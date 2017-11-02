package hlt

import "math"

func DegToRad(d float64) float64{
	return d / 180 * math.Pi
}

func RadToDeg(r float64) float64{
	return r / math.Pi * 180
}

func GetClosestPointOnLine(v1 *Position, v2 *Position, p *Position) Position {
	e1 := Position {
		X: v2.X - v1.X,
		Y: v2.Y - v1.Y,
	}
	e2 := Position {
		X: p.X - v1.X,
		Y: p.Y - v1.Y,
	}
	valDp := e1.DotProduct(&e2)
	len2 := e1.Magnitude()
	return Position {
		X: v1.X + (valDp * e1.X) / len2,
		Y: v1.Y + (valDp * e1.Y) / len2,
	}
}