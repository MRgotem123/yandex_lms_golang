package workWithSQL

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
)

var DB *sql.DB

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OutExpression struct {
	ExpressionID string
	Expression   string
	Result       string
	StatusID     int
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]) + "666"
}

func InsertUser(login, password string) error {
	userId, err := generateRandomID(10)
	if err != nil {
		return err
	}
	hashed := hashPassword(password)

	_, err = DB.Exec(`INSERT INTO users (userId, login, password) VALUES (?, ?, ?)`, userId, login, hashed)
	if err != nil {
		return err
	}
	return nil
}

func generateRandomID(length int) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	randomID := hex.EncodeToString(bytes)

	return randomID, nil
}

func UserVerification(DB *sql.DB, login, password string, task int) (string, error) {
	hashedPassword := hashPassword(password)

	var Password string
	var userID string
	err := DB.QueryRow(`SELECT userId, password FROM users WHERE login = ?`, login).Scan(&userID, &Password)
	if err != nil {
		if err == sql.ErrNoRows {
			if task == 1 {
				// Можем регестрироватся
				return userID, nil
			} else {
				return "", fmt.Errorf("Пользователь не найден.")
			}
		}
		return "", fmt.Errorf("Ошибка поиска, попробуйте ещё-раз.")
	}
	if hashedPassword == Password {
		if task == 1 {
			return "", fmt.Errorf("Вы уже зарегестрированы!")
		} else {
			return userID, nil
		}

	}

	if task == 1 {
		return "", fmt.Errorf("такой login уже существует, напишите другой login.")
	} else {
		return "", fmt.Errorf("Неверный пароль.")
	}

}
