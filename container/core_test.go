package container

import (
	"context"
	"testing"
)

func TestDefaultCore_Singleton(t *testing.T) {
	 c1 := DefaultCore()
	 c2 := DefaultCore()
	 if c1 == nil || c2 == nil {
	 	 t.Fatalf("DefaultCore returned nil")
	 }
	 if c1 != c2 {
	 	 t.Fatalf("DefaultCore should return the same instance (singleton)")
	 }
}

func TestCore_Init_Start_Shutdown_NoServices(t *testing.T) {
	 ctx := context.Background()
	 c := NewCore(ctx)
	 if c == nil {
	 	 t.Fatalf("NewCore returned nil")
	 }

	 if err := c.Init(); err != nil {
	 	 t.Fatalf("Init failed: %v", err)
	 }

	 if err := c.Start(); err != nil {
	 	 t.Fatalf("Start failed without services: %v", err)
	 }

	 // Shutdown is void; ensure it does not panic
	 defer func() {
	 	 if r := recover(); r != nil {
	 	 	 t.Fatalf("Shutdown panicked: %v", r)
	 	 }
	 }()
	 c.Shutdown()
}


