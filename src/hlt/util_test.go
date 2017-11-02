package hlt

import (
	"fmt"
)


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