package cycles_alg

import (
	"container/heap"
	"cycles/algs"
	"cycles/data_structs"
	"cycles/types"
	"encoding/json"
	"math"
	"slices"

	lammps_structs "github.com/Ivanestver/lammps-file-parser/structs"
)

type Data struct {
	points               []*types.Point
	edges                []*types.Edge
	graph                data_structs.Graph
	spanningTreeEdges    []*types.Edge
	nonSpanningTreeEdges []*types.Edge
	supportVectors       []types.SupportVector
}

type Cycle []*types.Edge

func CalculateCycles(jsonObj string) ([]Cycle, error) {
	// 1. Initialization step
	data, err := initialize(jsonObj)
	if err != nil {
		return nil, err
	}

	// 2. Iteration step
	supportVectorSize := len(data.edges)
	supportVectors := getSupportVectors(data.nonSpanningTreeEdges, len(data.points), data.edges)
	cycles := make([]Cycle, len(supportVectors))
	shift := len(data.points)
	for k := 0; k < len(cycles); k++ {
		supportVector := supportVectors[k]
		doubledGraph := createDoubledGraph(data.graph, data.edges, supportVector)
		for pointNumber := range data.points {
			cycle := getCycle(doubledGraph, pointNumber, shift+pointNumber)
			if cycles[k] == nil || len(cycles[k]) > len(cycle) {
				cycles[k] = cycle
			}
		}
		for j := k + 1; j < len(supportVectors); j++ {
			cycleSupportVector := turnCycleIntoSupportVector(cycles[k], supportVectorSize)
			if testScalarMultiplication(cycleSupportVector, supportVectors[j]) {
				supportVectors[j] = supportVectors[j].XOR(supportVectors[k])
			}
		}
	}

	return cycles, nil
}

func initialize(jsonObj string) (*Data, error) {
	graphJson, err := parseJson(jsonObj)
	if err != nil {
		return nil, err
	}
	// 1. Get a random spanning tree
	dfs := algs.MakeDFS(graphJson.Points, graphJson.Graph)
	spanningTree := dfs.Traverse(0)
	// 2. Get all the edges that are not in the spanning tree
	nonSpanningTreeEdges := getNonSpanningTreeEdges(spanningTree, graphJson.Edges, graphJson.Graph)
	// 3. Get support vectors
	supportVectors := getSupportVectors(nonSpanningTreeEdges, len(graphJson.Points), graphJson.Edges)
	return &Data{graphJson.Points, graphJson.Edges, graphJson.Graph, spanningTree, nonSpanningTreeEdges, supportVectors}, nil
}

func MakeGraphSmall() *GraphJson {
	/*
		1 *-* 1-2
		  |/| |/|
		0 *-* 0-3
		  0 1
	*/
	points := []*types.Point{
		types.NewPoint(0, 0, 0),
		types.NewPoint(0, 1, 0),
		types.NewPoint(1, 1, 0),
		types.NewPoint(1, 0, 0),
	}
	edges := []*types.Edge{
		{Number: 0, Edge: [2]int{0, 1}},
		{Number: 1, Edge: [2]int{1, 2}},
		{Number: 2, Edge: [2]int{2, 3}},
		{Number: 3, Edge: [2]int{3, 0}},
		{Number: 4, Edge: [2]int{0, 2}},
	}

	graph := data_structs.Graph{
		[]*types.Edge{nil, edges[0], edges[4], edges[3]},
		[]*types.Edge{edges[0], nil, edges[1], nil},
		[]*types.Edge{edges[4], edges[1], nil, edges[2]},
		[]*types.Edge{edges[3], nil, edges[2], nil},
	}
	return NewGraphJson(points, edges, graph)
}

func parseJson(jsonObj string) (*GraphJson, error) {
	var v lammps_structs.LammpsStruct
	if err := json.Unmarshal([]byte(jsonObj), &v); err == nil {
		graphJson := NewGraphJson(make([]*types.Point, len(v.Atoms)), make([]*types.Edge, len(v.Bonds)), make(data_structs.Graph, len(v.Atoms)))
		for i, atom := range v.Atoms {
			graphJson.Points[i] = types.NewPoint(atom.X, atom.Y, atom.Z)
		}
		for i, bond := range v.Bonds {
			graphJson.Edges[i] = &types.Edge{
				Number: bond.BondID - 1,
				Edge:   [2]int{bond.Ends[0].AtomID - 1, bond.Ends[1].AtomID - 1},
			}
		}
		graphJson.Graph = makeGraph(graphJson.Edges, len(graphJson.Points))
		return graphJson, nil
	} else {
		return nil, err
	}
}

func makeGraph(edges []*types.Edge, pointCount int) data_structs.Graph {
	graph := make(data_structs.Graph, pointCount)
	for i := 0; i < len(graph); i++ {
		graph[i] = make([]*types.Edge, len(graph))
	}
	for i, edge := range edges {
		graph[edge.Edge[0]][edge.Edge[1]] = edges[i]
		graph[edge.Edge[1]][edge.Edge[0]] = edges[i]
	}
	return graph
}

