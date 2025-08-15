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

// LockDefinitionRegistry locks the DefinitionRegistry and logs the action,
// preventing further modifications.
func LockDefinitionRegistry() {
	dr.Lock()
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

// GetDefinition returns all definitions with the given name, optionally filtered by filter.
// If filter is nil, returns all definitions for the name.
func (dr *DefinitionRegistry) GetDefinition(name string, filter DefinitionFilter) []*Definition {
	defs := dr.entries[name]
	if filter == nil {
		return defs
	}
	var result []*Definition
	for _, def := range defs {
		if filter(def) {
			result = append(result, def)
		}
	}
	return result
}

// GetDefinitionNames returns all definition names, optionally filtered by filter.
// If filter is nil, returns all names.
func (dr *DefinitionRegistry) GetDefinitionNames(filter DefinitionFilter) []string {
	var names []string
	for name, defs := range dr.entries {
		for _, def := range defs {
			if filter == nil || filter(def) {
				names = append(names, name)
				break // only add name once
			}
		}
	}
	return names
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
