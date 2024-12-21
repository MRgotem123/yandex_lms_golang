package main

func main() {

}
func SumOfValuesInMap(m map[int]int) int {
	sum := 0
	for v := range m {
		sum += v
	}
	return sum
}
