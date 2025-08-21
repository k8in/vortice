package object

import (
	"reflect"
	"strings"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

func newTestDefinition() *Definition {
	return &Definition{
		name:        "test",
		typ:         reflect.TypeOf(""),
		factory:     &Factory{},
		dependsOn:   []string{"a", "b"},
		methods:     &Methods{},
		scope:       Scope("singleton"),
		desc:        "desc",
		lazyInit:    true,
		autoStartup: true,
		tags:        []Tag{NewTag("tag1", "v1"), NewTag("tag2", "v2")},
	}
}

func TestDefinition_Getters(t *testing.T) {
	def := newTestDefinition()
	if def.Name() != "test" {
		t.Error("Name getter failed")
	}
	if def.Type() != reflect.TypeOf("") {
		t.Error("Type getter failed")
	}
	if def.Factory() == nil {
		t.Error("Factory getter failed")
	}
	if def.Methods() == nil {
		t.Error("Methods getter failed")
	}
	if def.Scope() != Scope("singleton") {
		t.Error("Scope getter failed")
	}
	if def.Desc() != "desc" {
		t.Error("Desc getter failed")
	}
	if !def.LazyInit() {
		t.Error("LazyInit getter failed")
	}
	if !def.AutoStartup() {
		t.Error("AutoStartup getter failed")
	}
}

func TestDefinition_SliceCopySafety(t *testing.T) {
	def := newTestDefinition()
	deps := def.DependsOn()
	tags := def.Tags()
	deps[0] = "changed"
	tags[0] = NewTag("changed", "v")
	if def.dependsOn[0] == "changed" {
		t.Error("DependsOn getter did not return a copy")
	}
	if def.tags[0].Key() == "changed" {
		t.Error("Tags getter did not return a copy")
	}
}

func TestDefinition_EmptyDependsOnTags(t *testing.T) {
	def := &Definition{}
	deps := def.DependsOn()
	tags := def.Tags()
	if deps == nil || tags == nil {
		t.Error("DependsOn/Tags should return empty slice, not nil")
	}
	if len(deps) != 0 {
		t.Error("DependsOn should return empty slice for nil")
	}
	if len(tags) != 0 {
		t.Error("Tags should return empty slice for nil")
	}
}

func TestDefinition_NilMethodsFactory(t *testing.T) {
	def := &Definition{}
	if def.Methods() != nil {
		t.Error("Methods getter should return nil if not set")
	}
	if def.Factory() != nil {
		t.Error("Factory getter should return nil if not set")
	}
}

func TestDefinition_ModifyReturnedSlice(t *testing.T) {
	def := newTestDefinition()
	a := def.DependsOn()
	b := def.DependsOn()
	a[0] = "x"
	if b[0] == "x" {
		t.Error("DependsOn getter returns shared slice")
	}
	ta := def.Tags()
	tb := def.Tags()
	ta[0] = NewTag("y", "v")
	if tb[0].Key() == "y" {
		t.Error("Tags getter returns shared slice")
	}
}

func TestDefinition_LazyInitDefault(t *testing.T) {
	def := &Definition{}
	if def.LazyInit() {
		t.Error("LazyInit default should be false")
	}
	if def.AutoStartup() {
		t.Error("AutoStartup default should be false")
	}
}

func TestDefinition_String(t *testing.T) {
	def := newTestDefinition()
	s := def.String()
	if s == "" || s == "test" {
		t.Error("String method should return a descriptive string")
	}
}

func TestDefinition_IsValid(t *testing.T) {
	def := &Definition{}
	// 所有字段为零值，应该为 false
	if def.IsValid() {
		t.Error("IsValid should be false for empty definition")
	}
	def.name = "abc"
	// 只设置 name，应该为 false
	if def.IsValid() {
		t.Error("IsValid should be false if type/factory/methods/tags/dependsOn not set")
	}
	def.typ = reflect.TypeOf(123)
	// 只设置 name 和 type，应该为 false
	if def.IsValid() {
		t.Error("IsValid should be false if factory/methods/tags/dependsOn not set")
	}
	def.factory = &Factory{}
	def.methods = &Methods{}
	def.tags = []Tag{}
	def.dependsOn = []string{}
	// 所有必需字段都设置后，应该为 true
	if !def.IsValid() {
		t.Error("IsValid should be true if all required fields are set")
	}
}

func TestProperty_GetTagsCopy(t *testing.T) {
	prop := NewProperty()
	prop.SetTags(NewTag("x", "y"), NewTag("a", "b"))
	tags := prop.GetTags()
	tags[0] = NewTag("z", "v")
	origTags := prop.GetTags()
	for _, tag := range origTags {
		if tag.Key() == "z" {
			t.Error("Property.GetTags should return a copy")
		}
	}
}

func TestProperty_DefaultValues(t *testing.T) {
	prop := NewProperty()
	if prop.Scope != Singleton {
		t.Error("Property default Scope should be Singleton")
	}
	if prop.LazyInit != true {
		t.Error("Property default LazyInit should be true")
	}
	if prop.AutoStartup != false {
		t.Error("Property default AutoStartup should be false")
	}
	if len(prop.GetTags()) != 0 {
		t.Error("Property default Tags should be empty")
	}
}

type dummyStruct struct{}
type dummyDep struct{}

func dummyFactory(a dummyDep) dummyStruct { return dummyStruct{} }
func dummyInvalidFactory(a int) int       { return a }

func TestParser_ParseValid(t *testing.T) {
	p := NewParser(dummyFactory)
	prop := NewProperty()
	def, err := p.Parse(prop)
	if err != nil {
		t.Errorf("Parse should succeed for valid factory, got error: %v", err)
	}
	if def == nil || !def.IsValid() {
		t.Error("Parse should return a valid Definition")
	}
}

func TestParser_ParseInvalidInputType(t *testing.T) {
	p := NewParser(dummyInvalidFactory)
	prop := NewProperty()
	_, err := p.Parse(prop)
	if err == nil {
		t.Error("Parse should fail for invalid input type")
	}
}

func TestParser_ParseInvalidOutputType(t *testing.T) {
	invalid := func(a dummyDep) int { return 1 }
	p := NewParser(invalid)
	prop := NewProperty()
	_, err := p.Parse(prop)
	if err == nil {
		t.Error("Parse should fail for invalid output type")
	}
}

func TestParser_ParseNotFunc(t *testing.T) {
	p := NewParser(123)
	prop := NewProperty()
	_, err := p.Parse(prop)
	if err == nil {
		t.Error("Parse should fail for non-function input")
	}
}

func TestDefinition_IsSingleton(t *testing.T) {
	def := &Definition{}
	def.scope = Singleton
	if !def.IsSingleton() {
		t.Error("IsSingleton should return true for Singleton scope")
	}
	def.scope = Prototype
	if def.IsSingleton() {
		t.Error("IsSingleton should return false for Prototype scope")
	}
}

func makeDef(scope Scope, lazy, auto bool, tags []Tag) *Definition {
	return &Definition{
		name:        "demo",
		typ:         reflect.TypeOf((*testing.T)(nil)),
		factory:     &Factory{name: "fac"},
		dependsOn:   []string{"A", "B"},
		methods:     &Methods{},
		scope:       scope,
		desc:        "desc",
		lazyInit:    lazy,
		tags:        tags,
		autoStartup: auto,
	}
}

func TestDefinition_MethodsAndCopies(t *testing.T) {
	def := makeDef(Singleton, false, false, []Tag{NewTag("k", "v")})
	if !def.IsValid() {
		t.Fatal("IsValid expected true")
	}
	if def.ID() != def.Factory().Name() {
		t.Fatal("ID mismatch")
	}
	if !def.IsSingleton() {
		t.Fatal("IsSingleton mismatch")
	}
	if def.Scope() != Singleton || def.Desc() != "desc" || def.LazyInit() {
		t.Fatal("getter mismatch")
	}
	// DependsOn 副本
	deps := def.DependsOn()
	deps[0] = "X"
	if def.DependsOn()[0] != "A" {
		t.Fatal("DependsOn must return copy")
	}
	// Tags 副本
	tags := def.Tags()
	tags[0] = NewTag("kk", "vv")
	if def.Tags()[0].Key() != "k" {
		t.Fatal("Tags must return copy")
	}
	if def.String() == "" {
		t.Fatal("String should be non-empty")
	}
}

func TestDefinition_IsValidFalseBranches(t *testing.T) {
	empty := &Definition{}
	if empty.IsValid() {
		t.Fatal("empty definition should be invalid")
	}
	// tags 为 nil
	def2 := &Definition{
		name:      "n",
		typ:       reflect.TypeOf(""),
		factory:   &Factory{name: "f"},
		dependsOn: []string{},
		methods:   &Methods{},
		tags:      nil,
	}
	if def2.IsValid() {
		t.Fatal("tags nil should make invalid")
	}
}

func TestDefinition_IsSingletonFalse(t *testing.T) {
	def := makeDef(Prototype, true, true, []Tag{})
	if def.IsSingleton() {
		t.Fatal("Prototype must not be singleton")
	}
}

func TestProperty_TagsLifecycle(t *testing.T) {
	prop := NewProperty()
	if len(prop.GetTags()) != 0 {
		t.Fatal("new property should have 0 tags")
	}
	// 空调用
	prop.SetTags()
	if len(prop.GetTags()) != 0 {
		t.Fatal("SetTags no-op failed")
	}
	t1 := NewTag("a", "1")
	t2 := NewTag("b", "2")
	prop.SetTags(t1, t2)
	if len(prop.GetTags()) != 2 {
		t.Fatal("expected 2 tags")
	}
	// 覆盖
	prop.SetTags(NewTag("a", "10"))
	got := map[string]string{}
	for _, tg := range prop.GetTags() {
		got[tg.Key()] = tg.Value()
	}
	if got["a"] != "10" || got["b"] != "2" {
		t.Fatalf("overwrite failed: %v", got)
	}
	// 副本
	cp := prop.GetTags()
	cp[0] = NewTag("x", "y")
	for _, tg := range prop.GetTags() {
		if tg.Key() == "x" {
			t.Fatal("GetTags must return copy")
		}
	}
}

func TestTag_Behavior(t *testing.T) {
	a := NewTag("k", "v")
	b := NewTag("k", "v")
	c := NewTag("k", "v2")
	if a.Key() != "k" || a.Value() != "v" {
		t.Fatal("Key/Value mismatch")
	}
	if !a.Equals(b) {
		t.Fatal("Equals should be true")
	}
	if a.Equals(c) {
		t.Fatal("Equals should be false")
	}
	if a.String() != "k=v" {
		t.Fatal("String mismatch")
	}
}

// Parser 相关
type depA struct{}
type depB struct{}
type retX struct{}
type retY struct{}

func factoryNoArg() *retX                   { return &retX{} }
func factoryWithArgs(a *depA, b depB) *retX { return &retX{} }

// 接口���数 & 结构实现
type ifaceArg interface{ M() }
type ifaceImpl struct{}

func (ifaceImpl) M() {}

func factoryWithInterfaceArg(i ifaceArg) *retX { return &retX{} }

// 返回指向非 struct 的指针（当前实现应失败：pointer 但 Elem 不是 struct）
func factoryReturnPtrInt() *int {
	x := 1
	return &x
}

// 无返回值函数（应报错）
func factoryNoReturn(a depA) { _ = a }

func TestParser_Parse_Success_NoArgs(t *testing.T) {
	p := NewParser(factoryNoArg)
	prop := NewProperty()
	def, err := p.Parse(prop)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(def.DependsOn()) != 0 {
		t.Fatal("expected 0 deps")
	}
	if def.Name() == "" {
		t.Fatal("name empty")
	}
	if !def.IsValid() {
		t.Fatal("definition invalid unexpectedly")
	}
}

func TestParser_Parse_Success_WithArgs(t *testing.T) {
	p := NewParser(factoryWithArgs)
	prop := NewProperty()
	def, err := p.Parse(prop)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(def.DependsOn()) != 2 {
		t.Fatalf("expected 2 deps, got %v", def.DependsOn())
	}
	// 依赖名格式（*Struct）
	for _, d := range def.DependsOn() {
		if d == "" {
			t.Fatal("dependency name empty")
		}
	}
}

func TestParser_Parse_Success_InterfaceArg(t *testing.T) {
	p := NewParser(factoryWithInterfaceArg)
	prop := NewProperty()
	def, err := p.Parse(prop)
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if len(def.DependsOn()) != 1 || !strings.Contains(def.DependsOn()[0], "ifaceArg") {
		t.Fatalf("interface arg dependency name unexpected: %v", def.DependsOn())
	}
}

func TestParser_Parse_Error_PointerNonStructReturn(t *testing.T) {
	p := NewParser(factoryReturnPtrInt)
	prop := NewProperty()
	if _, err := p.Parse(prop); err == nil {
		t.Fatalf("expected error for pointer-to-non-struct return type (*int)")
	}
}

func TestParser_Parse_Error_NoReturn(t *testing.T) {
	p := NewParser(factoryNoReturn)
	prop := NewProperty()
	if _, err := p.Parse(prop); err == nil {
		t.Fatal("expected error for function with zero return values")
	}
}

func TestGenerateDefinitionName_Wrapper(t *testing.T) {
	typPtr := reflect.TypeOf((*depA)(nil))
	n1 := GenerateDefinitionName(typPtr)
	n2 := generateReflectionName(typPtr)
	if n1 != n2 || !strings.Contains(n1, "*depA") {
		t.Fatalf("GenerateDefinitionName mismatch: %s vs %s", n1, n2)
	}
}

func TestParseDefinition_Wrapper(t *testing.T) {
	def, err := ParseDefinition(factoryNoArg, NewProperty())
	if err != nil {
		t.Fatalf("ParseDefinition failed: %v", err)
	}
	if !def.IsValid() {
		t.Fatalf("definition should be valid")
	}
}

func TestDefinition_StringContainsFactory(t *testing.T) {
	def := &Definition{
		name:      "compX",
		typ:       reflect.TypeOf((*retX)(nil)),
		factory:   &Factory{name: "facX"},
		dependsOn: []string{},
		methods:   &Methods{},
		scope:     Singleton,
		tags:      []Tag{},
	}
	s := def.String()
	if !strings.Contains(s, "facX") || !strings.Contains(s, "compX") {
		t.Fatalf("String should contain name & factory: %s", s)
	}
}

// 新增辅助工厂
func factoryWithPtrIntArg(*int) *retX           { return &retX{} }
func factoryMultiReturn() (*retX, error)        { return &retX{}, nil }
func factoryWithTwoArgs(a *depA, b *depB) *retX { return &retX{} }

// 新增：指针指向非 struct 的参数 (*int) => checkArgType 报错
func TestParser_Parse_Error_PtrToNonStructArg(t *testing.T) {
	p := NewParser(factoryWithPtrIntArg)
	prop := NewProperty()
	if _, err := p.Parse(prop); err == nil {
		t.Fatalf("expected error for pointer-to-non-struct argument (*int)")
	}
}

// 新增：多返回值 => checkOutputAndSet 报错
func TestParser_Parse_Error_MultiReturn(t *testing.T) {
	p := NewParser(factoryMultiReturn)
	prop := NewProperty()
	if _, err := p.Parse(prop); err == nil {
		t.Fatalf("expected error for multiple return values")
	}
}

// 新增：nil 函数输入 => 无效函数
func TestParser_Parse_Error_NilFunc(t *testing.T) {
	var fn any = nil
	p := NewParser(fn)
	if _, err := p.Parse(NewProperty()); err == nil {
		t.Fatalf("expected error for nil function input")
	}
}

// 新增：命名规则（指针接口 vs 值结构体）
func TestGenerateReflectionName_InterfacePointerAndValueStruct(t *testing.T) {
	it := reflect.TypeOf((*ifaceArg)(nil)) // pointer to interface
	nameIfacePtr := generateReflectionName(it)
	if !strings.Contains(nameIfacePtr, ".ifaceArg") || strings.Contains(nameIfacePtr, "*ifaceArg") {
		t.Fatalf("unexpected interface pointer name: %s", nameIfacePtr)
	}
	valStruct := reflect.TypeOf(depB{})
	nameValStruct := generateReflectionName(valStruct)
	if !strings.Contains(nameValStruct, "*depB") {
		t.Fatalf("value struct should have leading *: %s", nameValStruct)
	}
	if nameIfacePtr == nameValStruct {
		t.Fatalf("interface and struct names should differ: %s", nameIfacePtr)
	}
}

// 新增：多参数工厂解析校验依赖顺序 / Argn / 传递的标签
func TestParseDefinition_ArgCountAndTags(t *testing.T) {
	prop := NewProperty()
	prop.SetTags(NewTag("k1", "v1"), NewTag("k2", "v2"))
	def, err := ParseDefinition(factoryWithTwoArgs, prop)
	if err != nil {
		t.Fatalf("ParseDefinition failed: %v", err)
	}
	if !def.IsValid() {
		t.Fatalf("definition invalid")
	}
	if def.Factory().Argn() != 2 {
		t.Fatalf("expected Argn=2, got %d", def.Factory().Argn())
	}
	deps := def.DependsOn()
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %v", deps)
	}
	// 顺序应保持与参数声明一致
	if !strings.Contains(deps[0], "*depA") || !strings.Contains(deps[1], "*depB") {
		t.Fatalf("dependencies order/name mismatch: %v", deps)
	}
	// 标签透传
	tagKeys := map[string]bool{}
	for _, tg := range def.Tags() {
		tagKeys[tg.Key()] = true
	}
	if !tagKeys["k1"] || !tagKeys["k2"] {
		t.Fatalf("definition tags missing: %v", def.Tags())
	}
}

