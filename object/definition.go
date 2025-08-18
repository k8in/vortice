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
	// Singleton defines a scope where the component is instantiated once and shared throughout its lifecycle.
	Singleton Scope = "Singleton"
	// Prototype indicates that a component should be instantiated each time it is requested.
	Prototype Scope = "Prototype"
	// NSCore represents the core namespace used to identify the main or default category within the system.
	NSCore Namespace = "Core"
)

// ErrParseDefinition is the error returned when there's a failure in parsing the definition.
var ErrParseDefinition = errors.New("failed to parse definition")

// Definition encapsulates the details of a component including its namespace, name, type,
// factory, dependencies, methods, scope, description, and lazy initialization flag.
type Definition struct {
	name        string
	typ         reflect.Type
	factory     *Factory
	dependsOn   []string // dependsOn holds the list of component names that this component depends on.
	methods     *Methods
	ns          Namespace
	scope       Scope
	desc        string
	lazyInit    bool
	autoStartup bool
	tags        []string // tags holds a list of string tags associated with the component definition.
}

// IsValid checks if the Definition is valid, ensuring name, type, factory, dependsOn, and methods are set, and tags are not.
func (d *Definition) IsValid() bool {
	return d.name != "" && d.typ != nil && d.factory != nil &&
		d.dependsOn != nil && d.methods != nil && d.tags != nil
}

// Name returns the name of the component definition.
func (d *Definition) Name() string {
	return d.name
}

// Type returns the reflect.Type of the component associated with the Definition.
func (d *Definition) Type() reflect.Type {
	return d.typ
}

// Factory returns the factory associated with the component definition.
func (d *Definition) Factory() *Factory {
	return d.factory
}

// DependsOn returns a copy of the list of dependencies for the component.
// Always returns a non-nil slice (at least empty).
func (d *Definition) DependsOn() []string {
	if d.dependsOn == nil {
		return []string{}
	}
	deps := make([]string, len(d.dependsOn))
	copy(deps, d.dependsOn)
	return deps
}

// Methods returns the lifecycle and method information for the component.
func (d *Definition) Methods() *Methods {
	return d.methods
}

// Namespace returns the namespace of the component definition.
func (d *Definition) Namespace() Namespace {
	return d.ns
}

// Scope returns the scope of the component definition.
func (d *Definition) Scope() Scope {
	return d.scope
}

// Desc returns the description of the component definition.
func (d *Definition) Desc() string {
	return d.desc
}

// LazyInit returns whether the component should be lazily initialized.
func (d *Definition) LazyInit() bool {
	return d.lazyInit
}

// AutoStartup returns whether the component should automatically start up.
func (d *Definition) AutoStartup() bool {
	return d.autoStartup
}

// Tags returns a copy of the tags for the component definition.
// Always returns a non-nil slice (at least empty).
func (d *Definition) Tags() []string {
	if d.tags == nil {
		return []string{}
	}
	tags := make([]string, len(d.tags))
	copy(tags, d.tags)
	return tags
}

// String returns a string representation of the Definition,
// including its namespace, name, and tags.
func (d *Definition) String() string {
	return fmt.Sprintf("%s<%s %s %v>", d.ns, d.Name(), d.typ.Kind(), d.Tags())
}

// Property represents a configuration property with namespace, scope, description,
// and lazy initialization flag.
type Property struct {
	Namespace   Namespace
	Scope       Scope
	Desc        string
	LazyInit    bool
	AutoStartup bool
	tags        []string
}

// NewProperty creates a new Property instance with default values,
// including Core namespace and Prototype scope.
func NewProperty() *Property {
	return &Property{
		Namespace:   NSCore,
		Scope:       Singleton,
		Desc:        "",
		LazyInit:    true,
		AutoStartup: false,
		tags:        []string{},
	}
}

// SetTag adds a new key-value pair to the property's tags.
func (prop *Property) SetTag(key, val string) {
	prop.tags = append(prop.tags, key+"="+val)
}

// GetTags returns a copy of the tags associated with the property.
func (prop *Property) GetTags() []string {
	tags := make([]string, len(prop.tags))
	copy(tags, prop.tags)
	return tags
}

// Option is a function type for configuring Property with functional options.
type Option func(prop *Property)

// Parser is a struct used for parsing and validating function definitions,
// ensuring they meet certain criteria.
type Parser struct {
	fn   any
	rv   reflect.Value
	rt   reflect.Type
	rk   reflect.Kind
	argv []reflect.Value
	argn int
	deps []string
	obj  reflect.Type
}

// NewParser initializes a new parser instance for a given function,
// setting up reflection-based properties and dependencies.
func NewParser(fn any) *Parser {
	rv := reflect.ValueOf(fn)
	p := &Parser{
		fn:   fn,
		rt:   rv.Type(),
		rk:   rv.Kind(),
		rv:   rv,
		argv: []reflect.Value{},
		deps: []string{},
	}
	return p
}

// Parse initializes the parser and checks the input, returning a Definition or an error.
func (p *Parser) Parse(prop *Property) (*Definition, error) {
	if err := p.checkInputAndSet(); err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	if err := p.checkOutputAndSet(); err != nil {
		return nil, errors.Join(ErrParseDefinition, err)
	}
	def := p.newDefinition(prop)
	if !def.IsValid() {
		return nil, errors.Join(ErrParseDefinition, errors.New("invalid definition"))
	}
	return def, nil
}

// newDefinition creates and returns a new Definition based on the parsed function and properties.
func (p *Parser) newDefinition(prop *Property) *Definition {
	return &Definition{
		name:        generateReflectionName(p.obj),
		typ:         p.rt,
		factory:     newFactory(p.rv, p.argv, p.argn),
		dependsOn:   p.deps,
		methods:     newMethods(p.obj),
		ns:          prop.Namespace,
		scope:       prop.Scope,
		lazyInit:    prop.LazyInit,
		autoStartup: prop.AutoStartup,
		tags:        prop.GetTags(),
	}
}

// generateDefinitionName generates a unique name for the definition based
// on the namespace and argument type
func (p *Parser) generateDefinitionName(ns Namespace, argType reflect.Type) string {
	return fmt.Sprintf("%s:%s", ns, generateReflectionName(argType))
}

// checkInputAndSet checks the input function and sets the argument values
func (p *Parser) checkInputAndSet() error {
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
		p.deps = append(p.deps, p.generateDefinitionName(NSCore, argType))
	}
	return nil
}

// checkOutputAndSet verifies the output of the function,
// ensuring it has exactly one return value and sets the object type.
func (p *Parser) checkOutputAndSet() error {
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
func (p *Parser) checkArgType(rt reflect.Type) error {
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
func (p *Parser) checkReturnType(rt reflect.Type) error {
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
