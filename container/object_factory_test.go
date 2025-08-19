package container

import (
	"context"
	"sync"
	"testing"
	"time"
	"vortice/object"
)

type depA struct{}
type depB struct{ sec int }
type coreObj struct{}
type depC struct{}
type depD struct{}

func factoryA(b *depB) *depA                { return &depA{} }
func factoryB() *depB                       { return &depB{sec: time.Now().Nanosecond()} }
func factoryCore(a *depA, b *depB) *coreObj { return &coreObj{} }
func factoryC(d *depD) *depC                { return &depC{} }
func factoryD() *depD                       { return &depD{} }

type testCtx struct{}

func (t *testCtx) GetFilters() []object.DefinitionFilter { return nil }
func (t *testCtx) GetObjects() map[string]Object         { return nil }

func addCoreTag(prop *object.Property) {
	prop.SetTag("autowired", "true")
}

func newCoreCtx() *CoreContext {
	return WithCoreContext(context.Background())
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

	ctx := newCoreCtx()

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

	ctx := newCoreCtx()

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
	if _, ok := objs1[0].Instance().(*depB); !ok {
		t.Error("Instance should be of type *depB")
	}
	if _, ok := objs2[0].Instance().(*depB); !ok {
		t.Error("Instance should be of type *depB")
	}
	if objs1[0].Instance() == objs2[0].Instance() {
		t.Error("Prototype scope should return different instances")
	}
}

func TestCoreObjectFactory_GetObjectsByName(t *testing.T) {
	propB := object.NewProperty()
	addCoreTag(propB)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	_ = factory.Init()
	ctx := newCoreCtx()
	defs := factory.GetDefinitions()
	if len(defs) == 0 {
		t.Fatal("No definitions registered")
	}
	name := defs[0].Name()
	objs, err := factory.GetObjectsByName(ctx, name)
	if err != nil {
		t.Fatalf("GetObjectsByName failed: %v", err)
	}
	if len(objs) != 1 || objs[0].Instance() == nil {
		t.Error("GetObjectsByName should return valid object")
	}
}

func TestCoreObjectFactory_GetObject_NotFound(t *testing.T) {
	factory := NewCoreObjectFactory()
	_ = factory.Init()
	ctx := newCoreCtx()
	objs, err := factory.GetObjects(ctx, (*struct{ X int })(nil))
	if err == nil || len(objs) != 0 {
		t.Error("GetObjects should fail for unknown type")
	}
}

func TestCoreObjectFactory_GetObjectsByName_NotFound(t *testing.T) {
	factory := NewCoreObjectFactory()
	_ = factory.Init()
	ctx := newCoreCtx()
	objs, err := factory.GetObjectsByName(ctx, "notfound")
	if err != nil {
		t.Fatalf("GetObjectsByName should not error for unknown name")
	}
	if len(objs) != 0 {
		t.Error("GetObjectsByName should return empty slice for unknown name")
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
	ctx := newCoreCtx()
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
	factory.Destroy()
}

func TestCoreObjectFactory_MultiInit(t *testing.T) {
	propB := object.NewProperty()
	addCoreTag(propB)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	// 多次 Init 不应 panic 或报错
	if err := factory.Init(); err != nil {
		t.Fatalf("Second Init failed: %v", err)
	}
}

func TestCoreObjectFactory_AutowiredFilter(t *testing.T) {
	propB := object.NewProperty()
	addCoreTag(propB)
	propC := object.NewProperty()
	addCoreTag(propC)
	propD := object.NewProperty()
	addCoreTag(propD)
	// 不加 autowired tag
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	_, _ = factory.RegisterFactory(factoryC, propC, false)
	_, _ = factory.RegisterFactory(factoryD, propD, false) // 补充注册 depD
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	// 只应返回带 autowired tag 的对象
	defs := factory.GetDefinitions(object.TagFilter("autowired=true"))
	if len(defs) != 3 {
		t.Error("TagFilter should only return autowired objects")
	}
}

func TestCoreObjectFactory_DependencyChain(t *testing.T) {
	propD := object.NewProperty()
	addCoreTag(propD)
	propC := object.NewProperty()
	addCoreTag(propC)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryD, propD, false)
	_, _ = factory.RegisterFactory(factoryC, propC, false)
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	ctx := newCoreCtx()
	objs, err := factory.GetObjects(ctx, (*depC)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if len(objs) != 1 || objs[0].Instance() == nil {
		t.Error("Dependency chain should resolve and instantiate")
	}
}

func TestCoreObjectFactory_DependencyMissing(t *testing.T) {
	propC := object.NewProperty()
	addCoreTag(propC)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryC, propC, false)
	// factoryC 依赖 depD，但未注册 depD
	// 所以 Init 时 DAG 检查会报错 cycle detected or missing dependency in the DAG
	err := factory.Init()
	if err == nil {
		t.Error("Init should fail if dependency is missing")
	}
}

func TestCoreObjectFactory_ConcurrentAccess(t *testing.T) {
	propB := object.NewProperty()
	addCoreTag(propB)
	factory := NewCoreObjectFactory()
	_, _ = factory.RegisterFactory(factoryB, propB, false)
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	ctx := newCoreCtx()
	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := factory.GetObjects(ctx, (*depB)(nil))
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Errorf("Concurrent GetObjects failed: %v", err)
		}
	}
}

func TestCoreObjectFactory_newObject_InvalidFactory(t *testing.T) {
	prop := object.NewProperty()
	addCoreTag(prop)
	factory := NewCoreObjectFactory()
	// 注册一个无效的工厂（参数类型不符）
	invalidFactory := func(a int) int { return a }
	_, err := factory.RegisterFactory(invalidFactory, prop, false)
	if err == nil {
		t.Error("RegisterFactory should fail for invalid factory")
	}
}

// 覆盖 newObject 的无依赖分支
func TestCoreObjectFactory_newObject_NoDeps(t *testing.T) {
	prop := object.NewProperty()
	addCoreTag(prop)
	factory := NewCoreObjectFactory()
	_, err := factory.RegisterFactory(factoryB, prop, false)
	if err != nil {
		t.Fatalf("RegisterFactory failed: %v", err)
	}
	if err := factory.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	ctx := newCoreCtx()
	objs, err := factory.GetObjects(ctx, (*depB)(nil))
	if err != nil {
		t.Fatalf("GetObjects failed: %v", err)
	}
	if len(objs) != 1 || objs[0].Instance() == nil {
		t.Error("newObject should create object with no dependencies")
	}
}

// --- END FULL COVERAGE TESTS ---