func getNonSpanningTreeEdges(spanningTreeEdges types.Path, edges []*types.Edge, graph data_structs.Graph) []*types.Edge {
	nonSpanningTreeEdges := make([]*types.Edge, 0)
	for _, edge := range edges {
		if !slices.ContainsFunc(spanningTreeEdges, func(e *types.Edge) bool {
			return e.Equals(edge)
		}) {
			nonSpanningTreeEdges = append(nonSpanningTreeEdges, edge)
		}
	}
	return nonSpanningTreeEdges
}

// func splitSpanningTreeIntoEdges(spanningTree types.Path, graph data_structs.Graph) []*types.Edge {
// 	edges := make([]*types.Edge, len(spanningTree)-1)
// 	for i := 0; i < len(spanningTree)-1; i++ {
// 		edges[i] = graph[spanningTree[i]][spanningTree[i+1]]
// 	}
// 	return edges
// }

func getSupportVectors(nonSpanningTreeEdges []*types.Edge, pointsCount int, edges []*types.Edge) []types.SupportVector {
	edgesCount := len(edges)
	supportVectorsCount := edgesCount - (pointsCount - 1)
	supportVectors := make([]types.SupportVector, supportVectorsCount)
	i := 0
	for ; i < len(nonSpanningTreeEdges); i++ {
		supportVectors[i] = make(types.SupportVector, edgesCount)
		supportVectors[i][nonSpanningTreeEdges[i].Number] = 1
	}
	edgeI := len(edges) - 1
	for ; i < supportVectorsCount; i++ {
		supportVectors[i] = make(types.SupportVector, edgesCount)
		supportVectors[i][edges[edgeI].Number] = 1
		edgeI--
	}
	return supportVectors
}

func createDoubledGraph(originGraph data_structs.Graph, edges []*types.Edge, supportVector types.SupportVector) data_structs.Graph {
	originSize := len(originGraph)
	size := originSize * 2
	doubledGraph := make(data_structs.Graph, size)
	for i := 0; i < len(doubledGraph); i++ {
		doubledGraph[i] = make([]*types.Edge, size)
	}
	for _, edge := range edges {
		if supportVector[edge.Number] == 0 { // in the spanning tree{
			x := edge.Edge[0]
			y := edge.Edge[1]
			doubledGraph[x][y] = originGraph[x][y]
			doubledGraph[y][x] = originGraph[y][x]
			doubledGraph[originSize+x][originSize+y] = originGraph[x][y]
			doubledGraph[originSize+y][originSize+x] = originGraph[y][x]
		} else {
			x := edge.Edge[0]
			y := edge.Edge[1]
			doubledGraph[x][originSize+y] = originGraph[x][y]
			doubledGraph[originSize+y][x] = originGraph[y][x]
			doubledGraph[originSize+x][y] = originGraph[x][y]
			doubledGraph[y][originSize+x] = originGraph[y][x]
		}
	}
	return doubledGraph
}

func getCycle(doubledGraph data_structs.Graph, startingPoint, finishingPoint int) []*types.Edge {
	lengths := make([]float64, len(doubledGraph))
	prev := make([]int, len(doubledGraph))
	for i := range lengths {
		lengths[i] = math.MaxFloat64
		prev[i] = -1
	}
	lengths[startingPoint] = 0
	pq := &data_structs.PriorityQueue{data_structs.PQItem{Dist: 0, Num: startingPoint}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(data_structs.PQItem)
		for i, edge := range doubledGraph[it.Num] {
			if edge == nil {
				continue
			}
			if newDist := it.Dist + edge.Len(); newDist < lengths[i] {
				lengths[i] = newDist
				prev[i] = it.Num
				heap.Push(pq, data_structs.PQItem{
					Dist: newDist,
					Num:  i,
				})
			}
		}
	}
	curr := finishingPoint
	cycle := make([]*types.Edge, 0)
	for curr != startingPoint {
		prevPoint := prev[curr]
		cycle = append([]*types.Edge{doubledGraph[prevPoint][curr]}, cycle...)
		curr = prevPoint
	}
	return cycle
}

func turnCycleIntoSupportVector(cycle []*types.Edge, cycleSize int) types.SupportVector {
	cycleSupportVector := make(types.SupportVector, cycleSize)
	for i := range cycleSupportVector {
		if slices.ContainsFunc(cycle, func(edge *types.Edge) bool {
			return edge.Number == i
		}) {
			cycleSupportVector[i] = 1
		}
	}
	return cycleSupportVector
}

func testScalarMultiplication(cycleSupportVector, supportVector types.SupportVector) bool {
	return cycleSupportVector.GetScalarMultiplication(supportVector)%2 == 1
}
