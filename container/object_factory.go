package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

var (
	// TagAutowired is a Tag used to mark components that should be automatically wired.
	TagAutowired = object.NewTag("autowired", "true")
	// autowiredFilter is a DefinitionFilter that matches Definitions tagged with TagAutowired for automatic wiring.
	autowiredFilter = object.TagFilter(TagAutowired)
	// realizationSelectFunc is a function that selects and returns the first Definition from a list, or nil if the list is empty.
	realizationSelectFunc RealizationSelectFunc = func(defs []*object.Definition) *object.Definition {
		if defs == nil || len(defs) == 0 {
			return nil
		}
		return defs[0]
	}
)

type (
	// RealizationSelectFunc is a function type that selects and returns a single Definition from a list of Definitions.
	RealizationSelectFunc func(defs []*object.Definition) *object.Definition
	// RealizationSelector is an interface for selecting a single Definition from a list of Definitions.
	RealizationSelector interface {
		Select(defs []*object.Definition) *object.Definition
	}
	// ObjectFactory is an interface for creating and managing objects, including initialization and destruction.
	ObjectFactory interface {
		// DefinitionRegistry is an interface for managing and retrieving definitions, including initialization and registration.
		object.DefinitionRegistry
		// GetObjects retrieves a list of objects of the specified type from the factory, using the provided context.
		GetObjects(ctx Context, typ any) ([]Object, error)
		// GetObjectsByName retrieves a list of objects by name from the factory, using the provided context.
		GetObjectsByName(ctx Context, name string) ([]Object, error)
		// SetRealizationSelector sets the AutowiredSelector to be used for selecting Definitions during auto-wiring.
		SetRealizationSelector(selector RealizationSelector)
		// Destroy cleans up resources and finalizes the ObjectFactory, returning an error if the operation fails.
		Destroy()
	}
)

// Select invokes the AutowiredSelectFunc with a list of Definitions and returns a single selected Definition.
func (fn RealizationSelectFunc) Select(defs []*object.Definition) *object.Definition {
	return fn(defs)
}

// CoreObjectFactory is a factory for creating core objects, equipped with definition filters to
// customize object creation.
type CoreObjectFactory struct {
	object.DefinitionRegistry
	once     *sync.Once
	mutex    *sync.RWMutex
	selector RealizationSelector
	objs     map[string]Object
}

// NewCoreObjectFactory creates a new instance of CoreObjectFactory with a namespace filter for the core namespace.
func NewCoreObjectFactory() *CoreObjectFactory {
	return &CoreObjectFactory{
		DefinitionRegistry: object.NewDefinitionRegistry(),
		once:               &sync.Once{},
		mutex:              &sync.RWMutex{},
		selector:           realizationSelectFunc,
		objs:               map[string]Object{},
	}
}

// SetRealizationSelector sets the autowiring selector for choosing a definition when multiple are available.
func (c *CoreObjectFactory) SetRealizationSelector(selector RealizationSelector) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.selector = selector
}

// Init initializes the CoreObjectFactory and its singleton objects, returning an error if any occurs.
func (c *CoreObjectFactory) Init() error {
	var err error
	c.once.Do(func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if err = c.DefinitionRegistry.Init(); err != nil {
			return
		}
		l := util.Logger()
		for _, def := range c.GetDefinitions(object.ScopeFilter(object.Singleton)) {
			// 修复: 避免使用 := 遮蔽外层 err
			var obj Object
			obj, err = c.newObject(def, map[string]Object{})
			if err != nil {
				err = fmt.Errorf("newObject failed: %s: %w", def.Name(), err)
				return
			}
			l.Debug("creating object", zap.String("definition", def.String()))
			if !def.LazyInit() {
				if err = obj.Init(); err != nil {
					err = fmt.Errorf("object.Init failed: %s: %w", def.Name(), err)
					return
				}
				l.Debug("object initialized", zap.String("definition", def.String()))
			}
			c.objs[def.ID()] = obj
		}
	})
	return err
}

// Destroy cleans up all created objects by calling their Destroy method, ensuring proper resource release.
func (c *CoreObjectFactory) Destroy() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, obj := range c.objs {
		if err := obj.Destroy(); err != nil {
			// TODO warning
		}
	}
}

