package main

import "fmt"

func Factorial(n int) (int, error) {
	if n > 0 {
		count := 1
		for i := 0; i < n+1; i++ {
			count *= i
		}
		return count, nil
	} else {
		return 0, fmt.Errorf("factorial is not defined for negative numbers")
	}
}
