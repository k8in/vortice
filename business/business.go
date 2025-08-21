package business

import (
	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

// Option is a function type for configuring Property with functional options.
type Option object.Option

//// GetAbility retrieves the extension for a given ability type and target object within the provided context.
//func GetAbility[O Target, E Extension, T Ability[O, E]](ctx Context, typ T, obj O) E {
//	return nil
//}
//
//// GetAbilities retrieves the extension associated with a specific ability for a given object within a context.
//func GetAbilities[O Target, E Extension, T Ability[O, E]](ctx Context, typ T, obj O) []E {
//	return nil
//}

//// RegisterAbility registers a factory function for creating an Ability, with optional configuration.
//func RegisterAbility[O Target, E Extension, T Ability[O, E], FN AbilityFactory[O, E, T]](fn FN, opts ...Option) {
//	RegisterExtN(fn, opts...)
//}

// RegisterPlugin registers a plugin with the default core, panicking if an error occurs.
func RegisterPlugin(plugin *Plugin) {
	if err := DefaultCore().RegisterPlugin(plugin); err != nil {
		util.Logger().Panic("RegisterPlugin", zap.Error(err))
	}
}

// RegisterExt0 registers a factory function that takes no arguments and returns a value of type T, with optional configuration.
func RegisterExt0[T any, FN object.FactoryFunc0[T]](fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt1 registers a factory function with the system, allowing for creation of type T from type A, with optional configuration.
func RegisterExt1[T, A any, FN object.FactoryFunc1[T, A]](fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt2 registers a factory function with two arguments and optional configuration options.
func RegisterExt2[T, A, B any, FN object.FactoryFunc2[T, A, B]](fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt3 registers a factory function that creates an instance of T with three arguments A, B, and C.
func RegisterExt3[T, A, B, C any, FN object.FactoryFunc3[T, A, B, C]](fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt4 registers a factory function that creates an instance of T using four parameters.
func RegisterExt4[T, A, B, C, D any, FN object.FactoryFunc4[T, A, B, C, D]](fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt5 registers a factory function with five parameters, using provided options for configuration.
func RegisterExt5[T, A, B, C, D, E any, FN object.FactoryFunc5[T, A, B, C, D, E]](
	fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExt6 registers a factory function that creates an instance of T with six input parameters, using provided options.
func RegisterExt6[T, A, B, C, D, E, F any, FN object.FactoryFunc6[T, A, B, C, D, E, F]](
	fn FN, opts ...Option) {
	RegisterExtN(fn, opts...)
}

// RegisterExtN adds a factory function to the system with given options and sets an extension tag.
func RegisterExtN(fn any, opts ...Option) {
	prop := object.NewProperty()
	for _, option := range opts {
		option(prop)
	}
	if _, err := DefaultCore().RegisterExtension(fn, prop); err != nil {
		util.Logger().Panic("RegisterExtension", zap.Error(err))
	}
}
