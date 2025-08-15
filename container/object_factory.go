package container

import (
	"vortice/object"
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

func (c *CoreObjectFactory) newObject(def *object.Definition) (Object, error) {
	if def.Factory().Argn() == 0 {

	}
	//if def.ArgN() == 0 {
	//	rv := def.Factory().Call([]reflect.Value{})
	//	return component.NewObject(def, rv, rv.Interface())
	//}
	return nil, nil
}

func (c *CoreObjectFactory) new(def *object.Definition, cache map[string]Object) (Object, error) {
	return nil, nil
}
