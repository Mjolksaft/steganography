package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("error generating password: %w", err)
	}
	return hashedPassword, nil
}
