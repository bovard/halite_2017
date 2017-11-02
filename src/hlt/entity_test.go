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
