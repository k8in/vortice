package container

import (
	"context"
	"reflect"
	"testing"
	"vortice/object"
)

// --- AI GENERATED CODE BEGIN ---

type depA struct{}
type depB struct{}
type coreObj struct{}

func factoryA(b *depB) *depA                { return &depA{} }
func factoryB() *depB                       { return &depB{} }
func factoryCore(a *depA, b *depB) *coreObj { return &coreObj{} }

// TestContext 用于测试，模拟 Context 行为
type TestContext struct{ context.Context }

func (t *TestContext) GetFilters() []object.DefinitionFilter { return nil }
func (t *TestContext) GetObjects() map[string]Object         { return map[string]Object{} }

func TestCoreObjectFactory_GetObject_Singleton(t *testing.T) {
	propA := object.NewProperty()
	propB := object.NewProperty()
	propCore := object.NewProperty()

	_, _ = object.RegisterFactory(factoryB, propB, true)
	_, _ = object.RegisterFactory(factoryA, propA, true)
	_, _ = object.RegisterFactory(factoryCore, propCore, true)

	factory := NewObjectFactory()
	ctx := &TestContext{}

	obj, err := factory.GetObject(ctx, (*coreObj)(nil))
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	if obj == nil || obj.Instance() == nil {
		t.Error("GetObject should return a valid object")
	}
	// 再次获取应返回同一个对象（singleton）
	obj2, err := factory.GetObject(ctx, (*coreObj)(nil))
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	if obj.Instance() != obj2.Instance() {
		t.Error("Singleton scope should return same instance")
	}
}

func TestCoreObjectFactory_GetObject_Prototype(t *testing.T) {
	propB := object.NewProperty()
	propB.Scope = object.Prototype
	_, _ = object.RegisterFactory(factoryB, propB, true)

	factory := NewObjectFactory()
	ctx := &TestContext{}

	obj1, err := factory.GetObject(ctx, (*depB)(nil))
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	obj2, err := factory.GetObject(ctx, (*depB)(nil))
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	if obj1.Instance() == obj2.Instance() {
		t.Error("Prototype scope should return different instances")
	}
}

func TestCoreObjectFactory_GetObject_NotFound(t *testing.T) {
	factory := NewObjectFactory()
	ctx := &TestContext{}
	_, err := factory.GetObject(ctx, (*struct{ X int })(nil))
	if err == nil {
		t.Error("GetObject should fail for unknown type")
	}
}

func TestCoreObjectFactory_getType(t *testing.T) {
	factory := NewObjectFactory().(*CoreObjectFactory)
	if factory.getType((*depA)(nil)) == nil {
		t.Error("getType should return type for pointer to struct")
	}
	if factory.getType(depA{}) != nil {
		t.Error("getType should return nil for non-pointer")
	}
	if factory.getType(nil) != nil {
		t.Error("getType should return nil for nil")
	}
}

func TestCoreObjectFactory_newObject_DependencyInit(t *testing.T) {
	propB := object.NewProperty()
	_, _ = object.RegisterFactory(factoryB, propB, true)
	propA := object.NewProperty()
	_, _ = object.RegisterFactory(factoryA, propA, true)

	factory := NewObjectFactory()
	ctx := &TestContext{}
	obj, err := factory.GetObject(ctx, (*depA)(nil))
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	if !obj.Initialized() {
		t.Error("Dependency object should be initialized")
	}
}

func TestCoreObjectFactory_init_destroy(t *testing.T) {
	propB := object.NewProperty()
	_, _ = object.RegisterFactory(factoryB, propB, true)
	factory := NewObjectFactory().(*CoreObjectFactory)
	if err := factory.init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if err := factory.destroy(); err != nil {
		t.Fatalf("destroy failed: %v", err)
	}
}

func TestCoreObjectFactory_getDefinition(t *testing.T) {
	propB := object.NewProperty()
	_, _ = object.RegisterFactory(factoryB, propB, true)
	factory := NewObjectFactory().(*CoreObjectFactory)
	def, err := factory.getDefinition(object.GenerateDefinitionName(reflect.TypeOf((*depB)(nil))))
	if err != nil || def == nil {
		t.Error("getDefinition should return definition")
	}
	def2, err := factory.getDefinition("notfound")
	if err == nil || def2 != nil {
		t.Error("getDefinition should fail for unknown name")
	}
}

// --- AI GENERATED CODE END ---
