package vortice

import (
	"context"
	"reflect"
	"vortice/container"
	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

// GetElem retrieves an object of the specified pointer type from the container within the given context.
/*
	type Service interface{}
	var srv Service =  GetElem((*Service)(nil))
*/
func GetElem[T any](ctx context.Context, typ *T) T {
	coreCtx := container.WithCoreContext(ctx)
	objs, err := container.DefaultCore().GetObjects(coreCtx, typ)
	if err != nil || len(objs) == 0 {
		return zeroVal(typ).(T)
	}
	return objs[0].Instance().(T)
}

// Get retrieves an object of the specified type from the container within the given context.
/*
		type Object struct{}
	    var obj *Object = Get((*Object)(nil))
*/
func Get[T any](ctx context.Context, typ T) T {
	coreCtx := container.WithCoreContext(ctx)
	objs, err := container.DefaultCore().GetObjects(coreCtx, typ)
	if err != nil || len(objs) == 0 {
		return zeroVal(typ).(T)
	}
	return objs[0].Instance().(T)
}

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
	if _, err := container.DefaultCore().RegisterFactory(fn, prop, true); err != nil {
		util.Logger().Panic("register", zap.Error(err))
	}
}

func zeroVal(typ any) any {
	rt := reflect.TypeOf(typ)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return reflect.Zero(rt).Interface()
}
