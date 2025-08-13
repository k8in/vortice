package util

import "fmt"

// DAG 表示一个有向无环图（Directed Acyclic Graph）
type DAG struct {
	nodes map[string][]string
}

// NewDAG 创建一个新的 DAG 实例
func NewDAG() *DAG {
	return &DAG{nodes: make(map[string][]string)}
}

// AddNode 添加一个节点及其依赖关系
// node: 节点名称
// deps: 该节点依赖的其他节点名称
// 注意：不会去重，允许重复依赖；不会检测自环和缺失依赖
func (dag *DAG) AddNode(node string, deps ...string) {
	dag.nodes[node] = append(dag.nodes[node], deps...)
}

// Sort 对 DAG 进行拓扑排序
// 返回：排序后的节点列表（依赖优先），或检测到环时报错
func (dag *DAG) Sort() ([]string, error) {
	inDegree := map[string]int{}
	for node := range dag.nodes {
		inDegree[node] = 0
	}
	for _, deps := range dag.nodes {
		for _, dep := range deps {
			// 1.不会对dep去重，重复依赖会导致入度统计增加，但不影响最终结果
			// 2.如果 dep 没有被显式 AddNode，则此处会自动补全到 inDegree，
			// 但 dag.nodes[dep] 不会有依赖内容（即不会有自己的依赖列表）
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

	// 结果经过 reverse，保证依赖优先
	return dag.reverse(result), nil
}

// reverse 反转排序结果，使依赖优先
func (dag *DAG) reverse(result []string) []string {
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
