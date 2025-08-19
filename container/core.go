package container

import (
	"context"
)

type Core struct {
	context.Context
	factory ObjectFactory
}

func NewCore(ctx context.Context) *Core {
	return &Core{
		Context: ctx,
		factory: NewCoreObjectFactory(),
	}
}

func (c *Core) ObjectFactory() ObjectFactory {
	return c.factory
}

func (c *Core) Init() error {
	if err := c.factory.Init(); err != nil {
		return err
	}
	return nil
}

func (ctx *Core) Start() {
	//ctx.lcp.start(ctx, ctx.factory)
}

func (ctx *Core) Shutdown() {
	//ctx.lcp.stop(ctx)
	//ctx.factory.Destroy()
}
