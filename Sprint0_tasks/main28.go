package main

import (
	"fmt"
	"strings"
)

func main() {
	input := "А роза упала на лапу Азора"
	joined := strings.ReplaceAll(input, " ", "")
	lower := strings.ToLower(joined)
	fmt.Println(lower)
	for i := 0; i < len(lower)/2; i = i + 2 {
		if lower[i] != lower[len(lower)-i-2] {
			fmt.Println("false", lower[i], len(lower)-i-2)
			break
		}
	}
	fmt.Println("true")
}
