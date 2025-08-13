package util

import (
	"testing"
)

func TestDAG_SimpleSort(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A")
	dag.AddNode("B", "A")
	dag.AddNode("C", "B")

	order, err := dag.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 只验证依赖关系，不验证严格顺序
	index := map[string]int{}
	for i, v := range order {
		index[v] = i
	}
	if !(index["A"] < index["B"] && index["B"] < index["C"]) {
		t.Errorf("topological order incorrect: %v", order)
	}
}

func TestDAG_CycleDetection(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A", "B")
	dag.AddNode("B", "A")

	_, err := dag.Sort()
	if err == nil {
		t.Errorf("expected cycle error, got nil")
	}
}

func TestDAG_MissingDependency(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A", "B") // B 未显式添加

	_, err := dag.Sort()
	if err == nil {
		t.Errorf("expected missing dependency error, got nil")
	}
}

func TestDAG_SelfLoop(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A", "A")

	_, err := dag.Sort()
	if err == nil {
		t.Errorf("expected self-loop error, got nil")
	}
}

func TestDAG_DuplicateDependency(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A")
	dag.AddNode("B", "A", "A") // 重复依赖

	order, err := dag.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 只验证依赖关系，不验证严格顺序
	index := map[string]int{}
	for i, v := range order {
		index[v] = i
	}
	if !(index["A"] < index["B"]) {
		t.Errorf("topological order incorrect: %v", order)
	}
}

func TestDAG_Empty(t *testing.T) {
	dag := NewDAG()
	order, err := dag.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 0 {
		t.Errorf("expected empty order, got %v", order)
	}
}

func TestDAG_RepeatedDependency(t *testing.T) {
	dag := NewDAG()
	dag.AddNode("A")
	dag.AddNode("B", "A", "A", "A") // 多次重复依赖

	order, err := dag.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 检查排序结果只包含唯一节点
	seen := map[string]struct{}{}
	for _, v := range order {
		if _, ok := seen[v]; ok {
			t.Errorf("node %v appears more than once in result: %v", v, order)
		}
		seen[v] = struct{}{}
	}
	// 检查依赖关系
	index := map[string]int{}
	for i, v := range order {
		index[v] = i
	}
	if !(index["A"] < index["B"]) {
		t.Errorf("topological order incorrect with repeated dependency: %v", order)
	}
}
