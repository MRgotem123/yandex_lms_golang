package WorkWithSQL

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockUserRepository мок для тестирования
type MockUserRepository struct {
	UserVerificationFunc  func(login, password string, action int) (string, error)
	InsertUserFunc        func(login, password string) (string, error)
	InsertExpressionsFunc func(userID, expressionId, expression string, statusId int) error
}

func setupTestDB(t *testing.T) *sql.DB {
	// Создаем временную in-memory базу данных
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	return db
}

func TestCreateTables(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := &SQLiteUserRepository{db: db}

	t.Run("CreateTable", func(t *testing.T) {
		err := repo.CreateTabls() // Передаем db в функцию
		if err != nil {
			t.Errorf("CreateTable() error = %v, want nil", err)
		}

		// Проверяем существование таблиц
		_, err = db.Exec("SELECT * FROM users LIMIT 1")
		if err != nil {
			t.Errorf("Table users doesn't exist: %v", err)
		}

		_, err = db.Exec("SELECT * FROM expressions LIMIT 1")
		if err != nil {
			t.Errorf("Table expressions doesn't exist: %v", err)
		}
	})
}

func (m *MockUserRepository) UserVerification(login, password string, action int) (string, error) {
	return m.UserVerificationFunc(login, password, action)
}

func (m *MockUserRepository) InsertUser(login, password string) (string, error) {
	return m.InsertUserFunc(login, password)
}

func (m *MockUserRepository) InsertExpresions(userID, expressionId, expression string, statusId int) error {
	return m.InsertExpressionsFunc(userID, expressionId, expression, statusId)
}

