package cycles_alg

import (
	"cycles/algs"
	"cycles/data_structs"
	"cycles/types"
	"fmt"
	"slices"
	"testing"
)

func TestEdgeEquals(t *testing.T) {
	edge1 := types.Edge{
		Number: 0,
		Edge:   [2]int{1, 2},
	}
	edge2 := types.Edge{
		Number: 0,
		Edge:   [2]int{1, 2},
	}
	if !edge1.Equals(&edge2) {
		t.Errorf("Equal edges are not equal. Edge1 - %v, Edge2 - %v", edge1, edge2)
	}

	edge3 := types.Edge{
		Number: 0,
		Edge:   [2]int{2, 1},
	}
	if !edge1.Equals(&edge3) {
		t.Errorf("Equal edges are not equal. Edge1 - %v, Edge2 - %v", edge1, edge3)
	}

	edge4 := types.Edge{
		Number: 1,
		Edge:   [2]int{1, 2},
	}
	if edge1.Equals(&edge4) {
		t.Errorf("Equal edges are not equal. Edge1 - %v, Edge2 - %v", edge1, edge4)
	}

	edge5 := types.Edge{
		Number: 0,
		Edge:   [2]int{3, 1},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("No panic when it must be. Edge1 - %v, Edge2 - %v", edge1, edge5)
		}
	}()
	edge1.Equals(&edge5)
}

func TestDFS(t *testing.T) {
	graphJson := makeTestGraph()
	dfs := algs.MakeDFS(graphJson.Points, graphJson.Graph)
	real := dfs.Traverse(0)
	expected := getTestSpanningTree()
	equal := slices.CompareFunc(real, expected, func(p1, p2 *types.Edge) int {
		if p1.Equals(p2) {
			return 0
		} else {
			return -1
		}
	})
	if equal != 0 {
		exp := make([]types.Edge, len(expected))
		r := make([]types.Edge, len(real))
		for i := 0; i < len(real); i++ {
			exp[i] = *expected[i]
			r[i] = *real[i]
		}
		t.Errorf("expected: %v, got: %v", exp, r)
	}
}

func TestGetNonSpanningTreeEdges(t *testing.T) {
	real := getRealNonSpanningTreeEdges()
	expected := getTestNonSpanningTreeEdges()
	equal := slices.CompareFunc(expected, real, func(e1, e2 *types.Edge) int {
		if e1.Equals(e2) {
			return 0
		} else {
			return 1
		}
	})
	if equal != 0 {
		exp := make([]types.Edge, len(expected))
		r := make([]types.Edge, len(real))
		for i := 0; i < len(real); i++ {
			exp[i] = *expected[i]
			r[i] = *real[i]
		}
		t.Errorf("expected: %v, got: %v", exp, r)
	}
}

func getRealNonSpanningTreeEdges() []*types.Edge {
	spanningTree := getTestSpanningTree()
	graphJson := makeTestGraph()
	return getNonSpanningTreeEdges(spanningTree, graphJson.Edges)
}

func TestGetSupportVectors(t *testing.T) {
	graphJson := makeTestGraph()
	expected := getTestSupportVectors()
	nonSpanningTreeEdges := getRealNonSpanningTreeEdges()
	real := getSupportVectors(nonSpanningTreeEdges, len(graphJson.Points), graphJson.Edges)
	res := slices.CompareFunc(expected, real, func(e1, e2 types.SupportVector) int {
		return slices.Compare(e1, e2)
	})
	if res != 0 {
		t.Errorf("Expected: %v, got: %v", expected, real)
	}
}

func TestCreateDoubledGraph(t *testing.T) {
	graphJson := makeTestGraph()
	expected := getTestDoubledGraph()
	testSupportVector := types.SupportVector{0, 0, 0, 0, 1}
	real := createDoubledGraph(graphJson.Graph, graphJson.Edges, testSupportVector)
	res := slices.CompareFunc(expected, real, func(l1, l2 []*types.Edge) int {
		return slices.CompareFunc(l1, l2, func(e1, e2 *types.Edge) int {
			if (e1 == nil && e2 == nil) || e1.Equals(e2) {
				return 0
			} else {
				return 1
			}
		})
	})
	if res != 0 {
		exp := expected.GetOnlyNumbers()
		r := real.GetOnlyNumbers()
		t.Errorf("Wrong doubled graph. Expected: %v, got: %v", exp, r)
	}
}

