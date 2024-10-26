package main

import "fmt"

func main() {
	array := [10]int{1, 22, 4, 6, 56, 0, 23, 100, 9, 10}
	output := FindMinMaxInArray(array)
	fmt.Print(output)
}
func FindMinMaxInArray(array [10]int) (int, int) {
	mn := array[0]
	mx := array[0]
	for i := 0; i < 10; i++ {
		if mn >= array[i] {
			mn = array[i]
		} else if mx <= array[i] {
			mx = array[i]
		}
	}
	return mn, mx
}
