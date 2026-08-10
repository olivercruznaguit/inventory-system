package service

import (
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
