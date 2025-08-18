package business

import (
	"context"
)

type ctxKey string

var (
	nsKey  = ctxKey("namespace")
	objKey = ctxKey("object")
)

type Context[O any] interface {
	Namespace() string
	Object() O
}

type ctx[O any] struct {
	context.Context
}

func (ctx *ctx[O]) Namespace() string {
	return ctx.Value(nsKey).(string)
}

func (ctx *ctx[O]) Object() O {
	return ctx.Value(objKey).(O)
}

func WithNamespace(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, nsKey, string(ns))
}

func WithObject[O any](c context.Context, obj O) Context[O] {
	return &ctx[O]{Context: context.WithValue(c, objKey, obj)}
}
