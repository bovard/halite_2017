package hlt

import ("math")

type Entity struct {
	Point
	Radius   float64
	Health   float64
	Owner    int
	Id       int
	Distance float64
}

func (self *Entity) DistanceToCollision(target *Entity) float64 {
	return self.Point.DistanceTo(&target.Point) - self.Radius - target.Radius
}

func (self *Entity) ClosestPointTo(target *Entity, minDistance float64) Point {
	// returns closest point to self that is at least minDistance from target
	dist := self.Point.DistanceTo(&target.Point) - target.Radius - minDistance
	angle := target.Point.AngleTo(&self.Point)
	return target.Point.AddThrust(dist, angle)
}

func (self *Entity) WillCollideWith(target *Entity, vel *Vector) bool {
	mag := vel.Magnitude()
	if mag == 0 {
		return false
	}
	// if the object is too far away, return false
	if self.DistanceToCollision(target) > mag {
		return false
	}
	nextP := self.Point.AddVector(vel)
	projectedP := GetClosestPointOnLine(&self.Point, &nextP, &target.Point)
	// if the object isn't in the right direction, return false
	if math.Abs(self.AngleTo(&nextP) - self.AngleTo(&projectedP)) > .1 {
		return false
	} 	
	return projectedP.DistanceTo(&target.Point) - self.Radius - target.Radius <= 0
}