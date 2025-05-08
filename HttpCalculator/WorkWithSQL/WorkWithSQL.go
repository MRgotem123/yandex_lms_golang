package WorkWithSQL

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func CreateBD() (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}
	return DB, nil
}

func (r *SQLiteUserRepository) CreateTabls() error {
	var err error
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			userId TEXT PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS expressions (
			expressionId TEXT PRIMARY KEY,
			userId TEXT NOT NULL,
			expression TEXT NOT NULL,
			result TEXT,
			statusId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(userId),
			FOREIGN KEY (statusId) REFERENCES statuses(statusId)
		);`,
	}

	for _, q := range queries {
		if _, err = r.db.Exec(q); err != nil {
			log.Fatal("Failed to execute query:", err)
			return err
		}
	}

	fmt.Println("Database initialized successfully.")
	return nil
}

func (s *AuthService) Registrate(w http.ResponseWriter, r *http.Request) (string, error, int) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("only POST method allowed"), http.StatusMethodNotAllowed
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", fmt.Errorf("invalid JSON"), http.StatusBadRequest
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("login and password required"), http.StatusBadRequest
	}

	_, err := s.userRepo.UserVerification(req.Login, req.Password, 1)
	if err != nil {
		return "", fmt.Errorf("verification failed: %w", err), http.StatusInternalServerError
	}

	userID, err := s.userRepo.InsertUser(req.Login, req.Password)
	if err != nil {
		return "", fmt.Errorf("user creation failed: %w", err), http.StatusInternalServerError
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("userID: %s", userID)))
	return userID, nil, http.StatusOK
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) (string, error, int) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("only POST method allowed"), http.StatusMethodNotAllowed
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", fmt.Errorf("invalid JSON"), http.StatusBadRequest
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("login and password required"), http.StatusBadRequest
	}

	userID, err := s.userRepo.UserVerification(req.Login, req.Password, 2)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err), http.StatusUnauthorized
	}

	w.WriteHeader(http.StatusOK)
	return userID, nil, http.StatusOK
}

func (r *SQLiteUserRepository) InsertExpresions(userID, expressionId, expression string, statusId int) error {
	if userID == "" || expression == "" || expressionId == "" || statusId == 0 {
		return fmt.Errorf("Один из входных параметров пустой")
	}
	_, err := r.db.Exec(`INSERT INTO expressions (userId, expressionId, expression, statusId) VALUES (?, ?, ?, ?)`, userID, expressionId, expression, statusId)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteUserRepository) UpdateExpressionResult(expressionId, result string, statusId int) error {
	var err error
	if expressionId == "" {
		return fmt.Errorf("expressionId пустой")
	}

	switch {
	case result == "" && statusId != 0:
		_, err = r.db.Exec(`UPDATE expressions SET statusId = ? WHERE expressionId = ?`, statusId, expressionId)

	case result != "" && statusId == 0:
		_, err = r.db.Exec(`UPDATE expressions SET result = ? WHERE expressionId = ?`, result, expressionId)

	case result != "" && statusId != 0:
		_, err = r.db.Exec(`UPDATE expressions SET result = ?, statusId = ? WHERE expressionId = ?`, result, statusId, expressionId)

	default:
		return nil
	}

	return err
}

func (r *SQLiteUserRepository) GetExpression(expressionID string) (OutExpression, error) {
	var expression string
	var result string
	var statusId int
	err := r.db.QueryRow(`SELECT expression, result, statusId FROM expressions WHERE expressionId  = ?`, expressionID).Scan(&expression, &result, &statusId)
	if err != nil {
		if err == sql.ErrNoRows {
			return OutExpression{}, fmt.Errorf("expression не найден: %s", expressionID)
		}
		return OutExpression{}, err
	}

	return OutExpression{
		ExpressionID: expressionID,
		Expression:   expression,
		Result:       result,
		StatusID:     statusId,
	}, nil
}

func (r *SQLiteUserRepository) GetAllExpressions(userID string) ([]OutExpression, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID пустой")
	}

	rows, err := r.db.Query(`SELECT expressionId, expression, result, statusId FROM expressions WHERE userId = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе выражений пользователя: %v", err)
	}
	defer rows.Close()

	var expressions []OutExpression

	for rows.Next() {
		var expr OutExpression
		err = rows.Scan(&expr.ExpressionID, &expr.Expression, &expr.Result, &expr.StatusID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании выражения: %v", err)
		}
		expressions = append(expressions, expr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %v", err)
	}

	return expressions, nil
}

func (r *SQLiteUserRepository) FindSaymExpression(expression string) (string, error) {
	var result string
	err := r.db.QueryRow(`SELECT result FROM expressions WHERE expression = ?`, expression).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *SQLiteUserRepository) LenExpresions() (int, error) {
	if r == nil || r.db == nil {
		return 0, errors.New("repository is not initialized")
	}

	row := r.db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
