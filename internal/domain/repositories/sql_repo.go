package repositories

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
)

func Save(model interface{}) interface{} {
	err := database.DB.Create(model).Error
	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}
	return err
}

func Get(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}

func GetOne(model interface{}) interface{} {
	err := database.DB.Last(model).Error
	return err
}

func Update(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}
