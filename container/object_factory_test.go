package container

import (
	"context"
	"testing"
	"time"
	"vortice/object"
)

// --- AI GENERATED CODE BEGIN ---

type depA struct{}
type depB struct{ sec int }
type coreObj struct{}

func factoryA(b *depB) *depA                { return &depA{} }
func factoryB() *depB                       { return &depB{sec: time.Now().Nanosecond()} }
func factoryCore(a *depA, b *depB) *coreObj { return &coreObj{} }

type TestContext struct{ context.Context }

func (t *TestContext) GetFilters() []object.DefinitionFilter { return []object.DefinitionFilter{} }
func (t *TestContext) GetObjects() map[string]Object         { return map[string]Object{} }

func addCoreTag(prop *object.Property) {
	prop.SetTag("autowired", "true")
}

func TestCoreObjectFactory_GetObject_Singleton(t *testing.T) {
	propA := object.NewProperty()
	propB := object.NewProperty()
	propCore := object.NewProperty()
	addCoreTag(propA)
	addCoreTag(propB)
	addCoreTag(propCore)

	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	_, _ = factory.RegisterFactory(factoryA, propA, false)
	_, _ = factory.RegisterFactory(factoryCore, propCore, false)

	if err := factory.Init(); err != nil {
		t.Fatalf("factory Init failed: %v", err)
	}

	ctx := &TestContext{}

	objs, err := factory.GetObjects(ctx, (*coreObj)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if len(objs) != 1 || objs[0] == nil || objs[0].Instance() == nil {
		t.Error("GetObjects should return a valid singleton object")
	}
	// 再次获取应返回同一个对象（singleton）
	objs2, err := factory.GetObjects(ctx, (*coreObj)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if objs[0].Instance() != objs2[0].Instance() {
		t.Error("Singleton scope should return same instance")
	}
}

func TestCoreObjectFactory_GetObject_Prototype(t *testing.T) {
	propB := object.NewProperty()
	propB.Scope = object.Prototype
	addCoreTag(propB)

	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)

	if err := factory.Init(); err != nil {
		t.Fatalf("factory Init failed: %v", err)
	}

	ctx := &TestContext{}

	objs1, err := factory.GetObjects(ctx, (*depB)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	objs2, err := factory.GetObjects(ctx, (*depB)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if len(objs1) != 1 || len(objs2) != 1 {
		t.Fatal("Prototype should return one object per call")
	}
	// 补充类型断言，确保实例类型正确
	if _, ok := objs1[0].Instance().(*depB); !ok {
		t.Error("Instance should be of type *depB")
	}
	if _, ok := objs2[0].Instance().(*depB); !ok {
		t.Error("Instance should be of type *depB")
	}
	// Prototype模式下，每次获取的实例应该不同
	if objs1[0].Instance() == objs2[0].Instance() {
		//t.Logf("obj1:%#v, obj2:%#v\n", objs1[0].Instance(), objs2[0].Instance())
		t.Error("Prototype scope should return different instances")
	}
}

func TestCoreObjectFactory_GetObject_NotFound(t *testing.T) {
	factory := NewCoreObjectFactory()
	_ = factory.Init()
	ctx := &TestContext{}
	objs, err := factory.GetObjects(ctx, (*struct{ X int })(nil))
	if err == nil || len(objs) != 0 {
		t.Error("GetObjects should fail for unknown type")
	}
}

func TestCoreObjectFactory_getType(t *testing.T) {
	factory := NewCoreObjectFactory()
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
	addCoreTag(propB)

	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	propA := object.NewProperty()
	addCoreTag(propA)
	_, _ = factory.RegisterFactory(factoryA, propA, false)

	if err := factory.Init(); err != nil {
		t.Fatalf("factory Init failed: %v", err)
	}

	ctx := &TestContext{}
	objs, err := factory.GetObjects(ctx, (*depA)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if len(objs) != 1 || !objs[0].Initialized() {
		t.Error("Dependency object should be initialized")
	}
}

func TestCoreObjectFactory_init_destroy(t *testing.T) {
	propB := object.NewProperty()
	addCoreTag(propB)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if err := factory.Destroy(); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}
}

// --- AI GENERATED CODE END ---
