package container

import (
	"context"
	"time"

	"vortice/object"
	"vortice/util"
)

var (
	// serviceFilter is a DefinitionFilter that selects components with lifecycle methods, auto-startup, and singleton scope.
	serviceFilter = object.DefinitionFilter(func(def *object.Definition) bool {
		return def.Methods().IsLifeCycle() && def.AutoStartup() && def.Scope() == object.Singleton
	})
)

type lifecycleProcessor struct {
	comps   []Object
	tg      *util.TaskGroup
	timeout time.Duration
}

func newLifecycleProcessor(timeout time.Duration) *lifecycleProcessor {
	return &lifecycleProcessor{
		comps:   []Object{},
		tg:      util.NewTaskGroup("vortice.container.lifecycleProcessor"),
		timeout: timeout,
	}
}

func (p *lifecycleProcessor) start(ctx context.Context, factory ObjectFactory) {

	//for _, name := range component.GetLifecycleNames() {
	//	comp := factory.GetSingleton(name)
	//	if comp == nil {
	//		// bug 分支
	//		logger.Printf("Singleton %s not found", name)
	//	}
	//	if comp.Running() {
	//		continue
	//	}
	//	p.comps = append(p.comps, comp)
	//	logger.Printf("Starting %s......", comp)
	//	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	//	err := p.tg.GoAndWait(ctx, func(ctx context.Context) error {
	//		comp.Start()
	//		return nil
	//	})
	//	cancel()
	//	if err != nil {
	//		logger.Printf("Stopped while running %s, err: %s", comp, err.Error())
	//	} else if !comp.Running() {
	//		logger.Printf("Stopped while running %s, it wasn't running", comp)
	//	}
	//	logger.Printf("%s started successfully", comp)
	//}
}
