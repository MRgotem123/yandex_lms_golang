package main

import "fmt"

func main() {
	var number int
	fmt.Scanln(&number)
	if number < 0 {
		fmt.Println("Число отрицательное")
	} else if number > 0 {
		fmt.Println("Число положительное")
	} else if number == 0 {
		fmt.Print("Введен ноль")
	}
}
