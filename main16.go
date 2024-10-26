package main

import (
	"fmt"
)

func main() {
	input := [5]int{1, 2, 3, 4, 5}
	output := FiveSteps(input)
	for i := 0; i < len(output); i++ {
		fmt.Printf("%d ", output[i])
	}
}

func FiveSteps(input [5]int) []int {
	reversed := make([]int, len(input))
	for i := 0; i < len(input); i = i + 1 {
		reversed[i] = input[len(input)-1-i]
	}
	return reversed
}
