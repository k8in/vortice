package container

import (
	"context"
	"testing"
	"vortice/object"
)

func TestWithCoreContext_Basic(t *testing.T) {
	ctx := WithCoreContext(context.Background())
	if ctx == nil {
		t.Fatalf("WithCoreContext returned nil")
	}
	// GetFilters should return empty slice
	if fs := ctx.GetFilters(); fs == nil || len(fs) != 0 {
		t.Fatalf("GetFilters should return empty slice by default")
	}
	// GetObjects should return empty map
	if m := ctx.GetObjects(); m == nil || len(m) != 0 {
		t.Fatalf("GetObjects should return empty map by default")
	}
}

func TestCoreContext_SetFilter_NoPanic(t *testing.T) {
	ctx := WithCoreContext(context.Background())
	// SetFilter should accept empty and non-empty filters without panic
	ctx.SetFilter()
	ctx.SetFilter(nil)
	ctx.SetFilter(func(*object.Definition) bool { return true })
	// Note: GetFilters always returns empty slice per current implementation
	if fs := ctx.GetFilters(); len(fs) != 0 {
		t.Fatalf("GetFilters should still return empty slice")
	}
}
