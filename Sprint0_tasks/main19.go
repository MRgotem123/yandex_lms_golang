package main

import "fmt"

func main() {
	m := [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println(SumOfArray(m))
}
func SumOfArray(m [6]int) int {
	count := 0
	for i := 0; i < 6; i++ {
		count += m[i]
	}
	return count
}
