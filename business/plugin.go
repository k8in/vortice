package business

import "vortice/object"

// Plugin represents a modular component that can be registered and initialized within a Core, providing additional functionalities.
type Plugin struct {
	core      *Core
	name      string
	initFns   []func() error
	exts      map[string]*object.Definition
	abilities []any
}

// newPlugin creates a new Plugin instance with the specified name and associated Core.
func newPlugin(name string, core *Core) *Plugin {
	return &Plugin{
		core:      core,
		name:      name,
		initFns:   []func() error{},
		exts:      map[string]*object.Definition{},
		abilities: []any{},
	}
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string {
	return p.name
}

// Init appends initialization functions to the plugin, to be executed during the plugin's initialization.
func (p *Plugin) Init(fn ...func() error) {
	p.initFns = append(p.initFns, fn...)
}

// Register adds the plugin to the Core, making it available for initialization and use.
func (p *Plugin) Register() {
	p.core.RegisterPlugin(p)
}

// setDefinition adds a new object definition to the plugin's extensions map using its ID.
func (p *Plugin) setDefinition(def *object.Definition) {
	p.exts[def.ID()] = def
}

// init executes all initialization functions registered with the plugin, returning an error if any function fails.
func (p *Plugin) init() error {
	for _, initFunc := range p.initFns {
		if err := initFunc(); err != nil {
			return err
		}
	}
	return nil
}
