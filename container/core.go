package container

import (
	"context"
	"sync"
	"time"
)

var (
	// DefaultStartupTimeout defines the default duration to wait for services to start before timing out.
	DefaultStartupTimeout = 6 * time.Second

	core *Core
	once = &sync.Once{}
)

// DefaultCore returns a singleton instance of Core, initializing it with a background context if not already initialized.
func DefaultCore() *Core {
	once.Do(func() {
		core = NewCore(context.Background())
	})
	return core
}

// Core is the main structure that encapsulates the context, object factory, and lifecycle processor.
type Core struct {
	context.Context
	ObjectFactory
	lcp *lifecycleProcessor
}

// NewCore initializes and returns a new Core instance with the provided context, setting up an object factory and lifecycle processor.
func NewCore(ctx context.Context) *Core {
	return &Core{
		Context:       ctx,
		ObjectFactory: NewCoreObjectFactory(),
		lcp:           newLifecycleProcessor(DefaultStartupTimeout),
	}
}

// Init initializes the Core's ObjectFactory, preparing it for use and returns an error if initialization fails.
func (c *Core) Init() error {
	return c.ObjectFactory.Init()
}

// Start initiates the services defined by the factory, ensuring they are running and managing their lifecycle.
func (c *Core) Start() error {
	return c.lcp.start(c.Context, c.ObjectFactory)
}

// Shutdown stops all running services and cleans up resources, finalizing the Core.
func (c *Core) Shutdown() {
	c.lcp.stop(c.Context)
	c.ObjectFactory.Destroy()
}
