package controllers

import (
	"github.com/bonarizki-dat/Datatables-Gin/datatables"
	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/services"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetData(ctx *gin.Context) {
	var example []*models.Example
	repositories.Get(&example)
	utils.Ok(ctx, &example, "Data retrieved successfully")
}

func GetDataDatatables(ctx *gin.Context) {
	data, err := services.GetDataDatatables(ctx)

	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to retrieve data")
		return
	}

	datatables.JSON(ctx, data.(dto.Datatables))
}
