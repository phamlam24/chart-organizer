package auth

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const cost = 10

func AddNewUser(db *sql.DB, username string, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}

	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	// Try to add the user to the table
	_, err = db.Exec("INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)", username, string(hashedPassword), currentTime)
	if err != nil {
		return err
	}

	return nil
}

func CheckUsernameAndPassword(db *sql.DB, username string, password string) (bool, error) {
	var passwordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		// Passwords don't match
		return false, nil
	}

	return true, nil
}
