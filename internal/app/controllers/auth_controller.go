package controllers

import (
	"errors"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
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
