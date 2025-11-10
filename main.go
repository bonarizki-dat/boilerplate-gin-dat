package main

import (
	"time"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database/migrations"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/app/routers"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/metrics"

	"github.com/spf13/viper"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	// Initialize metrics tracking
	metrics.Init()
	logger.Infof("Metrics tracking initialized")

	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
