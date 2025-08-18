package util

import (
	"context"
	"errors"
	"testing"
	"time"
)

// --- AI GENERATED CODE BEGIN ---

func TestTaskGroup_GoAndWait_Success(t *testing.T) {
	g := NewTaskGroup("test")
	ctx := context.Background()
	err := g.GoAndWait(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestTaskGroup_GoAndWait_Error(t *testing.T) {
	g := NewTaskGroup("test")
	ctx := context.Background()
	expected := errors.New("fail")
	err := g.GoAndWait(ctx, func(ctx context.Context) error {
		return expected
	})
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("expected error %v, got %v", expected, err)
	}
}

func TestTaskGroup_GoAndWait_ContextCancel(t *testing.T) {
	g := NewTaskGroup("test")
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	err := g.GoAndWait(ctx, func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 100)
		return nil
	})
	if err == nil || ctx.Err() == nil || err.Error() != "test wait: context deadline exceeded" {
		t.Errorf("expected context deadline error, got %v", err)
	}
}

func TestTaskGroup_Go_Panic(t *testing.T) {
	g := NewTaskGroup("test")
	ctx := context.Background()
	ch := g.Go(ctx, func(ctx context.Context) error {
		panic("panic error")
	})
	err := <-ch
	if err == nil || err.Error() == "" || err.Error()[:11] != "panic error" {
		t.Errorf("expected panic error, got %v", err)
	}
}

// --- AI GENERATED CODE END ---
