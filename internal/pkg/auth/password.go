package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	cost int
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost, // 10
	}
}

func (s *PasswordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *PasswordService) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *PasswordService) ValidateStrength(password string) []string {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	
	// Add more validation rules as needed
	// e.g., require uppercase, lowercase, numbers, special characters
	
	return errors
}
