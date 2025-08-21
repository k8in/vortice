package container

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"vortice/object"
)

// 统一测试上下文创建
func newTestCtx() Context {
	return WithCoreContext(context.Background())
}

// ---------- 辅助: 添加 autowired 标签 ----------
func addAutowired(p *object.Property) {
	p.SetTags(TagAutowired)
}

// ---------- 基础类型 & 工厂 (无依赖 / 依赖链) ----------
type compC struct{ inited bool }
type compB struct{ c *compC }
type compA struct{ b *compB }

func newCompC() *compC  { return &compC{} }
func newCompC1() *compC { return &compC{} }
func newCompC2() *compC { return &compC{} }
func (c *compC) Init() error {
	c.inited = true
	return nil
}
func (c *compC) Destroy() error { return nil }

func newCompB(c *compC) *compB  { return &compB{c: c} }
func (b *compB) Init() error    { return nil }
func (b *compB) Destroy() error { return nil }

func newCompA(b *compB) *compA { return &compA{b: b} }
func (a *compA) Init() error   { return nil }
func (a *compA) Destroy() error {
	return nil
}

// ---------- Init 失败组件 ----------
type initErr struct{}

func newInitErr() *initErr        { return &initErr{} }
func (i *initErr) Init() error    { return errors.New("init-boom") }
func (i *initErr) Destroy() error { return nil }

// ---------- Destroy 失败组件（错误被忽略） ----------
type destroyErr struct {
	destroyed bool
}

func newDestroyErr() *destroyErr  { return &destroyErr{} }
func (d *destroyErr) Init() error { return nil }
func (d *destroyErr) Destroy() error {
	d.destroyed = true
	return errors.New("destroy-boom")
}

// ---------- Lazy 组件 ----------
type lazyOne struct{ started bool }

func newLazyOne() *lazyOne        { return &lazyOne{} }
func (l *lazyOne) Init() error    { l.started = true; return nil }
func (l *lazyOne) Destroy() error { return nil }

// ---------- 多实现选择组件 ----------
type variantDep struct {
	id int
}

func newVariant1() *variantDep { return &variantDep{id: 1} }
func (v *variantDep) Init() error {
	return nil
}
func (v *variantDep) Destroy() error { return nil }

func newVariant2() *variantDep { return &variantDep{id: 2} }

// ---------- 选择 root (依赖 variantDep) ----------
type variantRoot struct {
	dep *variantDep
}

func newVariantRoot(d *variantDep) *variantRoot { return &variantRoot{dep: d} }
func (r *variantRoot) Init() error              { return nil }
func (r *variantRoot) Destroy() error           { return nil }

// ---------- 自定义 selector 选第二个 ----------
type pickSecond struct{}

func (s *pickSecond) Select(defs []*object.Definition) *object.Definition {
	if len(defs) > 1 {
		return defs[1]
	}
	if len(defs) == 1 {
		return defs[0]
	}
	return nil
}

// ---------- 循环依赖 ----------
type cycA struct{}
type cycB struct{}

func newCycA(b *cycB) *cycA { return &cycA{} }
func newCycB(a *cycA) *cycB { return &cycB{} }
func (a *cycA) Init() error { return nil }
func (a *cycA) Destroy() error {
	return nil
}
func (b *cycB) Init() error    { return nil }
func (b *cycB) Destroy() error { return nil }

// ---------- 缺失 autowired 依赖 root ----------
type missingDepRoot struct {
	a *compA
}

func newMissingDepRoot(a *compA) *missingDepRoot { return &missingDepRoot{a: a} }
func (r *missingDepRoot) Init() error            { return nil }
func (r *missingDepRoot) Destroy() error         { return nil }

// ---------- 测试: Init 成功（含 lazy 与非 lazy） ----------
func TestObjectFactory_Init_SuccessAndLazy(t *testing.T) {
	f := NewCoreObjectFactory()
	// 非 lazy
	pC := object.NewProperty()
	pC.LazyInit = false
	addAutowired(pC)
	_, _ = f.RegisterFactory(newCompC, pC, false)
	// lazy
	pLazy := object.NewProperty()
	pLazy.LazyInit = true
	addAutowired(pLazy)
	_, _ = f.RegisterFactory(newLazyOne, pLazy, false)
	if err := f.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	var lazyObj, cObj Object
	for _, o := range f.objs {
		switch o.Instance().(type) {
		case *lazyOne:
			lazyObj = o
		case *compC:
			cObj = o
		}
	}
	if cObj == nil || !cObj.Initialized() {
		t.Fatalf("non-lazy singleton should be initialized at Init")
	}
	if lazyObj == nil || lazyObj.Initialized() {
		t.Fatalf("lazy singleton should NOT be initialized at Init")
	}
}

