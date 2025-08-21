package business

import (
	"context"
	"errors"
	"testing"

	"vortice/container"
	"vortice/object"
)

// 测试用类型
type mainExtType struct{}

func newMainExtType() mainExtType    { return mainExtType{} }
func newMainExtTypeDup() mainExtType { return mainExtType{} }

type pluginExtType struct{}

func newPluginExtType() pluginExtType    { return pluginExtType{} }
func newPluginExtTypeDup() pluginExtType { return pluginExtType{} }

type abilityType struct{}

func newAbilityA() abilityType { return abilityType{} }
func newAbilityB() abilityType { return abilityType{} }

// 工具
func hasTag(tags []object.Tag, k, v string) bool {
	for _, tg := range tags {
		if tg.Key() == k && tg.Value() == v {
			return true
		}
	}
	return false
}

func newCoreForTest() *Core { return NewCore(container.NewCore(context.Background())) }

func TestCore_RegisterPlugin(t *testing.T) {
	c := newCoreForTest()
	if err := c.RegisterPlugin(nil); !errors.Is(err, ErrNilPlugin) {
		t.Fatalf("expected ErrNilPlugin, got %v", err)
	}
	p := NewPlugin("p1")
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if err := c.RegisterPlugin(p); err == nil {
		t.Fatalf("duplicate same instance should fail")
	}
	if err := c.RegisterPlugin(NewPlugin("p1")); err == nil {
		t.Fatalf("duplicate same name should fail")
	}
}

func TestCore_RegisterExtension_Main(t *testing.T) {
	c := newCoreForTest()
	def, err := c.RegisterExtension(newMainExtType, object.NewProperty())
	if err != nil {
		t.Fatalf("register main ext failed: %v", err)
	}
	if !hasTag(def.Tags(), "namespace", MainNamespace) {
		t.Fatalf("missing namespace tag")
	}
	if !hasTag(def.Tags(), TagBizKindKey, "extension") {
		t.Fatalf("missing biz_kind=extension")
	}
	_, err = c.RegisterExtension(newMainExtTypeDup, object.NewProperty())
	if !errors.Is(err, ErrRegisterMainExt) {
		t.Fatalf("duplicate main ext should ErrRegisterMainExt, got %v", err)
	}
}

func TestCore_RegisterExtension_Plugin(t *testing.T) {
	c := newCoreForTest()
	p := NewPlugin("plugX")
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}
	// 人工设置 current 进入插件上下文
	c.current.Store(p)
	def, err := c.RegisterExtension(newPluginExtType, object.NewProperty())
	if err != nil {
		t.Fatalf("plugin ext register failed: %v", err)
	}
	if !hasTag(def.Tags(), "namespace", p.Name()) {
		t.Fatalf("plugin ext missing namespace tag")
	}
	if !hasTag(def.Tags(), TagBizKindKey, "extension") {
		t.Fatalf("plugin ext missing biz_kind tag")
	}
	if _, err := c.RegisterExtension(newPluginExtTypeDup, object.NewProperty()); !errors.Is(err, ErrRegisterPluginExt) {
		t.Fatalf("duplicate plugin ext expect ErrRegisterPluginExt got %v", err)
	}
	if p.GetExtension(def.Name()) != def {
		t.Fatalf("plugin should hold extension")
	}
	// 清除上下文
	c.current.Store(nil)
}

func TestCore_RegisterAbility(t *testing.T) {
	c := newCoreForTest()
	defA, err := c.RegisterAbility(newAbilityA, object.NewProperty())
	if err != nil {
		t.Fatalf("ability A register failed: %v", err)
	}
	defB, err := c.RegisterAbility(newAbilityB, object.NewProperty())
	if err != nil {
		t.Fatalf("ability B register failed: %v", err)
	}
	if !hasTag(defA.Tags(), TagBizKindKey, "ability") {
		t.Fatalf("missing ability biz_kind tag")
	}
	if !hasTag(defA.Tags(), "namespace", MainNamespace) {
		t.Fatalf("missing ability namespace tag")
	}
	list := c.abilities[defA.Name()]
	if len(list) != 2 {
		t.Fatalf("expected 2 ability defs got %d", len(list))
	}
	_ = defB
}

func TestCore_Init_Success_Readonly(t *testing.T) {
	c := newCoreForTest()
	p := NewPlugin("plugInit")
	p.Init(func() error { _, err := c.RegisterExtension(newPluginExtType, object.NewProperty()); return err })
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}
	if _, err := c.RegisterExtension(newMainExtType, object.NewProperty()); err != nil {
		t.Fatalf("pre-init main ext failed: %v", err)
	}
	if err := c.Init(); err != nil {
		t.Fatalf("core init failed: %v", err)
	}
	if err := c.Init(); !errors.Is(err, ErrInitialized) {
		t.Fatalf("second init expect ErrInitialized got %v", err)
	}
	if _, err := c.RegisterExtension(newMainExtTypeDup, object.NewProperty()); !errors.Is(err, ErrInReadonlyMode) {
		t.Fatalf("readonly register ext expect ErrInReadonlyMode got %v", err)
	}
	if _, err := c.RegisterAbility(newAbilityA, object.NewProperty()); !errors.Is(err, ErrInReadonlyMode) {
		t.Fatalf("readonly ability expect ErrInReadonlyMode got %v", err)
	}
	if err := c.RegisterPlugin(NewPlugin("newP")); !errors.Is(err, ErrInReadonlyMode) {
		t.Fatalf("readonly plugin expect ErrInReadonlyMode got %v", err)
	}
}

func TestCore_Init_PluginInitError(t *testing.T) {
	c := newCoreForTest()
	p := NewPlugin("bad")
	p.Init(func() error { return errors.New("boom") })
	if err := c.RegisterPlugin(p); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}
	if err := c.Init(); err == nil {
		t.Fatalf("expected init error")
	}
	// 已只读
	if _, err := c.RegisterExtension(newMainExtType, object.NewProperty()); !errors.Is(err, ErrInReadonlyMode) {
		t.Fatalf("after failed init expect readonly ext err got %v", err)
	}
}

func TestCore_Start_Shutdown(t *testing.T) {
	c := newCoreForTest()
	if err := c.core.Init(); err != nil {
		t.Fatalf("container core init failed: %v", err)
	}
	if err := c.Start(); err != nil {
		t.Fatalf("core start failed: %v", err)
	}
	c.Shutdown() // 只验证不 panic
}
