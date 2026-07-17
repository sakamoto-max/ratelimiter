package utils

import (
	"fmt"

	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	passInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", myErr.WrapErr(fmt.Errorf("failed to encrypt password : %w", err), myErr.InternalServerErr)
	}

	return string(passInBytes), nil
}

func ComparePassword(password string, encryptedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return myErr.WrapErr(fmt.Errorf("password is not correct : %w", err), myErr.UnauthorizedErr)
	}

	return nil
}
