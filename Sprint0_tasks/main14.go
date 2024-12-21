package main

import "fmt"

func main() {
	fmt.Println(SumDigitsRecursive(123))
}
func SumDigitsRecursive(n int) int {
	count := 0
	for n != 0 {
		count1 := n % 10
		count += count1
		n /= 10
	}
	return count
}
