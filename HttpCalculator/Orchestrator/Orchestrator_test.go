package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalSlicesString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAddToQueueTaskMap(t *testing.T) {
	clearQueue := func() {
		for k := range QueueTask {
			delete(QueueTask, k)
		}
	}

	tests := []struct {
		name           string
		id             string
		arg1           string
		arg2           string
		operation      string
		operation_time int
		wantErr        bool
		expectedErr    error
		expectedMap    map[string]SendValues
	}{
		{
			name:           "successful addition",
			id:             "task1",
			arg1:           "10",
			arg2:           "20",
			operation:      "+",
			operation_time: 100,
			wantErr:        false,
			expectedMap: map[string]SendValues{
				"task1": {
					Id:             "task1",
					Arg1:           "10",
					Arg2:           "20",
					Operation:      "+",
					Operation_time: 100,
				},
			},
		},
		{
			name:           "empty id",
			id:             "",
			arg1:           "10",
			arg2:           "20",
			operation:      "+",
			operation_time: 100,
			wantErr:        true,
			expectedErr:    fmt.Errorf("Один из элементов пустой!"),
			expectedMap:    map[string]SendValues{},
		},
		{
			name:           "empty arg1",
			id:             "task1",
			arg1:           "",
			arg2:           "20",
			operation:      "+",
			operation_time: 100,
			wantErr:        true,
			expectedErr:    fmt.Errorf("Один из элементов пустой!"),
			expectedMap:    map[string]SendValues{},
		},
		{
			name:           "empty arg2",
			id:             "task1",
			arg1:           "10",
			arg2:           "",
			operation:      "+",
			operation_time: 100,
			wantErr:        true,
			expectedErr:    fmt.Errorf("Один из элементов пустой!"),
			expectedMap:    map[string]SendValues{},
		},
		{
			name:           "empty operation",
			id:             "task1",
			arg1:           "10",
			arg2:           "20",
			operation:      "",
			operation_time: 100,
			wantErr:        true,
			expectedErr:    fmt.Errorf("Один из элементов пустой!"),
			expectedMap:    map[string]SendValues{},
		},
		{
			name:           "update existing task",
			id:             "task1",
			arg1:           "15",
			arg2:           "25",
			operation:      "*",
			operation_time: 200,
			wantErr:        false,
			expectedMap: map[string]SendValues{
				"task1": {
					Id:             "task1",
					Arg1:           "15",
					Arg2:           "25",
					Operation:      "*",
					Operation_time: 200,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearQueue()

			err := AddToQueueTaskMap(tt.id, tt.arg1, tt.arg2, tt.operation, tt.operation_time)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddToQueueTaskMap() error = %v, want %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("AddToQueueTaskMap() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if len(QueueTask) != len(tt.expectedMap) {
				t.Errorf("QueueTask length = %d, want %d", len(QueueTask), len(tt.expectedMap))
				return
			}

			for k, v := range tt.expectedMap {
				actual, ok := QueueTask[k]
				if !ok {
					t.Errorf("Expected key %s not found in QueueTask", k)
					continue
				}

				if actual != v {
					t.Errorf("QueueTask[%s] = %v, want %v", k, actual, v)
				}
			}
		})
	}
}

func TestNormalExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "successful expression",
			expression:  "-2+2-10*13.5",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "empty expression",
			expression:  "",
			wantErr:     true,
			expectedErr: fmt.Errorf("параметр 'calculate' отсутствует или пуст"),
		},
		{
			name:        "unknown symbol",
			expression:  "2-5!",
			wantErr:     true,
			expectedErr: fmt.Errorf("в параметре 'calculate' присутствуют невалидные знаки"),
		},
		{
			name:        "two characters in a row",
			expression:  "2+/5*10",
			wantErr:     true,
			expectedErr: fmt.Errorf("два оператора подряд недопустимы"),
		},
		{
			name:        "unknown dot",
			expression:  "2*10.",
			wantErr:     true,
			expectedErr: fmt.Errorf("некорректное использование десятичной точки"),
		},
		{
			name:        "two dot in a row",
			expression:  "2*10..5",
			wantErr:     true,
			expectedErr: fmt.Errorf("число не может содержать две точки подряд"),
		},
		{
			name:        "start with operator",
			expression:  "*2*10*10",
			wantErr:     true,
			expectedErr: fmt.Errorf("выражение не может начинаться с оператора (кроме унарного минуса)"),
		},
		{
			name:        "finish with operator",
			expression:  "2*10*10-",
			wantErr:     true,
			expectedErr: fmt.Errorf("выражение не может заканчиваться оператором"),
		},
		{
			name:        "start with operator after (",
			expression:  "(/2+10)*10",
			wantErr:     true,
			expectedErr: fmt.Errorf("'*' или '/' не могут идти сразу после '('"),
		},
		{
			name:        "two operator in a row",
			expression:  "2+10*/10",
			wantErr:     true,
			expectedErr: fmt.Errorf("два оператора подряд недопустимы"),
		},
		{
			name:        "unexted )",
			expression:  "2*(10-10))",
			wantErr:     true,
			expectedErr: fmt.Errorf("неверное количество закрывающих скобок"),
		},
		{
			name:        "unexted (",
			expression:  "2*((10-10)",
			wantErr:     true,
			expectedErr: fmt.Errorf("неверное количество открывающих скобок"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := NormalExpression(tt.expression)

			if (err != nil) != tt.wantErr {
				t.Errorf("NormalExpression() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("NormalExpression() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestAnalys(t *testing.T) {

	originalSlises := make([][]string, len(Slises_easyExpr))
	copy(originalSlises, Slises_easyExpr)
	defer func() {
		Slises_easyExpr = originalSlises
	}()

	tests := []struct {
		name           string
		rpnExpression  []string
		expectedDel    []int
		expectedID     []int
		expectedSlices [][]string
	}{
		{
			name:           "simple expression",
			rpnExpression:  []string{"2", "3", "+"},
			expectedDel:    []int{0, 1},
			expectedID:     []int{2},
			expectedSlices: [][]string{{"2", "3", "+"}},
		},
		{
			name:           "two operations",
			rpnExpression:  []string{"2", "3", "+", "5", "*"},
			expectedDel:    []int{0, 1},
			expectedID:     []int{2},
			expectedSlices: [][]string{{"2", "3", "+"}},
		},
		{
			name:           "complex expression",
			rpnExpression:  []string{"2", "3", "4", "*", "+"},
			expectedDel:    []int{1, 2},
			expectedID:     []int{3},
			expectedSlices: [][]string{{"3", "4", "*"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем Slises_easyExpr перед каждым тестом
			Slises_easyExpr = [][]string{}

			del, id := Analys(tt.rpnExpression)

			// Проверяем возвращаемые срезы
			if !equalSlices(del, tt.expectedDel) {
				t.Errorf("Analys() del = %v, expected %v", del, tt.expectedDel)
			}

			if !equalSlices(id, tt.expectedID) {
				t.Errorf("Analys() id = %v, expected %v", id, tt.expectedID)
			}

			// Проверяем Slises_easyExpr
			if len(Slises_easyExpr) != len(tt.expectedSlices) {
				t.Errorf("Slises_easyExpr length = %d, expected %d", len(Slises_easyExpr), len(tt.expectedSlices))
			} else {
				for i := range Slises_easyExpr {
					if !equalSlicesString(Slises_easyExpr[i], tt.expectedSlices[i]) {
						t.Errorf("Slises_easyExpr[%d] = %v, expected %v", i, Slises_easyExpr[i], tt.expectedSlices[i])
					}
				}
			}
		})
	}
}

func TestToRpn(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedRpn []string
	}{
		{
			name:        "normal",
			expression:  "2*10*10",
			expectedRpn: []string{"2", "10", "*", "10", "*"},
		},
		{
			name:        "with unary simbol",
			expression:  "2+-10*10+2",
			expectedRpn: []string{"2", "-10", "10", "*", "+", "2", "+"},
		},
		{
			name:        "with `( )`",
			expression:  "(2-10)*(3+8)",
			expectedRpn: []string{"2", "10", "-", "3", "8", "+", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			expressionRPN := ToRPN(tt.expression)

			if len(expressionRPN) != len(tt.expectedRpn) {
				t.Errorf("Expected RPN length %d, got %d", len(tt.expectedRpn), len(expressionRPN))
				return
			}

			for i := range expressionRPN {
				if expressionRPN[i] != tt.expectedRpn[i] {
					t.Errorf("At position %d: expected %q, got %q", i, tt.expectedRpn[i], expressionRPN[i])
				}
			}
		})
	}
}

func TestOrchestrator(t *testing.T) {
	// Сохраняем оригинальные глобальные переменные
	originalUserID := UserID
	originalExpressionToRPN := ExpressionToRPN
	originalSlisesEasyExpr := Slises_easyExpr
	//originalDB := DB
	defer func() {
		UserID = originalUserID
		ExpressionToRPN = originalExpressionToRPN
		Slises_easyExpr = originalSlisesEasyExpr
	}()

	tests := []struct {
		name         string
		userID       string
		requestBody  string
		expectedCode int
		expectedBody string
		checkGlobals func(t *testing.T)
	}{
		{
			name:         "unauthorized access",
			userID:       "",
			requestBody:  `{"expression": "2+2"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"Войдите в акаунт или харегестрируйтесь."}`,
		},
		{
			name:         "invalid JSON",
			userID:       "user1",
			requestBody:  "invalid json",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error": "Некорректное выражение"}`,
		},
		{
			name:         "simple addition",
			userID:       "user1",
			requestBody:  "3+4",
			expectedCode: http.StatusCreated,
			expectedBody: "Уникальный id на ваше выражение:",
			checkGlobals: func(t *testing.T) {
				if len(ExpressionToRPN) == 0 {
					t.Error("ExpressionToRPN должен содержать выражение")
				}
			},
		},
		{
			name:         "complex expression",
			userID:       "user1",
			requestBody:  "(3+4)*5/2",
			expectedCode: http.StatusCreated,
			expectedBody: "Уникальный id на ваше выражение:",
		},
		{
			name:         "division by zero",
			userID:       "user1",
			requestBody:  "2/0",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error": "деление на ноль!"}`,
		},
		{
			name:         "empty expression",
			userID:       "user1",
			requestBody:  "",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error": "Некорректное выражение"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка глобальных переменных
			UserID = tt.userID
			ExpressionToRPN = make(map[string][]string)
			Slises_easyExpr = nil

			// Создание тестового запроса
			req := httptest.NewRequest("POST", "/calculate", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Создание ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызов тестируемой функции
			Orchestrator(rr, req)

			// Проверка кода статуса
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// Проверка тела ответа
			body := rr.Body.String()
			if tt.expectedBody != "" && !bytes.Contains([]byte(body), []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}

			// Дополнительные проверки глобальных переменных
			if tt.checkGlobals != nil {
				tt.checkGlobals(t)
			}
		})
	}
}

// TestTaskHandlerIntegration - интеграционный тест

type Task struct {
	ID     string
	Result string
}

func TestOrchestratorReturnIntegration(t *testing.T) {
	// Создаем обработчик с recovery middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Вызов оригинального обработчика
		OrchestratorReturn(w, r)
	})

	// 1. Создаем тестовый сервер с нашим обработчиком
	server := httptest.NewServer(handler)
	defer server.Close()

	// 2. Тестируем создание новой задачи
	t.Run("Create task with invalid ID", func(t *testing.T) {
		newTask := Task{
			ID:     "idab123_id123abc",
			Result: "5.000",
		}

		// Сериализуем задачу в JSON
		body, err := json.Marshal(newTask)
		if err != nil {
			t.Fatalf("Failed to marshal task: %v", err)
		}

		// Создаем POST-запрос
		resp, err := http.Post(server.URL+"/internal/task", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Проверяем код статуса (ожидаем 404, так как ID нет в ExpressionToRPN)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
		}

		// Проверяем тело ответа
		var response map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if response["error"] != "Не найдено совпадение ID в ExpressionToRPN[ID]" {
			t.Errorf("Expected error message 'Не найдено совпадение ID в ExpressionToRPN[ID]', got '%s'", response["error"])
		}
	})

	// 3. Тестируем получение списка задач
	t.Run("Get tasks", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/internal/task")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound || resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d or %d, got %d", http.StatusNotFound, http.StatusOK, resp.StatusCode)
		}
	})
}
