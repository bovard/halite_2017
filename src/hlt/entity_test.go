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
		Id: 1,
	}
	p2 := Point{
		X: 5,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 2,
		Id: 2,
	}
	p3 := Point{
		X: 10,
		Y: 0,
	}
	e3 := Entity{
		Point:  p3,
		Radius: 2,
		Id: 3,
	}
	v0 := Vector {
		X: 0,
		Y: 0,
	}
	v1 := Vector {
		X: 5,
		Y: 0,
	}
	v2 := Vector {
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
		Id: 1,
	}
	p2 := Point{
		X: 4,
		Y: 1,
	}
	e2 := Entity{
		Point:  p2,
		Radius: 2,
		Id: 2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: 3,
		Id: 3,
	}
	v1 := Vector {
		X: 4,
		Y: 0,
	}
	v2 := Vector {
		X: 3,
		Y: 0,
	}
	v3 := Vector {
		X: 1,
		Y: 0,
	}
	v4 := Vector {
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
		Id: 1,
	}
	p2 := Point{
		X: 4,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id: 2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: .5,
		Id: 3,
	}
	v1 := Vector {
		X: 4,
		Y: 0,
	}
	v2 := Vector {
		X: 3,
		Y: 0,
	}
	v3 := Vector {
		X: 1,
		Y: 0,
	}
	v4 := Vector {
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
		Id: 1,
	}
	p2 := Point{
		X: 4,
		Y: 0,
	}
	e2 := Entity{
		Point:  p2,
		Radius: .5,
		Id: 2,
	}
	p3 := Point{
		X: 3,
		Y: 3,
	}
	e3 := Entity{
		Point:  p3,
		Radius: .5,
		Id: 3,
	}
	v1 := Vector {
		X: 4,
		Y: 0,
	}
	v1o := v1.Opposite()
	v2 := Vector {
		X: 3,
		Y: 0,
	}
	v1mv2 := v1.Subtract(&v2)
	v3 := Vector {
		X: 1,
		Y: 0,
	}
	v2mv3 := v2.Subtract(&v3)
	v3mv1o := v3.Subtract(&v1o)
	v4 := Vector {
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