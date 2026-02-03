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

	// Set timezone (must be valid IANA name, e.g. UTC, Asia/Jakarta)
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, err := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	if err != nil {
		logger.Fatalf("invalid SERVER_TIMEZONE %q: %v", viper.GetString("SERVER_TIMEZONE"), err)
	}
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
	if err := router.Run(config.ServerConfig()); err != nil {
		logger.Fatalf("server run error: %v", err)
	}
}
