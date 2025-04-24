package workWithSQL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func CreateTabls() error {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
		return err
	}

	queries := []string{
		// Таблица пользователей с userId как TEXT
		`CREATE TABLE IF NOT EXISTS users (
			userId TEXT PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`,
		// Таблица выражений
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
		if _, err = DB.Exec(q); err != nil {
			log.Fatal("Failed to execute query:", err)
			return err
		}
	}

	fmt.Println("Database initialized successfully.")
	return nil
}

func Registrate(w http.ResponseWriter, r *http.Request) (string, error, int) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Принимается только метод POST"), http.StatusMethodNotAllowed
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return "", fmt.Errorf("invalid JSON"), http.StatusBadRequest
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("логин и пароль пустые!"), http.StatusBadRequest
	}

	userID, err := UserVerification(DB, req.Login, req.Password, 1)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}

	err = InsertUser(req.Login, req.Password)
	if err != nil {
		return "", fmt.Errorf("Ошибка создания нового пользователя: %v", err), http.StatusInternalServerError
	}

	return userID, nil, http.StatusOK
}

func Login(w http.ResponseWriter, r *http.Request) (string, error, int) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Принимается только метод POST"), http.StatusMethodNotAllowed
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return "", fmt.Errorf("invalid JSON"), http.StatusBadRequest
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("логин и пароль пустые!"), http.StatusBadRequest
	}

	userID, err := UserVerification(DB, req.Login, req.Password, 2)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}

	return userID, nil, http.StatusOK
}

func InsertExpresions(userID, expressionId, expression string, statusId int) error {
	_, err := DB.Exec(`INSERT INTO expressions (userId, expressionId, expression, statusId) VALUES (?, ?, ?, ?)`, userID, expressionId, expression, statusId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateExpressionResult(expressionId, result string, statusId int) error {
	var err error

	switch {
	case result == "" && statusId != 0:
		_, err = DB.Exec(`UPDATE expressions SET statusId = ? WHERE expressionId = ?`, statusId, expressionId)

	case result != "" && statusId == 0:
		_, err = DB.Exec(`UPDATE expressions SET result = ? WHERE expressionId = ?`, result, expressionId)

	case result != "" && statusId != 0:
		_, err = DB.Exec(`UPDATE expressions SET result = ?, statusId = ? WHERE expressionId = ?`, result, statusId, expressionId)

	default:
		return nil
	}

	return err
}

func GetExpression(expressionID string) (OutExpression, error) {
	var expression string
	var result string
	var statusId int
	err := DB.QueryRow(`SELECT expression, result, statusId FROM expressions WHERE expressionId  = ?`, expressionID).Scan(&expression, &result, &statusId)
	if err != nil {
		if err == sql.ErrNoRows {
			return OutExpression{}, fmt.Errorf("expression не найден")
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

func FindSaymExpression(expression string) (string, error) {
	var expressionId string
	err := DB.QueryRow(`SELECT expressionId FROM expressions WHERE expression  = ?`, expression).Scan(&expressionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("expression не найден")
		}
		return "", err
	}
	return expressionId, nil
}

func LenExpresions() (int, error) {
	row := DB.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
