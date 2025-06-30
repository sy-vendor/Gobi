package infrastructure

import (
	"gobi/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl implements AuthService
type AuthServiceImpl struct{}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{}
}

// GenerateJWT generates a JWT token
func (s *AuthServiceImpl) GenerateJWT(userID uint, role string) (string, error) {
	return utils.GenerateJWT(userID, role)
}

// ValidateJWT validates a JWT token
func (s *AuthServiceImpl) ValidateJWT(token string) (uint, string, error) {
	return utils.ValidateJWT(token)
}

// HashPassword hashes a password using bcrypt
func (s *AuthServiceImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a password with its hash
func (s *AuthServiceImpl) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
