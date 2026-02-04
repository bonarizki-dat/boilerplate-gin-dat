package services

import (
	"context"
	"fmt"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
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

// GetData returns all example records. Used for simple list responses.
func (s *ExampleService) GetData(ctx context.Context) ([]*models.Example, error) {
	ctx, start := logger.LogStart(ctx, "ExampleService.GetData")

	var list []*models.Example
	ret := repositories.Get(&list)
	if ret != nil {
		if err, ok := ret.(error); ok {
			logger.LogFinish(ctx, "ExampleService.GetData", err, start)
			return nil, err
		}
		err := fmt.Errorf("unexpected Get result")
		logger.LogFinish(ctx, "ExampleService.GetData", err, start)
		return nil, err
	}
	logger.LogFinish(ctx, "ExampleService.GetData", nil, start)
	return list, nil
}

// GetDataDatatables returns data in DataTables format for server-side processing.
//
// Expects DataTables request query params from gin.Context; returns datatables result or error.
func (s *ExampleService) GetDataDatatables(ctx context.Context, c *gin.Context) (data interface{}, err error) {
	ctx, start := logger.LogStart(ctx, "ExampleService.GetDataDatatables")

	data, err = repositories.GetDataDatatables(c)
	logger.LogFinish(ctx, "ExampleService.GetDataDatatables", err, start)
	return data, err
}
