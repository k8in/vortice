package vortice

import (
	"vortice/object"
)

/*
GetInterface extracts and returns the value from a pointer to an interface of any type.

	type Service interface{}
	var srv Service =  GetInterface((*Service)(nil))
*/
func GetInterface[T any](itf *T) T {
	return *itf
}

/*
GetStruct returns the input interface as its underlying struct type.

		type Object struct{}
	    var obj *Object = GetStruct((*Object)(nil))
*/
func GetStruct[T any](itf T) T {
	return itf
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
