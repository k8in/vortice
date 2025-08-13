package object

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	Namespace string
	Scope     string
)

const (
	// Singleton 单例模式
	Singleton Scope = "Singleton"
	// Prototype 原型模式
	Prototype Scope = "Prototype"

	// Core 引擎命名空间
	Core Namespace = "Core"
)

type Definition struct {
	namespace Namespace
	name      string
	typ       reflect.Type
	factory   *factory
	dependsOn []string
	methods   *Methods
	scope     Scope
	desc      string
	lazyInit  bool
}

// Property ...
type Property struct {
	Namespace Namespace
	Scope     Scope
	Desc      string
	LazyInit  bool
}

func NewProperty() *Property {
	return &Property{
		Namespace: Core,
		Scope:     Prototype,
		Desc:      "",
		LazyInit:  true,
	}
}

// Parser ...
type Parser interface {
	Parse(fn any, prop *Property) (*Definition, error)
}

var ErrParseDefinition = errors.New("failed to parse definition")

type parser struct {
	fn   any
	rv   reflect.Value
	rt   reflect.Type
	rk   reflect.Kind
	prop *Property
	argv []reflect.Value
	argn int
	deps []string
	obj  reflect.Type
}

func (p *parser) Parse(fn any, prop *Property) (*Definition, error) {
	p.init(fn, prop)
	if err := p.checkInputAndSet(); err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	return nil, nil
}

func (p *parser) init(fn any, prop *Property) {
	rv := reflect.ValueOf(fn)
	p.fn = fn
	p.rt = rv.Type()
	p.rk = rv.Kind()
	p.rv = rv
	p.prop = prop
}

func (p *parser) checkInputAndSet() error {
	if p.rk != reflect.Func {
		return errors.New("input must be a function")
	}
	for i := 0; i < p.rt.NumIn(); i++ {
		argType := p.rt.In(i)
		if err := p.checkArgType(argType); err != nil {
			return err
		}
		p.argv = append(p.argv, reflect.ValueOf(argType))
		// p.deps = append(p.deps, FullCompName(GME, dp.generateObjectID(argType)))
	}
	return nil
}

func (p *parser) checkArgType(rt reflect.Type) error {
	switch rt.Kind() {
	case reflect.Struct, reflect.Interface:
		return nil
	case reflect.Ptr:
		if rt.Elem().Kind() == reflect.Struct {
			return nil
		}
	default:
	}
	return fmt.Errorf("invalid argument type: %v", rt.String())
}
