package container

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"vortice/object"
)

// ErrAlreadyBeenDestroyed is the error returned when an operation is attempted on an already destroyed object.
var ErrAlreadyBeenDestroyed = errors.New("object has already been destroyed")

// Object defines the interface for a managed object in the container,
// including lifecycle, identity, and reflection access.
type Object interface {
	// ID returns the unique identifier for the object.
	ID() string
	// Definition returns the object's definition, which includes details like its type,
	// dependencies, and lifecycle methods.
	Definition() *object.Definition
	// Instance returns the actual instance of the object.
	Instance() any
	// Value returns the reflect.Value of the object, allowing for type inspection
	// and manipulation.
	Value() reflect.Value
	// Initialized returns true if the object has been initialized,
	// indicating its readiness for use.
	Initialized() bool
	// Initializable is an interface for objects that can be initialized,
	// setting up their initial state.
	object.Initializable
	// Destroyable defines an interface for objects that can be destroyed,
	// typically releasing resources or resetting state.
	object.Destroyable
	// Lifecycle defines the methods for managing the lifecycle of a component,
	// including starting, stopping, and checking its running status.
	object.Lifecycle
}

// CoreObject is the default implementation of Object.
// It wraps a Definition and its instance, and manages lifecycle and thread safety.
type CoreObject struct {
	mux   *sync.RWMutex      // protects all field access
	def   *object.Definition // definition metadata
	value reflect.Value      // reflect.Value of instance
	ins   any                // raw instance
	init  *atomic.Bool       // initialization flag
}

// NewObject creates a new Object with the given definition, reflect value, and instance.
func NewObject(def *object.Definition, rv reflect.Value, ins any) Object {
	return &CoreObject{
		def:   def,
		value: rv,
		ins:   ins,
		mux:   &sync.RWMutex{},
		init:  &atomic.Bool{},
	}
}

// Value returns the reflect.Value of the instance.
func (obj *CoreObject) Value() reflect.Value {
	obj.mux.RLock()
	defer obj.mux.RUnlock()
	return obj.value
}

// Initialized returns true if the object has been initialized.
func (obj *CoreObject) Initialized() bool {
	return obj.init.Load()
}

// Init initializes the object if not already initialized.
func (obj *CoreObject) Init() error {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	if obj.init.Load() {
		return nil
	}
	if obj.def == nil {
		return ErrAlreadyBeenDestroyed
	}
	if err := obj.def.Methods().CallInit(obj.value); err != nil {
		return err
	}
	obj.init.Store(true)
	return nil
}

// Destroy destroys the object and releases resources.
func (obj *CoreObject) Destroy() error {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	if obj.def == nil {
		return ErrAlreadyBeenDestroyed
	}
	if err := obj.def.Methods().CallDestroy(obj.value); err != nil {
		return err
	}
	obj.def = nil
	obj.value = reflect.Value{}
	obj.ins = nil
	return nil
}

// Running returns true if the object is currently running.
func (obj *CoreObject) Running() bool {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	if obj.def == nil {
		return false
	}
	b, err := obj.def.Methods().CallRunning(obj.value)
	if err != nil {
		return false
	}
	return b
}

// Start starts the object.
func (obj *CoreObject) Start() error {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	if obj.def == nil {
		return ErrAlreadyBeenDestroyed
	}
	return obj.def.Methods().CallStart(obj.value)
}

// Stop stops the object.
func (obj *CoreObject) Stop() error {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	if obj.def == nil {
		return ErrAlreadyBeenDestroyed
	}
	return obj.def.Methods().CallStop(obj.value)
}

// ID returns the object's name from its definition.
func (obj *CoreObject) ID() string {
	obj.mux.RLock()
	defer obj.mux.RUnlock()
	if obj.def == nil {
		return ""
	}
	return obj.def.Name()
}

// Definition returns the object's definition.
func (obj *CoreObject) Definition() *object.Definition {
	obj.mux.RLock()
	defer obj.mux.RUnlock()
	return obj.def
}

// Instance returns the raw instance.
func (obj *CoreObject) Instance() any {
	obj.mux.RLock()
	defer obj.mux.RUnlock()
	return obj.ins
}