func TestGetFirstCycle(t *testing.T) {
	expected := getTestCycles()
	doubledGraph := getTestDoubledGraph()
	real := getCycle(doubledGraph, 0, len(doubledGraph)/2)
	for _, expectedCycle := range expected {
		res := slices.CompareFunc(expectedCycle, real, func(e1, e2 *types.Edge) int {
			if e1.Equals(e2) {
				return 0
			} else {
				return 1
			}
		})
		if res == 0 {
			return
		}
	}
	message := "Expected one of the following:"
	for _, expectedCycle := range expected {
		exp := make([]types.Edge, len(expectedCycle))
		for i := range exp {
			exp[i] = *expectedCycle[i]
		}
		message = fmt.Sprintf("%s %v, ", message, expectedCycle)
	}
	r := make([]types.Edge, len(real))
	for i := range real {
		r[i] = *real[i]
	}
	t.Errorf("%s, got: %v", message, r)
}

func TestTestScalarMultiplication(t *testing.T) {
	cycle := []*types.Edge{
		{Number: 0, Edge: [2]int{0, 1}},
		{Number: 1, Edge: [2]int{1, 2}},
		{Number: 4, Edge: [2]int{0, 2}},
	}
	supportVector := types.SupportVector{0, 0, 0, 0, 1}
	expected := true
	real := testScalarMultiplication(turnCycleIntoSupportVector(cycle, len(supportVector)), supportVector)
	if expected != real {
		t.Error("Expected: true, got: false")
	}
}

