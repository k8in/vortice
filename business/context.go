package business

import (
	"context"
	"vortice/container"
	"vortice/object"
)

// Context represents an interface for managing and accessing context-specific information,
// including namespace and objects.
type Context interface {
	container.Context
	// Namespace returns the namespace associated with the context.
	Namespace() string
}

type (
	// ctxKey is a type for context keys to store and retrieve values in a context.
	ctxKey string
)

var (
	// nsKey is a context key used to store the namespace in a context.
	nsKey = ctxKey("namespace")
)

// Ctx is a context type that extends CoreContext, providing additional functionality
// for managing application contexts with an associated object.
type Ctx struct {
	*container.CoreContext
}

// WithContext creates a new context with the specified namespace and returns a Context object.
func WithContext(ctx context.Context, ns string) Context {
	ctx = context.WithValue(ctx, nsKey, ns)
	coreCtx := container.WithCoreContext(ctx)
	filter := object.TagFilter(object.NewTag("namespace", ns))
	coreCtx.SetFilter(filter)
	return &Ctx{CoreContext: coreCtx}
}

// GetNamespace retrieves the namespace from the provided context, returning an empty string if not found or context is nil.
func GetNamespace(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(nsKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

// Namespace retrieves the namespace associated with the context.
func (ctx *Ctx) Namespace() string {
	return GetNamespace(ctx.Context)
}
