package services

import (
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/metrics"
)

// HealthService handles health check operations.
type HealthService struct {
	// Dependencies can be added here if needed
}

// NewHealthService creates a new HealthService instance.
func NewHealthService() *HealthService {
	return &HealthService{}
}

// CheckHealth performs health checks on application dependencies.
//
// Checks database connectivity and returns overall health status.
// Returns "healthy" only if all checks pass.
func (s *HealthService) CheckHealth() *dto.HealthResponse {
	checks := make(map[string]string)

	// Check database connectivity
	dbStatus := s.checkDatabase()
	checks["database"] = dbStatus

	// Determine overall status
	overallStatus := "healthy"
	if dbStatus != "ok" {
		overallStatus = "unhealthy"
	}

	return &dto.HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    checks,
		Uptime:    metrics.GetUptime(),
	}
}

// GetMetrics returns application metrics.
//
// Returns request counters and uptime statistics.
func (s *HealthService) GetMetrics() *dto.MetricsResponse {
	return &dto.MetricsResponse{
		TotalRequests:   metrics.GetTotalRequests(),
		SuccessRequests: metrics.GetSuccessRequests(),
		ErrorRequests:   metrics.GetErrorRequests(),
		UptimeSeconds:   metrics.GetUptime(),
		Timestamp:       time.Now(),
	}
}

// checkDatabase verifies database connectivity.
//
// Returns "ok" if database is reachable, "error" otherwise.
func (s *HealthService) checkDatabase() string {
	// Handle nil database (test environment)
	if database.DB == nil {
		logger.Warnf("health check: database not initialized")
		return "error"
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		logger.Errorf("health check: failed to get database instance: %v", err)
		return "error"
	}

	// Ping database with timeout
	if err := sqlDB.Ping(); err != nil {
		logger.Errorf("health check: database ping failed: %v", err)
		return "error"
	}

	return "ok"
}
