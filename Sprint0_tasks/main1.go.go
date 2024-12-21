package main

import "fmt"

func main() {
	var number int
	var number2 int
	fmt.Scanln(&number, &number2) // Fix: pass a pointer to number2
	if number > number2 {
		fmt.Print("Первое число больше второго")
	} else if number2 > number {
		fmt.Print("Второе число больше первого")
	} else {
		fmt.Print("Числа равны")
	}
}
