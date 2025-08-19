package vortice

import (
	"context"
	"testing"
)

// --- Types for Get test ---
type getObj struct{ id int }

// --- Types for GetElem test ---
type iSvc interface{ Name() string }
type svcImpl struct{}

func (s *svcImpl) Name() string { return "svc" }

func TestGet_ReturnsRegisteredPointerStruct(t *testing.T) {
	// unique type: getObj
	factory := func() *getObj { return &getObj{id: 42} }
	Register0(factory)
	obj := Get(context.Background(), (*getObj)(nil))
	if obj == nil {
		t.Fatalf("Get returned nil")
	}
	if obj.id != 42 {
		t.Fatalf("unexpected value: %v", obj.id)
	}
}

func TestGetElem_ReturnsRegisteredInterface(t *testing.T) {
	// register a factory that returns an interface value
	factory := func() iSvc { return &svcImpl{} }
	Register0(factory)
	var svc iSvc = GetElem(context.Background(), (*iSvc)(nil))
	if svc == nil {
		t.Fatalf("GetElem returned nil")
	}
	if svc.Name() != "svc" {
		t.Fatalf("unexpected interface impl: %v", svc.Name())
	}
}


