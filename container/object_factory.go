package container

import (
	"fmt"
	"reflect"

	"vortice/object"
	"vortice/util"
)

type ObjectFactory interface {
	Init() error
	Destroy() error
	GetObject(ns object.Namespace, typ any) ([]Object, error)
}

type CoreObjectFactory struct {
}

func (c *CoreObjectFactory) Init() error {
	//TODO implement me
	panic("implement me")
}

func (c *CoreObjectFactory) Destroy() error {
	//TODO implement me
	panic("implement me")
}

func (c *CoreObjectFactory) GetObject(ns object.Namespace, typ any) ([]Object, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoreObjectFactory) NewObject(def *object.Definition, ctx map[string]Object) (Object, error) {
	if def.Factory().Argn() == 0 {
		return c.new(def, ctx)
	}

	deps, err := c.getDependencies(def)
	if err != nil {
		return nil, err
	}

	//if def.ArgN() == 0 {
	//	rv := def.Factory().Call([]reflect.Value{})
	//	return component.NewObject(def, rv, rv.Interface())
	//}
	return nil, nil
}

// getDependencies retrieves and sorts the dependencies for a given object definition.
func (c *CoreObjectFactory) getDependencies(def *object.Definition) ([]string, error) {
	dag, deps := util.NewDAG(), def.DependsOn()
	dag.AddNode(def.Name(), deps...)
	queue := append([]string{}, deps...)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		def, err := c.getDefinition(node)
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
func (c *CoreObjectFactory) getDefinition(name string) (*object.Definition, error) {
	def := object.GetDefinitionRegistry().GetDefinition(name, object.NamespaceFilter(object.NSCore))
	if def == nil || len(def) == 0 {
		return nil, fmt.Errorf("definition not found: %s", name)
	}
	return def[0], nil
}

// new creates a new object based on the provided definition and context,
// handling dependencies and factory calls.
func (c *CoreObjectFactory) new(def *object.Definition, ctx map[string]Object) (Object, error) {
	if def.Factory().Argn() == 0 {
		rv := def.Factory().Call([]reflect.Value{})
		return NewObject(def, rv, rv.Interface()), nil
	}
	deps := def.DependsOn()
	argv := make([]reflect.Value, 0, len(deps))
	for _, dep := range deps {
		c, ok := ctx[dep]
		if !ok {
			return nil, fmt.Errorf("%s dependencies not found: %s", def, dep)
		}
		argv = append(argv, c.Value())
	}
	rv := def.Factory().Call(argv)
	return NewObject(def, rv, rv.Interface()), nil
}
