package container

import (
	"context"
	"vortice/object"
)

// Context represents an interface that extends the standard context.Context with additional methods
// to get filters and managed objects.
type Context interface {
	// Context represents the standard context.Context, providing a way to carry deadlines,
	// cancellation signals, and other request-scoped values.
	context.Context
	// SetFilter sets one or more filters to be used for filtering object definitions.
	SetFilter(...object.DefinitionFilter)
	// GetFilters returns a slice of DefinitionFilter that can be used to filter object definitions.
	GetFilters() []object.DefinitionFilter
	// GetObjects returns a map of managed objects, keyed by their unique identifiers.
	GetObjects() map[string]Object
}

// CoreContext extends the standard context.Context to include additional functionality for managing application contexts.
type CoreContext struct {
	context.Context
	dfs  []object.DefinitionFilter
	objs map[string]Object
}

// WithCoreContext creates a new CoreContext with the provided context, enhancing it for application-specific context management.
func WithCoreContext(ctx context.Context) *CoreContext {
	return &CoreContext{
		Context: ctx,
		dfs:     []object.DefinitionFilter{},
		objs:    map[string]Object{},
	}
}

// SetFilter updates the list of DefinitionFilter functions used for filtering component definitions.
func (c *CoreContext) SetFilter(filters ...object.DefinitionFilter) {
	if len(filters) == 0 {
		return
	}
	dfs := make([]object.DefinitionFilter, len(filters))
	copy(dfs, filters)
	c.dfs = dfs
}

// GetFilters returns a slice of DefinitionFilter functions that can be used to filter component definitions.
func (c *CoreContext) GetFilters() []object.DefinitionFilter {
	if len(c.dfs) == 0 {
		return []object.DefinitionFilter{}
	}
	cp := make([]object.DefinitionFilter, len(c.dfs))
	copy(cp, c.dfs)
	return cp
}

// GetObjects returns a map of objects managed by the container, currently returning an empty map.
func (c *CoreContext) GetObjects() map[string]Object {
	if c.objs == nil {
		c.objs = map[string]Object{}
	}
	return c.objs
}
