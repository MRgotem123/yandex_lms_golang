package main

import "fmt"

func main() {
	var number int
	var count int = 0
	fmt.Scanln(&number)
	for i := 1; i < number+1; i++ {
		count += i
	}
	fmt.Println(count)
}
