package controllers

import (
	"github.com/bonarizki-dat/Datatables-Gin/datatables"
	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ExampleController handles example/datatables-related HTTP requests.
type ExampleController struct {
	service *services.ExampleService
}

// NewExampleController creates a new ExampleController instance.
func NewExampleController(service *services.ExampleService) *ExampleController {
	return &ExampleController{
		service: service,
	}
}

// GetData handles GET for raw example data (standard JSON response).
func (ctrl *ExampleController) GetData(c *gin.Context) {
	ctx, start := logger.LogStart(c.Request.Context(), "ExampleController.GetData")

	var example []*models.Example
	repositories.Get(&example)
	logger.LogFinish(ctx, "ExampleController.GetData", nil, start)
	utils.Ok(c, &example, "Data retrieved successfully")
}

// GetDataDatatables handles GET for DataTables server-side format.
func (ctrl *ExampleController) GetDataDatatables(c *gin.Context) {
	ctx, start := logger.LogStart(c.Request.Context(), "ExampleController.GetDataDatatables")

	data, err := ctrl.service.GetDataDatatables(ctx, c)
	if err != nil {
		logger.LogFinish(ctx, "ExampleController.GetDataDatatables", err, start)
		utils.InternalServerError(c, err, "Failed to retrieve data")
		return
	}
	logger.LogFinish(ctx, "ExampleController.GetDataDatatables", nil, start)
	datatables.JSON(c, data.(dto.Datatables))
}
