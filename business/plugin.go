package business

type Namespace string

type PluginInitFunc func() error

type Plugin struct {
	ns   string
	fns  []PluginInitFunc
	used []any
}

func NewPlugin(ns string) *Plugin {
	return &Plugin{ns: ns, fns: []PluginInitFunc{}}
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
