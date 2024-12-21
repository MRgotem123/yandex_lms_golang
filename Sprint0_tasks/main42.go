package main

import "fmt"

func GetCharacterAtPosition(str string, position int) (rune, error) {
	if len(str) >= position {
		return rune(str[position]), nil
	} else {
		return rune(str[0]), fmt.Errorf("position out of range")
	}
}
