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

func ExampleCalculateAngleTo() {
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
	fmt.Println(first.CalculateAngleTo(&second))
	fmt.Println(RadToDeg(first.CalculateAngleTo(&third)))
	fmt.Println(RadToDeg(third.CalculateAngleTo(&second)))
	fmt.Println(RadToDeg(second.CalculateAngleTo(&third)))
	// Output:
	// 0
	// 90
	// -45
	// 135
}


