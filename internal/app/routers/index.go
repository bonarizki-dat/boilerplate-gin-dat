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

	// Public routes
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.GET("/datatables", controllers.GetDataDatatables)

	// Initialize services
	authService := services.NewAuthService()

	// Initialize controllers
	authController := controllers.NewAuthController(authService)

	// Auth routes (public - no authentication required)
	// Apply rate limiting to prevent brute force attacks
	authRoutes := route.Group("/auth")
	authRoutes.Use(middlewares.RateLimitMiddleware())
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}

	// Protected routes (require authentication)
	protectedRoutes := route.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware(authService))
	{
		// Example protected endpoint - get current user profile
		protectedRoutes.GET("/profile", func(ctx *gin.Context) {
			userID := ctx.GetUint("user_id")
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Profile retrieved successfully",
				"data": gin.H{
					"user_id": userID,
				},
			})
		})

		// Add more protected routes here
		// protectedRoutes.GET("/users", controllers.GetUsers)
		// protectedRoutes.POST("/users", controllers.CreateUser)
	}

	// Add All route
	// TestRoutes(route)
}
