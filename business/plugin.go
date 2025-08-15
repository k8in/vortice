package business

import "vortice/object"

type PluginInitFunc func() error

type Plugin struct {
	ns   object.Namespace
	fns  []PluginInitFunc
	used []any
}

func NewPlugin(ns string) *Plugin {
	return &Plugin{ns: object.Namespace(ns), fns: []PluginInitFunc{}}
}

func (p *Plugin) InitExtension(fn ...PluginInitFunc) {
	p.fns = append(p.fns, fn...)
}

func (p *Plugin) Use() {

}

func (p *Plugin) Register() {
	return
}

func (p *Plugin) init() error {
	for _, init := range p.fns {
		if err := init(); err != nil {
			return err
		}
	}
	return nil
}
