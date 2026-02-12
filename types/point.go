package types

import "math"

type Point struct {
	PointID int
	X, Y, Z float64
	State   State
}

func (point *Point) GetDistanceTo(other *Point) float64 {
	return math.Sqrt(
		(point.X-other.X)*(point.X*other.X) +
			(point.Y-other.Y)*(point.Y*other.Y) +
			(point.Z-other.Z)*(point.Z*other.Z))
}

func NewPoint(pointID int, X, Y, Z float64) *Point {
	return &Point{
		PointID: pointID,
		X:       X,
		Y:       Y,
		Z:       Z,
		State:   STATE_WHITE,
	}
}
