package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
)

// RefreshToken generates new access and refresh tokens using a valid refresh token.
//
// Returns ErrInvalidRefreshToken if the refresh token is invalid or not found.
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (resp *dto.RefreshTokenResponse, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.RefreshToken")

	user, err := s.userRepo.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Errorf("failed to get user by refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if user == nil {
		logger.Warnf("refresh token attempt with invalid token")
		logger.LogFinish(ctx, "AuthService.RefreshToken", ErrInvalidRefreshToken, start)
		return nil, ErrInvalidRefreshToken
	}

	// Generate new access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		logger.Errorf("failed to generate new access token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		logger.Errorf("failed to generate new refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Update refresh token in database
	user.RefreshToken = newRefreshToken
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Errorf("failed to update refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	logger.Infof("refresh token successful for user: %s", user.Email)

	// Build response
	response := &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
	}

	logger.LogFinish(ctx, "AuthService.RefreshToken", nil, start)
	return response, nil
}

// ForgotPassword generates a password reset token and returns it.
//
// In production, this token should be sent via email instead of returned in response.
// Returns ErrUserNotFound if email doesn't exist.
func (s *AuthService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (token string, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.ForgotPassword")

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logger.Errorf("failed to get user by email: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to process request: %w", err)
	}

	if user == nil {
		logger.Warnf("password reset attempt for non-existent email: %s", req.Email)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", ErrUserNotFound, start)
		return "", ErrUserNotFound
	}

	// Generate password reset token
	resetToken, err := s.generatePasswordResetToken()
	if err != nil {
		logger.Errorf("failed to generate reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Set token expiry (15 minutes from now)
	expiry := time.Now().Add(15 * time.Minute)
	user.PasswordResetToken = resetToken
	user.PasswordResetExpiry = &expiry

	// Save to database
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Errorf("failed to save reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	logger.Infof("password reset token generated for user: %s", user.Email)

	if s.mailer != nil {
		if err := s.mailer.SendPasswordResetEmail(user.Email, resetToken); err != nil {
			logger.Errorf("failed to send password reset email: %v", err)
			logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
			return "", fmt.Errorf("failed to send reset email: %w", err)
		}
		logger.LogFinish(ctx, "AuthService.ForgotPassword", nil, start)
		return "", nil
	}
	logger.LogFinish(ctx, "AuthService.ForgotPassword", nil, start)
	return resetToken, nil
}

// ResetPassword resets user password using a valid reset token.
//
// Returns ErrInvalidResetToken if token is invalid.
// Returns ErrResetTokenExpired if token has expired.
func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.ResetPassword")

	user, err := s.userRepo.GetUserByPasswordResetToken(req.Token)
	if err != nil {
		logger.Errorf("failed to get user by reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to validate reset token: %w", err)
	}

	if user == nil {
		logger.Warnf("password reset attempt with invalid token")
		logger.LogFinish(ctx, "AuthService.ResetPassword", ErrInvalidResetToken, start)
		return ErrInvalidResetToken
	}

	// Check if token has expired
	if user.PasswordResetExpiry == nil || time.Now().After(*user.PasswordResetExpiry) {
		logger.Warnf("password reset attempt with expired token for user: %s", user.Email)
		logger.LogFinish(ctx, "AuthService.ResetPassword", ErrResetTokenExpired, start)
		return ErrResetTokenExpired
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("failed to hash new password: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to process password: %w", err)
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.PasswordResetExpiry = nil

	// Save to database
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Errorf("failed to update password: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to update password: %w", err)
	}

	logger.Infof("password reset successful for user: %s", user.Email)

	logger.LogFinish(ctx, "AuthService.ResetPassword", nil, start)
	return nil
}

// generatePasswordResetToken creates a cryptographically secure random reset token.
//
// Returns a 64-character hexadecimal string.
func (s *AuthService) generatePasswordResetToken() (string, error) {
	// Generate 32 random bytes (will be 64 chars in hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
