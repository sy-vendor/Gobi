package services

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gobi/internal/models"
	"gobi/internal/repositories"
	"gobi/pkg/errors"
)

// UserService handles user-related business logic
// 只依赖接口，便于 mock 和扩展
type UserService struct {
	repo       repositories.UserRepository
	cache      CacheService
	auth       AuthService
	apiKeyRepo repositories.APIKeyRepository
}

// NewUserService creates a new UserService instance
func NewUserService(
	repo repositories.UserRepository,
	cache CacheService,
	auth AuthService,
	apiKeyRepo repositories.APIKeyRepository,
) *UserService {
	return &UserService{
		repo:       repo,
		cache:      cache,
		auth:       auth,
		apiKeyRepo: apiKeyRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(user *models.User) error {
	// Hash password
	hashedPassword, err := s.auth.HashPassword(user.Password)
	if err != nil {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeInternalServer,
			"Failed to hash password",
			err,
			errors.SeverityHigh,
			errors.CategorySystem,
		)
	}
	user.Password = hashedPassword
	user.Role = "user" // Default role

	// Check if username already exists
	existingUser, err := s.repo.FindByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeUserExists,
			"Username already exists",
			nil,
			errors.SeverityMedium,
			errors.CategoryBusiness,
		)
	}

	// Check if email already exists
	existingUser, err = s.repo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.NewErrorWithSeverity(
			errors.ErrCodeUserExists,
			"Email already exists",
			nil,
			errors.SeverityMedium,
			errors.CategoryBusiness,
		)
	}

	if err := s.repo.Create(user); err != nil {
		return errors.NewDatabaseError("Could not create user", err)
	}

	// Clear password before returning
	user.Password = ""
	return nil
}

// AuthenticateUser authenticates a user and returns JWT token
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.NewErrorWithSeverity(
			errors.ErrCodeInvalidCredentials,
			"Invalid username or password",
			err,
			errors.SeverityMedium,
			errors.CategoryAuth,
		)
	}
	if err := s.auth.ComparePassword(user.Password, password); err != nil {
		return "", errors.NewErrorWithSeverity(
			errors.ErrCodeInvalidCredentials,
			"Invalid username or password",
			err,
			errors.SeverityMedium,
			errors.CategoryAuth,
		)
	}
	// Generate JWT token
	token, err := s.auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.NewErrorWithSeverity(
			errors.ErrCodeInternalServer,
			"Could not generate token",
			err,
			errors.SeverityHigh,
			errors.CategorySystem,
		)
	}
	// Update last login time
	user.LastLogin = time.Now()
	if err := s.repo.Update(user); err != nil {
		// 记录更新失败但不影响登录
		errors.RecordError(errors.NewDatabaseError("Failed to update last login time", err))
	}
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

// CreateAPIKey creates a new API key for a user
func (s *UserService) CreateAPIKey(userID uint, name string, expiresAt *time.Time) (*models.APIKey, string, error) {
	// 生成安全的随机key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, "", errors.WrapError(err, "Failed to generate API key")
	}
	plainKey := base64.RawURLEncoding.EncodeToString(keyBytes)
	prefix := plainKey[:12]
	keyHash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.WrapError(err, "Failed to hash API key")
	}
	apiKey := &models.APIKey{
		UserID:    userID,
		Name:      name,
		KeyHash:   string(keyHash),
		Prefix:    prefix,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	if err := s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, "", errors.WrapError(err, "Could not create API key")
	}
	return apiKey, plainKey, nil
}

// ListAPIKeys lists all API keys for a user
func (s *UserService) ListAPIKeys(userID uint, isAdmin bool) ([]models.APIKey, error) {
	return s.apiKeyRepo.FindByUser(userID, isAdmin)
}

// RevokeAPIKey revokes an API key
func (s *UserService) RevokeAPIKey(userID uint, keyID uint, isAdmin bool) error {
	return s.apiKeyRepo.Revoke(keyID, userID, isAdmin)
}

// GetAPIKeyByPrefix retrieves an API key by its prefix
func (s *UserService) GetAPIKeyByPrefix(prefix string) (*models.APIKey, error) {
	return s.apiKeyRepo.FindByPrefix(prefix)
}

// ValidateAPIKey validates an API key
func (s *UserService) ValidateAPIKey(apiKey *models.APIKey, plainKey string) bool {
	if apiKey == nil || apiKey.Revoked {
		return false
	}
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(plainKey)) == nil
}