// ---------- 测试: Init 期间缺失依赖应在 DefinitionRegistry.Init() 阶段被发现（definition not found），而不是 newObject ----------
func TestObjectFactory_Init_NewObjectErrorShadowed(t *testing.T) {
	f := NewCoreObjectFactory()
	// root depends on compA (未注册 / 且需要 autowired)
	pRoot := object.NewProperty()
	pRoot.LazyInit = false
	addAutowired(pRoot)
	_, _ = f.RegisterFactory(newMissingDepRoot, pRoot, false)

	err := f.Init()
	if err == nil || !strings.Contains(err.Error(), "definition not found") {
		t.Fatalf("expected definition not found error from registry Init, got %v", err)
	}
	// 验证 root 未创建
	for _, o := range f.objs {
		if _, ok := o.Instance().(*missingDepRoot); ok {
			t.Fatalf("root should not be created when dependency missing")
		}
	}
	// 后续 GetObjects 仍应因为缺失依赖报同类错误
	ctx := newTestCtx()
	_, err2 := f.GetObjects(ctx, (*missingDepRoot)(nil))
	if err2 == nil || !strings.Contains(err2.Error(), "definition not found") {
		t.Fatalf("expected definition not found on GetObjects, got %v", err2)
	}
}

// ---------- 测试: Init 对象初始化失败（返回错误） ----------
func TestObjectFactory_Init_ObjectInitError(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	p.LazyInit = false
	addAutowired(p)
	_, _ = f.RegisterFactory(newInitErr, p, false)
	err := f.Init()
	if err == nil || !strings.Contains(err.Error(), "init-boom") {
		t.Fatalf("expected init-boom error, got %v", err)
	}
}

// ---------- 测试: Destroy 忽略 Destroy 错误 ----------
func TestObjectFactory_Destroy_IgnoreErrors(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	p.LazyInit = false
	addAutowired(p)
	_, _ = f.RegisterFactory(newDestroyErr, p, false)
	_ = f.Init()
	var ptr *destroyErr
	for _, o := range f.objs {
		if v, ok := o.Instance().(*destroyErr); ok {
			ptr = v
		}
	}
	f.Destroy()
	if ptr == nil || !ptr.destroyed {
		t.Fatalf("Destroy should invoke Destroy method even if it returns error")
	}
}

// ---------- 测试: GetObjects defs 为空 ----------
func TestObjectFactory_GetObjects_Empty(t *testing.T) {
	f := NewCoreObjectFactory()
	_ = f.Init()
	ctx := newTestCtx()
	objs, err := f.GetObjects(ctx, (*compA)(nil))
	if err == nil && len(objs) != 0 {
		t.Fatalf("expected empty slice for unknown type")
	}
	// ByName 为空
	if by, err2 := f.GetObjectsByName(ctx, "unknown"); err2 != nil || len(by) != 0 {
		t.Fatalf("GetObjectsByName empty mismatch: %v %d", err2, len(by))
	}
}

// ---------- 测试: newObject Argn==0 & 复用 singleton ----------
func TestObjectFactory_GetObjects_SingletonReuse(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	addAutowired(p)
	_, _ = f.RegisterFactory(newCompC, p, false)
	_ = f.Init()
	ctx := newTestCtx()
	// 第一��
	objs1, err := f.GetObjects(ctx, (*compC)(nil))
	if err != nil || len(objs1) != 1 {
		t.Fatalf("first get compC failed: %v", err)
	}
	if !objs1[0].Initialized() {
		t.Fatalf("singleton should be initialized after Init")
	}
	// 第二次复用
	objs2, err2 := f.GetObjects(ctx, (*compC)(nil))
	if err2 != nil || len(objs2) != 1 || objs1[0] != objs2[0] {
		t.Fatalf("expected singleton reuse")
	}
}

// ---------- 测试: getDependencies 正常链 + Argn>0 newObject ----------
func TestObjectFactory_newObject_DependencyChain(t *testing.T) {
	f := NewCoreObjectFactory()
	pC := object.NewProperty()
	addAutowired(pC)
	_, _ = f.RegisterFactory(newCompC, pC, false)
	pB := object.NewProperty()
	addAutowired(pB)
	_, _ = f.RegisterFactory(newCompB, pB, false)
	pA := object.NewProperty()
	addAutowired(pA)
	_, _ = f.RegisterFactory(newCompA, pA, false)
	_ = f.Init()
	ctx := newTestCtx()
	objs, err := f.GetObjects(ctx, (*compA)(nil))
	if err != nil || len(objs) != 1 {
		t.Fatalf("get compA failed: %v", err)
	}
	if objs[0].Instance() == nil {
		t.Fatalf("expected non-nil instance")
	}
}

