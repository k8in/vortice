package object

import (
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

func makeTestDefinition(name string, factoryName string, ns Namespace, tags []string) *Definition {
	return &Definition{
		name:      name,
		typ:       reflect.TypeOf(""),
		factory:   &Factory{name: factoryName},
		dependsOn: []string{},
		methods:   &Methods{},
		ns:        ns,
		scope:     Singleton,
		tags:      tags,
	}
}

func TestDefinitionRegistry_RegisterAndLock(t *testing.T) {
	reg := newDefinitionRegistry()
	def1 := makeTestDefinition("obj1", "obj1_factory", NSCore, []string{"tag"})
	def2 := makeTestDefinition("obj2", "obj2_factory", NSCore, []string{"tag"})

	// Register should succeed
	if err := reg.register(def1); err != nil {
		t.Fatalf("register def1 failed: %v", err)
	}
	if err := reg.register(def2); err != nil {
		t.Fatalf("register def2 failed: %v", err)
	}

	// Duplicate factory should fail
	dupDef := makeTestDefinition("obj3", "obj1_factory", NSCore, []string{"tag"})
	if err := reg.register(dupDef); err == nil {
		t.Error("register duplicate factory should fail")
	}

	// Lock registry
	reg.Lock()
	if !reg.readonly.Load() {
		t.Error("registry should be readonly after Lock")
	}

	// Register after lock should fail
	def4 := makeTestDefinition("obj4", "obj4_factory", NSCore, []string{"tag"})
	if err := reg.register(def4); err == nil {
		t.Error("register after lock should fail")
	}
}

func TestDefinitionRegistry_EntriesAndFactories(t *testing.T) {
	reg := newDefinitionRegistry()
	def := makeTestDefinition("obj", "obj_factory", NSCore, []string{"tag"})
	if err := reg.register(def); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	// Check entries
	list, ok := reg.entries[def.Name()]
	if !ok || len(list) != 1 || list[0] != def {
		t.Error("entries not updated correctly")
	}
	// Check factories
	fid := def.Factory().Name()
	name, ok := reg.factories[fid]
	if !ok || name != def.Name() {
		t.Error("factories not updated correctly")
	}
}

func TestRegisterFactory_Success(t *testing.T) {
	type dummyDep struct{}
	type dummyStruct struct{}
	factory := func(a dummyDep) dummyStruct { return dummyStruct{} }
	prop := NewProperty()
	def, err := RegisterFactory(factory, prop)
	if err != nil {
		t.Fatalf("RegisterFactory should succeed, got error: %v", err)
	}
	if def == nil || def.Name() == "" {
		t.Error("RegisterFactory should return valid definition")
	}
}

func TestRegisterFactory_Duplicate(t *testing.T) {
	type dummyDep struct{}
	type dummyStruct struct{}
	factory := func(a dummyDep) dummyStruct { return dummyStruct{} }
	prop := NewProperty()
	_, err := RegisterFactory(factory, prop)
	if err != nil {
		t.Fatalf("First RegisterFactory should succeed, got error: %v", err)
	}
	_, err = RegisterFactory(factory, prop)
	if err == nil {
		t.Error("RegisterFactory should fail for duplicate factory")
	}
}

func TestDefinitionRegistry_GetDefinition_Filter(t *testing.T) {
	reg := newDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", NSCore, []string{"tag1"})
	defB := makeTestDefinition("A", "fb", NSCore, []string{"tag2"})
	defC := makeTestDefinition("B", "fc", NSCore, []string{"tag1"})
	_ = reg.register(defA)
	_ = reg.register(defB)
	_ = reg.register(defC)

	// Filter: only factory name == "fa"
	filter := func(def *Definition) bool { return def.Factory().Name() == "fa" }
	defs := reg.GetDefinition("A", filter)
	if len(defs) != 1 || defs[0].factory.name != "fa" {
		t.Error("GetDefinition with filter failed")
	}

	// Multiple filters: factory name == "fa" and tag == "tag1"
	tagFilter := TagFilter("tag1")
	defsMulti := reg.GetDefinition("A", filter, tagFilter)
	if len(defsMulti) != 1 || defsMulti[0].factory.name != "fa" {
		t.Error("GetDefinition with multiple filters failed")
	}

	// No filter: should return all
	defsAll := reg.GetDefinition("A")
	if len(defsAll) != 2 {
		t.Error("GetDefinition without filter failed")
	}
}

func TestNamespaceFilter_TagFilter(t *testing.T) {
	reg := newDefinitionRegistry()
	defA := makeTestDefinition("A", "fa", "ns1", []string{"tag1"})
	defB := makeTestDefinition("A", "fb", "ns2", []string{"tag2"})
	_ = reg.register(defA)
	_ = reg.register(defB)

	nsFilter := NamespaceFilter("ns1")
	tagFilter := TagFilter("tag2")
	defsNS := reg.GetDefinition("A", nsFilter)
	if len(defsNS) != 1 || defsNS[0].ns != "ns1" {
		t.Error("NamespaceFilter failed")
	}
	defsTag := reg.GetDefinition("A", tagFilter)
	if len(defsTag) != 1 || defsTag[0].tags[0] != "tag2" {
		t.Error("TagFilter failed")
	}
}

// --- AI GENERATED CODE END ---
