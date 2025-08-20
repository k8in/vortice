package vortice

import (
	"testing"
	"vortice/object"
)

type Dummy struct{}
type Dummy0 struct{}
type Dummy1 struct{}
type Dummy2 struct{}
type Dummy3 struct{}
type Dummy4 struct{}
type Dummy5 struct{}
type Dummy6 struct{}
type DummyA struct{}
type DummyB struct{}
type DummyC struct{}
type DummyD struct{}
type DummyE struct{}
type DummyF struct{}

func dummyFactory0() Dummy0                                                 { return Dummy0{} }
func dummyFactory1(a DummyA) Dummy1                                         { return Dummy1{} }
func dummyFactory2(a DummyA, b DummyB) Dummy2                               { return Dummy2{} }
func dummyFactory3(a DummyA, b DummyB, c DummyC) Dummy3                     { return Dummy3{} }
func dummyFactory4(a DummyA, b DummyB, c DummyC, d DummyD) Dummy4           { return Dummy4{} }
func dummyFactory5(a DummyA, b DummyB, c DummyC, d DummyD, e DummyE) Dummy5 { return Dummy5{} }
func dummyFactory6(a DummyA, b DummyB, c DummyC, d DummyD, e DummyE, f DummyF) Dummy6 {
	return Dummy6{}
}

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
		prop.SetTags(object.NewTag("test", "v"))
	}
	// 使用唯一的 factory/type，避免与其它测试重复注册
	type UniqueDummy struct{}
	uniqueFactory := func() UniqueDummy { return UniqueDummy{} }
	Register0(uniqueFactory, opt)
	if !optCalled {
		t.Error("Option was not called in Register0")
	}
}
