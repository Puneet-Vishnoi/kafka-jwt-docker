package models

import (
	"time"

	"github.com/Puneet-Vishnoi/kafka-simple/database"
)

type User struct {
	UserID       string
	Email        string
	Password     string
	FirstName    string
	LastName     string
	UserType     string
	Token        string
	RefreshToken string
	UpdatedAt    time.Time
}

func UserExists(email string) (bool, error) {
	db := database.DB
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func GetUserByEmail(email string) (User, error) {
	db := database.DB

	var user User
	query := `
		SELECT user_id, email, password, first_name, last_name, user_type, token, refresh_token, updated_at
		FROM users WHERE email = $1
	`

	row := db.QueryRow(query, email)
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.UserType,
		&user.Token,
		&user.RefreshToken,
		&user.UpdatedAt,
	)

	return user, err
}

func CreateUser(user User) error {
	db := database.DB
	_, err := db.Exec(`
		INSERT INTO users (user_id, email, password, first_name, last_name, user_type, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.UserID, user.Email, user.Password, user.FirstName, user.LastName, user.UserType, time.Now())
	return err
}
