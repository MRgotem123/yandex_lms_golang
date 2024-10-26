package main

import "fmt"

func main() {
	var number int
	fmt.Scanln(&number)
	if number >= -1 {
		if number < 10 && number >= 0 {
			fmt.Println("Число меньше 10")
		} else if number >= 10 && number < 100 {
			fmt.Println("Число меньше 100")
		} else if number >= 100 && number < 1000 {
			fmt.Println("Число меньше 1000")
		} else {
			fmt.Println("Число больше или равно 1000")
		}
	} else {
		fmt.Println("Некорректный ввод")
	}
}
