package container

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"vortice/object"
)

type testComp struct {
	initialized bool
	running     bool
	destroyed   bool
}

func (c *testComp) Init() error {
	c.initialized = true
	return nil
}
func (c *testComp) Destroy() error {
	c.destroyed = true
	return nil
}
func (c *testComp) Start() error {
	c.running = true
	return nil
}
func (c *testComp) Stop() error {
	c.running = false
	return nil
}
func (c *testComp) Running() bool {
	return c.running
}

func newTestDefinitionForObject() *object.Definition {
	// 使用 Parser 创建 Definition
	factory := func() *testComp { return &testComp{} }
	prop := object.NewProperty()
	def, err := object.NewParser(factory).Parse(prop)
	if err != nil {
		panic(err)
	}
	return def
}

func TestCoreObject_Lifecycle(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp))

	if err := obj.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if !obj.Initialized() || !comp.initialized {
		t.Fatalf("component should be initialized")
	}

	// Start
	if err := obj.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !obj.Running() || !comp.running {
		t.Fatalf("Running flag or comp.running not set")
	}

	// Stop
	if err := obj.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if obj.Running() || comp.running {
		t.Fatalf("Running flag or comp.running not cleared")
	}

	// Destroy
	if err := obj.Destroy(); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}
	if obj.Alive() {
		t.Fatalf("object should not be alive after destroy")
	}
}

func TestCoreObject_ConcurrentAccess(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = obj.Init() // 安全：已初始化或已销毁返回 nil / ErrAlreadyBeenDestroyed
			_ = obj.Start()
			_ = obj.Stop()
			_ = obj.Destroy()
			_ = obj.ID()
			_ = obj.Definition()
			_ = obj.Instance()
			_ = obj.Value()
			_ = obj.Initialized()
			_ = obj.Running()
			_ = obj.Alive()
		}()
	}
	wg.Wait()
}

func TestCoreObject_DefinitionFields(t *testing.T) {
	def := newTestDefinitionForObject()
	if def.Name() == "" {
		t.Error("Definition.Name should not be empty")
	}
	if def.Type() == nil {
		t.Error("Definition.Type should not be nil")
	}
	if def.Factory() == nil {
		t.Error("Definition.Factory should not be nil")
	}
	if def.Methods() == nil {
		t.Error("Definition.Methods should not be nil")
	}
	if def.Scope() != object.Singleton {
		t.Error("Definition.Scope should be Singleton by default")
	}
	if def.Tags() == nil {
		t.Error("Definition.Tags should not be nil")
	}
	if def.DependsOn() == nil {
		t.Error("Definition.DependsOn should not be nil")
	}
	// LazyInit 语义可能改变（若实现调整，这里不强制要求 true）
	if def.Desc() != "" {
		t.Error("Definition.Desc should be empty by default")
	}
}

// 额外类型：用于制造 Running 反射调用错误（Definition 来源 testComp，但实例为 otherComp）
type otherComp struct {
	initialized bool
	destroyed   bool
}

func (o *otherComp) Init() error    { o.initialized = true; return nil }
func (o *otherComp) Destroy() error { o.destroyed = true; return nil }

// 重复初始化（通过 Methods 调用）不应报错
func TestCoreObject_ReInit_Noop(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp))

	if err := obj.Init(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}
	if err := obj.Init(); err != nil {
		t.Fatalf("second init (noop) failed: %v", err)
	}
}

func TestCoreObject_DestroyTwiceAndAfterEffects(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp))

	_ = obj.Init()
	if !obj.Alive() {
		t.Fatalf("alive expected before destroy")
	}
	if err := obj.Destroy(); err != nil {
		t.Fatalf("first destroy failed: %v", err)
	}
	if obj.Alive() {
		t.Fatalf("alive should be false after first destroy")
	}
	if err := obj.Destroy(); !errors.Is(err, ErrAlreadyBeenDestroyed) {
		t.Fatalf("expected ErrAlreadyBeenDestroyed on second destroy, got %v", err)
	}
	if err := obj.Start(); !errors.Is(err, ErrAlreadyBeenDestroyed) {
		t.Fatalf("expected ErrAlreadyBeenDestroyed on Start after destroy, got %v", err)
	}
	if err := obj.Stop(); !errors.Is(err, ErrAlreadyBeenDestroyed) {
		t.Fatalf("expected ErrAlreadyBeenDestroyed on Stop after destroy, got %v", err)
	}
	if id := obj.ID(); id != "" {
		t.Fatalf("expected empty ID after destroy, got %q", id)
	}
	if obj.Running() {
		t.Fatalf("expected Running false after destroy")
	}
	// Definition / Instance / Value 应为置空或零值
	if obj.Definition() != nil {
		t.Fatalf("definition should be nil after destroy")
	}
	if obj.Instance() != nil {
		t.Fatalf("instance should be nil after destroy")
	}
	if obj.Value().IsValid() {
		t.Fatalf("value should be zero after destroy")
	}
	// 额外：Destroy 后 Init 分支（def=nil）
	if err := obj.Init(); !errors.Is(err, ErrAlreadyBeenDestroyed) {
		t.Fatalf("expected ErrAlreadyBeenDestroyed on Init after destroy, got %v", err)
	}
}

