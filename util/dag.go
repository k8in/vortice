package util

import "fmt"

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
	if deps == nil || len(deps) == 0 {
		return
	}
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
			// 1. No deduplication for dep, duplicate dependencies will increase in-degree, but do not affect the final result.
			// 2. If dep is not explicitly added via AddNode, it will be automatically added to inDegree here,
			// but dag.nodes[dep] will not have its own dependency list.
			inDegree[dep]++
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

	if len(result) != len(dag.nodes) {
		return nil, fmt.Errorf("cycle detected or missing dependency in the DAG")
	}

	// The result is reversed to ensure dependency-first order.
	return dag.reverse(result), nil
}

// reverse reverses the sorted result to ensure dependency-first order.
func (dag *DAG) reverse(result []string) []string {
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
