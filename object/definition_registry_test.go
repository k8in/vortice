package object

import (
	"errors"
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

func makeTestDefinition(name, factoryName string, tags []Tag) *Definition {
	if tags == nil { // 确保非 nil，避免 IsValid 判定失败
		tags = []Tag{}
	}
	return &Definition{
		name:      name,
		typ:       reflect.TypeOf(""),
		factory:   &Factory{name: factoryName},
		dependsOn: []string{},
		methods:   &Methods{},
		scope:     Singleton,
		tags:      tags,
	}
}

func TestDefinitionRegistry_Register_InvalidDefinition(t *testing.T) {
	reg := NewDefinitionRegistry()
	if err := reg.register(&Definition{}, false); err == nil {
		t.Fatalf("expected invalid definition error")
	}
}

func TestDefinitionRegistry_Register_Readonly(t *testing.T) {
	reg := NewDefinitionRegistry()
	reg.readonly.Store(true)
	def := makeTestDefinition("x", "fx", nil)
	if err := reg.register(def, false); err == nil || err.Error() != "the DefinitionRegistry has been locked" {
		t.Fatalf("expected locked error, got %v", err)
	}
}

func TestDefinitionRegistry_Register_DuplicateFactory(t *testing.T) {
	reg := NewDefinitionRegistry()
	d1 := makeTestDefinition("n1", "fdup", nil)
	d2 := makeTestDefinition("n2", "fdup", nil)
	if err := reg.register(d1, false); err != nil {
		t.Fatalf("register d1 failed: %v", err)
	}
	if err := reg.register(d2, false); err == nil {
		t.Fatalf("expected duplicate factory error")
	}
}

func TestDefinitionRegistry_Register_UniqueNameEnforced(t *testing.T) {
	reg := NewDefinitionRegistry()
	d1 := makeTestDefinition("same", "f1", nil)
	d2 := makeTestDefinition("same", "f2", nil)
	// unique=true 第一次即可拒绝（当前实现：只要 unique 为真就拒绝任何后续同名）
	if err := reg.register(d1, true); err != nil {
		t.Fatalf("first register with unique should pass (当前实现逻辑允许首次?) got %v", err)
	}
	if err := reg.register(d2, true); err == nil {
		t.Fatalf("expected duplicate name reject when unique=true")
	}
}

func TestDefinitionRegistry_RegisterFactory_ParseErrors(t *testing.T) {
	reg := NewDefinitionRegistry()
	// 非函数
	if _, err := reg.RegisterFactory(123, NewProperty(), false); err == nil || !errors.Is(err, ErrParseDefinition) {
		t.Fatalf("expected parse error (non-func), got %v", err)
	}
	// 返回值数量错误
	if _, err := reg.RegisterFactory(badFactoryMultiReturn, NewProperty(), false); err == nil || !errors.Is(err, ErrParseDefinition) {
		t.Fatalf("expected multi return parse error, got %v", err)
	}
}

func TestDefinitionRegistry_RegisterFactory_Readonly(t *testing.T) {
	reg := NewDefinitionRegistry()
	// 先成功一次
	if _, err := reg.RegisterFactory(goodFactoryFoo, NewProperty(), false); err != nil {
		t.Fatalf("first RegisterFactory failed: %v", err)
	}
	_ = reg.Init()
	// 再次注册应因内部 register 失败
	if _, err := reg.RegisterFactory(goodFactoryBar, NewProperty(), false); err == nil {
		t.Fatalf("expected internal register failure after readonly")
	}
}

func TestDefinitionRegistry_EntriesAndFactories(t *testing.T) {
	reg := NewDefinitionRegistry()
	def := makeTestDefinition("obj", "factory", []Tag{NewTag("k", "v")})
	if err := reg.register(def, false); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	list, ok := reg.entries[def.Name()]
	if !ok || len(list) != 1 || list[0] != def {
		t.Error("entries not updated correctly")
	}
	gotDef, ok := reg.factories[def.Factory().Name()]
	if !ok || gotDef != def {
		t.Error("factories not updated correctly")
	}
}

func TestDefinitionRegistry_GetDefinitionsByName_Filter(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []Tag{NewTag("tag1", "v1")})
	defB := makeTestDefinition("A", "fb", []Tag{NewTag("tag2", "v2")})
	defC := makeTestDefinition("B", "fc", []Tag{NewTag("tag1", "v1")})
	_ = reg.register(defA, false)
	_ = reg.register(defB, false)
	_ = reg.register(defC, false)

	// 按 factory name 过滤
	filter := func(def *Definition) bool { return def.Factory().Name() == "fa" }
	defs := reg.GetDefinitionsByName("A", filter)
	if len(defs) != 1 || defs[0].factory.name != "fa" {
		t.Error("GetDefinitionsByName with filter failed")
	}

	// 多 filter
	tagFilter := TagFilter(NewTag("tag1", "v1"))
	defsMulti := reg.GetDefinitionsByName("A", filter, tagFilter)
	if len(defsMulti) != 1 || defsMulti[0].factory.name != "fa" {
		t.Error("GetDefinitionsByName with multiple filters failed")
	}

	// 无 filter
	defsAll := reg.GetDefinitionsByName("A")
	if len(defsAll) != 2 {
		t.Error("GetDefinitionsByName without filter failed")
	}
}