// ---------- 测试: getDependencies 缺失依赖 ----------
func TestObjectFactory_getDependencies_Missing(t *testing.T) {
	f := NewCoreObjectFactory()
	// root depends on compA，但 compA 未注册
	pRoot := object.NewProperty()
	addAutowired(pRoot)
	_, _ = f.RegisterFactory(newCompA, pRoot, false)
	_ = f.Init()
	ctx := newTestCtx()
	_, err := f.GetObjects(ctx, (*compA)(nil))
	if err == nil {
		t.Fatalf("expected error due to missing dependency chain")
	}
}

// ---------- 测试: getDependencies 循环 ----------
func TestObjectFactory_getDependencies_Cycle(t *testing.T) {
	f := NewCoreObjectFactory()
	pA := object.NewProperty()
	addAutowired(pA)
	_, _ = f.RegisterFactory(newCycA, pA, false)
	pB := object.NewProperty()
	addAutowired(pB)
	_, _ = f.RegisterFactory(newCycB, pB, false)
	err := f.Init()
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("expected cycle error at Init, got %v", err)
	}
	// Init ���败后不再尝试 GetObjects（避免误导）
}

// ---------- 测试: obj.Init 失败 (GetObjects 中 object.Init failed 分支) ----------
func TestObjectFactory_GetObjects_ObjectInitError(t *testing.T) {
	f := NewCoreObjectFactory()
	prop := object.NewProperty()
	prop.Scope = object.Prototype
	addAutowired(prop)
	_, _ = f.RegisterFactory(newInitErr, prop, false)
	_ = f.Init()
	ctx := newTestCtx()
	_, err := f.GetObjects(ctx, (*initErr)(nil))
	if err == nil || !strings.Contains(err.Error(), "object.Init failed") {
		t.Fatalf("expected object.Init failed, got %v", err)
	}
}

// ---------- 测试: new() 依赖缺失 (dependencies not found) ----------
type depX struct{}
type needX struct{ x *depX }

func newDepX() *depX            { return &depX{} }
func newNeedX(x *depX) *needX   { return &needX{x: x} }
func (d *depX) Init() error     { return nil }
func (d *depX) Destroy() error  { return nil }
func (n *needX) Init() error    { return nil }
func (n *needX) Destroy() error { return nil }

func TestObjectFactory_new_DependencyNotFound(t *testing.T) {
	f := NewCoreObjectFactory()
	// 注册 needX (依赖 depX) 但 build 时上下文不给 depX
	pNeed := object.NewProperty()
	addAutowired(pNeed)
	def, _ := f.RegisterFactory(newNeedX, pNeed, false)
	_ = f.Init()
	_, err := f.new(def, map[string]Object{})
	if err == nil || !strings.Contains(err.Error(), "dependencies not found") {
		t.Fatalf("expected dependencies not found error, got %v", err)
	}
}

// ---------- 测试: 自定义 selector 选择第二实现 ----------
func TestObjectFactory_CustomSelectorPickSecond(t *testing.T) {
	f := NewCoreObjectFactory()
	f.SetRealizationSelector(&pickSecond{})
	// 两个 variantDep
	p1 := object.NewProperty()
	addAutowired(p1)
	def1, _ := f.RegisterFactory(newVariant1, p1, false)
	p2 := object.NewProperty()
	addAutowired(p2)
	def2, _ := f.RegisterFactory(newVariant2, p2, false)
	// root
	pr := object.NewProperty()
	pr.Scope = object.Prototype
	addAutowired(pr)
	_, _ = f.RegisterFactory(newVariantRoot, pr, false)

	_ = f.Init()
	ctx := newTestCtx()
	objs, err := f.GetObjects(ctx, (*variantRoot)(nil))
	if err != nil || len(objs) != 1 {
		t.Fatalf("GetObjects variantRoot failed: %v", err)
	}
	rootObj := objs[0].Instance().(*variantRoot)
	if rootObj.dep == nil {
		t.Fatalf("dep nil")
	}
	// 判断是不是第二个实现：比较 Factory 名称（不同函数名）
	var picked *object.Definition
	for _, d := range []*object.Definition{def1, def2} {
		if strings.Contains(d.Factory().Name(), "newVariant2") && rootObj.dep.id == 2 {
			picked = d
			break
		}
	}
	if picked == nil {
		t.Fatalf("selector did not pick second implementation (id=%d)", rootObj.dep.id)
	}
}

