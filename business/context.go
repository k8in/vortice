package business

import (
	"context"
	"vortice/container"
)

// Context represents an interface that extends the standard context with additional
// methods for managing objects and namespaces.
type Context[O any] interface {
	container.Context
	Namespace() Namespace
	Object() O
}

type ctxKey string

var (
	// nsKey is a context key used to store the namespace in a context.
	nsKey = ctxKey("namespace")
	// objKey is a context key used to store the object in a context.
	objKey = ctxKey("object")
)

type Ctx[O any] struct {
	*container.CoreContext
	obj O
}

// WithContext creates a new context with the specified namespace and returns a Context object.
func WithContext[O any](ctx context.Context, ns string) Context[O] {
	ctx = context.WithValue(ctx, nsKey, Namespace(ns))
	return &Ctx[O]{CoreContext: container.WithCoreContext(ctx)}
}

// Namespace retrieves the namespace associated with the context.
func (ctx *Ctx[O]) Namespace() Namespace {
	return ctx.Value(nsKey).(Namespace)
}

// Object returns the object associated with the context.
func (ctx *Ctx[O]) Object() O {
	return ctx.obj
}
