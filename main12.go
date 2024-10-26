package main

import (
	"fmt"
	"math"
)

func main() {
	SqRoots()
}
func SqRoots() {
	var a, b, c float64
	fmt.Scanln(&a, &b, &c)
	D := b*b - 4*a*c
	if D > 0 {
		x1 := (-b + math.Sqrt(D)) / (2 * a)
		x2 := (-b - math.Sqrt(D)) / (2 * a)
		if x1 < x2 {
			fmt.Println(x1, x2)
		} else {
			fmt.Println(x2, x1)
		}
	} else if D == 0 {
		x1 := (-b + math.Sqrt(D)) / (2 * a)
		fmt.Println(x1)
	} else {
		fmt.Println("00")
	}
}
