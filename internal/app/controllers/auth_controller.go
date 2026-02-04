package controllers

import (
	"errors"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services/auth"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	service *auth.AuthService
}

// NewAuthController creates a new AuthController instance
func NewAuthController(service *auth.AuthService) *AuthController {
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
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Register")

	var req dto.RegisterRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid registration request: %v", err)
		logger.LogFinish(ctx, "AuthController.Register", err, start)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			logger.LogFinish(ctx, "AuthController.Register", err, start)
			utils.Conflict(c, err, "Email already exists")
			return
		}
		logger.Errorf("registration failed: %v", err)
		logger.LogFinish(ctx, "AuthController.Register", err, start)
		utils.InternalServerError(c, err, "Failed to register user")
		return
	}

	logger.LogFinish(ctx, "AuthController.Register", nil, start)
	utils.Created(c, response, "User registered successfully")
}

// Login handles user authentication endpoint.
//
// POST /auth/login
// Request body: LoginRequest (JSON)
// Response: AuthResponse with user info and JWT token
func (ctrl *AuthController) Login(c *gin.Context) {
	var err error
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Login")

	var req dto.LoginRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid login request: %v", err)
		logger.LogFinish(ctx, "AuthController.Login", err, start)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			logger.LogFinish(ctx, "AuthController.Login", err, start)
			utils.Unauthorized(c, err, "Invalid email or password")
			return
		}
		logger.Errorf("login failed: %v", err)
		logger.LogFinish(ctx, "AuthController.Login", err, start)
		utils.InternalServerError(c, err, "Failed to authenticate user")
		return
	}

	logger.LogFinish(ctx, "AuthController.Login", nil, start)
	utils.Ok(c, response, "Login successful")
}

// RefreshToken handles token refresh endpoint.
//
// POST /auth/refresh
// Request body: RefreshTokenRequest (JSON)
// Response: RefreshTokenResponse with new access and refresh tokens
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var err error
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.RefreshToken")

	var req dto.RefreshTokenRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid refresh token request: %v", err)
		logger.LogFinish(ctx, "AuthController.RefreshToken", err, start)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	response, err := ctrl.service.RefreshToken(ctx, &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			logger.LogFinish(ctx, "AuthController.RefreshToken", err, start)
			utils.Unauthorized(c, err, "Invalid or expired refresh token")
			return
		}
		logger.Errorf("token refresh failed: %v", err)
		logger.LogFinish(ctx, "AuthController.RefreshToken", err, start)
		utils.InternalServerError(c, err, "Failed to refresh token")
		return
	}

	logger.LogFinish(ctx, "AuthController.RefreshToken", nil, start)
	utils.Ok(c, response, "Token refreshed successfully")
}

// ForgotPassword handles forgot password endpoint.
//
// POST /auth/forgot-password
// Request body: ForgotPasswordRequest (JSON)
// Response: Success message (token sent via email in production)
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var err error
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.ForgotPassword")

	var req dto.ForgotPasswordRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid forgot password request: %v", err)
		logger.LogFinish(ctx, "AuthController.ForgotPassword", err, start)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	resetToken, err := ctrl.service.ForgotPassword(ctx, &req)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			logger.LogFinish(ctx, "AuthController.ForgotPassword", err, start)
			utils.Ok(c, nil, "If the email exists, a password reset link has been sent")
			return
		}
		logger.Errorf("forgot password failed: %v", err)
		logger.LogFinish(ctx, "AuthController.ForgotPassword", err, start)
		utils.InternalServerError(c, err, "Failed to process request")
		return
	}

	if config.IsProduction() {
		logger.LogFinish(ctx, "AuthController.ForgotPassword", nil, start)
		utils.Ok(c, map[string]string{
			"message": "Password reset instructions sent to email",
		}, "Password reset initiated")
		return
	}

	logger.LogFinish(ctx, "AuthController.ForgotPassword", nil, start)
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
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.ResetPassword")

	var req dto.ResetPasswordRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid reset password request: %v", err)
		logger.LogFinish(ctx, "AuthController.ResetPassword", err, start)
		utils.BadRequest(c, err, "Invalid request data")
		return
	}

	err = ctrl.service.ResetPassword(ctx, &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidResetToken) {
			logger.LogFinish(ctx, "AuthController.ResetPassword", err, start)
			utils.BadRequest(c, err, "Invalid reset token")
			return
		}
		if errors.Is(err, auth.ErrResetTokenExpired) {
			logger.LogFinish(ctx, "AuthController.ResetPassword", err, start)
			utils.BadRequest(c, err, "Reset token has expired")
			return
		}
		logger.Errorf("password reset failed: %v", err)
		logger.LogFinish(ctx, "AuthController.ResetPassword", err, start)
		utils.InternalServerError(c, err, "Failed to reset password")
		return
	}

	logger.LogFinish(ctx, "AuthController.ResetPassword", nil, start)
	utils.Ok(c, nil, "Password reset successfully")
}

// Profile returns the current authenticated user's profile.
//
// GET /api/profile (requires JWT)
// Response: user_id from context in standard response format
func (ctrl *AuthController) Profile(c *gin.Context) {
	ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Profile")

	userID := c.GetUint("user_id")
	logger.LogFinish(ctx, "AuthController.Profile", nil, start)
	utils.Ok(c, gin.H{"user_id": userID}, "Profile retrieved successfully")
}
