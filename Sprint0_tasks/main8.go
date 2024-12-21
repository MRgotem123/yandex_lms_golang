package main

import "fmt"

var number int
var count int = 1

func main() {
	fmt.Scanln(&number)
	for i := 1; i <= number; i++ {
		count *= i
	}
	fmt.Println(count)
}
