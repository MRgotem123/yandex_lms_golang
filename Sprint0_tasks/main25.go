package main

import "fmt"

func main() {
	nums1 := []int{1, 3, 5, 7, 99}
	nums2 := []int{2, 4, 5, 6, 44}

	icvel := make([]int, len(nums1)+len(nums2))
	icvel = append(icvel, nums1, nums2)
	fmt.Print(icvel)
}

/*func Join(nums1, nums2 []int) []int {
	icvel := make([]int, len(nums1) + len(nums2))
	copy(icvel, nums1)
	fmt.Print(icvel)
}*/
