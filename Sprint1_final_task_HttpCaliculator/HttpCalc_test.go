package HttpCaliculator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	handler := Input(CalcHandler)

	tests := []struct {
		name           string
		inp_task       string
		expectedStatus int
		out_task       string
	}{
		{
			name:           "Обычный пример",
			inp_task:       "?calculate=2%2B3*4",
			expectedStatus: http.StatusOK,
			out_task:       "Результат: 14.00",
		},
		{
			name:           "Отсутствие примера",
			inp_task:       "?calculate=",
			expectedStatus: http.StatusBadRequest,
			out_task:       "Параметр 'calculate' отсутствует или пуст",
		},
		{
			name:           "Пустой запрос",
			inp_task:       "",
			expectedStatus: http.StatusBadRequest,
			out_task:       "Параметр 'calculate' отсутствует или пуст",
		},
		{
			name:           "Деление на ноль",
			inp_task:       "?calculate=4/0",
			expectedStatus: http.StatusBadRequest,
			out_task:       "деление на ноль",
		},
		{
			name:           "Деление нуля",
			inp_task:       "?calculate=0/2",
			expectedStatus: http.StatusOK,
			out_task:       "Результат: 0.00",
		},
		{
			name:           "Нехватка скобок",
			inp_task:       "?calculate=(5+3*(2+3)",
			expectedStatus: http.StatusBadRequest,
			out_task:       "не удалось распознать число: (",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tt.inp_task, nil)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Проверка статуса
			if rr.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d, тест: %q", tt.expectedStatus, rr.Code, tt.name)
			}

			// Проверка тела ответа
			if tt.out_task != "" && !strings.Contains(rr.Body.String(), tt.out_task) {
				t.Errorf("ожидалось тело %q, получено %q", tt.out_task, rr.Body.String())
			}
		})
	}
}
