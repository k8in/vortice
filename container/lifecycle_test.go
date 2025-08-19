package container

import (
	"context"
	"errors"
	"testing"
	"time"
	"vortice/object"
)

// --- Test service types implementing object.Lifecycle ---

type svcOK struct{ running bool }

func (s *svcOK) Start() error  { s.running = true; return nil }
func (s *svcOK) Stop() error   { s.running = false; return nil }
func (s *svcOK) Running() bool { return s.running }

type svcAlreadyRunning struct{ running bool }

func (s *svcAlreadyRunning) Start() error  { return nil }
func (s *svcAlreadyRunning) Stop() error   { s.running = false; return nil }
func (s *svcAlreadyRunning) Running() bool { return true }

type svcStartErr struct{}

func (s *svcStartErr) Start() error  { return errors.New("start error") }
func (s *svcStartErr) Stop() error   { return nil }
func (s *svcStartErr) Running() bool { return false }

type svcNotRunningAfterStart struct{ running bool }

func (s *svcNotRunningAfterStart) Start() error  { return nil }
func (s *svcNotRunningAfterStart) Stop() error   { s.running = false; return nil }
func (s *svcNotRunningAfterStart) Running() bool { return false }

type svcStopErr struct{ running bool }

func (s *svcStopErr) Start() error  { s.running = true; return nil }
func (s *svcStopErr) Stop() error   { return errors.New("stop error") }
func (s *svcStopErr) Running() bool { return s.running }

// helper to create factory and register a service with AutoStartup+Singleton
func newFactoryWithService[T any](t *testing.T, fn any) *CoreObjectFactory {
	t.Helper()
	prop := object.NewProperty()
	prop.AutoStartup = true
	prop.Scope = object.Singleton
	// LazyInit 无关紧要，这里保持默认
	factory := NewCoreObjectFactory()
	if _, err := factory.RegisterFactory(fn, prop, false); err != nil {
		t.Fatalf("RegisterFactory failed: %v", err)
	}
	if err := factory.Init(); err != nil {
		t.Fatalf("factory Init failed: %v", err)
	}
	return factory
}

func newCoreCtx2() *CoreContext { return WithCoreContext(context.Background()) }

func TestLifecycleProcessor_Start_Success(t *testing.T) {
	factory := newFactoryWithService[*svcOK](t, func() *svcOK { return &svcOK{} })
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if len(lp.objs) != 1 {
		t.Fatalf("expected 1 started object, got %d", len(lp.objs))
	}
	if !lp.objs[0].Running() {
		t.Fatalf("object should be running after start")
	}
	if err := lp.stop(newCoreCtx2()); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestLifecycleProcessor_Start_AlreadyRunning(t *testing.T) {
	factory := newFactoryWithService[*svcAlreadyRunning](t, func() *svcAlreadyRunning {
		return &svcAlreadyRunning{running: true}
	})
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err != nil {
		t.Fatalf("start should not fail when already running: %v", err)
	}
	if len(lp.objs) != 0 {
		t.Fatalf("already-running service should not be added to running list, got %d", len(lp.objs))
	}
}

func TestLifecycleProcessor_Start_StartError(t *testing.T) {
	factory := newFactoryWithService[*svcStartErr](t, func() *svcStartErr { return &svcStartErr{} })
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err == nil {
		t.Fatalf("start should fail when service Start returns error")
	}
}

func TestLifecycleProcessor_Start_NotRunningAfterStart(t *testing.T) {
	factory := newFactoryWithService[*svcNotRunningAfterStart](t, func() *svcNotRunningAfterStart {
		return &svcNotRunningAfterStart{}
	})
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err == nil {
		t.Fatalf("start should fail when service not running after Start")
	}
}

func TestLifecycleProcessor_Stop_Success(t *testing.T) {
	factory := newFactoryWithService[*svcOK](t, func() *svcOK { return &svcOK{} })
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if len(lp.objs) != 1 || !lp.objs[0].Running() {
		t.Fatalf("object should be running after start")
	}
	if err := lp.stop(newCoreCtx2()); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestLifecycleProcessor_Stop_NotRunning(t *testing.T) {
	factory := newFactoryWithService[*svcOK](t, func() *svcOK { return &svcOK{} })
	// 获取未启动的对象，直接放入 lp.objs 以覆盖 stop 的未运行分支
	ctx := newCoreCtx2()
	defs := factory.GetDefinitions()
	if len(defs) == 0 {
		t.Fatalf("no definitions registered")
	}
	objs, err := factory.GetObjectsByName(ctx, defs[0].Name())
	if err != nil {
		t.Fatalf("GetObjectsByName failed: %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("expected one object")
	}
	if objs[0].Running() {
		t.Fatalf("object should not be running before lifecycle start")
	}
	lp := newLifecycleProcessor(200 * time.Millisecond)
	lp.objs = []Object{objs[0]}
	if err := lp.stop(ctx); err != nil {
		t.Fatalf("stop should not fail when object wasn't running: %v", err)
	}
}

func TestLifecycleProcessor_Stop_Error(t *testing.T) {
	factory := newFactoryWithService[*svcStopErr](t, func() *svcStopErr { return &svcStopErr{} })
	lp := newLifecycleProcessor(200 * time.Millisecond)
	if err := lp.start(newCoreCtx2(), factory); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if len(lp.objs) != 1 || !lp.objs[0].Running() {
		t.Fatalf("object should be running after start")
	}
	if err := lp.stop(newCoreCtx2()); err == nil {
		t.Fatalf("stop should fail when service Stop returns error")
	}
}
