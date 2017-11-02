package hlt

import (
	"fmt"
	"math"
)


func ExampleRotateAround() {
	center := Point {
		X: 10,
		Y: 10,
	}
	point := Point {
		X: 12,
		Y: 10,
	}
	result := point.RotateAround(&center, math.Pi)
	fmt.Println(result.X, result.Y)
	result = point.RotateAround(&center, 0)
	fmt.Println(result.X, result.Y)
	result = point.RotateAround(&center, math.Pi/2)
	fmt.Println(result.X, result.Y)
	// Output:
	// 8 10
	// 12 10
	// 10 12
}

func ExampleAddThrust() {
	center := Point {
		X: 10,
		Y: 10,
	}
	result := center.AddThrust(2.0, 0)
	fmt.Println(result.X, result.Y)
	result = center.AddThrust(2.0, math.Pi/2)
	fmt.Println(result.X, result.Y)
	result = center.AddThrust(2.0, math.Pi)
	fmt.Println(result.X, result.Y)
	// Output:
	// 12 10
	// 10 12
	// 8 10
}

func ExampleGetMidPoint() {
	first := Point {
		X: 10,
		Y: 8,
	}
	second := Point {
		X:  12,
		Y:  10,
	}
	result := first.GetMidPoint(&second)
	fmt.Println(result.X, result.Y)
	result = second.GetMidPoint(&first)
	fmt.Println(result.X, result.Y)
	// Output:
	// 11 9
	// 11 9
}


func ExampleDistanceTo() {
	first := Point {
		X: 10,
		Y: 10,
	}
	second := Point {
		X: 13,
		Y: 10,
	}
	third := Point {
		X: 10,
		Y: 14,
	}
	fmt.Println(first.DistanceTo(&second))
	fmt.Println(first.DistanceTo(&third))
	fmt.Println(second.DistanceTo(&third))
	// Output:
	// 3
	// 4
	// 5
}

func ExampleAngleTo() {
	first := Point {
		X: 10,
		Y: 10,
	}
	second := Point {
		X: 13,
		Y: 10,
	}
	third := Point {
		X: 10,
		Y: 13,
	}
	fmt.Println(first.AngleTo(&second))
	fmt.Println(RadToDeg(first.AngleTo(&third)))
	fmt.Println(RadToDeg(third.AngleTo(&second)))
	fmt.Println(RadToDeg(second.AngleTo(&third)))
	// Output:
	// 0
	// 90
	// -45
	// 135
}


func ExampleGetClosestPointOnLine() {
	v1 := Point {
		X: 0.0,
		Y: 0.0,
	}
	v2 := Point {
		X: 0.0,
		Y: 10.0,
	}
	v3 := Point {
		X: 10.0,
		Y: 0.0,
	}
	v4 := Point {
		X: 10.0,
		Y: 10.0,
	}
	p := Point {
		X: 5.0,
		Y: 5.0,
	}
		
	result := GetClosestPointOnLine(&v1, &v2, &p)
	fmt.Println(result.X, result.Y)
	result = GetClosestPointOnLine(&v1, &v3, &p)
	fmt.Println(result.X, result.Y)
	result = GetClosestPointOnLine(&v1, &v4, &p)
	fmt.Println(result.X, result.Y)
	// Output:
	// 0 5
	// 5 0 
	// 5 5
}


