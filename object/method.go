package object

import (
	"reflect"
)

type (
	// Initializable is an interface for objects that can be initialized,
	// setting up their initial state.
	Initializable interface {
		// Init initializes the object, setting up its initial state and preparing it for use.
		Init()
	}

	// Destroyable defines an interface for objects that can be destroyed,
	// typically releasing resources or resetting state.
	Destroyable interface {
		// Destroy releases any resources held by the object and resets its state.
		Destroy()
	}
	// Lifecycle defines the methods for managing the lifecycle of a component, including starting, stopping, and checking its running status.
	Lifecycle interface {
		// Start begins the operation of the component, initiating all necessary processes and services.
		Start()
		// Stop stops the component, ensuring all resources are properly released and operations are halted.
		Stop()
		// Running returns true if the component is currently running, otherwise false.
		Running() bool
	}
)

var (
	initMethodType    = reflect.TypeOf((*Initializable)(nil)).Elem()
	destroyMethodType = reflect.TypeOf((*Destroyable)(nil)).Elem()
	lifecycleType     = reflect.TypeOf((*Lifecycle)(nil)).Elem()

	initMethodName    = "Init"
	destroyMethodName = "Destroy"
	startMethodName   = "Start"
	stopMethodName    = "Stop"
	runningMethod     = "Running"
)

// Methods holds method pointers for initialization, destruction,
// and lifecycle management of a component.
type Methods struct {
	initMethod    *reflect.Method
	destroyMethod *reflect.Method
	startMethod   *reflect.Method
	stopMethod    *reflect.Method
	runningMethod *reflect.Method
}

// newMethods initializes and returns a Methods struct with methods
// for initialization, destruction, and lifecycle management.
func newMethods(obj reflect.Type) *Methods {
	return &Methods{
		initMethod:    newMethod(obj, initMethodType, initMethodName),
		destroyMethod: newMethod(obj, destroyMethodType, destroyMethodName),
		startMethod:   newMethod(obj, lifecycleType, startMethodName),
		stopMethod:    newMethod(obj, lifecycleType, stopMethodName),
		runningMethod: newMethod(obj, lifecycleType, runningMethod),
	}
}

// newMethod retrieves a method from a given type if it implements the specified interface
// and contains the named method.
func newMethod(obj reflect.Type, iter reflect.Type, method string) *reflect.Method {
	if obj.Implements(iter) {
		if method, ok := obj.MethodByName(method); ok {
			return &method
		}
	}
	return nil
}

// CallInit invokes the initialization method on the provided reflect.Value.
func (m *Methods) CallInit(rv reflect.Value) {
	m.call(rv, m.initMethod)
}

// CallDestroy invokes the destruction method on the provided reflect.Value.
func (m *Methods) CallDestroy(rv reflect.Value) {
	m.call(rv, m.destroyMethod)
}

// CallStart invokes the start method on the provided reflect.Value.
func (m *Methods) CallStart(rv reflect.Value) {
	m.call(rv, m.startMethod)
}

// CallStop invokes the stop method on the provided reflect.Value.
func (m *Methods) CallStop(rv reflect.Value) {
	m.call(rv, m.stopMethod)
}

// CallRunning invokes the running method on the provided reflect.Value.
func (m *Methods) CallRunning(rv reflect.Value) {
	m.call(rv, m.runningMethod)
}

// call invokes the specified method on the provided reflect.Value if the method is not nil.
func (m *Methods) call(rv reflect.Value, method *reflect.Method) {
	if method != nil {
		rv.MethodByName(method.Name).Call([]reflect.Value{})
	}
}
