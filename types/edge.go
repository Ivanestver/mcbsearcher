package types

type Edge struct {
	Number int
	Edge   [2]int
	State  State
}

func (edge *Edge) Equals(other *Edge) bool {
	if other == nil || edge.Number != other.Number {
		return false
	}
	if (edge.Edge[0] == other.Edge[0] &&
		edge.Edge[1] == other.Edge[1]) ||
		(edge.Edge[0] == other.Edge[1] &&
			edge.Edge[1] == other.Edge[0]) {
		return true
	} else {
		panic("Two edges with equal numbers but different points")
	}
}

func (e *Edge) Len() float64 {
	return 1.0
}

func (e *Edge) GetOtherSide(oneSide int) int {
	if e.Edge[0] == oneSide {
		return e.Edge[1]
	} else if e.Edge[1] == oneSide {
		return e.Edge[0]
	} else {
		return -1
	}
}
