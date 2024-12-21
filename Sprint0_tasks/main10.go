package main

import "fmt"

var number int

func main() {
	fmt.Scanln(&number)
	for i := 1; i < number+1; i++ {
		if i%3 != 0 {
			continue
		}
		fmt.Println(i)
	}
}
