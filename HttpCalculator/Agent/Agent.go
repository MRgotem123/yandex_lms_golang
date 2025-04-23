package main

import (
	"log"
	"net/http"
	"time"
)

const COMPUTING_POWER = 5

func Worker(i int) {
	defer wg.Done()
	log.Println("14. запусщинна горутина", i)
	for {
		client := &http.Client{}
		task, err := GetTask(client, "http://localhost:9090/internal/task")
		if err != nil {
			log.Printf("Ошибка получения задачи: %v\n", err)
			time.Sleep(5 * time.Second) // Подождать перед повторной попыткой
			continue
		}

		// Отправка подтверждения
		log.Println("Задача успешно запущена!")

		// Формируем выражение
		expression := []string{task.Arg1, task.Arg2, task.Operation}
		log.Println("ФОРМЕРУЕМ ВЫРАЖЕНИЕ", expression)

		// Выполнение вычислений
		result, err := EvaluateRPN(expression)
		if err != nil {
			log.Printf("Ошибка вычисления задачи: %v\n", err)
			continue
		}
		log.Println("ВЫЧИСЛЯЕМ ВЫРАЖЕНИЕ", result)

		time.Sleep(time.Duration(task.Operation_time) * time.Millisecond)

		// Отправка результата
		err = SendResultToOrchestrator("http://localhost:9090", task.Id, result)
		if err != nil {
			log.Printf("Ошибка отправки результата: %v\n", err)
			continue
		}
		log.Println("ОТПРАВЛЯЕМ ВЫРАЖЕНИЕ")
	}
}

func main() {
	wg.Add(COMPUTING_POWER)

	for i := 0; i < COMPUTING_POWER; i++ {
		go Worker(i)
	}
	wg.Wait()
}
