package main

import "fmt"

type Employee struct {
	name     string
	position string
	salary   float64
	bonus    float64
}

func (e Employee) Print() {
	fmt.Printf("Name: %s\nposition: %s\nsalary: %.2d\n, bonus: %.2d\n", e.name, e.position, e.salary+e.bonus)
}
