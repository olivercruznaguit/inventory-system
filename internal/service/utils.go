package service

import (
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func isValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	parts := strings.Split(address.Address, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]

	return strings.Contains(domain, ".") &&
		!strings.HasSuffix(domain, ".")
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
