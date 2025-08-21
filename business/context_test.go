package business

import (
	"context"
	"testing"
)

func TestWithContextNamespace(t *testing.T) {
	ns := "alpha"
	c := WithContext(context.Background(), ns)
	if got := c.Namespace(); got != ns {
		t.Fatalf("expected namespace %q, got %q", ns, got)
	}
	// 再创建一个不同命名空间，确保互不影响
	ns2 := "beta"
	c2 := WithContext(context.Background(), ns2)
	if got := c2.Namespace(); got != ns2 {
		t.Fatalf("expected namespace %q, got %q", ns2, got)
	}
	if c.Namespace() != ns {
		t.Fatalf("first context namespace mutated: want %q got %q", ns, c.Namespace())
	}
}

func TestGetNamespaceVariants(t *testing.T) {
	// nil context
	if got := GetNamespace(nil); got != "" {
		t.Fatalf("expected empty for nil context, got %q", got)
	}
	// 普通 context 无 namespace
	if got := GetNamespace(context.Background()); got != "" {
		t.Fatalf("expected empty for plain context, got %q", got)
	}
	// WithContext 包装
	ns := "gamma"
	c := WithContext(context.Background(), ns)
	if got := GetNamespace(c); got != ns {
		// 通过内部持有的原生 context 获取（利用 CoreContext 暴露的 Context() 兼容）
		t.Fatalf("expected %q, got %q", ns, got)
	}
}
