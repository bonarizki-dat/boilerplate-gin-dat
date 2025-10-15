package services

import(
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

func GetDataDatatables(c *gin.Context) (interface{},error){
	return repositories.GetDataDatatables(c)
}