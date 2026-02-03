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
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "Register")
	defer span.Finish(err)

	var req dto.RegisterRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid registration request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.Conflict(c, err, "Email already exists")
			return
		}
		logger.Errorf("registration failed: %v", err)
		utils.InternalServerError(c, err, "Failed to register user")
		return
	}

	utils.Created(c, response, "User registered successfully")
}

// Login handles user authentication endpoint.
//
// POST /auth/login
// Request body: LoginRequest (JSON)
// Response: AuthResponse with user info and JWT token
func (ctrl *AuthController) Login(c *gin.Context) {
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "Login")
	defer span.Finish(err)

	var req dto.LoginRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid login request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.Unauthorized(c, err, "Invalid email or password")
			return
		}
		logger.Errorf("login failed: %v", err)
		utils.InternalServerError(c, err, "Failed to authenticate user")
		return
	}

	utils.Ok(c, response, "Login successful")
}

// RefreshToken handles token refresh endpoint.
//
// POST /auth/refresh
// Request body: RefreshTokenRequest (JSON)
// Response: RefreshTokenResponse with new access and refresh tokens
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "RefreshToken")
	defer span.Finish(err)

	var req dto.RefreshTokenRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid refresh token request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			utils.Unauthorized(c, err, "Invalid or expired refresh token")
			return
		}
		logger.Errorf("token refresh failed: %v", err)
		utils.InternalServerError(c, err, "Failed to refresh token")
		return
	}

	utils.Ok(c, response, "Token refreshed successfully")
}

// ForgotPassword handles forgot password endpoint.
//
// POST /auth/forgot-password
// Request body: ForgotPasswordRequest (JSON)
// Response: Success message (token sent via email in production)
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "ForgotPassword")
	defer span.Finish(err)

	var req dto.ForgotPasswordRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid forgot password request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	resetToken, err := ctrl.service.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.Ok(c, nil, "If the email exists, a password reset link has been sent")
			return
		}
		logger.Errorf("forgot password failed: %v", err)
		utils.InternalServerError(c, err, "Failed to process request")
		return
	}

	if config.IsProduction() {
		utils.Ok(c, map[string]string{
			"message": "Password reset instructions sent to email",
		}, "Password reset initiated")
		return
	}

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
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "ResetPassword")
	defer span.Finish(err)

	var req dto.ResetPasswordRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid reset password request: %v", err)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	err = ctrl.service.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidResetToken) {
			utils.BadRequest(c, err, "Invalid reset token")
			return
		}
		if errors.Is(err, services.ErrResetTokenExpired) {
			utils.BadRequest(c, err, "Reset token has expired")
			return
		}
		logger.Errorf("password reset failed: %v", err)
		utils.InternalServerError(c, err, "Failed to reset password")
		return
	}

	utils.Ok(c, nil, "Password reset successfully")
}

// Profile returns the current authenticated user's profile.
//
// GET /api/profile (requires JWT)
// Response: user_id from context in standard response format
func (ctrl *AuthController) Profile(c *gin.Context) {
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "AuthController", "Profile")
	defer span.Finish(err)

	userID := c.GetUint("user_id")
	utils.Ok(c, gin.H{"user_id": userID}, "Profile retrieved successfully")
}
