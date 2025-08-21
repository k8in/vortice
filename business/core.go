package business

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"vortice/container"
	"vortice/object"
)

var (
	// ErrInitialized is the error returned when an attempt to initialize an already initialized BusinessCore is made.
	ErrInitialized = errors.New("the BusinessCore has already been initialized")
	// ErrInitContainer is the error returned when there's a failure to initialize the container.
	ErrInitContainer = errors.New("failed to init container")
	// ErrNilPlugin is the error returned when an operation is attempted with a nil plugin.
	ErrNilPlugin = errors.New("plugin cannot be nil")
	// ErrInReadonlyMode is the error returned when an operation is attempted while the BusinessCore is in readonly mode.
	ErrInReadonlyMode = errors.New("the BusinessCore is in readonly mode")

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
	core       *container.Core
	readonly   *atomic.Bool
	plugins    *sync.Map
	mutex      *sync.RWMutex
	current    *Plugin
	extensions map[string]*object.Definition
	abilities  map[string][]*object.Definition
}

// NewCore initializes a new Core instance with the provided container.Core, setting up a mutex and an empty plugin list.
func NewCore(core *container.Core) *Core {
	return &Core{
		core:       core,
		readonly:   &atomic.Bool{},
		plugins:    &sync.Map{},
		current:    nil,
		mutex:      &sync.RWMutex{},
		extensions: map[string]*object.Definition{},
		abilities:  map[string][]*object.Definition{},
	}
}

// Init initializes the Core and its plugins, setting it to readonly and preparing for operation.
func (c *Core) Init() error {
	if ok := c.readonly.CompareAndSwap(false, true); !ok {
		return ErrInitialized
	}
	var err error
	if err = c.core.Init(); err != nil {
		return errors.Join(ErrInitContainer, err)
	}
	c.plugins.Range(func(key any, value any) bool {
		plugin := value.(*Plugin)
		if err = c.openPlugin(plugin); err != nil {
			return false
		}
		if err = c.initPlugin(plugin); err != nil {
			return false
		}
		if err = c.closePlugin(plugin); err != nil {
			return false
		}
		return true
	})
	return err
}

// Start initiates the services, ensuring they are running and managing their lifecycle.
func (c *Core) Start() error {
	return c.core.Start()
}

// Shutdown stops all running services and cleans up resources, finalizing the Core.
func (c *Core) Shutdown() {
	c.core.Shutdown()
}

// RegisterExtension registers a factory function with the given property, setting extension and namespace tags.
func (c *Core) RegisterExtension(fn any, prop *object.Property) (*object.Definition, error) {
	if err := c.checkReadonlyMode(); err != nil {
		return nil, err
	}
	if plugin := c.current; plugin != nil {
		return c.registerPluginExt(fn, prop, plugin)
	}
	return c.registerMainExt(fn, prop)
}

//// RegisterAbility registers a new ability with the given function and property, returning its definition or an error.
//func (c *Core) RegisterAbility(fn any, prop *object.Property) (*object.Definition, error) {
//	if err := c.checkReadonlyMode(); err != nil {
//		return nil, err
//	}
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//	prop.SetTags(TagAbilityKind, newNamespaceTag(MainNamespace))
//	def, err := c.core.RegisterFactory(fn, prop, false)
//	if err != nil {
//		return nil, errors.Join(ErrRegisterAbility, err)
//	}
//	key := def.Name()
//	c.abilities[key] = append(c.abilities[key], def)
//	return def, nil
//}

// RegisterPlugin adds a plugin to the Core, checking for readonly mode and ensuring no current plugin or duplicate exists.
func (c *Core) RegisterPlugin(plugin *Plugin) error {
	if plugin == nil {
		return ErrNilPlugin
	}
	if err := c.checkReadonlyMode(); err != nil {
		return err
	}
	if _, ok := c.plugins.LoadOrStore(plugin.Name(), plugin); ok {
		return fmt.Errorf("%s already exists", plugin)
	}
	return nil
}

// registerMainExt registers an extension point with the given function and property, setting extension and main namespace tags.
func (c *Core) registerMainExt(fn any, prop *object.Property) (*object.Definition, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	prop.SetTags(TagExtensionKind, newNamespaceTag(MainNamespace))
	def, err := c.core.RegisterFactory(fn, prop, true)
	if err != nil {
		return nil, fmt.Errorf("register main extension failed: %w", err)
	}
	if def0, ok := c.extensions[def.Name()]; ok {
		err = fmt.Errorf("main extension %s already exists: %s", def.Name(), def0.Name())
		return nil, err
	}
	c.extensions[def.Name()] = def
	return def, nil
}

// registerPluginExt registers an extension for a plugin, setting appropriate tags and associating it with the plugin.
func (c *Core) registerPluginExt(fn any, prop *object.Property, plugin *Plugin) (*object.Definition, error) {
	prop.SetTags(TagExtensionKind, newNamespaceTag(plugin.Name()))
	def, err := c.core.RegisterFactory(fn, prop, false)
	if err != nil {
		return nil, fmt.Errorf("register plugin extension failed: %w", err)
	}
	if plugin != nil {
		if ok := plugin.addExtension(def); !ok {
			err := fmt.Errorf("plugin extension %s already exists", def.Name())
			return nil, err
		}
	}
	return def, nil
}

// checkReadonlyMode verifies if the Core is in readonly mode and panics if true.
func (c *Core) checkReadonlyMode() error {
	if c.readonly.Load() {
		return ErrInReadonlyMode
	}
	return nil
}

// openPlugin locks the Core, checks if a plugin is already open, and sets the given plugin as current.
func (c *Core) openPlugin(plugin *Plugin) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.current != nil {
		return fmt.Errorf("cannot open plugin %s, current is %s", plugin.Name(), c.current.Name())
	}
	c.current = plugin
	return nil
}

// closePlugin closes the specified plugin, ensuring it is the current active plugin and then setting current to nil.
func (c *Core) closePlugin(plugin *Plugin) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.current != plugin {
		return fmt.Errorf("cannot close plugin %s, current is %s", plugin.Name(), c.current.Name())
	}
	c.current = nil
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