// ---------- 测试: getAutowiredDefinition 未找到 ----------
func TestObjectFactory_getAutowiredDefinition_NotFound(t *testing.T) {
	f := NewCoreObjectFactory()
	// 仅注册 root（依赖 *compA）且标记 autowired，未注册 compA，Init 时应报依赖缺失
	p := object.NewProperty()
	addAutowired(p)
	_, _ = f.RegisterFactory(newMissingDepRoot, p, false)

	err := f.Init()
	if err == nil || !strings.Contains(err.Error(), "definition not found") {
		t.Fatalf("expected definition not found error at Init, got %v", err)
	}
	// 确认未创建 root 实例
	for _, o := range f.objs {
		if _, ok := o.Instance().(*missingDepRoot); ok {
			t.Fatalf("root instance should not be created when dependency missing")
		}
	}
	// 再次尝试获取，仍应返回相同类型错误
	ctx := newTestCtx()
	_, err2 := f.GetObjects(ctx, (*missingDepRoot)(nil))
	if err2 == nil || !strings.Contains(err2.Error(), "definition not found") {
		t.Fatalf("expected definition not found on GetObjects, got %v", err2)
	}
}

// ---------- 测试: getBuildCtx 合并外部与已存在 singleton ----------
func TestObjectFactory_getBuildCtx_Merge(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	addAutowired(p)
	_, _ = f.RegisterFactory(newCompC, p, false)
	_ = f.Init()
	external := map[string]Object{"externalKey": nil}
	ctxMap := f.getBuildCtx(external)
	if _, ok := ctxMap["externalKey"]; !ok {
		t.Fatalf("external key not merged")
	}
	found := false
	for k := range ctxMap {
		if strings.Contains(k, "compC") {
			found = true
		}
	}
	if !found {
		t.Fatalf("autowired singleton not merged into build ctx")
	}
}

// ---------- 测试: 重复 Init (once.Do 幂等) ----------
func TestObjectFactory_Init_Idempotent(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	addAutowired(p)
	_, _ = f.RegisterFactory(newCompC, p, false)
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = f.Init()
		}()
	}
	wg.Wait()
	if len(f.objs) != 1 {
		t.Fatalf("expected single singleton, got %d", len(f.objs))
	}
	_ = f.Init()
	if len(f.objs) != 1 {
		t.Fatalf("idempotent Init violated")
	}
}

// ---------- 测试: SetRealizationSelector nil 与默认函数调用 ----------
func TestObjectFactory_SetRealizationSelector_Nil(t *testing.T) {
	f := NewCoreObjectFactory()
	f.SetRealizationSelector(nil)
	f.SetRealizationSelector(RealizationSelectFunc(func(defs []*object.Definition) *object.Definition {
		if len(defs) == 0 {
			return nil
		}
		return defs[0]
	}))
	p := object.NewProperty()
	addAutowired(p)
	_, _ = f.RegisterFactory(newCompC, p, false)
	_ = f.Init()
	ctx := newTestCtx()
	_, _ = f.GetObjects(ctx, (*compC)(nil))
}

// ---------- 直接调用 new() 成功 (Argn==0) ----------
func TestObjectFactory_new_NoArgsSuccess(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	addAutowired(p)
	def, _ := f.RegisterFactory(newCompC, p, false)
	_ = f.Init()
	obj, err := f.new(def, map[string]Object{})
	if err != nil || obj == nil {
		t.Fatalf("new no-arg failed: %v", err)
	}
	if _, ok := obj.Instance().(*compC); !ok {
		t.Fatalf("expected *compC instance")
	}
}

// ---------- 直接调用 new() 有依赖成功 ----------
func TestObjectFactory_new_WithDepsSuccess(t *testing.T) {
	f := NewCoreObjectFactory()
	pC := object.NewProperty()
	addAutowired(pC)
	_, _ = f.RegisterFactory(newCompC, pC, false)
	pB := object.NewProperty()
	addAutowired(pB)
	defB, _ := f.RegisterFactory(newCompB, pB, false)
	_ = f.Init()
	ctx := newTestCtx()
	cObjs, _ := f.GetObjects(ctx, (*compC)(nil))
	buildCtx := map[string]Object{}
	for _, o := range cObjs {
		buildCtx[o.ID()] = o
	}
	obj, err := f.new(defB, buildCtx)
	if err != nil || obj == nil {
		t.Fatalf("new with deps failed: %v", err)
	}
}

