package routers

import (
	"net/http"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/middlewares"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	// Health check and metrics routes
	RegisterHealthRoutes(route)

	// Initialize services
	authService := services.NewAuthService()
	exampleService := services.NewExampleService()

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	exampleController := controllers.NewExampleController(exampleService)

	// Public routes
	route.GET("/datatables", exampleController.GetDataDatatables)

	// Auth routes (public - no authentication required)
	// Apply rate limiting to prevent brute force attacks
	authRoutes := route.Group("/auth")
	authRoutes.Use(middlewares.RateLimitMiddleware())
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
		authRoutes.POST("/forgot-password", authController.ForgotPassword)
		authRoutes.POST("/reset-password", authController.ResetPassword)
	}

	// Protected routes (require authentication)
	protectedRoutes := route.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware(authService))
	{
		// Example protected endpoint - get current user profile
		protectedRoutes.GET("/profile", authController.Profile)

		// Add more protected routes here
		// protectedRoutes.GET("/users", controllers.GetUsers)
		// protectedRoutes.POST("/users", controllers.CreateUser)
	}

	// Add All route
	// TestRoutes(route)
}
