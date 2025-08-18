package container

import (
	"context"

	"vortice/object"
)

type Core struct {
	context.Context
	factory ObjectFactory
	filters []object.DefinitionFilter
}

func NewCore(ctx context.Context) *Core {
	return &Core{
		Context: ctx,
		factory: NewObjectFactory(),
		filters: []object.DefinitionFilter{},
	}
}

func (c *Core) Init() {
	//object.GetDefinitionRegistry().Init()
	//c.factory.init()
}

func (ctx *Core) Start() {
	//ctx.lcp.start(ctx, ctx.factory)
}

func (ctx *Core) Shutdown() {
	//ctx.lcp.stop(ctx)
	//ctx.factory.Destroy()
}
