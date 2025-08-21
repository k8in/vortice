package util

import (
	"fmt"
	"strings"
)

// DAG represents a Directed Acyclic Graph.
type DAG struct {
	nodes map[string][]string
}

// NewDAG creates a new DAG instance.
func NewDAG() *DAG {
	return &DAG{nodes: make(map[string][]string)}
}

// AddNode adds a node and its dependencies.
// node: node name
// deps: names of nodes this node depends on
// Note: No deduplication, duplicate dependencies are allowed; no self-loop or missing dependency check.
func (dag *DAG) AddNode(node string, deps ...string) {
	dag.nodes[node] = append(dag.nodes[node], deps...)
}

// Sort performs topological sorting on the DAG.
// Returns: sorted node list (dependency first), or error if a cycle is detected.
func (dag *DAG) Sort() ([]string, error) {
	inDegree := map[string]int{}
	for node := range dag.nodes {
		inDegree[node] = 0
	}
	for _, deps := range dag.nodes {
		for _, dep := range deps {
			inDegree[dep]++ // 允许隐式节点：未显式 AddNode 的依赖会在此加入
		}
	}

	var queue []string
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}

	var result []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		for _, dep := range dag.nodes[node] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(result) != len(inDegree) {
		// 有残留节点 => 存在环
		cycle := dag.findCycle(inDegree)
		remaining := dag.remainingNodes(inDegree)
		if len(cycle) > 0 {
			return nil, fmt.Errorf("cycle detected: %s (all remaining nodes: %v)",
				strings.Join(cycle, " -> "), remaining)
		}
		return nil, fmt.Errorf("cycle detected among nodes: %v", remaining)
	}

	// The result is reversed to ensure dependency-first order.
	return dag.reverse(result), nil
}

// findCycle 尝试在残留节点子图中找到一条环路径
func (dag *DAG) findCycle(inDegree map[string]int) []string {
	// 只对 inDegree > 0 的节点进行 DFS
	subgraph := map[string]bool{}
	for n, deg := range inDegree {
		if deg > 0 {
			subgraph[n] = true
		}
	}
	visited := map[string]bool{}
	stack := map[string]bool{}
	path := []string{}
	var cycle []string

	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		stack[node] = true
		path = append(path, node)

		for _, dep := range dag.nodes[node] {
			// 只考虑仍在 subgraph 的节点
			if !subgraph[dep] {
				continue
			}
			if !visited[dep] {
				if dfs(dep) {
					return true
				}
			} else if stack[dep] {
				// 找到回边，截取 cycle
				cycle = extractCycle(path, dep)
				return true
			}
		}

		stack[node] = false
		path = path[:len(path)-1]
		return false
	}

	for n := range subgraph {
		if !visited[n] {
			if dfs(n) {
				break
			}
		}
	}
	return cycle
}

// extractCycle 从当前 DFS 路径中截取形成环的部分
func extractCycle(path []string, start string) []string {
	for i, n := range path {
		if n == start {
			cp := append([]string{}, path[i:]...)
			cp = append(cp, start) // 闭合
			return cp
		}
	}
	return nil
}

// remainingNodes 返回还未处理完的节点集合（inDegree>0）
func (dag *DAG) remainingNodes(inDegree map[string]int) []string {
	var rem []string
	for n, deg := range inDegree {
		if deg > 0 {
			rem = append(rem, n)
		}
	}
	return rem
}

// reverse reverses the sorted result to ensure dependency-first order.
func (dag *DAG) reverse(result []string) []string {
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
