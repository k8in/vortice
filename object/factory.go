package object

import (
	"reflect"
	"runtime"
)

// FactoryFunc0 is a type that represents a function which takes no arguments
// and returns a value of type T.
type FactoryFunc0[T any] interface {
	~func() T
}

// FactoryFunc1 is a function type that takes one argument of type A
// and returns a value of type T.
type FactoryFunc1[T, A any] interface {
	~func(A) T
}

// FactoryFunc2 represents a function type that takes two arguments of types A and B,
// and returns a value of type T.
type FactoryFunc2[T, A, B any] interface {
	~func(A, B) T
}

// FactoryFunc3 is a type that represents a function capable of creating an instance of T given
// three arguments A, B, and C.
type FactoryFunc3[T, A, B, C any] interface {
	~func(A, B, C) T
}

// FactoryFunc4 represents a function type that creates an instance of T given
// four parameters of types A, B, C, and D.
type FactoryFunc4[T, A, B, C, D any] interface {
	~func(A, B, C, D) T
}

// FactoryFunc5 represents a function type that takes five arguments of types A, B, C, D,
// and E, and returns a value of type T.
type FactoryFunc5[T, A, B, C, D, E any] interface {
	~func(A, B, C, D, E) T
}

// FactoryFunc6 is a function type that creates an instance of T,
// given six input parameters of types A, B, C, D, E, and F.
type FactoryFunc6[T, A, B, C, D, E, F any] interface {
	~func(A, B, C, D, E, F) T
}

// Factory represents a structure for creating and managing components,
// including their function and arguments.
type Factory struct {
	name string
	file string
	line int
	fn   reflect.Value
	argv []reflect.Value
	argn int
}

// NewFactory creates a new factory with the provided function, its arguments,
// and the number of expected arguments.
func NewFactory(rfn reflect.Value, argv []reflect.Value, argn int) *Factory {
	ptr := rfn.Pointer()
	fn := runtime.FuncForPC(ptr)
	file, line := fn.FileLine(ptr)
	return &Factory{
		name: fn.Name(), file: file, line: line,
		fn: rfn, argv: argv, argn: argn,
	}
}

// Call invokes the Factory function with the provided arguments
// and returns the result.
func (f *Factory) Call(argv []reflect.Value) reflect.Value {
	return f.fn.Call(argv)[0]
}

// Name returns the name of the Factory function.
func (f *Factory) Name() string {
	return f.name
}

// File returns the file path where the Factory function is defined.
func (f *Factory) File() string {
	return f.file
}

// Line returns the line number where the Factory function is defined.
func (f *Factory) Line() int {
	return f.line
}

// Func returns the reflect.Value of the Factory function.
func (f *Factory) Func() reflect.Value {
	return f.fn
}

// Argv returns the argument values for the Factory function.
func (f *Factory) Argv() []reflect.Value {
	return f.argv
}

// Argn returns the number of arguments expected by the Factory function.
func (f *Factory) Argn() int {
	return f.argn
}
