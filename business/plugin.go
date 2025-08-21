package business

import (
	"fmt"
	"vortice/object"
)

// Plugin represents a plugin with initialization functions, extensions, and abilities.
type Plugin struct {
	name       string
	inits      []func() error
	extensions map[string]*object.Definition
	abilities  map[string]*object.Definition
}

// NewPlugin creates a new Plugin instance with the specified name and associated Core.
func NewPlugin(name string) *Plugin {
	return &Plugin{
		name:       name,
		inits:      []func() error{},
		extensions: map[string]*object.Definition{},
		abilities:  map[string]*object.Definition{},
	}
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string {
	return p.name
}

// Init appends initialization functions to the plugin, to be executed during the plugin's initialization.
func (p *Plugin) Init(fn ...func() error) {
	p.inits = append(p.inits, fn...)
}

// GetExtension retrieves the object definition for a given extension name, returning nil if not found.
func (p *Plugin) GetExtension(name string) *object.Definition {
	if def, ok := p.extensions[name]; ok {
		return def
	}
	return nil
}

// String returns a string representation of the Plugin, including its name.
func (p *Plugin) String() string {
	return fmt.Sprintf("<Plugin %s>", p.name)
}

// addExtension adds a new extension to the plugin if it does not already exist, returning true if added.
func (p *Plugin) addExtension(def *object.Definition) bool {
	if _, ok := p.extensions[def.Name()]; !ok {
		p.extensions[def.Name()] = def
		return true
	}
	return false
}

// init executes all initialization functions registered with the plugin, returning an error if any function fails.
func (p *Plugin) init() error {
	for _, initFunc := range p.inits {
		if err := initFunc(); err != nil {
			return err
		}
	}
	return nil
}
