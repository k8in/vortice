package util

import (
	"context"
	"errors"
	"strings"
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

// 新增：多个并发任务（成功+错误混合），验证 channel 正确返回
func TestTaskGroup_Go_Multiple(t *testing.T) {
	g := NewTaskGroup("multi")
	ctx := context.Background()

	ch1 := g.Go(ctx, func(ctx context.Context) error {
		return nil
	})
	ch2 := g.Go(ctx, func(ctx context.Context) error {
		return errors.New("x")
	})
	ch3 := g.Go(ctx, func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	if err := <-ch1; err != nil {
		t.Fatalf("ch1 expected nil got %v", err)
	}
	if err := <-ch2; err == nil || err.Error() != "x" {
		t.Fatalf("ch2 expected x got %v", err)
	}
	if err := <-ch3; err != nil {
		t.Fatalf("ch3 expected nil got %v", err)
	}
}

// 新增：安全的上下文取消测试，避免 goroutine 永久阻塞（函数在取消后尽快返回）
func TestTaskGroup_GoAndWait_ContextCancel_NoLeak(t *testing.T) {
	g := NewTaskGroup("cancel-safe")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := g.GoAndWait(ctx, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			// 模拟原测试中的长耗时，但返回前让 goroutine 能发送结果
			time.Sleep(15 * time.Millisecond) // 保证主 select 已走 ctx.Done 分支
			return nil
		case <-time.After(200 * time.Millisecond):
			return nil
		}
	})
	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("expected deadline error, got %v", err)
	}
	if time.Since(start) > 150*time.Millisecond {
		t.Fatalf("cancel path took too long, possible leak")
	}
}

// 深度递归制造超长堆栈，触发截断逻辑 (len(stack) > 4096)
func deepPanic(n int) {
	if n == 0 {
		panic("deep panic")
	}
	deepPanic(n - 1)
}

// 新增：panic 堆栈截断覆盖
func TestTaskGroup_Go_Panic_LongStackTruncation(t *testing.T) {
	g := NewTaskGroup("long-stack")
	ctx := context.Background()
	ch := g.Go(ctx, func(ctx context.Context) error {
		deepPanic(600) // 产生深栈
		return nil
	})
	err := <-ch
	if err == nil {
		t.Fatalf("expected panic error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "deep panic") {
		t.Fatalf("expected panic message in error, got %s", msg)
	}
	if len(msg) > 5000 {
		t.Fatalf("expected truncated stack (<5000 chars), got len=%d", len(msg))
	}
}

// 新增：直接调用 Go 并立即读取，覆盖最短路径
func TestTaskGroup_Go_Immediate(t *testing.T) {
	g := NewTaskGroup("immediate")
	ctx := context.Background()
	ch := g.Go(ctx, func(ctx context.Context) error { return nil })
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout waiting for immediate task")
	}
}

// --- AI GENERATED CODE END ---
