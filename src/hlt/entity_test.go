package hlt

import (
	"fmt"
)

func ExampleClosestPointTo() {
	firstP := Point{
		X: 10,
		Y: 10,
	}
	first := Entity{
		Point:  firstP,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	secondP := Point{
		X: 13,
		Y: 10,
	}
	second := Entity{
		Point:  secondP,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	result := first.ClosestPointTo(&second, 1)
	fmt.Println(result.X, result.Y)
	result = second.ClosestPointTo(&first, .5)
	fmt.Println(result.X, result.Y)
	// Output:
	// 11 10
	// 12.5 10
}

func ExampleLinearWillCollideWith() {
	p1 := Point{
		X: 0,
		Y: 0,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 5,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 2,
		Id:     2,
	}
	p3 := Point{
		X: 10,
		Y: 0,
	}
	e3 := Entity{
		Point:  p3,
		Radius: 2,
		Id:     3,
	}
	v0 := Vector{
		X: 0,
		Y: 0,
	}
	v1 := Vector{
		X: 5,
		Y: 0,
	}
	v2 := Vector{
		X: 3.5,
		Y: 0,
	}
	fmt.Println(e1.WillCollideWith(&e2, &v0))
	fmt.Println(e1.WillCollideWith(&e3, &v1))
	fmt.Println(e1.WillCollideWith(&e2, &v1))
	fmt.Println(e1.WillCollideWith(&e2, &v2))
	// Output:
	// false
	// false
	// true
	// true
}

func ExamplePlanerWillCollideWith() {
	p1 := Point{
		X: 0,
		Y: 0,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 4,
		Y: 1,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 2,
		Id:     2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: 3,
		Id:     3,
	}
	v1 := Vector{
		X: 4,
		Y: 0,
	}
	v2 := Vector{
		X: 3,
		Y: 0,
	}
	v3 := Vector{
		X: 1,
		Y: 0,
	}
	v4 := Vector{
		X: -1,
		Y: -1,
	}
	fmt.Println(e1.WillCollideWith(&e2, &v1))
	fmt.Println(e1.WillCollideWith(&e2, &v2))
	fmt.Println(e1.WillCollideWith(&e2, &v3))
	fmt.Println(e1.WillCollideWith(&e2, &v4))
	fmt.Println(e1.WillCollideWith(&e3, &v1))
	fmt.Println(e1.WillCollideWith(&e3, &v2))
	fmt.Println(e1.WillCollideWith(&e3, &v3))
	fmt.Println(e1.WillCollideWith(&e3, &v4))
	// Output:
	// true
	// true
	// false
	// false
	// true
	// true
	// true
	// false
}

func ExampleShipsStaticWillCollideWith() {
	p1 := Point{
		X: 0,
		Y: 0,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 4,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id:     2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: .5,
		Id:     3,
	}
	v1 := Vector{
		X: 4,
		Y: 0,
	}
	v2 := Vector{
		X: 3,
		Y: 0,
	}
	v3 := Vector{
		X: 1,
		Y: 0,
	}
	v4 := Vector{
		X: -4,
		Y: -4,
	}
	fmt.Println(e1.WillCollideWith(&e2, &v1))
	fmt.Println(e1.WillCollideWith(&e2, &v2))
	fmt.Println(e1.WillCollideWith(&e2, &v3))
	fmt.Println(e1.WillCollideWith(&e2, &v4))
	fmt.Println(e1.WillCollideWith(&e3, &v1))
	fmt.Println(e1.WillCollideWith(&e3, &v2))
	fmt.Println(e1.WillCollideWith(&e3, &v3))
	fmt.Println(e1.WillCollideWith(&e3, &v4))
	// Output:
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
}

func ExampleShipsDyanmicWillCollideWith() {
	p1 := Point{
		X: 0,
		Y: 0,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 4,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id:     2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: .5,
		Id:     3,
	}
	v1 := Vector{
		X: 4,
		Y: 0,
	}
	v1o := v1.Opposite()
	v2 := Vector{
		X: 3,
		Y: 0,
	}
	v1mv2 := v1.Subtract(&v2)
	v3 := Vector{
		X: 1,
		Y: 0,
	}
	v2mv3 := v2.Subtract(&v3)
	v3mv1o := v3.Subtract(&v1o)
	v4 := Vector{
		X: -4,
		Y: -4,
	}
	v4o := v4.Opposite()
	v4omv1 := v4o.Subtract(&v1)

	fmt.Println(e1.WillCollideWith(&e2, &v1mv2))
	fmt.Println(e1.WillCollideWith(&e2, &v2mv3))
	fmt.Println(e1.WillCollideWith(&e2, &v3mv1o))
	fmt.Println(e1.WillCollideWith(&e2, &v4))
	fmt.Println(e1.WillCollideWith(&e3, &v4))
	fmt.Println(e1.WillCollideWith(&e3, &v4o))
	fmt.Println(e1.WillCollideWith(&e3, &v4omv1))
	fmt.Println(e1.WillCollideWith(&e3, &v1))
	// Output:
	// false
	// false
	// true
	// false
	// false
	// true
	// false
	// false
}

func ExampleShipsFromGame45158615WillCollideWith() {
	p1 := Point{
		X: 71.43,
		Y: 96.68,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 72.07,
		Y: 103.23,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id:     2,
	}
	h1 := Heading{
		Magnitude: 7,
		Angle:     93.0,
	}
	v1 := h1.ToVelocity()
	v2 := Vector{
		X: 0,
		Y: 0,
	}
	v1mv2 := v1.Subtract(&v2)

	fmt.Println(e1.WillCollideWith(&e2, &v1mv2))
	// Output:
	// true
}

func ExampleShipsFromGame3071260526WillCollideWith() {
	p1 := Point{
		X: 31.29,
		Y: 127.96,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 30.63,
		Y: 130.40,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id:     2,
	}
	h1 := Heading{
		Magnitude: 4,
		Angle:     128.0,
	}
	v1 := h1.ToVelocity()
	v2 := Vector{
		X: 0,
		Y: 0,
	}
	v1mv2 := v1.Subtract(&v2)

	fmt.Println(e1.WillCollideWith(&e2, &v1mv2))
	// Output:
	// true
}

func ExampleShipsFromGame1756470586WillCollideWith() {
	p1 := Point{
		X: 113.23,
		Y: 62.95,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 111.76,
		Y: 71.76,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 5.33,
		Id:     2,
	}
	h1 := Heading{
		Magnitude: 7,
		Angle:     59.0,
	}
	v1 := h1.ToVelocity()

	fmt.Println(e1.WillCollideWith(&e2, &v1))
	// Output:
	// true
}

func ExampleShipsFromGame2370134WillCollideWith() {
	p1 := Point{
		X: 184.64,
		Y: 63.06,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	p2 := Point{
		X: 183.99,
		Y: 50.34,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 11.39,
		Id:     2,
	}
	h1 := Heading{
		Magnitude: 7,
		Angle:     336.0,
	}
	v1 := h1.ToVelocity()

	fmt.Println(e1.WillCollideWith(&e2, &v1))
	// Output:
	// true
}

func ExampleShipsFromGame154201WillCollideWith() {
	// seed: 961841357
	p1 := Point{
		X: 133.9925,
		Y: 74.7725,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	h1 := Heading{
		Magnitude: 7,
		Angle:     148.0,
	}
	v1 := h1.ToVelocity()

	p2 := Point{
		X: 127.6954,
		Y: 72.3046,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 4.9415,
		Id:     2,
	}

	fmt.Println(e1.WillCollideWith(&e2, &v1))
	// Output:
	// true
}

func ExampleShipsFromGame2440091WillCollideWith() {
	// replay: 2440091
	p1 := Point{
		X: 78.98,
		Y: 93.03,
	}
	e1 := Entity{
		Point:  p1,
		Radius: 0.5,
		Id:     1,
	}
	h1 := Heading{
		Magnitude: 7,
		Angle:     263.0,
	}
	v1 := h1.ToVelocity()

	p2 := Point{
		X: 79.22,
		Y: 87.05,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id:     2,
	}

	fmt.Println(e1.WillCollideWith(&e2, &v1))
	// Output:
	// true
}
