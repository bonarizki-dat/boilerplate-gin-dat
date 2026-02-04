package routers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/middlewares"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes (register, login, refresh, forgot-password, reset-password).
// Applies rate limiting via RateLimitMiddleware; limits read from RATE_LIMIT_RPS / RATE_LIMIT_BURST in middleware.
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authController := controllers.NewAuthController(authService)
	authRoutes := router.Group("/auth")
	authRoutes.Use(middlewares.RateLimitMiddleware())
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
		authRoutes.POST("/forgot-password", authController.ForgotPassword)
		authRoutes.POST("/reset-password", authController.ResetPassword)
	}
}
