package object

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"

	"vortice/util"

	"go.uber.org/zap"
)

var (
	// ErrParseDefinition is the error returned when there's a failure in parsing a definition.
	ErrParseDefinition = errors.New("failed to parse definition")
)

type (
	// DefinitionFilter defines a filter function for Definition.
	DefinitionFilter func(*Definition) bool
	// DefinitionRegistry is an interface for managing and retrieving definitions.
	DefinitionRegistry interface {
		// Init initializes the DefinitionRegistry, preparing it for use and returns an error if initialization fails.
		Init() error
		// RegisterFactory registers a factory function with the given property and returns a new Definition,
		// or an error if registration fails.
		RegisterFactory(fn any, prop *Property, unique bool) (*Definition, error)
		// GetDefinitions returns a list of all Definitions, optionally filtered by the provided DefinitionFilter functions.
		GetDefinitions(filters ...DefinitionFilter) []*Definition
		// GetDefinitionsByName retrieves a list of Definitions by name, optionally filtered by the provided DefinitionFilter functions.
		GetDefinitionsByName(name string, filters ...DefinitionFilter) []*Definition
		// GetDefinitionsByType retrieves a list of Definitions matching the given type, optionally filtered by provided filters.
		GetDefinitionsByType(typ any, filters ...DefinitionFilter) ([]*Definition, error)
	}
)

// ScopeFilter returns a DefinitionFilter that matches Definitions with the specified scope.
func ScopeFilter(scope Scope) DefinitionFilter {
	return func(def *Definition) bool {
		return def.scope == scope
	}
}

// TagFilter returns a DefinitionFilter that matches Definitions with the specified tag.
func TagFilter(tags ...Tag) DefinitionFilter {
	return func(def *Definition) bool {
		if tags == nil || len(tags) == 0 {
			return false
		}
		for _, tag0 := range def.Tags() {
			for _, tag1 := range tags {
				if tag0.Equals(tag1) {
					return true
				}
			}
		}
		return false
	}
}

// DefaultDefRegistry manages a collection of component definitions and their associated factories,
// supporting read-only state.
type DefaultDefRegistry struct {
	readonly  *atomic.Bool
	entries   map[string][]*Definition
	factories map[string]*Definition
	inSeq     []string
}

// NewDefinitionRegistry creates and returns a new DefinitionRegistry with
// an initial read-write state.
func NewDefinitionRegistry() *DefaultDefRegistry {
	readonly := &atomic.Bool{}
	readonly.Store(false)
	return &DefaultDefRegistry{
		readonly:  readonly,
		entries:   map[string][]*Definition{},
		factories: map[string]*Definition{},
		inSeq:     []string{},
	}
}

// RegisterFactory registers a factory function with the given property, returning a new Definition.
func (dr *DefaultDefRegistry) RegisterFactory(fn any, prop *Property, unique bool) (*Definition, error) {
	def, err := ParseDefinition(fn, prop)
	if err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	if err := dr.register(def, unique); err != nil {
		errInnerRegister := fmt.Errorf("failed to register definition internally: %s", def.ID())
		return nil, errors.Join(errInnerRegister, err)
	}
	return def, nil
}

// Init locks the DefinitionRegistry, sorts and checks for circular dependencies, then logs the process.
func (dr *DefaultDefRegistry) Init() error {
	dr.readonly.Store(true)
	util.Logger().Info("the DefinitionRegistry has been locked")
	if err := dr.sortAndCheck(); err != nil {
		return err
	}
	return nil
}

// GetDefinitions returns a list of definitions that match all the provided filters.
func (dr *DefaultDefRegistry) GetDefinitions(filters ...DefinitionFilter) []*Definition {
	var result []*Definition
	for _, def := range dr.factories {
		matched := true
		for _, filter := range filters {
			if filter != nil && !filter(def) {
				matched = false
				break
			}
		}
		if matched {
			result = append(result, def)
		}
	}
	return result
}

