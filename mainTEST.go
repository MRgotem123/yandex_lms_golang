package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Calculator принимает строку с выражением и возвращает результат его вычисления.
func Calculator(expression string) (float64, error) {
	// Убираем пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")
	return parseExpression(expression)
}

// parseExpression рекурсивно обрабатывает выражение и вычисляет его результат.
func parseExpression(expression string) (float64, error) {
	// Если выражение начинается и заканчивается скобками, нужно обработать внутреннее выражение.
	if strings.HasPrefix(expression, "(") && strings.HasSuffix(expression, ")") {
		innerExpr := expression[1 : len(expression)-1]
		return parseExpression(innerExpr)
	}

	// Обработка сложения и вычитания (с учетом приоритета операций)
	for i := len(expression) - 1; i >= 0; i-- {
		if expression[i] == '+' || expression[i] == '-' {
			left, err := parseExpression(expression[:i])
			if err != nil {
				return 0, err
			}
			right, err := parseExpression(expression[i+1:])
			if err != nil {
				return 0, err
			}
			if expression[i] == '+' {
				return left + right, nil
			} else {
				return left - right, nil
			}
		}
	}

	// Обработка умножения и деления
	for i := len(expression) - 1; i >= 0; i-- {
		if expression[i] == '*' || expression[i] == '/' {
			left, err := parseExpression(expression[:i])
			if err != nil {
				return 0, err
			}
			right, err := parseExpression(expression[i+1:])
			if err != nil {
				return 0, err
			}
			if expression[i] == '*' {
				return left * right, nil
			} else {
				if right == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return left / right, nil
			}
		}
	}

	// Если ни один оператор не найден, преобразуем строку в число.
	return strconv.ParseFloat(expression, 64)
}

func main() {
	expressions := []string{
		"3 + 2 * (1 + 3) - 5 / (2 + 3)",
		"((2 + 3) * 4) - (5 * (1 + 1))",
		"4 * (3 + 2) - 6 / 3",
	}

	for _, expr := range expressions {
		result, err := Calculator(expr)
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Printf("Результат выражения '%s' = %f\n", expr, result)
		}
	}
}
