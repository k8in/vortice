package container

import (
	"fmt"
	"reflect"
	"sync"

	"vortice/object"
	"vortice/util"
)

// ObjectFactory is an interface for creating and managing objects, including initialization and destruction.
type ObjectFactory interface {
	// GetObject retrieves an object of the specified type within the given namespace, using the provided context.
	GetObject(ctx Context, typ any) (Object, error)
	// NewObject creates a new object based on the provided definition and context, returning the object and any error encountered.
	NewObject(ctx Context, def *object.Definition) (Object, error)
}

// CoreObjectFactory is a factory for creating core objects, equipped with definition filters to
// customize object creation.
type CoreObjectFactory struct {
	*sync.Mutex
	objs  map[string]Object
	comps []Object
}

// NewObjectFactory creates a new instance of CoreObjectFactory with a namespace filter for the core namespace.
func NewObjectFactory() ObjectFactory {
	return &CoreObjectFactory{
		Mutex: &sync.Mutex{},
		objs:  map[string]Object{},
		comps: []Object{},
	}
}

// GetObject retrieves an object based on the given namespace, type, and context, handling singleton scope and creation.
func (c *CoreObjectFactory) GetObject(ctx Context, typ any) (Object, error) {
	objType := c.getType(typ)
	if objType == nil {
		return nil, fmt.Errorf("object not found: %v", typ)
	}
	name := object.GenerateDefinitionName(objType)
	def, err := c.getDefinition(name, ctx.GetFilters()...)
	if err != nil {
		return nil, err
	}
	c.Lock()
	defer c.Unlock()
	if def.Scope() == object.Singleton {
		if obj, ok := c.objs[def.ID()]; ok {
			return obj, nil
		}
	}
	return c.NewObject(ctx, def)
}

// NewObject creates a new object based on the provided definition and context, handling dependencies.
func (c *CoreObjectFactory) NewObject(ctx Context, def *object.Definition) (Object, error) {
	if def.Factory().Argn() == 0 {
		return c.new(def, ctx.GetObjects())
	}
	deps, err := c.getDependencies(def)
	if err != nil {
		return nil, err
	}
	return c.build(ctx, def, deps)
}

// build constructs an object and its dependencies based on the provided definition, dependency list, and context.
func (c *CoreObjectFactory) build(ctx Context, def *object.Definition, deps []string) (Object, error) {
	cache := map[string]Object{}
	for k, v := range ctx.GetObjects() {
		cache[k] = v
	}
	var obj Object
	for _, name := range deps {
		def, err := c.getDefinition(name)
		if err != nil {
			return nil, err
		}
		obj = nil
		if def.Scope() == object.Singleton {
			obj, _ = c.objs[def.ID()]
		}
		if obj == nil {
			obj, err = c.new(def, ctx.GetObjects())
			if err != nil {
				return nil, err
			}
		}
		if !obj.Initialized() {
			if err = obj.Init(); err != nil {
				return nil, err
			}
		}
		cache[def.ID()] = obj
	}
	obj, _ = cache[def.ID()]
	return obj, nil
}

// getType returns the reflect.Type of the provided type if it is a pointer to an interface or struct, otherwise nil.
func (r *CoreObjectFactory) getType(typ any) reflect.Type {
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

// getDependencies retrieves and sorts the dependencies for a given object definition.
func (c *CoreObjectFactory) getDependencies(def *object.Definition, dfs ...object.DefinitionFilter) ([]string, error) {
	dag, deps := util.NewDAG(), def.DependsOn()
	dag.AddNode(def.Name(), deps...)
	queue := append([]string{}, deps...)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		def, err := c.getDefinition(node, dfs...)
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

// getDefinition retrieves a definition by name from the core namespace,
// returning an error if not found.
func (c *CoreObjectFactory) getDefinition(name string, dfs ...object.DefinitionFilter) (*object.Definition, error) {
	def := object.GetDefinitionRegistry().GetDefinition(name, dfs...)
	if def == nil || len(def) == 0 {
		return nil, fmt.Errorf("definition not found: %s", name)
	}
	return def[0], nil
}

// new creates a new object based on the provided definition and context,
// handling dependencies and factory calls.
func (c *CoreObjectFactory) new(def *object.Definition, objs map[string]Object) (Object, error) {
	if def.Factory().Argn() == 0 {
		rv := def.Factory().Call([]reflect.Value{})
		return NewObject(def, rv, rv.Interface()), nil
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
	return NewObject(def, rv, rv.Interface()), nil
}

func (c *CoreObjectFactory) init() error {
	//defNames := component.GetDefinitionNames()
	//r.mux.Lock()
	//defer r.mux.Unlock()
	//for _, name := range defNames {
	//	def := component.GetDefinition(name)
	//	if def.Scope() == component.SingletonScope {
	//		comp := r.newComponent(def)
	//		logger.Printf("%s created", def)
	//		if !def.LazyInit() {
	//			comp.Init()
	//		}
	//		r.singletons[name] = comp
	//	}
	//}
	//filter := object.ScopeFilter(object.Singleton)
	//defs := object.GetDefinitionRegistry().GetDefinitions(filter)
	//c.Lock()
	//defer c.Unlock()
	//for _, def := range defs {
	//	obj := c.new()
	//	if def.LazyInit() {
	//		continue
	//	}
	//}
	//return nil
	return nil
}

func (c *CoreObjectFactory) destroy() error {
	return nil
}
