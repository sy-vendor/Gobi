package repositories

import (
	"gobi/internal/models"
	"gobi/pkg/errors"
	"time"

	"gorm.io/gorm"
)

// UserRepositoryImpl implements UserRepository
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.WrapError(err, "Could not create user")
	}
	return nil
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &user, nil
}

// FindAll finds all users
func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, errors.WrapError(err, "Could not fetch users")
	}
	return users, nil
}

// Update updates a user
func (r *UserRepositoryImpl) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.WrapError(err, "Could not update user")
	}
	return nil
}

// Delete deletes a user
func (r *UserRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return errors.WrapError(err, "Could not delete user")
	}
	return nil
}

// UpdateLastLogin updates user's last login time
func (r *UserRepositoryImpl) UpdateLastLogin(id uint) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", time.Now()).Error; err != nil {
		return errors.WrapError(err, "Could not update last login")
	}
	return nil
}
