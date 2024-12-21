package main

import "fmt"

type Person struct {
	name   string
	age    int
	adress string
}

func (p Person) printData() {
	fmt.Printf("Name: %s, Age: %d\n, Adress: %d\n", p.name, p.age, p.adress)
}

func main() {

}
