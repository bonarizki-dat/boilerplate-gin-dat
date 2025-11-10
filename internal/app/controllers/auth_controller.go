package controllers

import (
	"errors"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	service *services.AuthService
}

// NewAuthController creates a new AuthController instance
func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

// Register handles user registration endpoint.
//
// POST /auth/register
// Request body: RegisterRequest (JSON)
// Response: AuthResponse with user info and JWT token
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid registration request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	// Call service
	response, err := ctrl.service.Register(&req)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.Conflict(c, err, "Email already exists")
			return
		}

		// Handle generic errors
		logger.Errorf("registration failed: %v", err)
		utils.InternalServerError(c, err, "Failed to register user")
		return
	}

	// Success response
	utils.Created(c, response, "User registered successfully")
}

// Login handles user authentication endpoint.
//
// POST /auth/login
// Request body: LoginRequest (JSON)
// Response: AuthResponse with user info and JWT token
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid login request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	// Call service
	response, err := ctrl.service.Login(&req)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.Unauthorized(c, err, "Invalid email or password")
			return
		}

		// Handle generic errors
		logger.Errorf("login failed: %v", err)
		utils.InternalServerError(c, err, "Failed to authenticate user")
		return
	}

	// Success response
	utils.Ok(c, response, "Login successful")
}

// RefreshToken handles token refresh endpoint.
//
// POST /auth/refresh
// Request body: RefreshTokenRequest (JSON)
// Response: RefreshTokenResponse with new access and refresh tokens
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid refresh token request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	// Call service
	response, err := ctrl.service.RefreshToken(&req)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			utils.Unauthorized(c, err, "Invalid or expired refresh token")
			return
		}

		// Handle generic errors
		logger.Errorf("token refresh failed: %v", err)
		utils.InternalServerError(c, err, "Failed to refresh token")
		return
	}

	// Success response
	utils.Ok(c, response, "Token refreshed successfully")
}

// ForgotPassword handles forgot password endpoint.
//
// POST /auth/forgot-password
// Request body: ForgotPasswordRequest (JSON)
// Response: Success message (token sent via email in production)
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid forgot password request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	// Call service
	resetToken, err := ctrl.service.ForgotPassword(&req)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrUserNotFound) {
			// Return success even if user not found (security best practice)
			// Don't reveal if email exists in system
			utils.Ok(c, nil, "If the email exists, a password reset link has been sent")
			return
		}

		// Handle generic errors
		logger.Errorf("forgot password failed: %v", err)
		utils.InternalServerError(c, err, "Failed to process request")
		return
	}

	// Success response
	// In production, don't return the token in response; send via email
	if config.IsProduction() {
		utils.Ok(c, map[string]string{
			"message": "Password reset instructions sent to email",
		}, "Password reset initiated")
		return
	}

	// Non-production: include token for development/testing convenience
	utils.Ok(c, map[string]string{
		"message": "Password reset instructions sent to email",
		"token":   resetToken,
	}, "Password reset initiated")
}

// ResetPassword handles password reset endpoint.
//
// POST /auth/reset-password
// Request body: ResetPasswordRequest (JSON)
// Response: Success message
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid reset password request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	// Call service
	err := ctrl.service.ResetPassword(&req)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrInvalidResetToken) {
			utils.BadRequest(c, err, "Invalid reset token")
			return
		}
		if errors.Is(err, services.ErrResetTokenExpired) {
			utils.BadRequest(c, err, "Reset token has expired")
			return
		}

		// Handle generic errors
		logger.Errorf("password reset failed: %v", err)
		utils.InternalServerError(c, err, "Failed to reset password")
		return
	}

	// Success response
	utils.Ok(c, nil, "Password reset successfully")
}
