package avtorisation

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
		// Таблица статусов
		`CREATE TABLE IF NOT EXISTS statuses (
			statusId INTEGER PRIMARY KEY,
			status TEXT NOT NULL UNIQUE
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

func Registrate(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Принимается только метод POST")
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return "", fmt.Errorf("invalid JSON")
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("логин и пароль пустые!")
	}

	userID, err := UserVerification(DB, req.Login, req.Password, 1)
	if err != nil {
		return "", err
	}

	err = InsertUser(req.Login, req.Password)
	if err != nil {
		return "", fmt.Errorf("Ошибка создания нового пользователя: %v", err)
	}

	return userID, nil
}

func Login(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Принимается только метод POST")
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return "", fmt.Errorf("invalid JSON")
	}

	if req.Login == "" || req.Password == "" {
		return "", fmt.Errorf("логин и пароль пустые!")
	}

	userID, err := UserVerification(DB, req.Login, req.Password, 2)
	if err != nil {
		return "", err
	}

	return userID, nil
}
