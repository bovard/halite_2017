package hlt

import "math"

func DegToRad(d float64) float64{
	return d / 180 * math.Pi
}

func RadToDeg(r float64) float64{
	return r / math.Pi * 180
}

func GetClosestPointOnLine(v1 *Point, v2 *Point, p *Point) Point {
	e1 := v1.VectorTo(v2)
	e2 := v1.VectorTo(p)
	valDp := e1.Dot(&e2)
	len2 := e1.SqMagnitude()
	return Point {
		X: v1.X + (valDp * e1.X) / len2,
		Y: v1.Y + (valDp * e1.Y) / len2,
	}
}