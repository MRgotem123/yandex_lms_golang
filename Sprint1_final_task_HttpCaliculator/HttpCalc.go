package HttpCaliculator

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

// проверка на пустое значение
func Input(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		calculate := queryParams.Get("calculate")

		if calculate == "" {
			http.Error(w, "Параметр 'calculate' отсутствует или пуст", http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}

// получаем пример и запускаем каликулятор
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	expression := queryParams.Get("calculate")

	expression = strings.ReplaceAll(expression, " ", "")

	rpn, err := toRPN(expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ansver, err := evaluateRPN(rpn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("Результат: %.2f", ansver)))
}

// переводим пример в обратную польскую нотацию (из 2 * 3 + 4  в  2 3 4 * +)
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
			if _, err := strconv.ParseFloat(number, 64); err != nil {
				return nil, fmt.Errorf("некорректный формат числа: %s", number)
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
	return output, nil
}

// вычисляем пример в обратной польской нотации
func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, fmt.Errorf("недостаточно цифр для операции %s", token)
			}
			if len(stack) >= len(tokens) {
				return 0, fmt.Errorf("Большое количество операторов %s", token)
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
					return 0, fmt.Errorf("деление на ноль")
				}
				stack = append(stack, a/b)
			}
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("не удалось распознать число: %s", token)
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("некорректное выражение: остаток в стеке %v", stack)
	}

	return stack[0], nil
}

func main() {
	mux := http.NewServeMux()

	handler := Input(CalcHandler)
	mux.HandleFunc("/", handler)

	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
