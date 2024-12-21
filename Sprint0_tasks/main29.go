package main

func main() {

}
func FindMaxKey(m map[int]int) int {
	first := true
	max := 0
	for k := range m {
		if first {
			first = false
			max = k
		} else if k > max {
			max = k
		}
	}
	return max
}
