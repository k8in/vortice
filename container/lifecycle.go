package container

import (
	"context"
	"errors"
	"time"

	"vortice/object"
	"vortice/util"

	"go.uber.org/zap"
)

const (
	// lifecycleTaskGroupName is the name assigned to the task group used for managing lifecycle tasks.
	lifecycleTaskGroupName = "container.lifecycleProcessor"
)

var (
	// serviceFilter is a DefinitionFilter that selects components with lifecycle methods, auto-startup, and singleton scope.
	serviceFilter = object.DefinitionFilter(func(def *object.Definition) bool {
		return def.Methods().IsLifeCycle() && def.AutoStartup() && def.Scope() == object.Singleton
	})
)

// lifecycleProcessor manages the lifecycle of a set of services, coordinating their start and stop operations.
type lifecycleProcessor struct {
	objs    []Object
	tasks   *util.TaskGroup
	timeout time.Duration
}

// newLifecycleProcessor creates a new lifecycleProcessor with the specified timeout for managing service lifecycles.
func newLifecycleProcessor(timeout time.Duration) *lifecycleProcessor {
	return &lifecycleProcessor{
		objs:    []Object{},
		tasks:   util.NewTaskGroup(lifecycleTaskGroupName),
		timeout: timeout,
	}
}

// start initiates the services defined by the factory, ensuring they are running and managing their lifecycle.
func (p *lifecycleProcessor) start(ctx context.Context, factory ObjectFactory) error {
	coreCtx := WithCoreContext(ctx)
	l := util.Logger()
	for _, def := range factory.GetDefinitions(serviceFilter) {
		objs, err := factory.GetObjectsByName(coreCtx, def.Name())
		if err != nil {
			return errors.Join(errors.New("GetObjectsByName failed"), err)
		}
		for _, obj := range objs {
			if obj.Running() {
				l.Warn("service has been started", zap.String("service", obj.ID()))
				continue
			}
			l.Info("starting service......", zap.String("service", obj.ID()))
			ctx, cancel := context.WithTimeout(ctx, p.timeout)
			err := p.tasks.GoAndWait(ctx, func(ctx context.Context) error {
				return obj.Start()
			})
			cancel()
			if err != nil {
				l.Error("stopped while running service", zap.String("service", obj.ID()),
					zap.Error(err))
				return err
			} else if !obj.Running() {
				err := errors.New("it wasn't running")
				l.Error("stopped while running service", zap.String("service", obj.ID()),
					zap.Error(err))
				return err
			}
			p.objs = append(p.objs, obj)
			l.Info("service started successfully", zap.String("service", obj.ID()))
		}
	}
	return nil
}

// stop stops all running services in reverse order, logging the status of each service.
func (p *lifecycleProcessor) stop(ctx context.Context) {
	n := len(p.objs)
	l := util.Logger()
	for i := n - 1; i >= 0; i-- {
		obj := p.objs[i]
		l.Info("stopping service......", zap.String("service", obj.ID()))
		if !obj.Running() {
			l.Info("service wasn't running", zap.String("service", obj.ID()))
			continue
		}
		ctx, cancel := context.WithTimeout(ctx, p.timeout)
		err := p.tasks.GoAndWait(ctx, func(ctx context.Context) error {
			return obj.Stop()
		})
		cancel()
		if err != nil {
			l.Error("service wasn't stopped",
				zap.String("service", obj.ID()), zap.Error(err))
		}
		l.Info("service stopped successfully", zap.String("service", obj.ID()))
	}
}
