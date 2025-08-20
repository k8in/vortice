package business

import (
	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

// GetAbilities retrieves the extension associated with a specific ability for a given object within a context.
func GetAbilities[O Target, E Extension, T Ability[O, E]](ctx Context, typ T, obj O) []E {
	return nil
}

// RegisterAbility registers a factory function for creating an Ability, with optional configuration.
func RegisterAbility[O Target, E Extension, T Ability[O, E], FN AbilityFactory[O, E, T]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register0 registers a factory function that takes no arguments and returns a value of type T, with optional configuration.
func Register0[T any, FN object.FactoryFunc0[T]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register1 registers a factory function with the system, allowing for creation of type T from type A, with optional configuration.
func Register1[T, A any, FN object.FactoryFunc1[T, A]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register2 registers a factory function with two arguments and optional configuration options.
func Register2[T, A, B any, FN object.FactoryFunc2[T, A, B]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register3 registers a factory function that creates an instance of T with three arguments A, B, and C.
func Register3[T, A, B, C any, FN object.FactoryFunc3[T, A, B, C]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register4 registers a factory function that creates an instance of T using four parameters.
func Register4[T, A, B, C, D any, FN object.FactoryFunc4[T, A, B, C, D]](fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register5 registers a factory function with five parameters, using provided options for configuration.
func Register5[T, A, B, C, D, E any, FN object.FactoryFunc5[T, A, B, C, D, E]](
	fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// Register6 registers a factory function that creates an instance of T with six input parameters, using provided options.
func Register6[T, A, B, C, D, E, F any, FN object.FactoryFunc6[T, A, B, C, D, E, F]](
	fn FN, opts ...Option) {
	RegisterN(fn, opts...)
}

// RegisterN adds a factory function to the system with given options and sets an extension tag.
func RegisterN(fn any, opts ...Option) {
	prop := object.NewProperty()
	for _, option := range opts {
		option(prop)
	}
	if _, err := DefaultCore().RegisterExtension(fn, prop); err != nil {
		util.Logger().Panic("RegisterExtension", zap.Error(err))
	}
}
