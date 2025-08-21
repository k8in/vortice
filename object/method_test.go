package object

import (
	"fmt"
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

type testComponent struct {
	initialized bool
	running     bool
	destroyed   bool
}

func (c *testComponent) Init() error {
	c.initialized = true
	return nil
}
func (c *testComponent) Destroy() error {
	c.destroyed = true
	return nil
}
func (c *testComponent) Start() error {
	c.running = true
	return nil
}
func (c *testComponent) Stop() error {
	c.running = false
	return nil
}
func (c *testComponent) Running() bool {
	return c.running
}

type noLifecycle struct{}

func (c *noLifecycle) Init() error    { return nil }
func (c *noLifecycle) Destroy() error { return nil }
func (c *noLifecycle) Start() error   { return nil }
func (c *noLifecycle) Stop() error    { return nil }
func (c *noLifecycle) Running() bool  { return false }

// 仅实现 Init/Destroy，不实现 Lifecycle
type partialComponent struct {
	inited    bool
	stopped   bool
	destroyed bool
}

func (p *partialComponent) Init() error {
	p.inited = true
	return nil
}
func (p *partialComponent) Destroy() error {
	p.destroyed = true
	return nil
}

// Init 返回 error
type errorInitComponent struct{ called bool }

func (e *errorInitComponent) Init() error {
	e.called = true
	return fmt.Errorf("init failed")
}
func (e *errorInitComponent) Destroy() error { return nil }

// Destroy 返回 error
type errorDestroyComponent struct{}

func (e *errorDestroyComponent) Init() error    { return nil }
func (e *errorDestroyComponent) Destroy() error { return fmt.Errorf("destroy failed") }

func TestMethods_Lifecycle(t *testing.T) {
	objType := reflect.TypeOf(&testComponent{})
	methods := newMethods(objType)
	comp := &testComponent{}
	val := reflect.ValueOf(comp)

	if err := methods.CallInit(val); err != nil || !comp.initialized {
		t.Errorf("Init method failed: err=%v, initialized=%v", err, comp.initialized)
	}

	if err := methods.CallStart(val); err != nil || !comp.running {
		t.Errorf("Start method failed: err=%v, running=%v", err, comp.running)
	}

	running, err := methods.CallRunning(val)
	if err != nil || !running {
		t.Errorf("Running method failed: err=%v, running=%v", err, running)
	}

	if err := methods.CallStop(val); err != nil || comp.running {
		t.Errorf("Stop method failed: err=%v, running=%v", err, comp.running)
	}

	if err := methods.CallDestroy(val); err != nil || !comp.destroyed {
		t.Errorf("Destroy method failed: err=%v, destroyed=%v", err, comp.destroyed)
	}
}

func TestMethods_NoLifecycle(t *testing.T) {
	objType := reflect.TypeOf(&noLifecycle{})
	methods := newMethods(objType)
	comp := &noLifecycle{}
	val := reflect.ValueOf(comp)

	if err := methods.CallInit(val); err != nil {
		t.Errorf("Init should succeed: err=%v", err)
	}
	if err := methods.CallStart(val); err != nil {
		t.Errorf("Start should succeed: err=%v", err)
	}
	running, err := methods.CallRunning(val)
	if err != nil {
		t.Errorf("Running should succeed: err=%v", err)
	}
	if running {
		t.Errorf("Running should be false for noLifecycle")
	}
	if err := methods.CallStop(val); err != nil {
		t.Errorf("Stop should succeed: err=%v", err)
	}
	if err := methods.CallDestroy(val); err != nil {
		t.Errorf("Destroy should succeed: err=%v", err)
	}
}

// 1. 部分接口实现：Start/Stop/Running 方法指针应为 nil
func TestMethods_PartialInterfaces(t *testing.T) {
	mt := newMethods(reflect.TypeOf(&partialComponent{}))
	pc := &partialComponent{}
	rv := reflect.ValueOf(pc)

	// Init 正常
	if err := mt.CallInit(rv); err != nil || !pc.inited {
		t.Fatalf("Init unexpected err=%v inited=%v", err, pc.inited)
	}
	// Start / Stop 方法指针为 nil -> simpleCall 走 (ok=false,err=nil) 分支直接返回 nil
	if err := mt.CallStart(rv); err != nil {
		t.Fatalf("CallStart should return nil when method pointer nil, got %v", err)
	}
	if err := mt.CallStop(rv); err != nil {
		t.Fatalf("CallStop should return nil when method pointer nil, got %v", err)
	}
	// Running 方法指针为 nil -> 返回 (false,nil)
	r, err := mt.CallRunning(rv)
	if err != nil {
		t.Fatalf("CallRunning unexpected err=%v", err)
	}
	if r {
		t.Fatalf("CallRunning should be false when method missing")
	}
	// Destroy 正常
	if err := mt.CallDestroy(rv); err != nil || !pc.destroyed {
		t.Fatalf("Destroy unexpected err=%v destroyed=%v", err, pc.destroyed)
	}
	// IsLifeCycle 应为 false
	if mt.IsLifeCycle() {
		t.Fatalf("IsLifeCycle should be false for partialComponent")
	}
}

// 2. Init 返回 error 分支覆盖(simpleCall 中 rv[0] 为非 nil error)
func TestMethods_ErrorInit(t *testing.T) {
	mt := newMethods(reflect.TypeOf(&errorInitComponent{}))
	ins := &errorInitComponent{}
	if err := mt.CallInit(reflect.ValueOf(ins)); err == nil {
		t.Fatalf("expected init error")
	} else if !ins.called {
		t.Fatalf("Init not invoked")
	}
}

// 3. Destroy 返回 error 分支覆盖
func TestMethods_ErrorDestroy(t *testing.T) {
	mt := newMethods(reflect.TypeOf(&errorDestroyComponent{}))
	ins := &errorDestroyComponent{}
	if err := mt.CallDestroy(reflect.ValueOf(ins)); err == nil {
		t.Fatalf("expected destroy error")
	}
}

// 4. 方法指针存在但实例类型不匹配 -> call 返回 not found error 分支
func TestMethods_MethodNotFoundOnInstance(t *testing.T) {
	// 基于 testComponent 类型生成的方法集合（包含生命周期方法指针）
	mt := newMethods(reflect.TypeOf(&testComponent{}))
	// 传入 partialComponent 实例（无 Start/Stop/Running 方法）触发 MethodByName not found
	pc := &partialComponent{}
	if err := mt.CallStart(reflect.ValueOf(pc)); err == nil {
		t.Fatalf("expected error when calling Start on incompatible instance")
	}
	if err := mt.CallStop(reflect.ValueOf(pc)); err == nil {
		t.Fatalf("expected error when calling Stop on incompatible instance")
	}
	if _, err := mt.CallRunning(reflect.ValueOf(pc)); err == nil {
		t.Fatalf("expected error when calling Running on incompatible instance")
	}
}

// 5. IsLifeCycle 正向再测（原用例已有，此处补充调用以统计）
func TestMethods_IsLifeCycle_Positive(t *testing.T) {
	mt := newMethods(reflect.TypeOf(&testComponent{}))
	if !mt.IsLifeCycle() {
		t.Fatalf("expected lifecycle true for testComponent")
	}
}

// 自定义类型：Init 返回 (error,int) -> 触发 simpleCall len(rv)!=1
type badInitMulti struct{}

func (b *badInitMulti) Init() (error, int) { return nil, 1 }

// 自定义类型：Init 返回 string -> 触发 simpleCall 中 ok=false 路径
type initString struct{}

func (i *initString) Init() string { return "ok" }

// 自定义类型：Running 返回 (bool,error) -> 触发 CallRunning len(rv)!=1
type runningMulti struct{}

func (r *runningMulti) Running() (bool, error) { return true, nil }

// 自定义类型：Running 返回 int -> 触发 CallRunning 非 bool 转换错误
type badRunning struct{}

func (b *badRunning) Running() int { return 1 }

// 覆盖 simpleCall 分支：len(rv)!=1
func TestMethods_CallInit_LengthMismatch(t *testing.T) {
	typ := reflect.TypeOf(&badInitMulti{})
	m, _ := typ.MethodByName("Init")
	mt := &Methods{obj: typ, initMethod: &m}
	ins := &badInitMulti{}
	if err := mt.CallInit(reflect.ValueOf(ins)); err == nil {
		t.Fatalf("expected error for len(rv)!=1 in Init")
	}
}

// 覆盖 simpleCall 分支：返回值非 error 类型 (ok=false)
func TestMethods_CallInit_NonErrorReturn(t *testing.T) {
	typ := reflect.TypeOf(&initString{})
	m, _ := typ.MethodByName("Init")
	mt := &Methods{obj: typ, initMethod: &m}
	ins := &initString{}
	if err := mt.CallInit(reflect.ValueOf(ins)); err != nil {
		t.Fatalf("expected nil error for non-error return, got %v", err)
	}
}

// 覆盖 CallRunning 分支：len(rv)!=1
func TestMethods_CallRunning_LengthMismatch(t *testing.T) {
	typ := reflect.TypeOf(&runningMulti{})
	m, _ := typ.MethodByName("Running")
	mt := &Methods{obj: typ, runningMethod: &m}
	ins := &runningMulti{}
	if _, err := mt.CallRunning(reflect.ValueOf(ins)); err == nil {
		t.Fatalf("expected error for len(rv)!=1 in Running")
	}
}

// 覆盖 CallRunning 分支：返回值非 bool
func TestMethods_CallRunning_NotBool(t *testing.T) {
	typ := reflect.TypeOf(&badRunning{})
	m, _ := typ.MethodByName("Running")
	mt := &Methods{obj: typ, runningMethod: &m}
	ins := &badRunning{}
	if _, err := mt.CallRunning(reflect.ValueOf(ins)); err == nil {
		t.Fatalf("expected error for non-bool Running return")
	}
}

// --- AI GENERATED CODE END ---
