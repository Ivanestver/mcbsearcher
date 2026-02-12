package cycles_alg

import (
	"cycles/data_structs"
	"cycles/types"
)

type GraphJson struct {
	Points []*types.Point
	Edges  []*types.Edge
	Graph  data_structs.Graph
}

func NewGraphJson(points []*types.Point, edges []*types.Edge, graph data_structs.Graph) *GraphJson {
	return &GraphJson{
		Points: points,
		Edges:  edges,
		Graph:  graph,
	}
}
