package main

import "fmt"

func main() {
	var number int
	var number1 int
	var number2 int
	fmt.Scanln(&number, &number1, &number2)
	if number != 0 && number1 != 0 && number2 != 0 {
		if number == number1 && number1 == number2 {
			fmt.Println("Все числа равны")
		} else if number == number1 || number == number2 || number1 == number2 {
			fmt.Println("Два числа равны")
		} else if number != number1 || number != number2 || number1 != number2 {
			fmt.Println("Все числа разные")
		}
	} else {
		fmt.Println("Некорректный ввод")
	}
}
