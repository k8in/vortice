package business

import (
	"errors"
	"sync"
	"sync/atomic"

	"vortice/container"
	"vortice/object"
	"vortice/util"
)

var (
	core *Core
	once = &sync.Once{}
)

// DefaultCore returns a singleton instance of Core, initializing it with a default container core if not already initialized.
func DefaultCore() *Core {
	once.Do(func() {
		core = NewCore(container.DefaultCore())
	})
	return core
}

// MainNamespace is the default namespace used when no specific plugin namespace is set.
const MainNamespace = "main"

// newNamespaceTag creates and returns a new Tag with the key "namespace" and the provided namespace string as value.
func newNamespaceTag(ns string) object.Tag {
	return object.NewTag("namespace", ns)
}

// Core encapsulates the container's core, a set of plugins, and a mutex for thread-safe operations.
type Core struct {
	core     *container.Core
	plugins  *sync.Map
	readonly *atomic.Bool
	current  *atomic.Value
}

// NewCore initializes a new Core instance with the provided container.Core, setting up a mutex and an empty plugin list.
func NewCore(core *container.Core) *Core {
	return &Core{
		core:     core,
		plugins:  &sync.Map{},
		readonly: &atomic.Bool{},
		current:  &atomic.Value{},
	}
}

// Init initializes the Core and its plugins, setting it to readonly and preparing for operation.
func (c *Core) Init() error {
	if c.readonly.Load() {
		return nil
	}
	c.readonly.Store(true)
	if err := c.core.Init(); err != nil {
		return err
	}
	var err error
	c.plugins.Range(func(key any, value any) bool {
		if !c.current.CompareAndSwap(nil, value) {
			err = errors.New("")
			return false
		}
		if err := c.initPlugin(value.(*Plugin)); err != nil {
			err = errors.New("")
			return false
		}
		if !c.current.CompareAndSwap(value, nil) {
			err = errors.New("")
			return false
		}
		return true
	})
	return err
}

// RegisterExtension registers a factory function with the given property, setting extension and namespace tags.
func (c *Core) RegisterExtension(fn any, prop *object.Property) (*object.Definition, error) {
	if err := c.checkReadonly(); err != nil {
		return nil, err
	}
	if v := c.current.Load(); v != nil {
		return c.registerPluginExt(fn, prop, v.(*Plugin))
	}
	return c.registerDefaultExt(fn, prop)
}

// NewPlugin creates and returns a new Plugin instance with the given name, associated with the current Core.
func (c *Core) NewPlugin(name string) *Plugin {
	return newPlugin(name, c)
}

// Start initiates the services, ensuring they are running and managing their lifecycle.
func (c *Core) Start() error {
	return c.core.Start()
}

// Shutdown stops all running services and cleans up resources, finalizing the Core.
func (c *Core) Shutdown() {
	c.core.Shutdown()
}

// RegisterPlugin adds a new plugin to the Core, ensuring thread safety by locking the mutex.
func (c *Core) RegisterPlugin(plugin *Plugin) {
	if err := c.checkReadonly(); err != nil {
		util.Logger().Panic("1")
	}
	if c.current != nil {
		util.Logger().Panic("2")
	}
	if _, ok := c.plugins.Load(plugin.Name()); ok {
		util.Logger().Panic("3")
	}
	c.plugins.Store(plugin.Name(), plugin)
}

// registerDefaultExt registers an extension point with the given function and property, setting extension and main namespace tags.
func (c *Core) registerDefaultExt(fn any, prop *object.Property) (*object.Definition, error) {
	prop.SetTags(TagExtension, newNamespaceTag(MainNamespace))
	return c.core.RegisterFactory(fn, prop, true)
}

// registerPluginExt registers an extension for a plugin, setting appropriate tags and associating it with the plugin.
func (c *Core) registerPluginExt(fn any, prop *object.Property, plugin *Plugin) (*object.Definition, error) {
	prop.SetTags(TagExtension, newNamespaceTag(plugin.Name()))
	def, err := c.core.RegisterFactory(fn, prop, false)
	if err != nil {
		return nil, err
	}
	if plugin != nil {
		plugin.setDefinition(def)
	}
	return def, nil
}

// checkReadonly verifies if the Core is in readonly mode and panics if true.
func (c *Core) checkReadonly() error {
	if c.readonly.Load() {
		return errors.New("readonly")
	}
	return nil
}

// initPlugin initializes a given plugin, ensuring it's ready for use and checking its abilities.
func (c *Core) initPlugin(plugin *Plugin) error {
	if err := plugin.init(); err != nil {
		return errors.New("")
	}
	// check ability
	return nil
}
