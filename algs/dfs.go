package algs

import (
	"cycles/data_structs"
	"cycles/types"
)

type DFS struct {
	points []*types.Point
	graph  data_structs.Graph
}

func (dfs *DFS) Traverse(startingPoint int) types.Path {
	stack := data_structs.Stack[*types.Edge]{}
	for i := len(dfs.graph[startingPoint]) - 1; i >= 0; i-- {
		if dfs.graph[startingPoint][i] != nil {
			stack.Push(dfs.graph[startingPoint][i])
		}
	}
	dfs.points[startingPoint].State = types.STATE_BLACK
	tree := types.Path{}
	for !stack.IsEmpty() {
		edge := stack.Pop()
		if edge.State == types.STATE_BLACK {
			continue
		}
		edge.State = types.STATE_BLACK
		otherPoint := 0
		if dfs.points[edge.Edge[0]].State == types.STATE_WHITE {
			otherPoint = edge.Edge[0]
		} else if dfs.points[edge.Edge[1]].State == types.STATE_WHITE {
			otherPoint = edge.Edge[1]
		} else {
			continue
		}
		tree = append(tree, edge)
		for i := len(dfs.graph[otherPoint]) - 1; i >= 0; i-- {
			if nextEdge := dfs.graph[otherPoint][i]; nextEdge != nil && nextEdge.State != types.STATE_BLACK {
				stack.Push(nextEdge)
			}
		}
		dfs.points[otherPoint].State = types.STATE_BLACK
	}
	dfs.reset()
	return tree
}

func (dfs *DFS) reset() {
	for i := 0; i < len(dfs.points); i++ {
		dfs.points[i].State = types.STATE_WHITE
	}
}

func MakeDFS(points []*types.Point, graph data_structs.Graph) *DFS {
	return &DFS{
		points: points,
		graph:  graph,
	}
}
