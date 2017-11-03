package hlt

import (
	"math"
)

type Point struct {
	X, Y float64
}

func (self *Point) RotateAround(target *Point, angle float64) Point {
	x1 := self.X - target.X
	y1 := self.Y - target.Y
	x2 := x1*math.Cos(angle) - y1*math.Sin(angle)
	y2 := x1*math.Sin(angle) - y1*math.Cos(angle)
	return Point{
		X: x2 + target.X,
		Y: y2 + target.Y,
	}
}

func (self *Point) AddThrust(magnitude float64, angle float64) Point {
	return Point{
		X: self.X + magnitude*math.Cos(angle),
		Y: self.Y + magnitude*math.Sin(angle),
	}
}

func (self *Point) AddVector(v *Vector) Point {
	return Point {
		X: self.X + v.X,
		Y: self.Y + v.Y,
	}
}

func (self *Point) GetMidPoint(target *Point) Point {
	return Point{
		X: (self.X + target.X) / 2,
		Y: (self.Y + target.Y) / 2,
	}
}

func (self *Point) DistanceTo(target *Point) float64 {
	// returns euclidean distance to target
	dx := target.X - self.X
	dy := target.Y - self.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func (self *Point) SqDistanceTo(target *Point) float64 {
	// returns euclidean distance to target
	dx := target.X - self.X
	dy := target.Y - self.Y
	return dx*dx + dy*dy
}

func (self *Point) AngleTo(target *Point) float64 {
	// returns angle in radians from self to target
	dx := target.X - self.X
	dy := target.Y - self.Y

	return math.Atan2(dy, dx)
}

func (self *Point) VectorTo(other *Point) Vector {
	return Vector{
		X: other.X - self.X,
		Y: other.Y - self.Y,
	}
}

func GetClosestPointOnLine(v1 *Point, v2 *Point, p *Point) Point {
	e1 := v1.VectorTo(v2)
	e2 := v1.VectorTo(p)
	valDp := e1.Dot(&e2)
	len2 := e1.SqMagnitude()
	return Point{
		X: v1.X + (valDp*e1.X)/len2,
		Y: v1.Y + (valDp*e1.Y)/len2,
	}
}
