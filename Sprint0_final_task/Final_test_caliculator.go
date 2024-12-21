package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func toRPN(expression string) ([]string, error) {
	var output []string
	var operators []rune

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' {
			number := string(char)
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				i++
				number += string(expression[i])
			}
			output = append(output, number)
		} else if char == '(' {
			operators = append(operators, char)
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) > 0 && operators[len(operators)-1] == '(' {
				operators = operators[:len(operators)-1]
			}
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[char] {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
		}
	}

	for len(operators) > 0 {
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}
	fmt.Println("toPRN: ", output)
	return output, nil
}

func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, fmt.Errorf("функция выдает ошибку")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("функция выдает ошибку")
				}
				stack = append(stack, a/b)
			}
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("функция выдает ошибку")
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("функция выдает ошибку")
	}

	return stack[0], nil
}

func Calc(expression string) (float64, error) {

	expression = strings.ReplaceAll(expression, " ", "")

	rpn, err := toRPN(expression)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpn)
}
func main() {
	fmt.Println(Calc("2+3*4"))
}
