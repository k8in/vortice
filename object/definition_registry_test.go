package object

import (
	"reflect"
	"testing"
)

func makeTestDefinition(name, factoryName string, tags []string) *Definition {
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

func TestDefinitionRegistry_Register_DuplicateFactory(t *testing.T) {
	reg := NewDefinitionRegistry()
	def1 := makeTestDefinition("obj1", "factory1", []string{"tag"})
	def2 := makeTestDefinition("obj2", "factory1", []string{"tag"})
	if err := reg.register(def1, false); err != nil {
		t.Fatalf("register def1 failed: %v", err)
	}
	if err := reg.register(def2, false); err == nil {
		t.Error("register should fail for duplicate factory")
	}
}

func TestDefinitionRegistry_Register_Unique(t *testing.T) {
	reg := NewDefinitionRegistry()
	def1 := makeTestDefinition("obj1", "factory1", []string{"tag"})
	def2 := makeTestDefinition("obj1", "factory2", []string{"tag"})
	if err := reg.register(def1, true); err != nil {
		t.Fatalf("register def1 failed: %v", err)
	}
	// unique=true，name重复，应该报错
	if err := reg.register(def2, true); err == nil {
		t.Error("register should fail for duplicate name when unique is true")
	}
	// unique=false，name重复，允许
	if err := reg.register(def2, false); err != nil {
		t.Errorf("register should allow duplicate name when unique is false, got %v", err)
	}
}

func TestDefinitionRegistry_RegisterAndLock(t *testing.T) {
	reg := NewDefinitionRegistry()
	def := makeTestDefinition("obj", "factory", []string{"tag"})
	if err := reg.register(def, false); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	reg.Init()
	if !reg.readonly.Load() {
		t.Error("registry should be readonly after Init")
	}
	def2 := makeTestDefinition("obj2", "factory2", []string{"tag"})
	if err := reg.register(def2, false); err == nil {
		t.Error("register should fail after Init")
	}
}

func TestDefinitionRegistry_EntriesAndFactories(t *testing.T) {
	reg := NewDefinitionRegistry()
	def := makeTestDefinition("obj", "factory", []string{"tag"})
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
	defA := makeTestDefinition("A", "fa", []string{"tag1"})
	defB := makeTestDefinition("A", "fb", []string{"tag2"})
	defC := makeTestDefinition("B", "fc", []string{"tag1"})
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
	tagFilter := TagFilter("tag1")
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
	defA := makeTestDefinition("A", "fa", []string{"tag1"})
	defB := makeTestDefinition("B", "fb", []string{"tag2"})
	defC := makeTestDefinition("C", "fc", []string{"tag1"})
	_ = reg.register(defA, false)
	_ = reg.register(defB, false)
	_ = reg.register(defC, false)

	// 无 filter
	allDefs := reg.GetDefinitions()
	if len(allDefs) != 3 {
		t.Errorf("GetDefinitions without filter failed, got %d", len(allDefs))
	}

	// tag filter
	tagFilter := TagFilter("tag1")
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
	defA := makeTestDefinition("A", "fa", []string{})
	defB := makeTestDefinition("B", "fb", []string{})
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
	defA := makeTestDefinition("A", "fa", []string{})
	defA.dependsOn = []string{"B"}
	_ = reg.register(defA, false)
	err := reg.sortAndCheck()
	if err == nil {
		t.Error("sortAndCheck should fail for missing dependency")
	}
}

func TestDefinitionRegistry_Init_SortAndCheck(t *testing.T) {
	reg := NewDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", []string{})
	defB := makeTestDefinition("B", "fb", []string{})
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
	defA1 := makeTestDefinition("A", "fa1", []string{})
	defA2 := makeTestDefinition("A", "fa2", []string{})
	defB := makeTestDefinition("B", "fb", []string{})
	defC := makeTestDefinition("C", "fc", []string{})
	defD := makeTestDefinition("D", "fd", []string{})

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
