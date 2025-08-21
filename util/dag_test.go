package util

import (
	"strings"
	"testing"
)

// --- AI GENERATED CODE BEGIN ---

// 空图
func TestDAG_Empty(t *testing.T) {
	d := NewDAG()
	order, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 0 {
		t.Fatalf("expected empty result, got %v", order)
	}
}

// 隐式依赖（B 未显式添加），验证依赖在前（reverse 后 B 应排在 A 前）
func TestDAG_ImplicitDependency(t *testing.T) {
	d := NewDAG()
	d.AddNode("A", "B") // B 未显式 AddNode
	order, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 {
		t.Fatalf("expected 2 nodes, got %v", order)
	}
	idx := indexMap(order)
	if !(idx["B"] < idx["A"]) {
		t.Fatalf("expected B before A (dependency-first), got %v", order)
	}
}

// 简单链 A -> B -> C
func TestDAG_SimpleChain(t *testing.T) {
	d := NewDAG()
	d.AddNode("A")
	d.AddNode("B", "A")
	d.AddNode("C", "B")
	o, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := indexMap(o)
	if !(m["A"] < m["B"] && m["B"] < m["C"]) {
		t.Fatalf("chain order invalid: %v", o)
	}
}

// 自环 A -> A
func TestDAG_SelfLoop(t *testing.T) {
	d := NewDAG()
	d.AddNode("A", "A")
	_, err := d.Sort()
	if err == nil || !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected self loop cycle error, got %v", err)
	}
	if !strings.Contains(err.Error(), "A") {
		t.Fatalf("error should mention node A, got %v", err)
	}
}

// 复杂环: A -> B, B -> C, C -> D, D -> B  (环 B-C-D-B)
func TestDAG_ComplexCycle(t *testing.T) {
	d := NewDAG()
	d.AddNode("A", "B")
	d.AddNode("B", "C")
	d.AddNode("C", "D")
	d.AddNode("D", "B")
	_, err := d.Sort()
	if err == nil {
		t.Fatalf("expected cycle error")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("error should declare cycle, got %v", err)
	}
	// 环中的节点
	for _, n := range []string{"B", "C", "D"} {
		if !strings.Contains(err.Error(), n) {
			t.Fatalf("cycle error should mention %s, got %v", n, err)
		}
	}
}

// 重复依赖 B -> A,A
func TestDAG_DuplicateDependencies(t *testing.T) {
	d := NewDAG()
	d.AddNode("A")
	d.AddNode("B", "A", "A", "A")
	o, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := indexMap(o)
	if !(m["A"] < m["B"]) {
		t.Fatalf("dependency order wrong: %v", o)
	}
	// 确保不重复节点
	seen := map[string]bool{}
	for _, v := range o {
		if seen[v] {
			t.Fatalf("node %s appears twice in %v", v, o)
		}
		seen[v] = true
	}
}

// 大型无环图（多源多汇）
func TestDAG_LargeAcyclic(t *testing.T) {
	d := NewDAG()
	// 层级:
	//   L0: R
	//   L1: A,B
	//   L2: C(D), D(E,F)
	//   L3: G(H)
	d.AddNode("R")
	d.AddNode("A", "R")
	d.AddNode("B", "R")
	d.AddNode("C", "A")
	d.AddNode("D", "A", "B")
	d.AddNode("E", "C")
	d.AddNode("F", "C")
	d.AddNode("G", "D", "E")
	d.AddNode("H", "G")

	o, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := indexMap(o)
	// 基础校验：父级在子级前
	assertBefore(t, m, "R", "A", "B")
	assertBefore(t, m, "A", "C", "D")
	assertBefore(t, m, "B", "D")
	assertBefore(t, m, "C", "E", "F")
	assertBefore(t, m, "D", "G")
	assertBefore(t, m, "E", "G")
	assertBefore(t, m, "G", "H")
}

// 隐式+自环+未声明依赖节点同时出现：A 依赖 (B, A)
func TestDAG_ImplicitCycle(t *testing.T) {
	d := NewDAG()
	d.AddNode("A", "B", "A") // B 未显式添加；A 自环
	_, err := d.Sort()
	if err == nil || !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected cycle error, got %v", err)
	}
	// 应包含 A 与 B
	for _, n := range []string{"A", "B"} {
		if !strings.Contains(err.Error(), n) {
			t.Fatalf("error should mention %s, got %v", n, err)
		}
	}
}