// Running 反射调用失败分支：Definition 针对 *testComp，实例为 *otherComp
func TestCoreObject_RunningErrorPath(t *testing.T) {
	def := newTestDefinitionForObject()
	oth := &otherComp{}
	obj := NewObject(def, reflect.ValueOf(oth))

	// Running 缺失方法 -> false
	if r := obj.Running(); r {
		t.Fatalf("expected Running false on method mismatch")
	}

	// Init 存在于 otherComp 中，应成功
	if err := obj.Init(); err != nil {
		t.Fatalf("did not expect init error, got %v", err)
	}
	if !oth.initialized || !obj.Initialized() {
		t.Fatalf("expected initialized flags true")
	}

	// Start 缺失 -> 返回错误
	if err := obj.Start(); err == nil {
		t.Fatalf("expected start error due to missing Start method on instance")
	}

	// Stop 同样缺失 -> 返回错误（忽略是否已经 start 失败）
	if err := obj.Stop(); err == nil {
		t.Fatalf("expected stop error due to missing Stop method on instance")
	}

	// Alive 仍为 true（未 Destroy）
	if !obj.Alive() {
		t.Fatalf("object should be alive (not destroyed) during mismatch tests")
	}
}

// 新增组件：Init 返回错误
type initErrComp struct{}

func (c *initErrComp) Init() error    { return errors.New("init fail") }
func (c *initErrComp) Destroy() error { return nil }

// 新增组件：Destroy 返回错误
type destroyErrComp struct{}

func (c *destroyErrComp) Init() error    { return nil }
func (c *destroyErrComp) Destroy() error { return errors.New("destroy fail") }

// 工厂生成 Definition
func newInitErrDefinition() *object.Definition {
	f := func() *initErrComp { return &initErrComp{} }
	prop := object.NewProperty()
	d, err := object.NewParser(f).Parse(prop)
	if err != nil {
		panic(err)
	}
	return d
}
func newDestroyErrDefinition() *object.Definition {
	f := func() *destroyErrComp { return &destroyErrComp{} }
	prop := object.NewProperty()
	d, err := object.NewParser(f).Parse(prop)
	if err != nil {
		panic(err)
	}
	return d
}

// 新增：Init 错误路径（CallInit 返回错误，不置 init 标志）
func TestCoreObject_InitErrorPath(t *testing.T) {
	def := newInitErrDefinition()
	comp := &initErrComp{}
	obj := NewObject(def, reflect.ValueOf(comp))
	err := obj.Init()
	if err == nil || err.Error() != "init fail" {
		t.Fatalf("expected init fail error, got %v", err)
	}
	if obj.Initialized() {
		t.Fatalf("initialized flag should remain false on init error")
	}
	if !obj.Alive() {
		t.Fatalf("object should still be alive after init error")
	}
}

// 新增：Destroy 错误路径（CallDestroy 返回错误，不清理 def）
func TestCoreObject_DestroyErrorPath(t *testing.T) {
	def := newDestroyErrDefinition()
	comp := &destroyErrComp{}
	obj := NewObject(def, reflect.ValueOf(comp))
	if err := obj.Init(); err != nil {
		t.Fatalf("init should succeed, got %v", err)
	}
	if err := obj.Destroy(); err == nil || err.Error() != "destroy fail" {
		t.Fatalf("expected destroy fail error, got %v", err)
	}
	// 仍存活
	if !obj.Alive() {
		t.Fatalf("object should remain alive after destroy error")
	}
	// 再次 Destroy 仍错误（重复覆盖错误路径）
	if err := obj.Destroy(); err == nil {
		t.Fatalf("expected destroy fail again")
	}
	// 成员仍未被清空
	if obj.Definition() == nil {
		t.Fatalf("definition should remain after destroy error")
	}
	if obj.Instance() == nil {
		t.Fatalf("instance should remain after destroy error")
	}
}

// 新增：ID 正常值与销毁后空值已覆盖，这里补一次直接读取
func TestCoreObject_IDBeforeDestroy(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp))
	if obj.ID() == "" {
		t.Fatalf("expected non-empty ID before destroy")
	}
}
