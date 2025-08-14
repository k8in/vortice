package object

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	// Namespace represents a unique identifier for a specific domain or category within the system.
	Namespace string
	// Scope defines the lifecycle and visibility boundaries of a component within the system.
	Scope string
)

const (
	// Singleton 单例模式
	Singleton Scope = "Singleton"
	// Prototype 原型模式
	Prototype Scope = "Prototype"

	// Core 引擎命名空间
	Core Namespace = "Core"
)

// ErrParseDefinition is the error returned when there's a failure in parsing the definition.
var ErrParseDefinition = errors.New("failed to parse definition")

// Definition encapsulates the details of a component including its namespace, name, type,
// factory, dependencies, methods, scope, description, and lazy initialization flag.
type Definition struct {
	name      string
	typ       reflect.Type
	factory   *factory
	dependsOn []string
	methods   *Methods
	ns        Namespace
	scope     Scope
	desc      string
	lazyInit  bool
}

// Property represents a configuration property with namespace, scope, description,
// and lazy initialization flag.
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

// Parser defines an interface for parsing a function and property to produce
// a component definition.
type Parser interface {
	Parse(fn any, prop *Property) (*Definition, error)
}

// parser is a struct used for parsing and validating function definitions,
// ensuring they meet certain criteria.
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

// Parse initializes the parser and checks the input, returning a Definition or an error.
func (p *parser) Parse(fn any, prop *Property) (*Definition, error) {
	p.initFunc(fn)
	if err := p.checkInputAndSet(); err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	if err := p.checkOutputAndSet(); err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	return p.newDefinition(prop), nil
}

// newDefinition creates and returns a new Definition based on the parsed function and properties.
func (p *parser) newDefinition(prop *Property) *Definition {
	return &Definition{
		name:      generateReflectionName(p.obj),
		typ:       p.obj,
		factory:   newFactory(p.rv, p.argv, p.argn),
		dependsOn: p.deps,
		methods:   newMethods(p.obj),
		ns:        prop.Namespace,
		scope:     prop.Scope,
		desc:      prop.Desc,
		lazyInit:  prop.LazyInit,
	}
}

// init initializes the parser with the function and property
func (p *parser) initFunc(fn any) {
	rv := reflect.ValueOf(fn)
	p.fn = fn
	p.rt = rv.Type()
	p.rk = rv.Kind()
	p.rv = rv
}

// generateDefinitionName generates a unique name for the definition based
// on the namespace and argument type
func (p *parser) generateDefinitionName(ns Namespace, argType reflect.Type) string {
	return fmt.Sprintf("%s:%s", ns, generateReflectionName(argType))
}

// checkInputAndSet checks the input function and sets the argument values
func (p *parser) checkInputAndSet() error {
	if p.rk != reflect.Func {
		return errors.New("input must be a function")
	}
	for i := 0; i < p.rt.NumIn(); i++ {
		argType := p.rt.In(i)
		if err := p.checkArgType(argType); err != nil {
			return err
		}
		p.argn = p.rt.NumIn()
		p.argv = append(p.argv, reflect.ValueOf(argType))
		p.deps = append(p.deps, p.generateDefinitionName(Core, argType))
	}
	return nil
}

// checkOutputAndSet verifies the output of the function,
// ensuring it has exactly one return value and sets the object type.
func (p *parser) checkOutputAndSet() error {
	if p.rt.NumOut() != 1 {
		return errors.New("invalid factory function output")
	}
	p.obj = p.rt.Out(0)
	if err := p.checkReturnType(p.obj); err != nil {
		return err
	}
	return nil
}

// checkArgType checks the argument type and returns an error if it is invalid
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

// checkReturnType validates the return of a function,
// ensuring it is an interface, struct, or pointer to a struct.
func (p *parser) checkReturnType(rt reflect.Type) error {
	switch rt.Kind() {
	case reflect.Interface, reflect.Struct:
		return nil
	case reflect.Ptr:
		if rt.Elem().Kind() == reflect.Struct {
			return nil
		}
	default:
		return fmt.Errorf("invalid output type: %v", rt.String())
	}
	return nil
}

// generateReflectionName generates a string representation of the type name
func generateReflectionName(rt reflect.Type) string {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	name := rt.Name()
	if rt.Kind() == reflect.Struct {
		name = "*" + name
	}
	return rt.PkgPath() + "." + name
}
