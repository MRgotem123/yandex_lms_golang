package main

import "fmt"

func main() {
	nums := [10]int{6, 5, 4, 6, 2}
	b := make([]int, len(nums), cap(nums))
	fmt.Println(b)
	//a := (len(nums))
	//nums2 := [a]nums[:a]
	//fmt.Println(nums2)
	/*for i := 0; i < len(nums); i++ {
		fmt.Print(nums[i])
	}
	nums1 := cap(nums)
	nums2 := len(nums)

	fmt.Println("", nums2, nums1)*/
}

/*func SliceCopy(nums []int) []int {
	a := (len(nums))
	nums2 := nums[:a]
	return nums2
}*/
