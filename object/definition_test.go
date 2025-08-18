package object

import (
	"reflect"
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
		ns:          Namespace("default"),
		scope:       Scope("singleton"),
		desc:        "desc",
		lazyInit:    true,
		autoStartup: true,
		tags:        []string{"tag1", "tag2"},
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
	if def.Namespace() != Namespace("default") {
		t.Error("Namespace getter failed")
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
	tags[0] = "changed"
	if def.dependsOn[0] == "changed" {
		t.Error("DependsOn getter did not return a copy")
	}
	if def.tags[0] == "changed" {
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
	ta[0] = "y"
	if tb[0] == "y" {
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
	def.tags = []string{}
	def.dependsOn = []string{}
	// 所有必需字段都设置后，应该为 true
	if !def.IsValid() {
		t.Error("IsValid should be true if all required fields are set")
	}
}

func TestProperty_GetTagsCopy(t *testing.T) {
	prop := &Property{tags: []string{"x", "y"}}
	tags := prop.GetTags()
	tags[0] = "z"
	if prop.tags[0] == "z" {
		t.Error("Property.GetTags should return a copy")
	}
}

func TestProperty_DefaultValues(t *testing.T) {
	prop := NewProperty()
	if prop.Namespace != NSCore {
		t.Error("Property default Namespace should be Core")
	}
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

// --- AI GENERATED CODE END ---
