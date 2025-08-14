package object

import (
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

type testComponent struct {
	initialized bool
	running     bool
	destroyed   bool
}

func (c *testComponent) Init() {
	c.initialized = true
}
func (c *testComponent) Destroy() {
	c.destroyed = true
}
func (c *testComponent) Start() {
	c.running = true
}
func (c *testComponent) Stop() {
	c.running = false
}
func (c *testComponent) Running() bool {
	return c.running
}

type noLifecycle struct{}

func TestMethods_Lifecycle(t *testing.T) {
	objType := reflect.TypeOf(&testComponent{})
	methods := newMethods(objType)
	comp := &testComponent{}
	val := reflect.ValueOf(comp)

	ok, _, err := methods.CallInit(val)
	if !ok || err != nil || !comp.initialized {
		t.Errorf("Init method failed: ok=%v, err=%v, initialized=%v", ok, err, comp.initialized)
	}

	ok, _, err = methods.CallStart(val)
	if !ok || err != nil || !comp.running {
		t.Errorf("Start method failed: ok=%v, err=%v, running=%v", ok, err, comp.running)
	}

	ok, _, err = methods.CallRunning(val)
	if !ok || err != nil || !comp.running {
		t.Errorf("Running method failed: ok=%v, err=%v, running=%v", ok, err, comp.running)
	}

	ok, _, err = methods.CallStop(val)
	if !ok || err != nil || comp.running {
		t.Errorf("Stop method failed: ok=%v, err=%v, running=%v", ok, err, comp.running)
	}

	ok, _, err = methods.CallDestroy(val)
	if !ok || err != nil || !comp.destroyed {
		t.Errorf("Destroy method failed: ok=%v, err=%v, destroyed=%v", ok, err, comp.destroyed)
	}
}

func TestMethods_NoLifecycle(t *testing.T) {
	objType := reflect.TypeOf(&noLifecycle{})
	methods := newMethods(objType)
	comp := &noLifecycle{}
	val := reflect.ValueOf(comp)

	ok, _, err := methods.CallInit(val)
	if ok || err != nil {
		t.Errorf("Init should not exist: ok=%v, err=%v", ok, err)
	}
	ok, _, err = methods.CallStart(val)
	if ok || err != nil {
		t.Errorf("Start should not exist: ok=%v, err=%v", ok, err)
	}
	ok, _, err = methods.CallRunning(val)
	if ok || err != nil {
		t.Errorf("Running should not exist: ok=%v, err=%v", ok, err)
	}
	ok, _, err = methods.CallStop(val)
	if ok || err != nil {
		t.Errorf("Stop should not exist: ok=%v, err=%v", ok, err)
	}
	ok, _, err = methods.CallDestroy(val)
	if ok || err != nil {
		t.Errorf("Destroy should not exist: ok=%v, err=%v", ok, err)
	}
}

// --- AI GENERATED CODE END ---
