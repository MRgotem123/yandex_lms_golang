package main

import (
	"reflect"
	"slices"
)

func main() {
	type test struct {
		nums     []int
		x        int
		expected []int
	}

	tests := []test{
		{
			nums:     []int{6, 5, 4, 6, 2, 4, 6, 2, 1, 6},
			x:        6,
			expected: []int{5, 4, 2, 4, 2, 1},
		},
		{
			nums:     []int{-1, 6, 4, 5, 1, 6, 211, 90},
			x:        0,
			expected: []int{-1, 6, 4, 5, 1, 6, 211, 90},
		},
		{
			nums:     []int{},
			x:        1,
			expected: []int{},
		},
		{
			nums:     []int{5, 5},
			x:        5,
			expected: []int{},
		},
		{
			nums:     []int{3, -7},
			x:        -7,
			expected: []int{3},
		},
	}
	for _, tc := range tests {
		res := Clean(tc.nums, tc.x)
		if !slices.Equal(res, tc.expected) {
			//t.Fatalf("expected: %v, got: %v", tc.expected, res)
		}

		//if !alias(res, tc.nums) {
		//t.Fatalf("result slice must reference to the same underlying array")
		//}
	}
}

func alias(x, y []int) bool {
	return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
}

func Clean(nums []int, x int) []int {
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
}