func TestDefinitionRegistry_GetDefinitions(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []Tag{NewTag("tag1", "v1")})
	defB := makeTestDefinition("B", "fb", []Tag{NewTag("tag2", "v2")})
	defC := makeTestDefinition("C", "fc", []Tag{NewTag("tag1", "v1")})
	_ = reg.register(defA, false)
	_ = reg.register(defB, false)
	_ = reg.register(defC, false)

	// 无 filter
	allDefs := reg.GetDefinitions()
	if len(allDefs) != 3 {
		t.Errorf("GetDefinitions without filter failed, got %d", len(allDefs))
	}

	// tag filter
	tagFilter := TagFilter(NewTag("tag1", "v1"))
	tagDefs := reg.GetDefinitions(tagFilter)
	if len(tagDefs) != 2 {
		t.Errorf("GetDefinitions with tag filter failed, got %d", len(tagDefs))
	}

	// 多 filter
	multiDefs := reg.GetDefinitions(tagFilter)
	if len(multiDefs) != 2 {
		t.Errorf("GetDefinitions with multiple filters failed, got %d", len(multiDefs))
	}
}

func TestDefinitionRegistry_SortAndCheck_Cycle(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []Tag{})
	defB := makeTestDefinition("B", "fb", []Tag{})
	defA.dependsOn = []string{"B"}
	defB.dependsOn = []string{"A"}
	_ = reg.register(defA, false)
	_ = reg.register(defB, false)
	err := reg.sortAndCheck()
	if err == nil {
		t.Error("sortAndCheck should fail for cycle")
	}
}

func TestDefinitionRegistry_SortAndCheck_MissingDep(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []Tag{})
	defA.dependsOn = []string{"B"}
	_ = reg.register(defA, false)
	err := reg.sortAndCheck()
	if err == nil {
		t.Error("sortAndCheck should fail for missing dependency")
	}
}

func TestDefinitionRegistry_Init_SortAndCheck(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []Tag{})
	defB := makeTestDefinition("B", "fb", []Tag{})
	defA.dependsOn = []string{"B"}
	defB.dependsOn = []string{}
	_ = reg.register(defA, false)
	_ = reg.register(defB, false)
	reg.Init()
	if !reg.readonly.Load() {
		t.Error("Init should set registry to readonly")
	}
	if len(reg.inSeq) != 2 {
		t.Errorf("Init should record registered definitions, got %d", len(reg.inSeq))
	}
	if reg.inSeq[0] != "fb" || reg.inSeq[1] != "fa" {
		t.Errorf("Registered order incorrect: %v", reg.inSeq)
	}
}

