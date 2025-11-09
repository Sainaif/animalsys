package security

import (
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength is the minimum password length
	MinPasswordLength = 8

	// DefaultCost is the default bcrypt cost
	DefaultCost = 12
)

// PasswordService handles password hashing and verification
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", errors.NewBadRequest("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", errors.Wrap(err, 500, "failed to hash password")
	}

	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength validates password strength
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return errors.NewBadRequest("password must be at least 8 characters")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	// Check for at least one lowercase letter
	hasLower := false
	// Check for at least one digit
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.NewBadRequest("password must contain at least one uppercase letter")
	}

	if !hasLower {
		return errors.NewBadRequest("password must contain at least one lowercase letter")
	}

	if !hasDigit {
		return errors.NewBadRequest("password must contain at least one digit")
	}

	return nil
}
