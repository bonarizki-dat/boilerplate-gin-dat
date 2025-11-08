package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Common errors for auth service
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthService handles authentication-related business logic
type AuthService struct {
	// Dependencies can be added here if needed
}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register creates a new user account with validation and password hashing.
//
// Returns ErrEmailAlreadyExists if email is already registered.
// Password is hashed using bcrypt before storage.
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	existingUser, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if existingUser != nil {
		logger.Warnf("registration attempt with existing email: %s", req.Email)
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	// Create user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save to database
	if err := repositories.CreateUser(user); err != nil {
		logger.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Infof("user registered successfully: %s", user.Email)

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.Errorf("failed to generate token: %v", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken: token,
		TokenType:   "Bearer",
	}

	return response, nil
}

// Login authenticates a user with email and password.
//
// Returns ErrInvalidCredentials if email or password is incorrect.
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		logger.Errorf("failed to get user: %v", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if user == nil {
		logger.Warnf("login attempt with non-existent email: %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		logger.Warnf("login attempt with invalid password: %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	logger.Infof("user logged in successfully: %s", user.Email)

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.Errorf("failed to generate token: %v", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken: token,
		TokenType:   "Bearer",
	}

	return response, nil
}

// ValidateToken validates a JWT token and returns the user ID.
//
// Returns error if token is invalid, expired, or malformed.
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	// Get JWT secret from config
	secret := viper.GetString("SECRET")
	if secret == "" {
		return 0, errors.New("JWT secret not configured")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid user_id in token")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}

// Private helper methods

// hashPassword hashes a plain text password using bcrypt.
func (s *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// verifyPassword compares a hashed password with a plain text password.
func (s *AuthService) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateToken creates a JWT token for authenticated user.
//
// Token contains user ID and email in claims.
// Expiry time is 24 hours from creation.
func (s *AuthService) generateToken(user *models.User) (string, error) {
	// Get JWT secret from config
	secret := viper.GetString("SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not configured")
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiry
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
