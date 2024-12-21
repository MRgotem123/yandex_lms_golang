package main

import "fmt"

func main() {
	nums := []int{5, 5}
	x := 5
	nums_k := len(nums)
	for i := 0; i < nums_k; i++ {
		if nums[i] == x {
			nums = append(nums[:i], nums[i+1:]...)
			nums_k--
		}
	}
	fmt.Println(nums)
}

/*func Clean(nums []int, x int) []int {
	count := 0
	for i := 0; i < len(nums); i++ {
		count += 1
		if nums[i] == x {
			nums[i] = nums[i+1]
			for a := count; a <= len(nums)-count; {
				nums[i] = nums[i+1]
			}
		} else {
			continue
		}
	}
	return nums
}*/
