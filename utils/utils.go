package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt password hashing
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing failed. Error: %v", err)
		return "", errors.New("Internal Server Error")
	}
	return string(hashedPassword), nil
}

// Check if input password is correct
func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

// Set response content to json and write any data that is sent
func EncodeJson(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("JSON encoding error")
		return errors.New("Internal Server Error")
	}
	return nil
}
