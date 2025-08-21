package business

import (
	"context"
	"fmt"
	"testing"
	"vortice/container"
	"vortice/object"
)

// helper 创建一个全新 Core（不使用单例，避免全局状态污染）
func newCore() *Core {
	return NewCore(container.NewCore(context.Background()))
}

func TestDefaultCoreSingleton(t *testing.T) {
	a := DefaultCore()
	b := DefaultCore()
	if a != b {
		t.Fatalf("DefaultCore should return singleton instance")
	}
}

func TestRegisterPluginBeforeAndAfterInit(t *testing.T) {
	c := newCore()
	p := NewPlugin("p1")
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin before init failed: %v", err)
	}
	// 重复注册
	if err := c.RegisterPlugin(p); err == nil {
		t.Fatalf("expected duplicate plugin registration error")
	}
	// 初始化
	if err := c.Init(); err != nil {
		t.Fatalf("core init failed: %v", err)
	}
	// 只读后再注册
	if err := c.RegisterPlugin(NewPlugin("p2")); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode, got %v", err)
	}
}

func TestCoreInitReadonlyAndReinit(t *testing.T) {
	c := newCore()
	if err := c.Init(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}
	// 再次 Init
	if err := c.Init(); err != ErrInitialized {
		t.Fatalf("expected ErrInitialized, got %v", err)
	}
}

func TestPluginInitHooksExecuted(t *testing.T) {
	c := newCore()
	p := NewPlugin("hooks")
	calls := 0
	p.Init(func() error {
		calls++
		return nil
	}, func() error {
		calls++
		return nil
	})
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}
	if err := c.Init(); err != nil {
		t.Fatalf("core init failed: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 init hook calls, got %d", calls)
	}
}

func TestPluginInitErrorPropagates(t *testing.T) {
	c := newCore()
	p := NewPlugin("bad")
	p.Init(func() error {
		return fmt.Errorf("boom")
	})
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}
	err := c.Init()
	if err == nil {
		t.Fatalf("expected error from plugin init")
	}
	// 出错后仍为只读
	if err2 := c.RegisterPlugin(NewPlugin("later")); err2 != ErrInReadonlyMode {
		t.Fatalf("expected readonly after failed init, got %v", err2)
	}
}

// 新增：RegisterAbility 在只读模式下应返回 ErrInReadonlyMode
func TestRegisterAbilityReadonly(t *testing.T) {
	c := newCore()
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if _, err := c.RegisterAbility(nil, nil); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode, got %v", err)
	}
}

// 新增：核心生命周期 Start / Shutdown 正常流程
func TestCoreStartAndShutdown(t *testing.T) {
	c := newCore()
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if err := c.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	// 不期望 crash
	c.Shutdown()
}

// 新增：Shutdown 可重复调用（幂等性 / 不 panic）
func TestCoreShutdownIdempotent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("shutdown should be idempotent, but panicked: %v", r)
		}
	}()
	c := newCore()
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	c.Shutdown()
	c.Shutdown() // 第二次调用不应 panic
}

// 新增：���证 Start 后及 Shutdown 后 RegisterExtension / RegisterAbility 均被只读模式拒绝
func TestCoreRegistersDisallowedAfterStartAndShutdown(t *testing.T) {
	c := newCore()

	// Init 前：应允许注册 Extension / Ability
	type preInitExt struct{}
	newPreInitExt := func() *preInitExt { return &preInitExt{} }
	if _, err := c.RegisterExtension(newPreInitExt, object.NewProperty()); err != nil {
		t.Fatalf("pre-init RegisterExtension failed: %v", err)
	}
	type preInitAbility struct{}
	newPreInitAbility := func() *preInitAbility { return &preInitAbility{} }
	if _, err := c.RegisterAbility(newPreInitAbility, object.NewProperty()); err != nil {
		t.Fatalf("pre-init RegisterAbility failed: %v", err)
	}

	// 进入只读
	if err := c.Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if err := c.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	// Start 后：应被拒绝
	if _, err := c.RegisterExtension(func() *struct{} { return &struct{}{} }, object.NewProperty()); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode for RegisterExtension after start, got %v", err)
	}
	if _, err := c.RegisterAbility(func() *struct{} { return &struct{}{} }, object.NewProperty()); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode for RegisterAbility after start, got %v", err)
	}

	// Shutdown 后：仍应被拒绝
	c.Shutdown()
	if _, err := c.RegisterExtension(func() *struct{} { return &struct{}{} }, object.NewProperty()); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode for RegisterExtension after shutdown, got %v", err)
	}
	if _, err := c.RegisterAbility(func() *struct{} { return &struct{}{} }, object.NewProperty()); err != ErrInReadonlyMode {
		t.Fatalf("expected ErrInReadonlyMode for RegisterAbility after shutdown, got %v", err)
	}
}
