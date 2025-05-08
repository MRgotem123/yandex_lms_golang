package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
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
		t.Error("Ожидалась ошибка при сетевой проблеме, но её нет")
	}
}

func TestWorker(t *testing.T) {
	tests := []struct {
		name           string
		getTask        func(url string) (*SendValues2, error)
		evaluateRPN    func(tokens []string) (float64, error)
		sendResultFunc func(serverURL, taskID string, result float64) error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "normal worker",
			getTask: func(url string) (*SendValues2, error) {
				return &SendValues2{
					Id:             "123",
					Arg1:           "5",
					Arg2:           "3",
					Operation:      "+",
					Operation_time: 200,
				}, nil
			},
			evaluateRPN: func(tokens []string) (float64, error) {
				return 8.000, nil
			},
			sendResultFunc: func(serverURL, taskID string, result float64) error {
				if taskID != "123" {
					t.Errorf("expected taskID 123, got %s", taskID)
				}
				if result != 8.0 {
					t.Errorf("expected result 8.0, got %f", result)
				}
				return nil
			},
			expectError: false,
		},
		{
			name: "getTask error",
			getTask: func(url string) (*SendValues2, error) {
				return nil, errors.New("getTask error")
			},
			evaluateRPN: func(tokens []string) (float64, error) {
				t.Error("evaluateRPN should not be called when getTask fails")
				return 0, nil
			},
			sendResultFunc: func(serverURL, taskID string, result float64) error {
				t.Error("sendResultFunc should not be called when getTask fails")
				return nil
			},
			expectError: false, // Worker продолжает работать после ошибки
		},
		{
			name: "evaluateRPN error",
			getTask: func(url string) (*SendValues2, error) {
				return &SendValues2{
					Id:             "123",
					Arg1:           "5",
					Arg2:           "3",
					Operation:      "+",
					Operation_time: 200,
				}, nil
			},
			evaluateRPN: func(tokens []string) (float64, error) {
				return 0, errors.New("evaluateRPN error")
			},
			sendResultFunc: func(serverURL, taskID string, result float64) error {
				t.Error("sendResultFunc should not be called when evaluateRPN fails")
				return nil
			},
			expectError: true,
		},
		{
			name: "sendResult error",
			getTask: func(url string) (*SendValues2, error) {
				return &SendValues2{
					Id:             "123",
					Arg1:           "5",
					Arg2:           "3",
					Operation:      "+",
					Operation_time: 200,
				}, nil
			},
			evaluateRPN: func(tokens []string) (float64, error) {
				return 8.0, nil
			},
			sendResultFunc: func(serverURL, taskID string, result float64) error {
				return errors.New("sendResult error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			done := make(chan struct{})
			defer close(done)

			go func() {
				defer wg.Done()
				Worker(1)
			}()

			time.Sleep(300 * time.Millisecond)
		})
	}
}
