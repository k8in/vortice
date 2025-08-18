package object

import (
	"errors"
	"fmt"
	"sync/atomic"

	"vortice/util"

	"go.uber.org/zap"
)

var (
	// ErrRegisterFactory is the error returned when a factory registration fails.
	ErrRegisterFactory = errors.New("failed to register factory")

	dr = newDefinitionRegistry()
)

// RegisterFactory registers a factory function with the given property,
// returning a new Definition.
func RegisterFactory(fn any, prop *Property, unique bool) (*Definition, error) {
	parser := NewParser(fn)
	def, err := parser.Parse(prop)
	if err != nil {
		return nil, errors.Join(ErrRegisterFactory, err)
	}
	if err := dr.register(def, unique); err != nil {
		return nil, errors.Join(ErrRegisterFactory, err)
	}
	return def, nil
}

// GetDefinitionRegistry returns the global DefinitionRegistry instance.
func GetDefinitionRegistry() *DefinitionRegistry {
	return dr
}

// DefinitionFilter defines a filter function for Definition.
type DefinitionFilter func(*Definition) bool

// ScopeFilter returns a DefinitionFilter that matches Definitions with the specified scope.
func ScopeFilter(scope Scope) DefinitionFilter {
	return func(def *Definition) bool {
		return def.scope == scope
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

// DefinitionRegistry manages a collection of component definitions and their associated factories,
// supporting read-only state.
type DefinitionRegistry struct {
	readonly  *atomic.Bool
	entries   map[string][]*Definition
	factories map[string]*Definition
	inSeq     []string
}

// newDefinitionRegistry creates and returns a new DefinitionRegistry with
// an initial read-write state.
func newDefinitionRegistry() *DefinitionRegistry {
	readonly := &atomic.Bool{}
	readonly.Store(false)
	return &DefinitionRegistry{
		readonly:  readonly,
		entries:   map[string][]*Definition{},
		factories: map[string]*Definition{},
		inSeq:     []string{},
	}
}

// Init locks the DefinitionRegistry, sorts and checks definitions for cycles, and prepares it for use.
func (dr *DefinitionRegistry) Init() {
	dr.readonly.Store(true)
	util.Logger().Info("The DefinitionRegistry has already been locked")
	if err := dr.sortAndCheck(); err != nil {
		util.Logger().Panic("init", zap.Error(err))
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

// GetDefinitions returns a list of definitions that match all the provided filters.
func (dr *DefinitionRegistry) GetDefinitions(filters ...DefinitionFilter) []*Definition {
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

// register adds a new Definition to the registry, ensuring it's unique if required and not in read-only mode.
func (dr *DefinitionRegistry) register(def *Definition, unique bool) error {
	if dr.readonly.Load() {
		return errors.New("the DefinitionRegistry has been locked")
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
理论上此时所有的组件依赖项应已被注册，并且彼此没有依赖环（构造注入模式）

系统初始化阶段要提前检测：
1）依赖组件忘记注册
2）存在依赖环（GO 包依赖机制在编译期会检查这种情况，但是包内依赖环无法避免）

	type ServiceA interface {
		New(ServiceB)
	}

	type ServiceB interface {
		New(ServiceA)
	}

// 为了简化系统设计，不允许依赖环这种代码设计
*/
func (dr *DefinitionRegistry) sortAndCheck() error {
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
		// 检查是否��注册
		defs, ok := dr.entries[name]
		if !ok {
			return fmt.Errorf("%s validation failed: definition not found", name)
		}
		for _, def := range defs {
			util.Logger().Info("validation passed",
				zap.String("definition", def.Name()),
				zap.String("factory", def.Factory().Name()),
				zap.Int("dependsOn", len(def.DependsOn())))
			inSeq = append(inSeq, def.Factory().Name())
		}
	}
	dr.inSeq = inSeq
	return nil
}
