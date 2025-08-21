package business

import (
	"errors"
	"testing"
	"vortice/object"
)

// 示例扩展类型
type DemoExt struct{}

type AnotherExt struct{}

func NewDemoExt() DemoExt       { return DemoExt{} }
func NewDemoExtV2() DemoExt     { return DemoExt{} }
func NewAnotherExt() AnotherExt { return AnotherExt{} }

// 构造 Definition（使用 Parser）
func mustParse(t *testing.T, fn any) *object.Definition {
	t.Helper()
	prop := object.NewProperty()
	def, err := object.NewParser(fn).Parse(prop)
	if err != nil {
		t.Fatalf("parse definition failed: %v", err)
	}
	return def
}

func Test_newPlugin_basic(t *testing.T) {
	p := NewPlugin("demo")
	if p.name != "demo" {
		t.Fatalf("name mismatch: %s", p.name)
	}
	if len(p.inits) != 0 {
		t.Fatalf("expect 0 inits")
	}
	if len(p.extensions) != 0 {
		t.Fatalf("expect 0 extensions")
	}
	if len(p.abilities) != 0 {
		t.Fatalf("expect 0 abilities")
	}
}

func TestPlugin_Name(t *testing.T) {
	if NewPlugin("x").Name() != "x" {
		t.Fatalf("Name() mismatch")
	}
}

func TestPlugin_Init_append(t *testing.T) {
	p := NewPlugin("demo")
	p.Init(func() error { return nil })
	p.Init(func() error { return nil }, func() error { return nil })
	if len(p.inits) != 3 {
		t.Fatalf("want 3 inits got %d", len(p.inits))
	}
}

func TestPlugin_addExtension_and_GetExtension(t *testing.T) {
	p := NewPlugin("demo")
	if p.GetExtension("not-exist") != nil {
		t.Fatalf("expect nil for absent extension")
	}

	def1 := mustParse(t, NewDemoExt)
	def2 := mustParse(t, NewDemoExtV2)  // 同返回类型，应被拒绝（第二次）
	def3 := mustParse(t, NewAnotherExt) // 不同返回类型，应允许

	// 第一次添加 DemoExt 成功
	if !p.addExtension(def1) {
		t.Fatalf("first addExtension def1 should succeed")
	}
	if got := p.GetExtension(def1.Name()); got != def1 {
		t.Fatalf("GetExtension mismatch for def1")
	}

	// 第二次添加同类型 DemoExt 工厂，应返回 false
	if p.addExtension(def2) {
		t.Fatalf("adding second factory with same return type should fail")
	}
	// 仍然只能拿到第一个定义
	if got := p.GetExtension(def1.Name()); got != def1 {
		t.Fatalf("original definition should remain after duplicate add attempt")
	}

	// 添加不同返回类型 AnotherExt 成功
	if !p.addExtension(def3) {
		t.Fatalf("adding different return type factory should succeed")
	}
	if got := p.GetExtension(def3.Name()); got != def3 {
		t.Fatalf("GetExtension mismatch for def3")
	}
}

func TestPlugin_init_success(t *testing.T) {
	p := NewPlugin("demo")
	order := []int{}
	p.Init(
		func() error { order = append(order, 1); return nil },
		func() error { order = append(order, 2); return nil },
		func() error { order = append(order, 3); return nil },
	)
	if err := p.init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	exp := []int{1, 2, 3}
	if len(order) != 3 {
		t.Fatalf("order length mismatch")
	}
	for i, v := range exp {
		if order[i] != v {
			t.Fatalf("order[%d]=%d want %d", i, order[i], v)
		}
	}
}

func TestPlugin_init_stopOnError(t *testing.T) {
	p := NewPlugin("demo")
	order := []int{}
	p.Init(
		func() error { order = append(order, 1); return nil },
		func() error { order = append(order, 2); return errors.New("boom") },
		func() error { order = append(order, 3); return nil },
	)
	err := p.init()
	if err == nil {
		t.Fatalf("expect error")
	}
	if len(order) != 2 {
		t.Fatalf("should stop at second init, got order len %d", len(order))
	}
	if order[0] != 1 || order[1] != 2 {
		t.Fatalf("execution order wrong: %+v", order)
	}
}

func TestPlugin_noInits(t *testing.T) {
	p := NewPlugin("demo")
	if err := p.init(); err != nil {
		t.Fatalf("no init funcs should succeed: %v", err)
	}
}

func TestPlugin_Integration(t *testing.T) {
	p := NewPlugin("integration-test")

	t.Run("完整插件生命周期", func(t *testing.T) {
		if p.Name() != "integration-test" {
			t.Fatalf("Plugin name = %v, want integration-test", p.Name())
		}
		initCalled := false
		p.Init(func() error { initCalled = true; return nil })
		def1 := mustParse(t, NewDemoExt)
		defDup := mustParse(t, NewDemoExtV2)
		defOther := mustParse(t, NewAnotherExt)
		if !p.addExtension(def1) {
			t.Fatalf("add ext1 failed")
		}
		// 重复类型应失败
		if p.addExtension(defDup) {
			t.Fatalf("duplicate return type added unexpectedly")
		}
		// 不同类型成功
		if !p.addExtension(defOther) {
			t.Fatalf("add different type failed")
		}
		if err := p.init(); err != nil {
			t.Fatalf("init error: %v", err)
		}
		if !initCalled {
			t.Fatalf("init func not called")
		}
		if p.GetExtension(def1.Name()) != def1 {
			t.Fatalf("ext1 mismatch")
		}
		if p.GetExtension(defOther.Name()) != defOther {
			t.Fatalf("extOther mismatch")
		}
		if p.GetExtension("not-exist") != nil {
			t.Fatalf("expect nil for non-exist")
		}
	})
}
