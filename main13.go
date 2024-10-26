package main

import (
	"fmt"
)

func main() {
	Add := Add(10, 5)
	fmt.Println("Сумма:", Add)
}

func Add(a, b float64) float64 {
	Add := a + b
	return Add
}
func Multiply(a, b float64) float64 {
	Multiply := a * b
	return Multiply
}
func PrintNumbersAscending(n int) {
	for i := 1; i < n+1; i++ {
		fmt.Print(i, " ")
	}
}
