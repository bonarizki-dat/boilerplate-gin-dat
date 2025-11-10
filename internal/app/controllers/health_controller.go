package controllers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

// HealthController handles health check and metrics endpoints.
type HealthController struct {
	service *services.HealthService
}

// NewHealthController creates a new HealthController instance.
func NewHealthController(service *services.HealthService) *HealthController {
	return &HealthController{
		service: service,
	}
}

// Health performs application health check.
//
// GET /health
// Returns health status of the application and its dependencies.
func (ctrl *HealthController) Health(c *gin.Context) {
	response := ctrl.service.CheckHealth()

	// Return 503 if unhealthy, 200 if healthy
	if response.Status == "unhealthy" {
		c.JSON(503, gin.H{
			"success": false,
			"message": "Service is unhealthy",
			"data":    response,
			"errors":  nil,
		})
		return
	}

	utils.Ok(c, response, "Service is healthy")
}

// Metrics returns application metrics.
//
// GET /metrics
// Returns basic request counters and uptime statistics.
func (ctrl *HealthController) Metrics(c *gin.Context) {
	response := ctrl.service.GetMetrics()
	utils.Ok(c, response, "Metrics retrieved successfully")
}
