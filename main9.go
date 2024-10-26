package main

import "fmt"

var number int
var count int = 0

func main() {
	fmt.Scanln(&number)
	if number > 0 {
		for i := 0; i < number+1; i++ {
			if i%2 == 0 {
				continue
			}
			count += i
		}
		fmt.Println(count)
	} else {
		fmt.Println("Некорректный ввод")
	}
}