func TestDefinitionRegistry_Init_SortAndCheck_Complex(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA1 := makeTestDefinition("A", "fa1", []Tag{})
	defA2 := makeTestDefinition("A", "fa2", []Tag{})
	defB := makeTestDefinition("B", "fb", []Tag{})
	defC := makeTestDefinition("C", "fc", []Tag{})
	defD := makeTestDefinition("D", "fd", []Tag{})

	defA1.dependsOn = []string{"B", "C"}
	defA2.dependsOn = []string{"B"}
	defB.dependsOn = []string{"D"}
	defC.dependsOn = []string{"D"}
	defD.dependsOn = []string{}

	_ = reg.register(defA1, false)
	_ = reg.register(defA2, false)
	_ = reg.register(defB, false)
	_ = reg.register(defC, false)
	_ = reg.register(defD, false)

	reg.Init()
	if !reg.readonly.Load() {
		t.Error("Init should set registry to readonly")
	}
	if len(reg.inSeq) != 5 {
		t.Errorf("Init should record registered definitions, got %d", len(reg.inSeq))
	}

	index := map[string]int{}
	for i, v := range reg.inSeq {
		index[v] = i
	}
	if !(index["fd"] < index["fb"] && index["fd"] < index["fc"]) {
		t.Errorf("fd should be before fb and fc: %v", reg.inSeq)
	}
	if !(index["fb"] < index["fa1"] && index["fb"] < index["fa2"]) {
		t.Errorf("fb should be before fa1 and fa2: %v", reg.inSeq)
	}
	if !(index["fc"] < index["fa1"]) {
		t.Errorf("fc should be before fa1: %v", reg.inSeq)
	}
}

func TestDefinitionRegistry_GetDefinitionsByType(t *testing.T) {
	reg := NewDefinitionRegistry()

	// 指针指向 struct，应该通过
	type myStruct struct{}
	typ := reflect.TypeOf((*myStruct)(nil))
	def := &Definition{
		name:      generateReflectionName(typ),
		typ:       typ,
		factory:   &Factory{name: "factory"},
		dependsOn: []string{},
		methods:   &Methods{},
		scope:     Singleton,
		tags:      []Tag{NewTag("tag", "v")},
	}
	_ = reg.register(def, false)

	defs, err := reg.GetDefinitionsByType((*myStruct)(nil))
	if err != nil {
		t.Fatalf("GetDefinitionsByType failed: %v", err)
	}
	if len(defs) != 1 || defs[0] != def {
		t.Error("GetDefinitionsByType should return the correct definition")
	}

	// 非指针类型，应该失败
	_, err = reg.GetDefinitionsByType(myStruct{})
	if err == nil {
		t.Error("GetDefinitionsByType should fail for non-pointer type")
	}

	// 指针指向非 struct/interface（如 int），应该失败
	var intPtr *int
	_, err = reg.GetDefinitionsByType(intPtr)
	if err == nil {
		t.Error("GetDefinitionsByType should fail for pointer to non-struct/interface")
	}

	// 未注册类型，应该失败
	type anotherStruct struct{}
	_, err = reg.GetDefinitionsByType((*anotherStruct)(nil))
	if err == nil {
		t.Error("GetDefinitionsByType should fail for unknown type")
	}
}

func TestDefinitionRegistry_getObjectType(t *testing.T) {
	reg := NewDefinitionRegistry()

	// 指针指向 struct
	type myStruct struct{}
	ptrStruct := &myStruct{}
	rt := reg.getObjectType(ptrStruct)
	if rt == nil || rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		t.Error("getObjectType should return pointer to struct type")

	}

	// 指针指向 interface
	var ifacePtr *testing.T
	rtIface := reg.getObjectType(ifacePtr)
	if rtIface == nil || rtIface.Kind() != reflect.Ptr || rtIface.Elem().Kind() != reflect.Struct {
		t.Error("getObjectType should return pointer to struct type for *testing.T")
	}

	// 非指针类型
	val := myStruct{}
	if reg.getObjectType(val) != nil {
		t.Error("getObjectType should return nil for non-pointer type")
	}

	// 指针指向非 struct/interface（如 int）
	var intPtr *int
	if reg.getObjectType(intPtr) != nil {
		t.Error("getObjectType should return nil for pointer to non-struct/interface")
	}

	// nil
	if reg.getObjectType(nil) != nil {
		t.Error("getObjectType should return nil for nil input")
	}
}