// GetDefinitionsByName retrieves definitions by name, optionally filtered by provided DefinitionFilter functions.
func (dr *DefaultDefRegistry) GetDefinitionsByName(name string, filters ...DefinitionFilter) []*Definition {
	defs := dr.entries[name]
	if len(filters) == 0 {
		return defs
	}
	var result []*Definition
	for _, def := range defs {
		matched := true
		for _, filter := range filters {
			if filter != nil && !filter(def) {
				matched = false
				break
			}
		}
		if matched {
			result = append(result, def)
		}
	}
	return result
}

// GetDefinitionsByType retrieves definitions by type, optionally filtered by provided DefinitionFilter functions.
func (dr *DefaultDefRegistry) GetDefinitionsByType(typ any, filters ...DefinitionFilter) ([]*Definition, error) {
	objType := dr.getObjectType(typ)
	if objType == nil {
		return nil, errors.Join(errors.New("getObjectType failed"),
			fmt.Errorf("invalid object type: %#v", typ))
	}
	defs := dr.GetDefinitionsByName(generateReflectionName(objType), filters...)
	if defs == nil || len(defs) == 0 {
		return nil, fmt.Errorf("object type not found: %v", typ)
	}
	return defs, nil
}

// getType returns the reflect.Type of the provided type if it is a pointer to an interface or struct, otherwise nil.
func (r *DefaultDefRegistry) getObjectType(typ any) reflect.Type {
	if typ == nil {
		return nil
	}
	rt := reflect.TypeOf(typ)
	if rt.Kind() != reflect.Ptr {
		return nil
	}
	rek := rt.Elem().Kind()
	if rek != reflect.Interface && rek != reflect.Struct {
		return nil
	}
	return rt
}

// register adds a new Definition to the registry, ensuring it's unique if required and not in read-only mode.
func (dr *DefaultDefRegistry) register(def *Definition, unique bool) error {
	if dr.readonly.Load() {
		return errors.New("the DefinitionRegistry has been locked")
	}
	if !def.IsValid() {
		return errors.New("definition not valid")
	}
	if unique {
		if _, ok := dr.entries[def.Name()]; ok {
			return fmt.Errorf("object type %s does not allow duplicate definition", def.Name())
		}
	}
	fid := def.Factory().Name()
	if _, ok := dr.factories[fid]; ok {
		return fmt.Errorf("definition's factory function %s already exists", fid)
	}
	dr.factories[fid] = def
	dr.entries[def.Name()] = append(dr.entries[def.Name()], def)
	dr.inSeq = append(dr.inSeq, fid)
	return nil
}

/*
In theory, at this point all component dependencies should have been registered, and there should be
no dependency cycles (constructor injection mode).

During system initialization, the following should be checked in advance:

 1. Dependent components forgotten to register

 2. Existence of dependency cycles (GO package dependency mechanism checks this at compile time,
    but intra-package cycles cannot be avoided)

    type ServiceA interface {
    New(ServiceB)
    }

    type ServiceB interface {
    New(ServiceA)
    }

To simplify system design, code with dependency cycles is not allowed.
*/
func (dr *DefaultDefRegistry) sortAndCheck() error {
	dag, defs := util.NewDAG(), dr.factories
	for _, def := range defs {
		dag.AddNode(def.Name(), def.DependsOn()...)
	}
	sorted, err := dag.Sort()
	if err != nil {
		return err
	}
	inSeq := []string{}
	for _, name := range sorted {
		defs, ok := dr.entries[name]
		if !ok {
			err := fmt.Errorf("definition not found: %s", name)
			util.Logger().Error("validation failed", zap.String("name", name), zap.Error(err))
			return err
		}
		for _, def := range defs {
			util.Logger().Debug("validation passed",
				zap.String("name", def.Name()),
				zap.String("factory", def.Factory().Name()),
				zap.Int("dependsOn", len(def.DependsOn())),
				zap.Bool("isSingleton", def.IsSingleton()),
				zap.String("desc", def.Desc()),
				zap.Bool("lazyInit", def.LazyInit()),
				zap.Bool("autoStartup", def.AutoStartup()),
				zap.Any("tags", def.Tags()))
			inSeq = append(inSeq, def.Factory().Name())
		}
	}
	dr.inSeq = inSeq
	return nil
}
