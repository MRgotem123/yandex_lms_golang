package main

import "fmt"

func main() {
	numsData := []int{1, 2, 3}
	numsCap := 10
	nums := make([]int, len(numsData), numsCap)
	if len(nums) < cap(nums) {
		nums = nums[:len(nums)]
		fmt.Println(nums)
	} else {
		fmt.Print(nums)
	}
}

/*func SliceCopy(nums []int) []int {
	if len(nums) < cap(nums) {
		nums = nums[:len(nums)]
		for i := 0; i < len(nums); i++ {
			nums_ret[i] = nums[i]
		}
	} else {
		return nums
	}
	return nums
}*/