// --- AI GENERATED CODE END ---

// 返回空接口的工厂；若 newMethods 对 interface 返回 nil，则触发 IsValid=false ���支
type emptyIface interface{}

func factoryReturnEmptyIface() emptyIface { return struct{}{} }

func TestParser_Parse_Error_InterfaceMissingMethods(t *testing.T) {
	p := NewParser(factoryReturnEmptyIface)
	prop := NewProperty()
	def, err := p.Parse(prop)
	if err == nil {
		// 如果实现改变（interface 也生成 methods），确保 def 有效以避免测试误判
		if def == nil || !def.IsValid() {
			t.Fatalf("expected valid definition when methods generated")
		}
		return
	}
	if err != ErrMissingRequiredField {
		t.Fatalf("expected ErrMissingRequiredField, got %v", err)
	}
}

// 覆盖 SetTags 分支：传入 nil slice 而非零参数调用
func TestProperty_SetTagsNilSlice(t *testing.T) {
	prop := NewProperty()
	var nilSlice []Tag
	prop.SetTags(nilSlice...) // 应走早返回
	if len(prop.GetTags()) != 0 {
		t.Fatalf("expected no tags after SetTags(nilSlice)")
	}
}

// 覆盖 Parse -> newDefinition 对 Property 各字段赋值与解析后 Definition 的 ID/String 访问
func factoryWithOverrides() *retX { return &retX{} }

