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

func (self *Vector) Add(other *Vector) Vector {
	return Vector {
		X: self.X + other.X,
		Y: self.Y + other.Y,
	}
}

func (self *Vector) Opposite() Vector {
	return Vector {
		X: -self.X,
		Y: -self.Y,
	}
}

func CreateVector(mag int, angle float64) Vector {
	return Vector{
		X: float64(mag) * math.Cos(angle),
		Y: float64(mag) * math.Sin(angle),
	}
}

//func (self *Vector) ToHeading() Heading {
	//mag := math.Sqrt(self.X * self.X + self.Y * self.Y)
	//ang := math.Atan2(self.Y / self.X)
	//return CreateHeading(int(mag), ang)
//}

func (self *Vector) Subtract(other *Vector) Vector {
	return Vector {
		X: self.X - other.X,
		Y: self.Y - other.Y,
	}
}