package main

import "fmt"

func main() {
	nums := []int{0, 1, 2, 3, 4, 5}
	for i := 0; i < len(nums)/2; i++ {
		fmt.Print(nums[i], nums[len(nums)/2+i])
	}
}
