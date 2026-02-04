package routers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/middlewares"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services/auth"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds all routing list here and delegates to Register* per feature.
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		utils.NotFound(ctx, nil, "Route not found")
	})

	RegisterHealthRoutes(route)

	authService := auth.NewAuthService()
	exampleService := services.NewExampleService()

	RegisterAuthRoutes(route, authService)
	RegisterExampleRoutes(route, exampleService)

	// Protected routes (require authentication)
	authController := controllers.NewAuthController(authService)
	protectedRoutes := route.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware(authService))
	{
		protectedRoutes.GET("/profile", authController.Profile)
	}
}
