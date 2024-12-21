package main

import "fmt"

func main() {
	var number int
	var number1 int
	fmt.Scanln(&number, &number1)
	if number == 0 || number1 == 0 {
		fmt.Println("Одно из чисел равно нулю")
	} else {
		if number > 0 && number1 < 0 || number1 > 0 && number < 0 {
			fmt.Println("Одно число положительное, а другое отрицательное")
		} else if number > 0 && number1 > 0 {
			fmt.Println("Оба числа положительные")
		} else {
			fmt.Println("Оба числа отрицательные")
		}
	}
}
