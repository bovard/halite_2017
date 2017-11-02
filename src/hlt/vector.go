package hlt

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func (self *Vector) Dot(other *Vector) float64 {
	return self.X*other.X + self.Y*other.Y
}

func (self *Vector) Magnitude() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}

func (self *Vector) SqMagnitude() float64 {
	return self.X*self.X + self.Y*self.Y
}