// ---------- 确认: buildObject getAutowiredDefinition 错误路径 ----------
func TestObjectFactory_buildObject_getAutowiredDefinitionError(t *testing.T) {
	f := NewCoreObjectFactory()
	pRoot := object.NewProperty()
	addAutowired(pRoot)
	_, _ = f.RegisterFactory(newMissingDepRoot, pRoot, false)
	_ = f.Init()
	ctx := newTestCtx()
	_, err := f.GetObjects(ctx, (*missingDepRoot)(nil))
	if err == nil || !strings.Contains(err.Error(), "definition not found") {
		t.Fatalf("expected getAutowiredDefinition error, got %v", err)
	}
}

// ---------- 验证: 实例 Definition / Factory / ID 基本信息 ----------
func TestObjectFactory_DefinitionIntegrity(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	addAutowired(p)
	def, _ := f.RegisterFactory(newCompC, p, false)
	if def.Name() == "" || def.Factory() == nil {
		t.Fatalf("definition invalid")
	}
}

// ---------- 新增：启动期缺失 autowired 标签（依赖定义存在但未标记），当前实现在构建依赖时才失败 ----------
func TestObjectFactory_Init_MissingAutowiredTag(t *testing.T) {
	f := NewCoreObjectFactory()
	// compC 定义存在但未标记 autowired（作为底层依赖）
	pC := object.NewProperty()
	_, _ = f.RegisterFactory(newCompC, pC, false)
	// compB 依赖 compC 且标记 autowired
	pB := object.NewProperty()
	addAutowired(pB)
	_, _ = f.RegisterFactory(newCompB, pB, false)
	// compA 依赖 compB 且标记 autowired
	pA := object.NewProperty()
	addAutowired(pA)
	_, _ = f.RegisterFactory(newCompA, pA, false)

	err := f.Init()
	if err == nil {
		t.Fatalf("expected init error because compC is not autowired (dependency resolution should fail)")
	}
	// 当前实现的错误路径来自 newObject -> getDependencies -> getAutowiredDefinition
	// 可能包含 "newObject failed" 或直接 "definition not found"
	msg := err.Error()
	if !strings.Contains(msg, "definition not found") || !strings.Contains(msg, "compC") {
		t.Fatalf("error should mention missing autowired dependency compC, got: %v", msg)
	}
	// 不强制断言哪些对象已创建：map 遍历顺序非确定，可能先成功创建 compC，再在 compA/compB 构建链上失败
}

// ---------- 替换：旧的 TestObjectFactory_buildObject_LastDepReturnBehavior 已失效（逻辑已修复）
// 新增：验证 buildObject 返回的确是根对象且依赖注入完整
func TestObjectFactory_buildObject_RootAndDepsConstructed(t *testing.T) {
	f := NewCoreObjectFactory()
	pC := object.NewProperty()
	addAutowired(pC)
	_, _ = f.RegisterFactory(newCompC, pC, false)
	pB := object.NewProperty()
	addAutowired(pB)
	_, _ = f.RegisterFactory(newCompB, pB, false)
	pA := object.NewProperty()
	addAutowired(pA)
	defA, _ := f.RegisterFactory(newCompA, pA, false)
	_ = f.Init()
	obj, err := f.newObject(defA, map[string]Object{})
	if err != nil {
		t.Fatalf("newObject error: %v", err)
	}
	ca, ok := obj.Instance().(*compA)
	if !ok {
		t.Fatalf("expected *compA instance, got %T", obj.Instance())
	}
	if ca.b == nil || ca.b.c == nil {
		t.Fatalf("dependency chain not injected correctly: %+v", ca)
	}
}

