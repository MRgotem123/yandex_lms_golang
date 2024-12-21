package main

import "fmt"

var number int
var count int

func main() {
	fmt.Scanln(&number)
	for i := 1; i <= number; i++ {
		if i%3 == 0 || i%5 == 0 {
			continue
		}
		count += i
	}
	fmt.Println(count)
}
