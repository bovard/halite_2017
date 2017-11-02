package hlt

import (
	"math"
)

type Entity struct {
	Point
	Radius   float64
	Health   float64
	Owner    int
	Id       int
	Distance float64
}

func (self *Entity) DistanceToCollision(target *Entity) float64 {
	// returns euclidean distance to target
	dx := target.Point.X - self.Point.X
	dy := target.Point.Y - self.Point.Y

	return math.Sqrt(dx*dx+dy*dy) - self.Radius - target.Radius
}

func (self *Entity) ClosestPointTo(target *Entity, minDistance float64) Point {
	// returns closest point to self that is at least minDistance from target
	dist := self.Point.DistanceTo(&target.Point) - target.Radius - minDistance
	angle := target.Point.AngleTo(&self.Point)
	return target.Point.AddThrust(dist, angle)
}
