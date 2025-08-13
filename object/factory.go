package object

import (
	"reflect"
)

type FactoryFunc0[T any] interface {
	~func() T
}

type FactoryFunc1[T, A any] interface {
	~func(A) T
}

type FactoryFunc2[T, A, B any] interface {
	~func(A, B) T
}

type FactoryFunc3[T, A, B, C any] interface {
	~func(A, B, C) T
}

type FactoryFunc4[T, A, B, C, D any] interface {
	~func(A, B, C, D) T
}

type FactoryFunc5[T, A, B, C, D, E any] interface {
	~func(A, B, C, D, E) T
}

type FactoryFunc6[T, A, B, C, D, E, F any] interface {
	~func(A, B, C, D, E, F) T
}

type factory struct {
	name string
	file string
	line int
	argV []reflect.Value
	argN int
	fn   reflect.Value
}
