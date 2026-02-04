package repositories

import (
	"errors"
	"fmt"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"gorm.io/gorm"
)

// UserRepository defines data access for user entity (used by auth and others).
type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUserByRefreshToken(token string) (*models.User, error)
	GetUserByPasswordResetToken(token string) (*models.User, error)
}

// userRepo is the default implementation of UserRepository.
type userRepo struct{}

// NewUserRepository returns a new UserRepository implementation.
func NewUserRepository() UserRepository {
	return &userRepo{}
}

var defaultUserRepo UserRepository = NewUserRepository()

// CreateUser creates a new user in the database.
//
// Returns error if email already exists or database operation fails.
func CreateUser(user *models.User) error {
	return defaultUserRepo.CreateUser(user)
}

func (r *userRepo) CreateUser(user *models.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		logger.Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email address.
//
// Returns nil if user not found.
func GetUserByEmail(email string) (*models.User, error) {
	return defaultUserRepo.GetUserByEmail(email)
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID.
//
// Returns error if user not found.
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		logger.Errorf("failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(user *models.User) error {
	return defaultUserRepo.UpdateUser(user)
}

func (r *userRepo) UpdateUser(user *models.User) error {
	if err := database.DB.Save(user).Error; err != nil {
		logger.Errorf("failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user by ID.
func DeleteUser(id uint) error {
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		logger.Errorf("failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUserByRefreshToken retrieves a user by refresh token.
//
// Returns nil if user not found or token is invalid.
func GetUserByRefreshToken(token string) (*models.User, error) {
	return defaultUserRepo.GetUserByRefreshToken(token)
}

func (r *userRepo) GetUserByRefreshToken(token string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("refresh_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by refresh token: %v", err)
		return nil, fmt.Errorf("failed to get user by refresh token: %w", err)
	}
	return &user, nil
}

// GetUserByPasswordResetToken retrieves a user by password reset token.
//
// Returns nil if user not found or token is invalid.
func GetUserByPasswordResetToken(token string) (*models.User, error) {
	return defaultUserRepo.GetUserByPasswordResetToken(token)
}

func (r *userRepo) GetUserByPasswordResetToken(token string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by password reset token: %v", err)
		return nil, fmt.Errorf("failed to get user by password reset token: %w", err)
	}
	return &user, nil
}
