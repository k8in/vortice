package util

import (
	"errors"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestLogger_DefaultAndSet(t *testing.T) {
	// 记录初始 logger
	l1 := Logger()
	if l1 == nil {
		t.Fatalf("expected non-nil default logger")
	}
	// 设置自定义 logger
	custom, _ := zap.NewDevelopment()
	SetLogger(custom)
	l2 := Logger()
	if l2 != custom {
		t.Fatalf("SetLogger not applied")
	}
	// 置为 nil 触发内部回退
	SetLogger(nil)
	l3 := Logger()
	if l3 == nil || l3 == custom {
		t.Fatalf("expected new default logger after nil, got %v", l3)
	}
}

func TestLogger_FallbackOnBuildError(t *testing.T) {
	// 保存原构建函数
	orig := buildLogger
	defer func() { buildLogger = orig }()

	called := false
	buildLogger = func(cfg zap.Config) (*zap.Logger, error) {
		called = true
		return nil, errors.New("boom")
	}
	l := defaultLogger()
	if !called {
		t.Fatalf("expected buildLogger override to be called")
	}
	// 失败后应返回 Nop（不可直接比较指针，但可以比较核心是否启用高等级日志）
	if l == nil {
		t.Fatalf("expected fallback logger not nil")
	}
	// 再次调用 Logger() 时仍可正常工作
	SetLogger(nil)
	if Logger() == nil {
		t.Fatalf("Logger fallback after build error should be non-nil")
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	orig := Logger()
	defer SetLogger(orig)

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			if i%5 == 0 {
				SetLogger(zap.NewNop())
			} else {
				_ = Logger()
			}
		}(i)
	}
	close(start)
	wg.Wait()
	// 最终再调用一次保证无竞态
	_ = Logger()
}

func TestDefaultLogger_Idempotent(t *testing.T) {
	l1 := defaultLogger()
	l2 := defaultLogger()
	if l1 == nil || l2 == nil {
		t.Fatalf("defaultLogger should not return nil")
	}
	// 不要求相同实例，但至少可用
	_ = l1.Sugar()
	_ = l2.Sugar()
}
