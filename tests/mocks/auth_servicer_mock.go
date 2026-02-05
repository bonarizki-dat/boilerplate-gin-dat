package mocks

import (
	"context"
	"errors"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services/auth"
)

// MockAuthServicer is a test double for auth.AuthServicer.
// Set fields to control behaviour; nil means return zero value and nil error.
var _ auth.AuthServicer = (*MockAuthServicer)(nil)

// MockAuthServicer implements auth.AuthServicer for controller tests.
type MockAuthServicer struct {
	RegisterFunc       func(context.Context, *dto.RegisterRequest) (*dto.AuthResponse, error)
	LoginFunc          func(context.Context, *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshTokenFunc   func(context.Context, *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	ForgotPasswordFunc func(context.Context, *dto.ForgotPasswordRequest) (string, error)
	ResetPasswordFunc  func(context.Context, *dto.ResetPasswordRequest) error
	ValidateTokenFunc  func(string) (uint, error)
}

func (m *MockAuthServicer) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (string, error) {
	if m.ForgotPasswordFunc != nil {
		return m.ForgotPasswordFunc(ctx, req)
	}
	return "", nil
}

func (m *MockAuthServicer) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(ctx, req)
	}
	return nil
}

func (m *MockAuthServicer) ValidateToken(token string) (uint, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return 0, errors.New("invalid token")
}
