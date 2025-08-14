package object

import (
	"reflect"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

func factory0() int {
	return 42
}
func factory1(a string) string {
	return a + "_ok"
}
func factory2(a int, b int) int {
	return a + b
}
func factory3(a, b, c string) string {
	return a + b + c
}

func TestFactory_Call0(t *testing.T) {
	f := newFactory(reflect.ValueOf(factory0), nil, 0)
	res := f.Call([]reflect.Value{})
	if res.Int() != 42 {
		t.Errorf("expected 42, got %v", res.Int())
	}
}

func TestFactory_Call1(t *testing.T) {
	f := newFactory(reflect.ValueOf(factory1), nil, 1)
	arg := reflect.ValueOf("hello")
	res := f.Call([]reflect.Value{arg})
	if res.String() != "hello_ok" {
		t.Errorf("expected hello_ok, got %v", res.String())
	}
}

func TestFactory_Call2(t *testing.T) {
	f := newFactory(reflect.ValueOf(factory2), nil, 2)
	arg1 := reflect.ValueOf(10)
	arg2 := reflect.ValueOf(32)
	res := f.Call([]reflect.Value{arg1, arg2})
	if res.Int() != 42 {
		t.Errorf("expected 42, got %v", res.Int())
	}
}

func TestFactory_Call3(t *testing.T) {
	f := newFactory(reflect.ValueOf(factory3), nil, 3)
	args := []reflect.Value{
		reflect.ValueOf("a"),
		reflect.ValueOf("b"),
		reflect.ValueOf("c"),
	}
	res := f.Call(args)
	if res.String() != "abc" {
		t.Errorf("expected abc, got %v", res.String())
	}
}

func TestFactory_Getters(t *testing.T) {
	f := newFactory(reflect.ValueOf(factory1), nil, 1)
	if f.Name() == "" {
		t.Error("Name() should not be empty")
	}
	if f.File() == "" {
		t.Error("File() should not be empty")
	}
	if f.Line() <= 0 {
		t.Error("Line() should be positive")
	}
	if f.Func().Pointer() != reflect.ValueOf(factory1).Pointer() {
		t.Error("Func() should return correct reflect.Value")
	}
	if f.Argv() != nil {
		t.Error("Argv() should be nil for this test")
	}
	if f.Argn() != 1 {
		t.Errorf("Argn() should be 1, got %d", f.Argn())
	}
}

// --- AI GENERATED CODE END ---
