package services

import (
	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/errors"
	"time"
)

// UserService handles user-related business logic
// 只依赖接口，便于 mock 和扩展
type UserService struct {
	repo  repositories.UserRepository
	cache CacheService
	auth  AuthService
}

// NewUserService creates a new UserService instance
func NewUserService(
	repo repositories.UserRepository,
	cache CacheService,
	auth AuthService,
) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
		auth:  auth,
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(user *models.User) error {
	// Hash password
	hashedPassword, err := s.auth.HashPassword(user.Password)
	if err != nil {
		return errors.WrapError(err, "Failed to hash password")
	}
	user.Password = hashedPassword
	user.Role = "user" // Default role

	// Check if username already exists
	existingUser, err := s.repo.FindByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.NewConflictError("Username already exists", nil)
	}

	// Check if email already exists
	existingUser, err = s.repo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.NewConflictError("Email already exists", nil)
	}

	if err := s.repo.Create(user); err != nil {
		return errors.WrapError(err, "Could not create user")
	}

	// Clear password before returning
	user.Password = ""
	return nil
}

// AuthenticateUser authenticates a user and returns JWT token
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.ErrUnauthorized
	}
	if err := s.auth.ComparePassword(user.Password, password); err != nil {
		return "", errors.ErrUnauthorized
	}
	// Generate JWT token
	token, err := s.auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.WrapError(err, "Could not generate token")
	}
	// Update last login time
	user.LastLogin = time.Now()
	s.repo.Update(user)
	return token, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	user.Password = "" // Never return password
	return user, nil
}

// ListUsers retrieves all users (admin only)
func (s *UserService) ListUsers() ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.WrapError(err, "Could not fetch users")
	}
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

// UpdateUser updates a user's email and role
func (s *UserService) UpdateUser(userID uint, newEmail, newRole string) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if newEmail != "" {
		user.Email = newEmail
	}
	if newRole != "" {
		user.Role = newRole
	}
	if err := s.repo.Update(user); err != nil {
		return nil, errors.WrapError(err, "Could not update user")
	}
	user.Password = ""
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(userID uint) error {
	if err := s.repo.Delete(userID); err != nil {
		return errors.WrapError(err, "Could not delete user")
	}
	return nil
}

// ResetPassword updates a user's password
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.ErrNotFound
	}
	hashed, err := s.auth.HashPassword(newPassword)
	if err != nil {
		return errors.WrapError(err, "Failed to hash password")
	}
	user.Password = hashed
	if err := s.repo.Update(user); err != nil {
		return errors.WrapError(err, "Could not update password")
	}
	return nil
}

// GetAPIKeyByPrefix retrieves an API key by its prefix
func (s *UserService) GetAPIKeyByPrefix(prefix string) (*models.APIKey, error) {
	// This would typically use an APIKeyRepository
	// For now, return a mock implementation
	return nil, errors.ErrNotFound
}

// ValidateAPIKey validates an API key
func (s *UserService) ValidateAPIKey(apiKey *models.APIKey, plainKey string) bool {
	// This would typically validate the API key hash
	// For now, return false
	return false
}

// CreateAPIKey creates a new API key for a user
func (s *UserService) CreateAPIKey(userID uint, name string) (*models.APIKey, error) {
	// This would typically use an APIKeyRepository
	// For now, return a mock implementation
	return nil, errors.NewError(500, "API key creation not implemented", nil)
}

// ListAPIKeys lists all API keys for a user
func (s *UserService) ListAPIKeys(userID uint) ([]models.APIKey, error) {
	// This would typically use an APIKeyRepository
	// For now, return a mock implementation
	return nil, errors.NewError(500, "API key listing not implemented", nil)
}

// RevokeAPIKey revokes an API key
func (s *UserService) RevokeAPIKey(userID uint, keyID uint) error {
	// This would typically use an APIKeyRepository
	// For now, return a mock implementation
	return errors.NewError(500, "API key revocation not implemented", nil)
}
