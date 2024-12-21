package main

import (
	"errors"
	"fmt"
)

func main() {
	type test struct {
		nums      []int
		n         int
		limit     int
		expected  []int
		wantError bool
	}

	tests := []test{
		{
			nums:      []int{4, 7, 89, 3, 21, 2, 5, 7, 32, 4, 6, 8, 0, 3, 4, 6, 2, 115, 12},
			n:         5,
			limit:     3,
			expected:  []int{2, 0, 2},
			wantError: false,
		},
		{
			nums:      nil,
			wantError: true,
		},
		{
			nums:      []int{},
			n:         5,
			limit:     3,
			expected:  []int{},
			wantError: false,
		},
		{
			nums:      []int{3, 5, 6},
			n:         5,
			limit:     10,
			expected:  []int{3, 5, 6},
			wantError: false,
		},
		{
			nums:      []int{-13, 0, 6},
			n:         1,
			limit:     -5,
			expected:  []int{-13},
			wantError: false,
		},
		{
			nums:      []int{},
			n:         -1,
			limit:     5,
			wantError: true,
		},
	}

	for _, tc := range tests {
		fmt.Println(UnderLimit(tc.nums, tc.limit, tc.n))
	}
}

func UnderLimit(nums []int, limit int, n int) ([]int, error) {
	if nums == nil {
		return nil, errors.New("nums = nil")
	}
	if n < 0 {
		return nil, errors.New("n < 0")
	}
	cusok2 := make([]int, len(nums))
	count := 0
	for i := 0; i < len(nums) && count <= n; i++ {
		if nums[i] < limit {
			cusok2[count] = nums[i]
			count++
		} else {
			continue
		}
	}
	cusok2 = cusok2[:count]
	return cusok2, nil
}
