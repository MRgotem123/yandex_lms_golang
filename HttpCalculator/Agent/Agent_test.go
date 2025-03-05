package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetTask(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		responseStatus int
		expectError    bool
		expectedTask   *SendValues2
	}{
		{
			name: "Успешное получение задачи",
			responseBody: `{
				"id": "123",
				"arg1": "5",
				"arg2": "3",
				"operation": "+",
				"operation_time": 200
			}`,
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedTask: &SendValues2{
				Id:             "123",
				Arg1:           "5",
				Arg2:           "3",
				Operation:      "+",
				Operation_time: 200,
			},
		},
		{
			name:           "Ошибка сервера (500)",
			responseBody:   `{"error": "internal error"}`,
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
			expectedTask:   nil,
		},
		{
			name:           "Некорректный JSON",
			responseBody:   `{"id": 123, "arg1": "5", "arg2": "3", "operation": "+", "operation_time": "200"`, // Ломанный JSON
			responseStatus: http.StatusOK,
			expectError:    true,
			expectedTask:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-сервер
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.responseBody))
			}))
			defer mockServer.Close()

			client := &http.Client{}
			task, err := GetTask(client, mockServer.URL)

			// Проверяем, должна ли быть ошибка
			if (err != nil) != tt.expectError {
				t.Errorf("Ожидалась ошибка: %v, но получено: %v", tt.expectError, err)
			}

			// Проверяем корректность полученной задачи
			if tt.expectedTask != nil {
				if task == nil {
					t.Errorf("Ожидалась задача, но получено nil")
				} else {
					if task.Id != tt.expectedTask.Id || task.Arg1 != tt.expectedTask.Arg1 ||
						task.Arg2 != tt.expectedTask.Arg2 || task.Operation != tt.expectedTask.Operation ||
						task.Operation_time != tt.expectedTask.Operation_time {
						t.Errorf("Ожидалось %+v, но получено %+v", tt.expectedTask, task)
					}
				}
			}
		})
	}
}

func TestSendResultToOrchestrator(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		result         float64
		responseStatus int
		expectError    bool
	}{
		{
			name:           "Успешная отправка",
			taskID:         "123",
			result:         10.5,
			responseStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Ошибка сервера (500)",
			taskID:         "456",
			result:         7.25,
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				if r.Method != http.MethodPost {
					t.Errorf("Ожидался метод POST, но получен: %s", r.Method)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Ожидался Content-Type application/json, но получен: %s", r.Header.Get("Content-Type"))
				}

				// Читаем тело запроса
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Ошибка чтения тела запроса: %v", err)
				}

				var receivedData resultData2
				err = json.Unmarshal(body, &receivedData)
				if err != nil {
					t.Fatalf("Ошибка декодирования JSON: %v", err)
				}

				// Проверяем корректность данных
				expectedResult := strconv.FormatFloat(tt.result, 'f', 3, 64)
				if receivedData.ID != tt.taskID || receivedData.Result != expectedResult {
					t.Errorf("Ожидались данные %+v, но получены %+v", resultData2{ID: tt.taskID, Result: expectedResult}, receivedData)
				}

				w.WriteHeader(tt.responseStatus)
			}))
			defer mockServer.Close()

			err := SendResultToOrchestrator(mockServer.URL, tt.taskID, tt.result)

			if (err != nil) != tt.expectError {
				t.Errorf("Ожидалась ошибка: %v, но получено: %v", tt.expectError, err)
			}
		})
	}
}

// Тест для сетевой ошибки
func TestSendResultToOrchestrator_NetworkError(t *testing.T) {
	// Передаем некорректный URL, чтобы вызвать ошибку соединения
	err := SendResultToOrchestrator("http://invalid-url", "789", 15.0)

	if err == nil {
		t.Error("Ожидалась ошибка при сетевой проблеме, но ее нет")
	}
}