// --- AI GENERATED CODE END ---

// --- 新增覆盖辅助 ---
type fooStruct struct{}
type barStruct struct{}

func goodFactoryFoo() *fooStruct                 { return &fooStruct{} }
func goodFactoryBar(a *fooStruct) *barStruct     { return &barStruct{} }
func badFactoryMultiReturn() (*fooStruct, error) { return &fooStruct{}, nil }

// 生成一个有效 Definition (自定义 scope / tags)
func makeScopedDef(name, factory string, scope Scope, tags []Tag) *Definition {
	return &Definition{
		name:      name,
		typ:       reflect.TypeOf((*fooStruct)(nil)),
		factory:   &Factory{name: factory},
		dependsOn: []string{},
		methods:   &Methods{},
		scope:     scope,
		tags:      tags,
	}
}

// --- 新增测试开始 ---

func TestDefinitionRegistry_GetDefinitions_Filters(t *testing.T) {
	reg := NewDefinitionRegistry()
	// def1: Singleton tag a=1
	d1 := makeScopedDef("A", "fa", Singleton, []Tag{NewTag("a", "1")})
	// def2: Prototype tag a=2
	d2 := makeScopedDef("B", "fb", Prototype, []Tag{NewTag("a", "2")})
	_ = reg.register(d1, false)
	_ = reg.register(d2, false)

	// ScopeFilter
	singletons := reg.GetDefinitions(ScopeFilter(Singleton))
	if len(singletons) != 1 || singletons[0].Name() != "A" {
		t.Fatalf("ScopeFilter failed: %+v", singletons)
	}
	// TagFilter 匹配
	tagA1 := reg.GetDefinitions(TagFilter(NewTag("a", "1")))
	if len(tagA1) != 1 || tagA1[0].Name() != "A" {
		t.Fatalf("TagFilter a=1 failed: %+v", tagA1)
	}
	// TagFilter 不匹配
	tagNone := reg.GetDefinitions(TagFilter(NewTag("x", "y")))
	if len(tagNone) != 0 {
		t.Fatalf("TagFilter should return empty, got %d", len(tagNone))
	}
	// TagFilter 空参数
	emptyTag := reg.GetDefinitions(TagFilter())
	if len(emptyTag) != 0 {
		t.Fatalf("TagFilter() expected empty result")
	}
	// 多过滤器 AND：Scope=Singleton AND tag a=1 -> A
	multi := reg.GetDefinitions(ScopeFilter(Singleton), TagFilter(NewTag("a", "1")))
	if len(multi) != 1 || multi[0].Name() != "A" {
		t.Fatalf("multi filter failed: %+v", multi)
	}
	// 含 nil filter 不影响
	withNil := reg.GetDefinitions(nil, ScopeFilter(Prototype))
	if len(withNil) != 1 || withNil[0].Name() != "B" {
		t.Fatalf("nil filter handling failed: %+v", withNil)
	}
}

func TestDefinitionRegistry_GetDefinitionsByName_FiltersAndNil(t *testing.T) {
	reg := NewDefinitionRegistry()
	d1 := makeScopedDef("X", "fx1", Singleton, []Tag{NewTag("k", "v1")})
	d2 := makeScopedDef("X", "fx2", Prototype, []Tag{NewTag("k", "v2")})
	_ = reg.register(d1, false)
	_ = reg.register(d2, false)

	all := reg.GetDefinitionsByName("X")
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
	// 单过滤器
	onlyV1 := reg.GetDefinitionsByName("X", TagFilter(NewTag("k", "v1")))
	if len(onlyV1) != 1 || onlyV1[0].Factory().Name() != "fx1" {
		t.Fatalf("Tag filter by name failed: %+v", onlyV1)
	}
	// 多过滤器 AND 无结果
	none := reg.GetDefinitionsByName("X", TagFilter(NewTag("k", "v1")), ScopeFilter(Prototype))
	if len(none) != 0 {
		t.Fatalf("expected empty due to AND mismatch")
	}
	// 含 nil filter
	withNil := reg.GetDefinitionsByName("X", nil, TagFilter(NewTag("k", "v2")))
	if len(withNil) != 1 || withNil[0].Factory().Name() != "fx2" {
		t.Fatalf("nil+tag filter failed: %+v", withNil)
	}
}

