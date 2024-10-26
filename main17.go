package main

import (
	"fmt"
)

func main() {
	input := [6]int{1, 2, 3, 4, 5, 6}
	output := ThirdElementInArray(input)
	fmt.Print(output)
}

func ThirdElementInArray(input [6]int) int {
	return input[2]
}
