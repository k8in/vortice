package container

import (
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
	obj := NewObject(def, reflect.ValueOf(comp), comp)

	// Test Init
	if err := obj.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if !obj.Initialized() || !comp.initialized {
		t.Fatalf("Init flag or comp.initialized not set")
	}

	// Test Start
	if err := obj.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !obj.Running() || !comp.running {
		t.Fatalf("Running flag or comp.running not set")
	}

	// Test Stop
	if err := obj.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if obj.Running() || comp.running {
		t.Fatalf("Running flag or comp.running not cleared")
	}

	// Test Destroy
	if err := obj.Destroy(); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}
	if !comp.destroyed {
		t.Fatalf("Destroy flag not set")
	}
}

func TestCoreObject_ConcurrentAccess(t *testing.T) {
	comp := &testComp{}
	def := newTestDefinitionForObject()
	obj := NewObject(def, reflect.ValueOf(comp), comp)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = obj.Init()
			_ = obj.Start()
			_ = obj.Stop()
			_ = obj.Destroy()
			_ = obj.ID()
			_ = obj.Definition()
			_ = obj.Instance()
			_ = obj.Value()
			_ = obj.Initialized()
			_ = obj.Running()
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
	if def.Namespace() != object.NSCore {
		t.Error("Definition.Namespace should be Core by default")
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
	if !def.LazyInit() {
		t.Error("Definition.LazyInit should be true by default")
	}
	if def.Desc() != "" {
		t.Error("Definition.Desc should be empty by default")
	}
}