func TestParser_Parse_WithPropertyOverrides(t *testing.T) {
	prop := NewProperty()
	prop.Scope = Prototype
	prop.LazyInit = false
	prop.AutoStartup = true
	prop.SetTags(NewTag("ov", "1"))

	def, err := ParseDefinition(factoryWithOverrides, prop)
	if err != nil {
		t.Fatalf("ParseDefinition overrides failed: %v", err)
	}
	if def.Scope() != Prototype {
		t.Fatalf("expected scope Prototype got %s", def.Scope())
	}
	if def.LazyInit() {
		t.Fatalf("expected LazyInit false")
	}
	if !def.AutoStartup() {
		t.Fatalf("expected AutoStartup true")
	}
	// 确认标签透传
	found := false
	for _, tg := range def.Tags() {
		if tg.Key() == "ov" && tg.Value() == "1" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("override tag not found in definition tags: %v", def.Tags())
	}
	// 访问 ID / String 以确保相关行被统计（解析型 definition）
	if def.ID() == "" {
		t.Fatalf("expected non-empty ID for parsed definition")
	}
	if s := def.String(); s == "" || !strings.Contains(s, def.ID()) || !strings.Contains(s, def.Name()) {
		t.Fatalf("definition String should contain name & id, got %s", s)
	}
	// Desc 未由 Property 复制，这里应为空字符串
	if def.Desc() != "" {
		t.Fatalf("expected empty desc (property.Desc not propagated), got %q", def.Desc())
	}
}
