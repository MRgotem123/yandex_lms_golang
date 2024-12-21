package main

import "fmt"

type Animal interface {
	MakeSound() string
}

func MakeSound(s Animal) string {
	return s.MakeSound()
}

type Dog struct {
	sound string
}
type Cat struct {
	sound string
}

func (d Dog) MakeSound() string {
	fmt.Print("Гав!")
}
func (c Cat) MakeSound() float64 {
	fmt.Print("Мяу!")
}
