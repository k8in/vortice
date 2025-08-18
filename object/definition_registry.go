package object

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	// ErrRegisterFactory is the error returned when a factory registration fails.
	ErrRegisterFactory = errors.New("failed to register factory")

	dr = newDefinitionRegistry()
)

// RegisterFactory registers a factory function with the given property,
// returning a new Definition.
func RegisterFactory(fn any, prop *Property) (*Definition, error) {
	parser := NewParser(fn)
	def, err := parser.Parse(prop)
	if err != nil {
		return nil, errors.Join(ErrRegisterFactory, err)
	}
	if err := dr.register(def); err != nil {
		return nil, errors.Join(ErrRegisterFactory, err)
	}
	return def, nil
}

// GetDefinitionRegistry returns the global DefinitionRegistry instance.
func GetDefinitionRegistry() *DefinitionRegistry {
	return dr
}

// DefinitionRegistry manages a collection of component definitions and their associated factories,
// supporting read-only state.
type DefinitionRegistry struct {
	entries   map[string][]*Definition
	factories map[string]string
	readonly  *atomic.Bool
}

// newDefinitionRegistry creates and returns a new DefinitionRegistry with
// an initial read-write state.
func newDefinitionRegistry() *DefinitionRegistry {
	readonly := &atomic.Bool{}
	readonly.Store(false)
	return &DefinitionRegistry{
		entries:   map[string][]*Definition{},
		factories: map[string]string{},
		readonly:  readonly,
	}
}

// Lock sets the DefinitionRegistry to a read-only state, preventing further modifications.
func (dr *DefinitionRegistry) Lock() {
	dr.readonly.Store(true)
}

// DefinitionFilter defines a filter function for Definition.
type DefinitionFilter func(*Definition) bool

// NamespaceFilter returns a DefinitionFilter that matches Definitions with the specified namespace.
func NamespaceFilter(ns Namespace) DefinitionFilter {
	return func(def *Definition) bool {
		return def.Namespace() == ns
	}
}

// TagFilter returns a DefinitionFilter that matches Definitions with the specified tag.
func TagFilter(match string) DefinitionFilter {
	return func(def *Definition) bool {
		for _, tag := range def.Tags() {
			if tag == match {
				return true
			}
		}
		return false
	}
}

// GetDefinition retrieves definitions by name, optionally applying filters to refine the results.
func (dr *DefinitionRegistry) GetDefinition(name string, filters ...DefinitionFilter) []*Definition {
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

// register adds a new Definition to the registry, ensuring no duplicate factory functions and not in read-only mode.
func (dr *DefinitionRegistry) register(def *Definition) error {
	if dr.readonly.Load() {
		return errors.New("the DefinitionRegistry has been locked")
	}
	fid := def.Factory().Name()
	if _, ok := dr.factories[fid]; ok {
		return fmt.Errorf("object's factory function %s already exists", fid)
	}
	dr.factories[fid] = def.Name()
	dr.entries[def.Name()] = append(dr.entries[def.Name()], def)
	return nil
}
