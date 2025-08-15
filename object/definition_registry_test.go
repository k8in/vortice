package object

import (
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

func TestDefinitionRegistry_RegisterAndLock(t *testing.T) {
	newTestDefinition := func(name string) *Definition {
		return &Definition{
			name:      name,
			typ:       reflect.TypeOf(""),
			factory:   &Factory{name: name + "_factory"},
			dependsOn: []string{},
			methods:   &Methods{},
			ns:        Core,
			scope:     Singleton,
			tags:      []string{"tag"},
		}
	}

	reg := newDefinitionRegistry()
	def1 := newTestDefinition("obj1")
	def2 := newTestDefinition("obj2")

	// Register should succeed
	if err := reg.register(def1); err != nil {
		t.Fatalf("register def1 failed: %v", err)
	}
	if err := reg.register(def2); err != nil {
		t.Fatalf("register def2 failed: %v", err)
	}

	// Duplicate factory should fail
	dupDef := newTestDefinition("obj3")
	dupDef.factory = def1.factory
	if err := reg.register(dupDef); err == nil {
		t.Error("register duplicate factory should fail")
	}

	// Lock registry
	reg.Lock()
	if !reg.readonly.Load() {
		t.Error("registry should be readonly after Lock")
	}

	// Register after lock should fail
	def4 := newTestDefinition("obj4")
	if err := reg.register(def4); err == nil {
		t.Error("register after lock should fail")
	}
}

func TestDefinitionRegistry_EntriesAndFactories(t *testing.T) {
	newTestDefinition := func(name string) *Definition {
		return &Definition{
			name:      name,
			typ:       reflect.TypeOf(""),
			factory:   &Factory{name: name + "_factory"},
			dependsOn: []string{},
			methods:   &Methods{},
			ns:        Core,
			scope:     Singleton,
			tags:      []string{"tag"},
		}
	}

	reg := newDefinitionRegistry()
	def := newTestDefinition("obj")
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

func TestLockDefinitionRegistry(t *testing.T) {
	type dummyDep struct{}
	type dummyStruct struct{}
	factory := func(a dummyDep) dummyStruct { return dummyStruct{} }
	prop := NewProperty()
	LockDefinitionRegistry()
	_, err := RegisterFactory(factory, prop)
	if err == nil {
		t.Error("RegisterFactory should fail after registry is locked")
	}
}

func TestDefinitionRegistry_GetDefinition_Filter(t *testing.T) {
	reg := newDefinitionRegistry()
	defA := &Definition{name: "A", typ: reflect.TypeOf(""), factory: &Factory{name: "fa"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	defB := &Definition{name: "A", typ: reflect.TypeOf(""), factory: &Factory{name: "fb"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	defC := &Definition{name: "B", typ: reflect.TypeOf(""), factory: &Factory{name: "fc"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	_ = reg.register(defA)
	_ = reg.register(defB)
	_ = reg.register(defC)

	// Filter: only factory name == "fa"
	filter := func(def *Definition) bool { return def.Factory().Name() == "fa" }
	defs := reg.GetDefinition("A", filter)
	if len(defs) != 1 || defs[0].factory.name != "fa" {
		t.Error("GetDefinition with filter failed")
	}

	// No filter: should return all
	defsAll := reg.GetDefinition("A", nil)
	if len(defsAll) != 2 {
		t.Error("GetDefinition without filter failed")
	}
}

func TestDefinitionRegistry_GetDefinitionNames_Filter(t *testing.T) {
	reg := newDefinitionRegistry()
	defA := &Definition{name: "A", typ: reflect.TypeOf(""), factory: &Factory{name: "fa"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	defB := &Definition{name: "B", typ: reflect.TypeOf(""), factory: &Factory{name: "fb"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	defC := &Definition{name: "C", typ: reflect.TypeOf(""), factory: &Factory{name: "fc"}, dependsOn: []string{}, methods: &Methods{}, ns: Core, scope: Singleton, tags: []string{"tag"}}
	_ = reg.register(defA)
	_ = reg.register(defB)
	_ = reg.register(defC)

	// Filter: only factory name == "fb"
	filter := func(def *Definition) bool { return def.Factory().Name() == "fb" }
	names := reg.GetDefinitionNames(filter)
	if len(names) != 1 || names[0] != "B" {
		t.Error("GetDefinitionNames with filter failed")
	}

	// No filter: should return all names
	namesAll := reg.GetDefinitionNames(nil)
	if len(namesAll) != 3 {
		t.Error("GetDefinitionNames without filter failed")
	}
}

// --- AI GENERATED CODE END ---