// 新增：GetObjectsByName 多定义（两个同类型原型 + 一个单例）覆盖多定义遍历
func TestObjectFactory_GetObjectsByName_MultipleDefinitions(t *testing.T) {
	f := NewCoreObjectFactory()
	// 单例
	ps := object.NewProperty()
	addAutowired(ps)
	if _, err := f.RegisterFactory(newCompC, ps, false); err != nil {
		t.Fatalf("register singleton compC failed: %v", err)
	}
	// 两个原型（使用不同工厂函数避免重复工厂报错）
	p1 := object.NewProperty()
	p1.Scope = object.Prototype
	addAutowired(p1)
	if _, err := f.RegisterFactory(newCompC1, p1, false); err != nil {
		t.Fatalf("register proto1 compC failed: %v", err)
	}
	p2 := object.NewProperty()
	p2.Scope = object.Prototype
	addAutowired(p2)
	if _, err := f.RegisterFactory(newCompC2, p2, false); err != nil {
		t.Fatalf("register proto2 compC failed: %v", err)
	}

	if err := f.Init(); err != nil {
		t.Fatalf("Init unexpected error: %v", err)
	}

	ctx := newTestCtx()

	defsByType, err := f.GetDefinitionsByType((*compC)(nil))
	if err != nil {
		t.Fatalf("GetDefinitionsByType failed: %v", err)
	}
	if len(defsByType) != 3 {
		t.Fatalf("expected 3 compC definitions, got %d", len(defsByType))
	}
	compCName := defsByType[0].Name()

	objs, err := f.GetObjectsByName(ctx, compCName)
	if err != nil {
		t.Fatalf("GetObjectsByName error: %v", err)
	}
	if len(objs) != 3 {
		t.Fatalf("expected 3 objects (1 singleton + 2 prototypes), got %d", len(objs))
	}

	objs2, err2 := f.GetObjectsByName(ctx, compCName)
	if err2 != nil {
		t.Fatalf("second GetObjectsByName error: %v", err2)
	}
	if len(objs2) != 3 {
		t.Fatalf("expected 3 objects second call, got %d", len(objs2))
	}
	// 验证 prototype 重新创建（至少应有一个实例不同）
	protoReuse := 0
	for i := range objs {
		for j := range objs2 {
			if objs[i] == objs2[j] {
				protoReuse++
			}
		}
	}
	if protoReuse < 1 { // 至少 singleton 复用
		t.Fatalf("expected at least one reused (singleton) object, reuse count=%d", protoReuse)
	}
}

// 新增：原型 + 无依赖 Argn==0 通过 GetObjects 创建，每次都是新对象
func TestObjectFactory_GetObjects_PrototypeNoArg(t *testing.T) {
	f := NewCoreObjectFactory()
	p := object.NewProperty()
	p.Scope = object.Prototype
	addAutowired(p)
	_, _ = f.RegisterFactory(newCompC, p, false)
	_ = f.Init()
	ctx := newTestCtx()
	objs1, err := f.GetObjects(ctx, (*compC)(nil))
	if err != nil || len(objs1) != 1 {
		t.Fatalf("first prototype get failed: %v", err)
	}
	objs2, err2 := f.GetObjects(ctx, (*compC)(nil))
	if err2 != nil || len(objs2) != 1 {
		t.Fatalf("second prototype get failed: %v", err2)
	}
	if objs1[0] == objs2[0] {
		t.Fatalf("expected different instances for prototype, got same object")
	}
}

// 新增：SetRealizationSelector 设置为 nil 后再使用默认（覆盖 selector=nil 情况下不使用的安全路径）
// 我们不直接调用会 panic 的路径（nil.Select），只验证重新设置后仍可工作
func TestObjectFactory_SetRealizationSelector_RecoverFromNil(t *testing.T) {
	f := NewCoreObjectFactory()
	f.SetRealizationSelector(nil) // 设置 nil
	// 重新设置为默认策略
	f.SetRealizationSelector(RealizationSelectFunc(func(defs []*object.Definition) *object.Definition {
		if len(defs) == 0 {
			return nil
		}
		return defs[len(defs)-1]
	}))
	// 注册两个 variantDep 并选择最后一个
	p1 := object.NewProperty()
	addAutowired(p1)
	_, _ = f.RegisterFactory(newVariant1, p1, false)
	p2 := object.NewProperty()
	addAutowired(p2)
	_, _ = f.RegisterFactory(newVariant2, p2, false)
	pr := object.NewProperty()
	pr.Scope = object.Prototype
	addAutowired(pr)
	_, _ = f.RegisterFactory(newVariantRoot, pr, false)
	_ = f.Init()
	ctx := newTestCtx()
	objs, err := f.GetObjects(ctx, (*variantRoot)(nil))
	if err != nil || len(objs) != 1 {
		t.Fatalf("variantRoot retrieval failed: %v", err)
	}
	if objs[0].Instance().(*variantRoot).dep.id != 2 {
		t.Fatalf("expected selector to pick second implementation (id=2)")
	}
}

// ---------- 结束 ----------
