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
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "ExampleController", "GetData")
	defer span.Finish(err)

	var example []*models.Example
	repositories.Get(&example)
	utils.Ok(c, &example, "Data retrieved successfully")
}

// GetDataDatatables handles GET for DataTables server-side format.
func (ctrl *ExampleController) GetDataDatatables(c *gin.Context) {
	var err error
	requestID := c.GetString("request_id")
	span := logger.StartWithRequestID(requestID, "ExampleController", "GetDataDatatables")
	defer span.Finish(err)

	data, err := ctrl.service.GetDataDatatables(c.Request.Context(), c)
	if err != nil {
		utils.InternalServerError(c, err, "Failed to retrieve data")
		return
	}
	datatables.JSON(c, data.(dto.Datatables))
}
