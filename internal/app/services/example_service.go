package services

import (
	"context"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ExampleService handles example/datatables-related business logic.
type ExampleService struct{}

// NewExampleService creates a new ExampleService instance.
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// GetDataDatatables returns data in DataTables format for server-side processing.
func (s *ExampleService) GetDataDatatables(ctx context.Context, c *gin.Context) (data interface{}, err error) {
	ctx, start := logger.LogStart(ctx, "ExampleService.GetDataDatatables")

	data, err = repositories.GetDataDatatables(c)
	logger.LogFinish(ctx, "ExampleService.GetDataDatatables", err, start)
	return data, err
}
