package vortice

import (
	"testing"
	"vortice/object"
)

type Dummy struct{}
type DummyA struct{}
type DummyB struct{}
type DummyC struct{}
type DummyD struct{}
type DummyE struct{}
type DummyF struct{}

func dummyFactory0() Dummy                                                           { return Dummy{} }
func dummyFactory1(a DummyA) Dummy                                                   { return Dummy{} }
func dummyFactory2(a DummyA, b DummyB) Dummy                                         { return Dummy{} }
func dummyFactory3(a DummyA, b DummyB, c DummyC) Dummy                               { return Dummy{} }
func dummyFactory4(a DummyA, b DummyB, c DummyC, d DummyD) Dummy                     { return Dummy{} }
func dummyFactory5(a DummyA, b DummyB, c DummyC, d DummyD, e DummyE) Dummy           { return Dummy{} }
func dummyFactory6(a DummyA, b DummyB, c DummyC, d DummyD, e DummyE, f DummyF) Dummy { return Dummy{} }

func TestRegister0To6(t *testing.T) {
	t.Run("Register0", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register0 panicked: %v", r)
			}
		}()
		Register0(dummyFactory0)
	})
	t.Run("Register1", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register1 panicked: %v", r)
			}
		}()
		Register1(dummyFactory1)
	})
	t.Run("Register2", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register2 panicked: %v", r)
			}
		}()
		Register2(dummyFactory2)
	})
	t.Run("Register3", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register3 panicked: %v", r)
			}
		}()
		Register3(dummyFactory3)
	})
	t.Run("Register4", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register4 panicked: %v", r)
			}
		}()
		Register4(dummyFactory4)
	})
	t.Run("Register5", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register5 panicked: %v", r)
			}
		}()
		Register5(dummyFactory5)
	})
	t.Run("Register6", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Register6 panicked: %v", r)
			}
		}()
		Register6(dummyFactory6)
	})
}

func TestRegister_WithOption(t *testing.T) {
	optCalled := false
	opt := func(prop *object.Property) {
		optCalled = true
		prop.SetTag("test", "v")
	}
	// 使用唯一的 factory，避免重复注册 panic
	uniqueFactory := func() Dummy { return Dummy{} }
	Register0(uniqueFactory, opt)
	if !optCalled {
		t.Error("Option was not called in Register0")
	}
}
