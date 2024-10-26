package main

import "fmt"

func DivideIntegers(a, b int) (float64, error) {
	a1 := float64(a)
	b1 := float64(b)
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a1 / b1, nil
}
