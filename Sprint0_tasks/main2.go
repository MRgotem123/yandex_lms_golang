package main

import "fmt"

func main() {
	var number int
	fmt.Scanln(&number)
	if number == 0 {
		fmt.Print("Число равно нулю")
	} else {
		if number > 0 {
			if number%2 == 0 {
				fmt.Println("Число положительное и четное")
			} else {
				fmt.Println("Число положительное и нечетное")
			}
		} else if number%2 == 0 {
			fmt.Println("Число отрицательное и четное")
		} else {
			fmt.Println("Число отрицательное и нечетное")
		}
	}
}
