package object

import (
	"fmt"
	"reflect"
)

type (
	// Initializable is an interface for objects that can be initialized,
	// setting up their initial state.
	Initializable interface {
		// Init initializes the object, setting up its initial state and preparing it for use.
		Init() error
	}
	// Destroyable defines an interface for objects that can be destroyed,
	// typically releasing resources or resetting state.
	Destroyable interface {
		// Destroy releases any resources held by the object and resets its state.
		Destroy() error
	}
	// Lifecycle defines the methods for managing the lifecycle of a component, including starting, stopping, and checking its running status.
	Lifecycle interface {
		// Start begins the operation of the component, initiating all necessary processes and services.
		Start() error
		// Stop stops the component, ensuring all resources are properly released and operations are halted.
		Stop() error
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
	obj           reflect.Type
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
		obj:           obj,
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

// CallInit invokes the initialization method on the provided reflect.Value
// and returns the result along with any error.
func (m *Methods) CallInit(ins reflect.Value) error {
	return m.simpleCall(ins, m.initMethod)
}

// CallDestroy invokes the destruction method on the provided reflect.Value
// and returns the result along with any error.
func (m *Methods) CallDestroy(ins reflect.Value) error {
	return m.simpleCall(ins, m.destroyMethod)
}

// CallStart invokes the start method on the provided reflect.Value
// and returns the result along with any error.
func (m *Methods) CallStart(ins reflect.Value) error {
	return m.simpleCall(ins, m.startMethod)
}

// CallStop invokes the stop method on the provided reflect.Value
// and returns the result along with any error.
func (m *Methods) CallStop(ins reflect.Value) error {
	return m.simpleCall(ins, m.stopMethod)
}

// CallRunning invokes the running method on the provided reflect.Value
// and returns the result along with any error.
func (m *Methods) CallRunning(ins reflect.Value) (bool, error) {
	ok, rv, err := m.call(ins, m.runningMethod)
	if err != nil || !ok {
		return false, err
	}
	if len(rv) != 1 {
		return false, fmt.Errorf("instance %s: method %s return value is invalid",
			reflect.TypeOf(ins), m.runningMethod.Name)
	}
	b, ok := rv[0].Interface().(bool)
	if !ok {
		return false, fmt.Errorf("instance %s: method %s return value is not bool",
			reflect.TypeOf(ins), m.runningMethod.Name)
	}
	return b, nil
}

// simpleCall invokes a method on the provided instance and returns an error
// if the first return value is of error type.
func (m *Methods) simpleCall(ins reflect.Value, method *reflect.Method) error {
	ok, rv, err := m.call(ins, method)
	if err != nil || !ok {
		return err
	}
	if len(rv) != 1 {
		return fmt.Errorf("instance %s: method %s return value is invalid",
			reflect.TypeOf(ins), method.Name)
	}
	if e, ok := rv[0].Interface().(error); ok && e != nil {
		return e
	}
	return nil
}

// call attempts to invoke a specified method on the given instance
// and returns the result along with any error.
func (m *Methods) call(ins reflect.Value, method *reflect.Method) (bool, []reflect.Value, error) {
	if method == nil {
		return false, nil, nil
	}
	if m := ins.MethodByName(method.Name); m.IsValid() {
		return true, m.Call([]reflect.Value{}), nil
	}
	return false, nil, fmt.Errorf("instance %s: method %s not found",
		reflect.TypeOf(ins), method.Name)
}
