package container

import (
	"context"
	"vortice/object"
)

type Context interface {
	context.Context
	SetFilter(dfs ...object.DefinitionFilter)
	GetFilters() []object.DefinitionFilter
	GetObjects() map[string]Object
}
