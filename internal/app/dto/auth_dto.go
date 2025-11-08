package dto

// RegisterRequest represents the payload for user registration.
//
// All fields are required and validated using go-playground/validator tags.
type RegisterRequest struct {
	// Name is the full name of the user
	Name string `json:"name" binding:"required,min=3,max=255"`

	// Email must be unique and valid format
	Email string `json:"email" binding:"required,email"`

	// Password must be at least 8 characters
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the payload for user authentication.
type LoginRequest struct {
	// Email of the registered user
	Email string `json:"email" binding:"required,email"`

	// Password for authentication
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response after successful authentication.
//
// Contains user information and JWT access token.
type AuthResponse struct {
	// User contains the authenticated user's basic information
	User UserResponse `json:"user"`

	// AccessToken is the JWT token for API authentication
	AccessToken string `json:"access_token"`

	// TokenType is always "Bearer" for JWT
	TokenType string `json:"token_type"`
}

// UserResponse represents user information in API responses.
//
// Password is never included in this response.
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
