package hlt

import (
	"fmt"
	"math"
)


func ExampleRotateAround() {
	center := Entity{
		X: 10,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	point := Entity{
		X: 12,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	result := point.RotateAround(center, math.Pi)
	fmt.Println(result.X, result.Y)
	result = point.RotateAround(center, 0)
	fmt.Println(result.X, result.Y)
	result = point.RotateAround(center, math.Pi/2)
	fmt.Println(result.X, result.Y)
	// Output:
	// 8 10
	// 12 10
	// 10 12
}

func ExampleAddThrust() {
	center := Entity{
		X: 10,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
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
	first := Entity{
		X: 10,
		Y: 8,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	second := Entity{
		X: 12,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	result := first.GetMidPoint(second)
	fmt.Println(result.X, result.Y)
	result = second.GetMidPoint(first)
	fmt.Println(result.X, result.Y)
	// Output:
	// 11 9
	// 11 9
}


func ExampleCalculateDistanceTo() {
	first := Entity{
		X: 10,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	second := Entity{
		X: 13,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	third := Entity{
		X: 10,
		Y: 14,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	fmt.Println(first.CalculateDistanceTo(second))
	fmt.Println(first.CalculateDistanceTo(third))
	fmt.Println(second.CalculateDistanceTo(third))
	// Output:
	// 3
	// 4
	// 5
}

func ExampleCalculateAngleTo() {
	first := Entity{
		X: 10,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	second := Entity{
		X: 13,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	third := Entity{
		X: 10,
		Y: 13,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	fmt.Println(first.CalculateAngleTo(second))
	fmt.Println(RadToDeg(first.CalculateAngleTo(third)))
	fmt.Println(RadToDeg(third.CalculateAngleTo(second)))
	fmt.Println(RadToDeg(second.CalculateAngleTo(third)))
	// Output:
	// 0
	// 90
	// -45
	// 135
}

func ExampleClosestPointTo() {
	first := Entity{
		X: 10,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	second := Entity{
		X: 13,
		Y: 10,
		Radius: 0,
		Health: 0,
		Owner:  -1,
		Id:     -1,
	}
	result := first.ClosestPointTo(second, 1)
	fmt.Println(result.X, result.Y)
	result = second.ClosestPointTo(first, .5)
	fmt.Println(result.X, result.Y)
	// Output:
	// 11 10
	// 12.5 10
}
