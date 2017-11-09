package hlt

import (
	"fmt"
	"math"
	"strconv"
)

type Heading struct {
	Magnitude int
	Angle     int
}

func (self *Heading) ToMoveCmd(ship *Ship, message int) string {
	angle := ((message + 1) * 360) + self.Angle
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.Id), strconv.Itoa(self.Magnitude), strconv.Itoa(angle))
}

func (self *Heading) ToVelocity() Vector {
	angle := DegToRad(float64(self.Angle))
	return Vector{
		X: float64(self.Magnitude) * math.Cos(angle),
		Y: float64(self.Magnitude) * math.Sin(angle),
	}
}

func (self *Heading) GetAngleInRads() float64 {
	return DegToRad(float64(self.Angle))
}

func CreateHeading(magnitude int, angle float64) Heading {
	var boundedAngle int
	angle = RadToDeg(angle)
	if angle > 0.0 {
		boundedAngle = int(math.Floor(angle + .5))
	} else {
		boundedAngle = int(math.Ceil(angle - .5))
	}
	boundedAngle = ((boundedAngle % 360) + 360) % 360
	return Heading{
		Magnitude: magnitude,
		Angle:     boundedAngle,
	}
}
