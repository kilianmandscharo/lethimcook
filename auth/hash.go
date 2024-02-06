package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	if len(passwordBytes) > 72 {
		return "", errors.New("the maximum password length is 72 bytes")
	}

	bytes, err := bcrypt.GenerateFromPassword(passwordBytes, 13)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func validatePassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
