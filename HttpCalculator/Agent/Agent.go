package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const COMPUTING_POWER = 5

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

// var resultData2 = make(map[string]string)
var mutex sync.Mutex

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

func SendResultToOrchestrator(serverURL, taskID string, result float64) error {
	resultStr := strconv.FormatFloat(result, 'f', 3, 64)

	//resultData2[taskID] = resultStr
	sendResult := resultData2{
		ID:     taskID,
		Result: resultStr,
	}

	jsonData, err := json.Marshal(sendResult)
	if err != nil {
		fmt.Println("Ошибка маршалинга JSON:", err)
		return err
	}

	fmt.Println("Отправляем результат на Оркестратор:", string(jsonData))
	// Отправляем результат обратно на тот же адрес
	resp, err := http.Post(serverURL+"/internal/task", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Ошибка отправки POST-запроса:", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Оркестратор ответил статусом:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Ошибка от Оркестратора:", string(body))
		return fmt.Errorf("не удалось отправить результат: код %d", resp.StatusCode)
	}

	fmt.Println("Результат успешно отправлен!")
	return nil
}

func Worker(i int) {
	defer wg.Done()
	fmt.Println("136. запусщинна горутина", i)
	for {
		client := &http.Client{}
		task, err := GetTask(client, "http://localhost:9090/internal/task")
		if err != nil {
			fmt.Printf("Ошибка получения задачи: %v\n", err)
			time.Sleep(5 * time.Second) // Подождать перед повторной попыткой
			continue
		}

		// Отправка подтверждения
		fmt.Println("Задача успешно запущена!")

		//if task.Arg2 == "0" && task.Operation == "/" { }
		// Формируем выражение
		expression := []string{task.Arg1, task.Arg2, task.Operation}
		fmt.Println("ФОРМЕРУЕМ ВЫРАЖЕНИЕ", expression)

		// Выполнение вычислений
		result, err := evaluateRPN(expression)
		if err != nil {
			fmt.Printf("Ошибка вычисления задачи: %v\n", err)
			continue
		}
		fmt.Println("ВЫЧИСЛЯЕМ ВЫРАЖЕНИЕ", result)

		time.Sleep(time.Duration(task.Operation_time) * time.Millisecond)

		// Отправка результата
		err = SendResultToOrchestrator("http://localhost:9090", task.Id, result)
		if err != nil {
			fmt.Printf("Ошибка отправки результата: %v\n", err)
			continue
		}
		fmt.Println("ОТПРАВЛЯЕМ ВЫРАЖЕНИЕ")
	}
}

func main() {
	wg.Add(COMPUTING_POWER)

	for i := 0; i < COMPUTING_POWER; i++ {
		go Worker(i)
	}
	wg.Wait()
}
