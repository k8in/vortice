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
	// ErrNilPlugin is the error returned when an operation is attempted with a nil plugin.
	ErrNilPlugin = errors.New("plugin cannot be nil")
	// ErrInReadonlyMode is the error returned when an operation is attempted while the BusinessCore is in readonly mode.
	ErrInReadonlyMode = errors.New("the BusinessCore is in readonly mode")
	// ErrRegisterPlugin is the error returned when a plugin registration fails.
	ErrRegisterPlugin = errors.New("failed to register plugin")
	// ErrRegisterExtension is the error returned when an extension registration fails.
	ErrRegisterExtension = errors.New("failed to register extension")
	// ErrRegisterMainExt is the error returned when a factory registration fails.
	ErrRegisterMainExt = errors.New("failed to register main extension")
	// ErrRegisterPluginExt is the error returned when a plugin extension registration fails.
	ErrRegisterPluginExt = errors.New("failed to register plugin extension")
	// ErrRegisterAbility is the error returned when an ability registration fails.
	ErrRegisterAbility = errors.New("failed to register ability")
	// ErrInitialized is the error returned when an attempt to initialize an already initialized BusinessCore is made.
	ErrInitialized = errors.New("the BusinessCore has already been initialized")

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
	current    *atomic.Value
	mutex      *sync.Mutex
	extensions map[string]*object.Definition
	abilities  map[string][]*object.Definition
}

// NewCore initializes a new Core instance with the provided container.Core, setting up a mutex and an empty plugin list.
func NewCore(core *container.Core) *Core {
	return &Core{
		core:       core,
		readonly:   &atomic.Bool{},
		plugins:    &sync.Map{},
		current:    &atomic.Value{},
		mutex:      &sync.Mutex{},
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
		return errors.Join(errors.New("failed to container core"), err)
	}
	c.plugins.Range(func(key any, value any) bool {
		if !c.current.CompareAndSwap(nil, value) {
			err = errors.New("another plugin is already registered")
			return false
		}
		plugin := value.(*Plugin)
		if err := c.initPlugin(plugin); err != nil {
			err = fmt.Errorf("failed to initialize plugin: %s", plugin.Name())
			return false
		}
		if !c.current.CompareAndSwap(value, nil) {
			err = errors.New("failed to clear current plugin")
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
		return nil, errors.Join(ErrRegisterExtension, err)
	}
	if v := c.current.Load(); v != nil {
		return c.registerPluginExt(fn, prop, v.(*Plugin))
	}
	return c.registerMainExt(fn, prop)
}

// RegisterAbility registers a new ability with the given function and property, returning its definition or an error.
func (c *Core) RegisterAbility(fn any, prop *object.Property) (*object.Definition, error) {
	if err := c.checkReadonlyMode(); err != nil {
		return nil, errors.Join(ErrRegisterAbility, err)
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	prop.SetTags(TagAbilityKind, newNamespaceTag(MainNamespace))
	def, err := c.core.RegisterFactory(fn, prop, false)
	if err != nil {
		return nil, errors.Join(ErrRegisterAbility, err)
	}
	key := def.Name()
	c.abilities[key] = append(c.abilities[key], def)
	return def, nil
}

// RegisterPlugin adds a plugin to the Core, checking for readonly mode and ensuring no current plugin or duplicate exists.
func (c *Core) RegisterPlugin(plugin *Plugin) error {
	if plugin == nil {
		return errors.Join(ErrRegisterPlugin, ErrNilPlugin)
	}
	if err := c.checkReadonlyMode(); err != nil {
		return errors.Join(ErrRegisterPlugin, err)
	}
	if ok := c.plugins.CompareAndSwap(plugin.Name(), nil, plugin); !ok {
		err := fmt.Errorf("plugin %s already exists", plugin.Name())
		return errors.Join(ErrRegisterPlugin, err)
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
		return nil, errors.Join(ErrRegisterMainExt, err)
	}
	if def0, ok := c.extensions[def.Name()]; ok {
		err = fmt.Errorf("extension %s already exists: %s", def.Name(), def0.Name())
		return nil, errors.Join(ErrRegisterMainExt, err)
	}
	c.extensions[def.Name()] = def
	return def, nil
}

// registerPluginExt registers an extension for a plugin, setting appropriate tags and associating it with the plugin.
func (c *Core) registerPluginExt(fn any, prop *object.Property, plugin *Plugin) (*object.Definition, error) {
	prop.SetTags(TagExtensionKind, newNamespaceTag(plugin.Name()))
	def, err := c.core.RegisterFactory(fn, prop, false)
	if err != nil {
		return nil, errors.Join(ErrRegisterPluginExt, err)
	}
	if plugin != nil {
		if ok := plugin.addExtension(def); !ok {
			err := fmt.Errorf("extension %s already exists", def.Name())
			return nil, errors.Join(ErrRegisterPluginExt, err)
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

// initPlugin initializes a given plugin, ensuring it's ready for use and checking its abilities.
func (c *Core) initPlugin(plugin *Plugin) error {
	if err := plugin.init(); err != nil {
		return errors.New("")
	}
	// check ability
	return nil
}
