package auth

import (
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error("Failed to hash password", slog.String("err", err.Error()))
		return nil, fmt.Errorf("error generating password: %w", err)
	}
	return hashedPassword, nil
}

func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("Password check failed", slog.String("err", err.Error()))
		return err
	}
	return nil
}
