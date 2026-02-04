package routers

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/controllers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterExampleRoutes registers example/datatables public routes.
func RegisterExampleRoutes(router *gin.Engine, exampleService *services.ExampleService) {
	exampleController := controllers.NewExampleController(exampleService)
	router.GET("/datatables", exampleController.GetDataDatatables)
}
