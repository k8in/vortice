package business

import (
	"context"
	"vortice/object"
)

type ctxKey string

var (
	nsKey  = ctxKey("namespace")
	objKey = ctxKey("object")
)

type Context[O any] interface {
	Namespace() object.Namespace
	Object() O
}

type ctx[O any] struct {
	context.Context
}

func (ctx *ctx[O]) Namespace() object.Namespace {
	return ctx.Value(nsKey).(object.Namespace)
}

func (ctx *ctx[O]) Object() O {
	return ctx.Value(objKey).(O)
}

func WithNamespace(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, nsKey, object.Namespace(ns))
}

func WithObject[O any](c context.Context, obj O) Context[O] {
	return &ctx[O]{Context: context.WithValue(c, objKey, obj)}
}
