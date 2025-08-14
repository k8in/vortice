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

// factory represents a structure for creating and managing components,
// including their function and arguments.
type factory struct {
	name string
	file string
	line int
	fn   reflect.Value
	argv []reflect.Value
	argn int
}

// newFactory creates a new factory with the provided function, its arguments,
// and the number of expected arguments.
func newFactory(rfn reflect.Value, argv []reflect.Value, argn int) *factory {
	ptr := rfn.Pointer()
	fn := runtime.FuncForPC(ptr)
	file, line := fn.FileLine(ptr)
	return &factory{
		name: fn.Name(), file: file, line: line,
		fn: rfn, argv: argv, argn: argn,
	}
}
