package routers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/middlewares"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	cfg := config.Get()
	debug := cfg != nil && cfg.Server.Debug
	allowedHosts := "0.0.0.0"
	if cfg != nil && cfg.Server.AllowedHosts != "" {
		allowedHosts = cfg.Server.AllowedHosts
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.SetTrustedProxies([]string{allowedHosts})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.RequestLogMiddleware())
	router.Use(middlewares.MetricsMiddleware())

	RegisterRoutes(router) //routes register

	return router
}
