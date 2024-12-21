package main

import "fmt"

func main() {

}
func IsPowerOfTwoRecursive(N int) {
	if N == 1 {
		fmt.Println("YES")
	} else if N < 1 || N%2 != 0 {
		fmt.Println("NO")
	} else {
		IsPowerOfTwoRecursive(N / 2)
	}
}
