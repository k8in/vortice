package business

import (
	"vortice/object"
)

type Option object.Option

func GetAbility[E, O any, C Context[O]](ctx C, ext E) E {
	return ext
}

func Register0[T any, FN object.FactoryFunc0[T]](fn FN, opts ...Option) {
	// 实现略
}

func Register1[T, A any, FN object.FactoryFunc1[T, A]](fn FN, opts ...Option) {
	// 实现略
}

func Register2[T, A, B any, FN object.FactoryFunc2[T, A, B]](fn FN, opts ...Option) {
	// 实现略
}

func Register3[T, A, B, C any, FN object.FactoryFunc3[T, A, B, C]](fn FN, opts ...Option) {
	// 实现略
}

func Register4[T, A, B, C, D any, FN object.FactoryFunc4[T, A, B, C, D]](fn FN, opts ...Option) {
	// 实现略
}

func Register5[T, A, B, C, D, E any, FN object.FactoryFunc5[T, A, B, C, D, E]](
	fn FN, opts ...Option) {
	// 实现略
}

func Register6[T, A, B, C, D, E, F any, FN object.FactoryFunc6[T, A, B, C, D, E, F]](
	fn FN, opts ...Option) {
	// 实现略
}
