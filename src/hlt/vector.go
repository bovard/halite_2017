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
	return Vector{
		X: self.X + other.X,
		Y: self.Y + other.Y,
	}
}

func (self *Vector) Opposite() Vector {
	return Vector{
		X: -self.X,
		Y: -self.Y,
	}
}

func (self *Vector) RescaleToMag(mag int) Vector {
	scaler := float64(mag) / self.Magnitude()
	return Vector{
		X: scaler * self.X,
		Y: scaler * self.Y,
	}
}

func (self *Vector) RescaleToMagFloat(mag float64) Vector {
	scaler := mag / self.Magnitude()
	return Vector{
		X: scaler * self.X,
		Y: scaler * self.Y,
	}
}

func CreateVector(mag int, angle float64) Vector {
	return Vector{
		X: float64(mag) * math.Cos(angle),
		Y: float64(mag) * math.Sin(angle),
	}
}

func CreateRoundedVector(mag int, angle float64) Vector {
	h := CreateHeading(mag, angle)
	return h.ToVelocity()
}

//func (self *Vector) ToHeading() Heading {
//mag := math.Sqrt(self.X * self.X + self.Y * self.Y)
//ang := math.Atan2(self.Y / self.X)
//return CreateHeading(int(mag), ang)
//}

func (self *Vector) Subtract(other *Vector) Vector {
	return Vector{
		X: self.X - other.X,
		Y: self.Y - other.Y,
	}
}
