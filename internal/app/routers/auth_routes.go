package routers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/middlewares"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services/auth"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes under the given group (e.g. /api/v1).
// Full path will be groupPrefix/auth/... Applies rate limiting via RateLimitMiddleware.
func RegisterAuthRoutes(group *gin.RouterGroup, authService *auth.AuthService) {
	authController := controllers.NewAuthController(authService)
	authRoutes := group.Group("/auth")
	authRoutes.Use(middlewares.RateLimitMiddleware())
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
		authRoutes.POST("/forgot-password", authController.ForgotPassword)
		authRoutes.POST("/reset-password", authController.ResetPassword)
	}
}
