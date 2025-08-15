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

func (c *testComponent) Init() error {
	c.initialized = true
	return nil
}
func (c *testComponent) Destroy() error {
	c.destroyed = true
	return nil
}
func (c *testComponent) Start() error {
	c.running = true
	return nil
}
func (c *testComponent) Stop() error {
	c.running = false
	return nil
}
func (c *testComponent) Running() bool {
	return c.running
}

type noLifecycle struct{}

func (c *noLifecycle) Init() error    { return nil }
func (c *noLifecycle) Destroy() error { return nil }
func (c *noLifecycle) Start() error   { return nil }
func (c *noLifecycle) Stop() error    { return nil }
func (c *noLifecycle) Running() bool  { return false }

func TestMethods_Lifecycle(t *testing.T) {
	objType := reflect.TypeOf(&testComponent{})
	methods := newMethods(objType)
	comp := &testComponent{}
	val := reflect.ValueOf(comp)

	if err := methods.CallInit(val); err != nil || !comp.initialized {
		t.Errorf("Init method failed: err=%v, initialized=%v", err, comp.initialized)
	}

	if err := methods.CallStart(val); err != nil || !comp.running {
		t.Errorf("Start method failed: err=%v, running=%v", err, comp.running)
	}

	running, err := methods.CallRunning(val)
	if err != nil || !running {
		t.Errorf("Running method failed: err=%v, running=%v", err, running)
	}

	if err := methods.CallStop(val); err != nil || comp.running {
		t.Errorf("Stop method failed: err=%v, running=%v", err, comp.running)
	}

	if err := methods.CallDestroy(val); err != nil || !comp.destroyed {
		t.Errorf("Destroy method failed: err=%v, destroyed=%v", err, comp.destroyed)
	}
}

func TestMethods_NoLifecycle(t *testing.T) {
	objType := reflect.TypeOf(&noLifecycle{})
	methods := newMethods(objType)
	comp := &noLifecycle{}
	val := reflect.ValueOf(comp)

	if err := methods.CallInit(val); err != nil {
		t.Errorf("Init should succeed: err=%v", err)
	}
	if err := methods.CallStart(val); err != nil {
		t.Errorf("Start should succeed: err=%v", err)
	}
	running, err := methods.CallRunning(val)
	if err != nil {
		t.Errorf("Running should succeed: err=%v", err)
	}
	if running {
		t.Errorf("Running should be false for noLifecycle")
	}
	if err := methods.CallStop(val); err != nil {
		t.Errorf("Stop should succeed: err=%v", err)
	}
	if err := methods.CallDestroy(val); err != nil {
		t.Errorf("Destroy should succeed: err=%v", err)
	}
}

// --- AI GENERATED CODE END ---
