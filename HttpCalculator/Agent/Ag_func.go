package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type SendValues2 struct {
	Id             string `json:"id"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

type resultData2 struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

var wg sync.WaitGroup

func GetTask(client *http.Client, url string) (*SendValues2, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("нет задач, статус: %d", resp.StatusCode)
	}

	var task SendValues2
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func EvaluateRPN(tokens []string) (float64, error) {
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

func SendResultToOrchestrator(serverURL, taskID string, result float64) error {
	resultStr := strconv.FormatFloat(result, 'f', 3, 64)

	sendResult := resultData2{
		ID:     taskID,
		Result: resultStr,
	}

	jsonData, err := json.Marshal(sendResult)
	if err != nil {
		log.Println("Ошибка маршалинга JSON:", err)
		return err
	}

	fmt.Println("Отправляем результат на Оркестратор:", string(jsonData))
	// Отправляем результат обратно на тот же адрес
	resp, err := http.Post(serverURL+"/internal/task", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Println("Ошибка отправки POST-запроса:", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Оркестратор ответил статусом:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Ошибка от Оркестратора:", string(body))
		return fmt.Errorf("не удалось отправить результат: код %d", resp.StatusCode)
	}

	log.Println("Результат успешно отправлен!")
	return nil
}
