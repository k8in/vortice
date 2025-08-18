package container

import (
	"context"
	"vortice/object"
)

type Context interface {
	context.Context
	GetFilters() []object.DefinitionFilter
	GetObjects() map[string]Object
}
