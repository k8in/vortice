package object

import (
	"sync/atomic"
	"vortice/util"

	"go.uber.org/zap"
)

var (
	dr = newDefinitionRegistry()
)

// RegisterFactory registers a factory function with the given property,
// returning a new Definition.
func RegisterFactory(fn any, prop *Property) *Definition {
	parser := newParser(fn)
	def, err := parser.Parse(prop)
	if err != nil {
		util.Logger().Panic("failed to register factory: %v", zap.Error(err))
	}
	dr.register(def)
	return def
}

// LockDefinitionRegistry locks the DefinitionRegistry and logs the action,
// preventing further modifications.
func LockDefinitionRegistry() {
	dr.Lock()
	util.Logger().Info("The DefinitionRegistry has already been locked")
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

func (dr *DefinitionRegistry) register(def *Definition) {
	return
}

func (dr *DefinitionRegistry) GetDefinition(name string) []*Definition {
	return nil
}

func (dr *DefinitionRegistry) GetDefinitionNames() []string {
	return nil
}
