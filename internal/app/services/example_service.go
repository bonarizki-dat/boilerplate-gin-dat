package services

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
	"github.com/gin-gonic/gin"
)

// ExampleService handles example/datatables-related business logic.
type ExampleService struct{}

// NewExampleService creates a new ExampleService instance.
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// GetDataDatatables returns data in DataTables format for server-side processing.
func (s *ExampleService) GetDataDatatables(c *gin.Context) (interface{}, error) {
	return repositories.GetDataDatatables(c)
}
