package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testValues struct {
	Expression string
	Status     string
	Result     string
}

func TestOrchestrator(t *testing.T) {

	tests := []struct {
		name           string
		expression     string
		expectedStatus int
		output         string
	}{
		{
			name:           "Корректное выражение",
			expression:     "2+2-5*2",
			expectedStatus: http.StatusCreated,
			output:         "Уникальный id на ваше выражение:",
		},
		{
			name:           "Некорректное выражение",
			expression:     "!2+2-5*2",
			expectedStatus: http.StatusUnprocessableEntity,
			output:         "Ошибка некоректное выражение: в параметре 'calculate' присутствуют невалидные знаки\n",
		},
		{
			name:           "Неверное количество скобок",
			expression:     "(2+2-5*2",
			expectedStatus: http.StatusUnprocessableEntity,
			output:         "Ошибка некоректное выражение: неверное количество открывающих скобок\n",
		},
	}

	//handler := http.HandlerFunc(Orchestrator)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "http://localhost:9090/api/v1/calculate", bytes.NewReader([]byte(tt.expression)))

			rr := httptest.NewRecorder()
			Orchestrator(rr, req)

			res := rr.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			fmt.Println(string(data))

			// Проверяем, что тело ответа содержит ожидаемый результат
			if !strings.Contains(string(data), tt.output) {
				t.Errorf("ожидалось тело %q, получено %q", tt.output, rr.Body.String())
			}
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Ожидался статус: %d, получен статус: %d", tt.expectedStatus, res.StatusCode)
			}

			req2 := httptest.NewRequest("GET", "http://localhost:9090/api/v1/expressions", nil)

			rr2 := httptest.NewRecorder()
			Expressions(rr2, req2)

			res2 := rr2.Result()
			defer res2.Body.Close()
			body, err := io.ReadAll(res2.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if res2.StatusCode != http.StatusOK {
				t.Errorf("Ожидался статус: %d, получен статус: %d", tt.expectedStatus, res.StatusCode)
			}

			bodystr := string(body)

			if !strings.Contains(bodystr, "2+2-5*2") {
				t.Errorf("ожидалось тело 2+2-5*2, получено %q", bodystr)
			}
		})
	}
}