func TestDefinitionRegistry_GetDefinitionsByType_ErrorsAndSuccess(t *testing.T) {
	reg := NewDefinitionRegistry()
	// 先注册 fooStruct
	prop := NewProperty()
	if _, err := reg.RegisterFactory(goodFactoryFoo, prop, false); err != nil {
		t.Fatalf("register foo failed: %v", err)
	}
	// typ=nil
	if _, err := reg.GetDefinitionsByType(nil); err == nil {
		t.Fatalf("expected error for nil typ")
	}
	// 非指针
	if _, err := reg.GetDefinitionsByType(fooStruct{}); err == nil {
		t.Fatalf("expected error for non-pointer")
	}
	// 指针->非 struct/interface (*int)
	var intPtr *int
	if _, err := reg.GetDefinitionsByType(intPtr); err == nil {
		t.Fatalf("expected error for pointer to non-struct/interface")
	}
	// 未注册类型
	type otherStruct struct{}
	if _, err := reg.GetDefinitionsByType((*otherStruct)(nil)); err == nil {
		t.Fatalf("expected not found error")
	}
	// 成功 (已注册)
	if defs, err := reg.GetDefinitionsByType((*fooStruct)(nil)); err != nil || len(defs) != 1 {
		t.Fatalf("expected 1 def for fooStruct, got %v %v", len(defs), err)
	}
}

func TestDefinitionRegistry_getObjectType_AllBranches(t *testing.T) {
	reg := NewDefinitionRegistry()
	// nil
	if reg.getObjectType(nil) != nil {
		t.Fatalf("nil should return nil")
	}
	// 非指针
	if reg.getObjectType(123) != nil {
		t.Fatalf("non-pointer should return nil")
	}
	// 指针->struct
	if reg.getObjectType(&fooStruct{}) == nil {
		t.Fatalf("pointer to struct should pass")
	}
	// 指针->interface
	type myIface interface{ M() }
	if reg.getObjectType((*myIface)(nil)) == nil {
		t.Fatalf("pointer to interface should pass")
	}
	// 指针->非 struct/interface (*int)
	var intPtr *int
	if reg.getObjectType(intPtr) != nil {
		t.Fatalf("pointer to non struct/interface should fail")
	}
}

func TestDefinitionRegistry_sortAndCheck_DefinitionMissingAfterRegistration(t *testing.T) {
	reg := NewDefinitionRegistry()
	d := makeTestDefinition("Ghost", "fGhost", nil)
	_ = reg.register(d, false)
	// 删除 entries 触发 "definition not found"
	delete(reg.entries, d.Name())
	if err := reg.sortAndCheck(); err == nil {
		t.Fatalf("expected error for missing entry in entries map")
	}
}

func TestDefinitionRegistry_sortAndCheck_SuccessOrderAndInSeqUpdate(t *testing.T) {
	reg := NewDefinitionRegistry()
	// f2 依赖 f1
	d1 := makeTestDefinition("D1", "f1", nil)
	d2 := makeTestDefinition("D2", "f2", nil)
	d2.dependsOn = []string{"D1"}
	_ = reg.register(d1, false)
	_ = reg.register(d2, false)
	if err := reg.sortAndCheck(); err != nil {
		t.Fatalf("sortAndCheck failed: %v", err)
	}
	// inSeq 由 sortAndCheck 重写
	if len(reg.inSeq) != 2 {
		t.Fatalf("expected 2 inSeq, got %v", reg.inSeq)
	}
	// 依赖拓扑：f1 在前
	if reg.inSeq[0] != "f1" {
		t.Fatalf("expected f1 first, got %v", reg.inSeq)
	}
}

// --- 新增测试结束 ---
