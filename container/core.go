package container

import (
	"context"
	"time"
)

var (
	// DefaultStartupTimeout defines the default duration to wait for services to start before timing out.
	DefaultStartupTimeout = 6 * time.Second
)

// Core is the main structure that encapsulates the context, object factory, and lifecycle processor.
type Core struct {
	context.Context
	factory ObjectFactory
	lcp     *lifecycleProcessor
}

// NewCore initializes and returns a new Core instance with the provided context, setting up an object factory and lifecycle processor.
func NewCore(ctx context.Context) *Core {
	return &Core{
		Context: ctx,
		factory: NewCoreObjectFactory(),
		lcp:     newLifecycleProcessor(DefaultStartupTimeout),
	}
}

// Init initializes the Core's ObjectFactory, preparing it for use and returns an error if initialization fails.
func (c *Core) Init() error {
	return c.factory.Init()
}

// Start initiates the services defined by the factory, ensuring they are running and managing their lifecycle.
func (c *Core) Start() error {
	return c.lcp.start(c.Context, c.factory)
}

// Shutdown stops all running services and cleans up resources, finalizing the Core.
func (c *Core) Shutdown() {
	c.lcp.stop(c.Context)
	c.factory.Destroy()
}
