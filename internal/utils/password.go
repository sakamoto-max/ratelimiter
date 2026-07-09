package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	passInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt the password : %w", err)
	}

	return string(passInBytes), nil
}

func ComparePassword(password string, encryptedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password is not correct : %w", err)
	}

	return nil
}