func makeTestGraph() *GraphJson {
	/*
		1 *-* 1-2
		  |/| |/|
		0 *-* 0-3
		  0 1
	*/
	points := []*types.Point{
		types.NewPoint(0, 0, 0, 0),
		types.NewPoint(1, 0, 1, 0),
		types.NewPoint(2, 1, 1, 0),
		types.NewPoint(3, 1, 0, 0),
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

func getTestSpanningTree() types.Path {
	return []*types.Edge{
		{Number: 0, Edge: [2]int{0, 1}}, {Number: 1, Edge: [2]int{1, 2}}, {Number: 2, Edge: [2]int{2, 3}},
	}

}

func getTestNonSpanningTreeEdges() []*types.Edge {
	return []*types.Edge{
		{Number: 3, Edge: [2]int{3, 0}},
		{Number: 4, Edge: [2]int{0, 2}},
	}
}

func getTestSupportVectors() []types.SupportVector {
	return []types.SupportVector{
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1},
	}
}

func getTestDoubledGraph() data_structs.Graph {
	graphJson := makeTestGraph()
	originSize := len(graphJson.Graph)
	doubledSize := originSize * 2
	doubledGraph := make(data_structs.Graph, doubledSize)
	for i := 0; i < len(doubledGraph); i++ {
		doubledGraph[i] = make([]*types.Edge, doubledSize)
	}
	doubledGraph[0][1] = graphJson.Graph[0][1]
	doubledGraph[1][0] = graphJson.Graph[1][0]
	doubledGraph[originSize+0][originSize+1] = graphJson.Graph[0][1]
	doubledGraph[originSize+1][originSize+0] = graphJson.Graph[1][0]

	doubledGraph[1][2] = graphJson.Graph[1][2]
	doubledGraph[2][1] = graphJson.Graph[2][1]
	doubledGraph[originSize+1][originSize+2] = graphJson.Graph[1][2]
	doubledGraph[originSize+2][originSize+1] = graphJson.Graph[2][1]

	doubledGraph[2][3] = graphJson.Graph[2][3]
	doubledGraph[3][2] = graphJson.Graph[3][2]
	doubledGraph[originSize+2][originSize+3] = graphJson.Graph[2][3]
	doubledGraph[originSize+3][originSize+2] = graphJson.Graph[3][2]

	doubledGraph[3][0] = graphJson.Graph[3][0]
	doubledGraph[0][3] = graphJson.Graph[0][3]
	doubledGraph[originSize+3][originSize+0] = graphJson.Graph[3][0]
	doubledGraph[originSize+0][originSize+3] = graphJson.Graph[0][3]

	doubledGraph[0][originSize+2] = graphJson.Graph[0][2]
	doubledGraph[originSize+2][0] = graphJson.Graph[2][0]
	doubledGraph[originSize+0][2] = graphJson.Graph[0][2]
	doubledGraph[2][originSize+0] = graphJson.Graph[2][0]
	return doubledGraph
}

func getTestCycles() [][]*types.Edge {
	return [][]*types.Edge{
		[]*types.Edge{
			{Number: 0, Edge: [2]int{0, 1}},
			{Number: 1, Edge: [2]int{1, 2}},
			{Number: 4, Edge: [2]int{0, 2}},
		},
		[]*types.Edge{
			{Number: 4, Edge: [2]int{0, 2}},
			{Number: 2, Edge: [2]int{2, 3}},
			{Number: 3, Edge: [2]int{3, 0}},
		},
	}
}

func getTestSecondCycle() []*types.Edge {
	return []*types.Edge{
		{Number: 3, Edge: [2]int{3, 0}},
		{Number: 4, Edge: [2]int{0, 2}},
		{Number: 2, Edge: [2]int{2, 3}},
	}
}

func makeGraphSmall() *GraphJson {
	/*
		1 *-* 1-2
		  |/| |/|
		0 *-* 0-3
		  0 1
	*/
	points := []*types.Point{
		types.NewPoint(0, 0, 0, 0),
		types.NewPoint(1, 0, 1, 0),
		types.NewPoint(2, 1, 1, 0),
		types.NewPoint(3, 1, 0, 0),
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

func makeGraphBig() ([]*types.Point, []*types.Edge, data_structs.Graph) {
	/*
		2 *-*-* 6-7-8
		| | | | | | |
		1 *-*-* 3-4-5
		| | | | | | |
		0 *-*-* 0-1-2
		  0-1-2
	*/
	points := []*types.Point{
		types.NewPoint(0, 0, 0, 0),
		types.NewPoint(1, 1, 0, 0),
		types.NewPoint(2, 2, 0, 0),
		types.NewPoint(3, 0, 1, 0),
		types.NewPoint(4, 1, 1, 0),
		types.NewPoint(5, 2, 1, 0),
		types.NewPoint(6, 0, 2, 0),
		types.NewPoint(7, 1, 2, 0),
		types.NewPoint(8, 2, 2, 0),
	}

	edges := []*types.Edge{
		{Number: 0, Edge: [2]int{0, 1}},
		{Number: 1, Edge: [2]int{0, 3}},
		{Number: 2, Edge: [2]int{1, 2}},
		{Number: 3, Edge: [2]int{1, 4}},
		{Number: 4, Edge: [2]int{2, 5}},
		{Number: 5, Edge: [2]int{3, 4}},
		{Number: 6, Edge: [2]int{3, 6}},
		{Number: 7, Edge: [2]int{4, 5}},
		{Number: 8, Edge: [2]int{4, 7}},
		{Number: 9, Edge: [2]int{5, 8}},
		{Number: 10, Edge: [2]int{6, 7}},
		{Number: 11, Edge: [2]int{7, 8}},
	}

	graph := make(data_structs.Graph, len(points))
	for i := 0; i < len(graph); i++ {
		graph[i] = make([]*types.Edge, len(graph))
	}
	for i, edge := range edges {
		graph[edge.Edge[0]][edge.Edge[1]] = edges[i]
		graph[edge.Edge[1]][edge.Edge[0]] = edges[i]
	}

	return points, edges, graph
}
