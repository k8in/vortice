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

type Ctx struct {
	context.Context
	factory ObjectFactory
	filters []object.DefinitionFilter
	objects map[string]Object
}
