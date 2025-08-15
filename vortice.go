package vortice

import (
	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

/*
GetInterface extracts and returns the value from a pointer to an interface of any type.

	type Service interface{}
	var srv Service =  GetInterface((*Service)(nil))
*/
//func GetInterface[T any](itf *T) T {
//	return *itf
//}

/*
GetStruct returns the input interface as its underlying struct type.

		type Object struct{}
	    var obj *Object = GetStruct((*Object)(nil))
*/
//func GetStruct[T any](itf T) T {
//	return itf
//}

// Register0 registers a factory function that takes no arguments, with optional configuration options.
func Register0[T any, FN object.FactoryFunc0[T]](fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register1 registers a factory function with the object system, allowing for
// the creation of objects with type T from A.
func Register1[T, A any, FN object.FactoryFunc1[T, A]](fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register2 registers a factory function that takes two arguments and returns a value,
// with optional configuration.
func Register2[T, A, B any, FN object.FactoryFunc2[T, A, B]](fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register3 registers a factory function that creates an instance of T using three arguments A, B, and C.
func Register3[T, A, B, C any, FN object.FactoryFunc3[T, A, B, C]](fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register4 registers a factory function that creates an instance of T using four parameters, with optional configuration.
func Register4[T, A, B, C, D any, FN object.FactoryFunc4[T, A, B, C, D]](fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register5 registers a factory function that takes five arguments and returns a value, with optional configuration.
func Register5[T, A, B, C, D, E any, FN object.FactoryFunc5[T, A, B, C, D, E]](
	fn FN, opts ...Option) {
	register(fn, opts...)
}

// Register6 registers a factory function that creates an instance of T with six input parameters, applying given options.
func Register6[T, A, B, C, D, E, F any, FN object.FactoryFunc6[T, A, B, C, D, E, F]](
	fn FN, opts ...Option) {
	register(fn, opts...)
}

func register(fn any, opts ...Option) {
	prop := object.NewProperty()
	for _, option := range opts {
		option(prop)
	}
	if _, err := object.RegisterFactory(fn, prop); err != nil {
		util.Logger().Panic("register", zap.Error(err))
	}
}
