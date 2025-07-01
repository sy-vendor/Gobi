package repositories

import (
	"gobi/internal/models"
	"time"

	"gorm.io/gorm"
)

// UserRepositoryImpl implements UserRepository
type UserRepositoryImpl struct {
	*BaseRepository
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.BaseRepository.Create(user)
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, r.WrapDBError(err, "Could not find user by username")
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, r.WrapDBError(err, "Could not find user by email")
	}
	return &user, nil
}

// FindAll finds all users
func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, r.WrapDBError(err, "Could not fetch users")
	}
	return users, nil
}

// Update updates a user
func (r *UserRepositoryImpl) Update(user *models.User) error {
	return r.BaseRepository.Update(user)
}

// Delete deletes a user
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.User{}, id)
}

// UpdateLastLogin updates user's last login time
func (r *UserRepositoryImpl) UpdateLastLogin(id uint) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", time.Now()).Error; err != nil {
		return r.WrapDBError(err, "Could not update last login")
	}
	return nil
}
