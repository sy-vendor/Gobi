package services

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new UserService instance
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(user *models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WrapError(err, "Failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.Role = "user" // Default role

	// Check if username or email already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		return errors.NewConflictError("User or email already exists", nil)
	} else if err != gorm.ErrRecordNotFound {
		return errors.WrapError(err, "Database error")
	}

	if err := s.db.Create(user).Error; err != nil {
		return errors.WrapError(err, "Could not create user")
	}

	// Clear password before returning
	user.Password = ""
	return nil
}

// AuthenticateUser authenticates a user and returns JWT token
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.ErrUnauthorized
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.WrapError(err, "Could not generate token")
	}

	// Update last login time
	user.LastLogin = time.Now()
	s.db.Save(&user)

	return token, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	user.Password = "" // Never return password
	return &user, nil
}

// ListUsers retrieves all users (admin only)
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch users")
	}

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// UpdateUser updates a user's email and role
func (s *UserService) UpdateUser(userID uint, newEmail, newRole string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.ErrNotFound
	}

	if newEmail != "" {
		user.Email = newEmail
	}
	if newRole != "" {
		user.Role = newRole
	}
	if err := s.db.Save(&user).Error; err != nil {
		return nil, errors.WrapError(err, "Could not update user")
	}

	user.Password = ""
	return &user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(userID uint) error {
	if err := s.db.Delete(&models.User{}, userID).Error; err != nil {
		return errors.WrapError(err, "Could not delete user")
	}
	return nil
}

// ResetPassword updates a user's password
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.ErrNotFound
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.WrapError(err, "Failed to hash password")
	}
	user.Password = string(hashed)
	return s.db.Save(&user).Error
}
