package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gemdivk/Crowdfunding-system/internal/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID       int       `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func Register(user *User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	query := `
    INSERT INTO "User" (name, email, password_hash, role, is_verified, created_at, updated_at) 
    VALUES ($1, $2, $3, $4, FALSE, NOW(), NOW()) RETURNING user_id, created_at, updated_at
`
	err = db.DB.QueryRow(query, user.Name, user.Email, string(hashedPassword), user.Role).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func Authenticate(email, password string) (*User, error) {
	var user User
	var hashedPassword string

	// Query the database for the user by email
	err := db.DB.QueryRow(`SELECT user_id, name, email, password_hash, role, is_verified, created_at, updated_at 
        FROM "User" WHERE email = $1`, email).Scan(&user.UserID, &user.Name, &user.Email, &hashedPassword, &user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	user.Password = ""
	return &user, nil
}

func GetUserIDbyEmail(email string) (int, error) {
	var userID int
	query := `Select user_id from "User" where email = $1`
	err := db.DB.QueryRow(query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("No user was found by this email: %s", email)
		}
		return 0, err
	}
	return userID, nil

}
func VerifyUserEmail(userID int) error {
	query := `UPDATE "User" SET is_verified = TRUE WHERE user_id = $1`
	_, err := db.DB.Exec(query, userID)
	return err
}