// GetObjects retrieves and initializes objects of the specified type, returning them along with any error.
func (c *CoreObjectFactory) GetObjects(ctx Context, typ any) ([]Object, error) {
	defs, err := c.GetDefinitionsByType(typ, ctx.GetFilters()...)
	if err != nil {
		return nil, err
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.getObjects(defs, ctx.GetObjects())
}

// GetObjectsByName retrieves and initializes objects by name, returning them along with any error.
func (c *CoreObjectFactory) GetObjectsByName(ctx Context, name string) ([]Object, error) {
	defs := c.GetDefinitionsByName(name, ctx.GetFilters()...)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.getObjects(defs, ctx.GetObjects())
}

// getObjects processes a list of object definitions, creates and initializes the objects, and returns them.
func (c *CoreObjectFactory) getObjects(defs []*object.Definition, getCtx map[string]Object) ([]Object, error) {
	objs := []Object{}
	if defs == nil || len(defs) == 0 {
		return objs, nil
	}
	for _, def := range defs {
		var (
			obj Object
			err error
		)
		if def.IsSingleton() {
			if v, ok := c.objs[def.ID()]; ok {
				obj = v
			}
		}
		if obj == nil {
			obj, err = c.newObject(def, getCtx)
			if err != nil {
				return nil, errors.Join(errors.New("newObject failed"), err)
			}
		}
		if !obj.Initialized() {
			if err := obj.Init(); err != nil {
				return nil, errors.Join(errors.New("object.Init failed"), err)
			}
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

// NewObject creates a new object based on the provided definition and context, handling dependencies.
func (c *CoreObjectFactory) newObject(def *object.Definition, objs map[string]Object) (Object, error) {
	if def.Factory().Argn() == 0 {
		return c.new(def, objs)
	}
	deps, err := c.getDependencies(def)
	if err != nil {
		return nil, err
	}
	return c.buildObject(def, deps, c.getBuildCtx(objs))
}

// getBuildCtx returns a context map containing the provided objects and core tagged objects from the factory.
func (c *CoreObjectFactory) getBuildCtx(objs map[string]Object) map[string]Object {
	ctx := map[string]Object{}
	for k, v := range objs {
		ctx[k] = v
	}
	for _, v := range c.objs {
		def := v.Definition()
		for _, tag0 := range def.Tags() {
			if TagAutowired.Equals(tag0) {
				ctx[def.Name()] = v
			}
		}
	}
	return ctx
}

// buildObject constructs an object based on its definition, resolving and initializing its dependencies.
func (c *CoreObjectFactory) buildObject(def *object.Definition,
	deps []string, objs map[string]Object) (Object, error) {
	var obj Object
	for _, name := range deps {
		obj = nil
		def, err := c.getAutowiredDefinition(name)
		if err != nil {
			return nil, err
		}
		if def.IsSingleton() {
			obj, _ = c.objs[def.ID()]
		}
		if obj == nil {
			obj, err = c.new(def, objs)
			if err != nil {
				return nil, err
			}
		}
		objs[obj.ID()] = obj
	}
	obj, err := c.new(def, objs)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// getDependencies retrieves and sorts the dependencies for a given object definition.
func (c *CoreObjectFactory) getDependencies(def *object.Definition) ([]string, error) {
	dag, deps := util.NewDAG(), def.DependsOn()
	dag.AddNode(def.Name(), deps...)
	queue := append([]string{}, deps...)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		def, err := c.getAutowiredDefinition(node)
		if err != nil {
			return nil, err
		}
		deps := def.DependsOn()
		dag.AddNode(node, deps...)
		queue = append(queue, deps...)
	}
	sorted, err := dag.Sort()
	if err != nil {
		return nil, err
	}
	return sorted, nil
}

// getAutowiredDefinition retrieves a definition by name with auto-wired filter, returning an error if not found.
func (c *CoreObjectFactory) getAutowiredDefinition(name string) (*object.Definition, error) {
	defs := c.GetDefinitionsByName(name, autowiredFilter)
	if defs == nil || len(defs) == 0 {
		return nil, fmt.Errorf("definition not found: %s", name)
	}
	return c.selector.Select(defs), nil
}

// new creates a new object based on the provided definition and context,
// handling dependencies and factory calls.
func (c *CoreObjectFactory) new(def *object.Definition, objs map[string]Object) (Object, error) {
	if def.Factory().Argn() == 0 {
		rv := def.Factory().Call([]reflect.Value{})
		return NewObject(def, rv), nil
	}
	deps := def.DependsOn()
	argv := make([]reflect.Value, 0, len(deps))
	for _, dep := range deps {
		c, ok := objs[dep]
		if !ok {
			return nil, fmt.Errorf("%s dependencies not found: %s", def, dep)
		}
		argv = append(argv, c.Value())
	}
	rv := def.Factory().Call(argv)
	return NewObject(def, rv), nil
}