func TestAuthService_Registrate(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		mockVerify     func(login, password string, action int) (string, error)
		mockInsert     func(login, password string) (string, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name:   "Successful registration",
			method: http.MethodPost,
			body:   `{"login":"test","password":"pass"}`,
			mockVerify: func(_, _ string, _ int) (string, error) {
				return "", nil
			},
			mockInsert: func(_, _ string) (string, error) {
				return "123", nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid HTTP method",
			method:         http.MethodGet,
			body:           `{"login":"test","password":"pass"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:           "Empty login and password",
			method:         http.MethodPost,
			body:           `{"login":"","password":""}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "Verification error",
			method: http.MethodPost,
			body:   `{"login":"test","password":"pass"}`,
			mockVerify: func(_, _ string, _ int) (string, error) {
				return "", errors.New("verification error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockUserRepository{
				UserVerificationFunc: tt.mockVerify,
				InsertUserFunc:       tt.mockInsert,
			}
			service := NewAuthService(repo)

			req := httptest.NewRequest(tt.method, "/registrate", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			_, err, status := service.Registrate(w, req)

			if status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}

			if (err != nil) != tt.expectError {
				t.Errorf("expected error %v, got %v", tt.expectError, err != nil)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		mockVerify     func(login, password string, action int) (string, error)
		expectedStatus int
		expectError    bool
	}{
		{
			name:   "Successful login",
			method: http.MethodPost,
			body:   `{"login":"test","password":"pass"}`,
			mockVerify: func(_, _ string, _ int) (string, error) {
				return "123", nil
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:   "Login failed",
			method: http.MethodPost,
			body:   `{"login":"test","password":"wrong"}`,
			mockVerify: func(_, _ string, _ int) (string, error) {
				return "", errors.New("invalid credentials")
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Empty login and password",
			method:         http.MethodPost,
			body:           `{"login":"","password":""}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "Verification error",
			method: http.MethodPost,
			body:   `{"login":"test","password":"pass"}`,
			mockVerify: func(_, _ string, _ int) (string, error) {
				return "", errors.New("verification error")
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Invalid HTTP method",
			method:         http.MethodGet,
			body:           `{"login":"test","password":"pass"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockUserRepository{
				UserVerificationFunc: tt.mockVerify,
			}
			service := NewAuthService(repo)

			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			_, err, status := service.Login(w, req)

			if status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}

			if (err != nil) != tt.expectError {
				t.Errorf("expected error %v, got %v", tt.expectError, err != nil)
			}
		})
	}
}

func TestInsertExpresions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &SQLiteUserRepository{db: db}

	tests := []struct {
		name         string
		userID       string
		expressionId string
		expression   string
		statusId     int
		wantErr      bool
	}{
		{
			name:         "successful insert",
			userID:       "user1",
			expressionId: "expr1",
			expression:   "2+2",
			statusId:     1,
			wantErr:      false,
		},
		{
			name:         "empty userID",
			userID:       "",
			expressionId: "expr1",
			expression:   "2+2",
			statusId:     1,
			wantErr:      true,
		},
		{
			name:         "empty expressionID",
			userID:       "user1",
			expressionId: "",
			expression:   "2+2",
			statusId:     1,
			wantErr:      true,
		},
		{
			name:         "empty expression",
			userID:       "user1",
			expressionId: "expr1",
			expression:   "",
			statusId:     1,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTabls()
			if err != nil {
				t.Errorf("CreateTables() error = %v, want nil", err)
			}
			err = repo.InsertExpresions(tt.userID, tt.expressionId, tt.expression, tt.statusId)
			if tt.wantErr && err == nil {
				t.Errorf("InsertExpresions() want error, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("InsertExpresions() want no error, got %v", err)
			}
		})
	}
}

func TestUpdateExpressionResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &SQLiteUserRepository{db: db}

	tests := []struct {
		name         string
		expressionId string
		result       string
		statusId     int
		wantErr      bool
	}{
		{
			name:         "successful insert",
			expressionId: "expr1",
			result:       "4.45",
			statusId:     3,
			wantErr:      false,
		},
		{
			name:         "empty expressionID",
			expressionId: "",
			result:       "4.45",
			statusId:     3,
			wantErr:      true,
		},
		{
			name:         "empty result",
			expressionId: "expr1",
			result:       "",
			statusId:     3,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTabls()
			if err != nil {
				t.Errorf("CreateTables() error = %v, want nil", err)
			}
			err = repo.UpdateExpressionResult(tt.expressionId, tt.result, tt.statusId)
			if tt.wantErr && err == nil {
				t.Errorf("InsertExpresions() want error, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("InsertExpresions() want no error, got %v", err)
			}
		})
	}
}

func TestGetExpression(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &SQLiteUserRepository{db: db}

	tests := []struct {
		name         string
		expressionId string
		wantErr      bool
	}{
		{
			name:         "successful insert",
			expressionId: "expr1",
			wantErr:      false,
		},
		{
			name:         "empty expressionID",
			expressionId: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTabls()
			if err != nil {
				t.Errorf("CreateTables() error = %v, want nil", err)
			}

			if tt.expressionId != "" {
				if err = repo.InsertExpresions("123", tt.expressionId, "2+2", 1); err != nil {
					t.Errorf("InsertExpresions() error = %v, want nil", err)
				}
				if err = repo.UpdateExpressionResult(tt.expressionId, "4.000", 3); err != nil {
					t.Errorf("UpdateExpressionResult() error = %v, want nil", err)
				}
			}

			_, err = repo.GetExpression(tt.expressionId)
			if tt.wantErr && err == nil {
				t.Errorf("Want error, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Want no error, got %v", err)
			}
		})
	}
}

func TestGetAllExpression(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &SQLiteUserRepository{db: db}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "successful insert",
			userID:  "expr1",
			wantErr: false,
		},
		{
			name:    "empty userID",
			userID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTabls()
			if err != nil {
				t.Errorf("CreateTables() error = %v, want nil", err)
			}

			if tt.userID != "" {
				if err = repo.InsertExpresions("userID", "expr1", "2+2", 1); err != nil {
					t.Errorf("InsertExpresions() error = %v, want nil", err)
				}
				if err = repo.UpdateExpressionResult("expr1", "4.000", 3); err != nil {
					t.Errorf("UpdateExpressionResult() error = %v, want nil", err)
				}
			}

			_, err = repo.GetAllExpressions(tt.userID)
			if tt.wantErr && err == nil {
				t.Errorf("Want error, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Want no error, got %v", err)
			}
		})
	}
}