// 带外部已拓扑节点的环：X 独立；A-B-C-A 且 C 额外依赖 X，覆盖 findCycle 中对非 subgraph 依赖跳过的分支
func TestDAG_CycleWithExternalEdge(t *testing.T) {
	d := NewDAG()
	d.AddNode("X")      // 外部无依赖节点
	d.AddNode("A", "B") // 环起点
	d.AddNode("B", "C")
	d.AddNode("C", "A", "X") // 指向已可拓扑节点 X，用于触发 subgraph 跳过
	_, err := d.Sort()
	if err == nil {
		t.Fatalf("expected cycle error")
	}
	if !strings.Contains(err.Error(), "cycle detected") {
		t.Fatalf("expected cycle detected message, got %v", err)
	}
	for _, n := range []string{"A", "B", "C"} {
		if !strings.Contains(err.Error(), n) {
			t.Fatalf("cycle error should list %s, got %v", n, err)
		}
	}
}

// 多个零入度节点 + 一个依赖链，验证 reverse 后依赖优先
func TestDAG_MultiZeroIndegree(t *testing.T) {
	d := NewDAG()
	d.AddNode("K")
	d.AddNode("L")
	d.AddNode("M", "K")
	order, err := d.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := indexMap(order)
	if !(m["K"] < m["M"]) {
		t.Fatalf("expected K before M (dependency-first), got %v", order)
	}
	if m["L"] == m["K"] || m["L"] == m["M"] {
		t.Fatalf("positions should be distinct: %v", order)
	}
	// 再次调用确保无副作用
	order2, err2 := d.Sort()
	if err2 != nil || len(order2) != len(order) {
		t.Fatalf("second sort unstable: %v %v", order2, err2)
	}
}

// 人工构造：所有节点入度>0 但无环（真实拓扑不可能），用于覆盖 findCycle 返回空分支
func TestDAG_FindCycleNoRealCycle(t *testing.T) {
	d := NewDAG()
	d.AddNode("A", "B")
	d.AddNode("B") // 没有回边
	// 手工伪造 inDegree（把 A,B 都设为 >0），让 findCycle 进入 DFS 但找不到环
	inDegree := map[string]int{"A": 1, "B": 1}
	cycle := d.findCycle(inDegree)
	if len(cycle) != 0 {
		t.Fatalf("expected no cycle (contrived state), got %v", cycle)
	}
}

// 直接调用 remainingNodes 覆盖
func TestDAG_RemainingNodesDirect(t *testing.T) {
	d := NewDAG()
	inDegree := map[string]int{"A": 0, "B": 2, "C": 1}
	rem := d.remainingNodes(inDegree)
	expect := map[string]bool{"B": true, "C": true}
	if len(rem) != 2 {
		t.Fatalf("expected 2 remaining nodes, got %v", rem)
	}
	for _, n := range rem {
		if !expect[n] {
			t.Fatalf("unexpected remaining node %s in %v", n, rem)
		}
	}
}

// 直接调用 reverse 覆盖多长度场景
func TestDAG_ReverseVariants(t *testing.T) {
	d := NewDAG()
	cases := [][]string{
		{},
		{"A"},
		{"A", "B"},
		{"A", "B", "C"},
	}
	for _, c := range cases {
		out := d.reverse(append([]string{}, c...)) // 复制防止原地修改影响后续
		// 再次 reverse 应还原
		restore := d.reverse(append([]string{}, out...))
		if len(restore) != len(c) {
			t.Fatalf("length mismatch after double reverse: %v -> %v -> %v", c, out, restore)
		}
		for i := range c {
			if restore[i] != c[i] {
				t.Fatalf("double reverse failed: %v -> %v -> %v", c, out, restore)
			}
		}
	}
}

// 直接调用 extractCycle 覆盖路径截取
func TestDAG_ExtractCycleDirect(t *testing.T) {
	path := []string{"X", "Y", "Z", "W"}
	cycle := extractCycle(path, "Y")
	// 期望: Y,Z,W,Y
	if len(cycle) != 4 || cycle[0] != "Y" || cycle[len(cycle)-1] != "Y" {
		t.Fatalf("unexpected cycle slice: %v", cycle)
	}
	if strings.Join(cycle, ",") != "Y,Z,W,Y" {
		t.Fatalf("cycle content mismatch: %v", cycle)
	}
}

// 工具: 构造索引
func indexMap(order []string) map[string]int {
	m := map[string]int{}
	for i, v := range order {
		m[v] = i
	}
	return m
}

// 工具: 断言 a 在多 b 之前
func assertBefore(t *testing.T, idx map[string]int, a string, bs ...string) {
	t.Helper()
	for _, b := range bs {
		if idx[a] >= idx[b] {
			t.Fatalf("expected %s before %s (idx=%d,%d)", a, b, idx[a], idx[b])
		}
	}
}

// --- AI GENERATED CODE END ---
