package WorkWithSQL

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

type UserRepository interface {
	UserVerification(login, password string, action int) (string, error)
	InsertUser(login, password string) (string, error)
}

// AuthService обработчик аутентификации
type AuthService struct {
	userRepo UserRepository
}

// NewAuthService конструктор сервиса
func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

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

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]) + "666"
}

func (r *SQLiteUserRepository) InsertUser(login, password string) (string, error) {
	userId, err := generateRandomID(10)
	if err != nil {
		return "", err
	}
	hashed := hashPassword(password)

	_, err = r.db.Exec(`INSERT INTO users (userId, login, password) VALUES (?, ?, ?)`, userId, login, hashed)
	if err != nil {
		return "", err
	}
	return userId, nil
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

func (r *SQLiteUserRepository) UserVerification(login, password string, task int) (string, error) {
	hashedPassword := hashPassword(password)

	var Password string
	var userID string
	err := r.db.QueryRow(`SELECT userId, password FROM users WHERE login = ?`, login).Scan(&userID, &Password)
	if err != nil {
		if err == sql.ErrNoRows {
			if task == 1 {
				// Можем регестрироватся
				log.Println("Можем регестрироватся")
				return "", nil
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
