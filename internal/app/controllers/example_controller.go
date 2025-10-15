package controllers

import (
	"net/http"

	"github.com/bonarizki-dat/Datatables-Gin/datatables"
	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

func GetData(ctx *gin.Context) {
	var example []*models.Example
	repositories.Get(&example)
	ctx.JSON(http.StatusOK, &example)

}

func GetDataDatatables(ctx *gin.Context) {
	data, err := services.GetDataDatatables(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, 
			"message": "oppsss . . .",
		})
	}
	
	datatables.JSON(ctx, data.(dto.Datatables))
}
