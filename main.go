package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	viper.SetDefault("SERVER_TIMEZONE", "UTC")
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
	// later separate migration
	migrations.Migrate()

	router := routers.SetupRoute()

	addr := config.ServerConfig()
	requestTimeout := 30 * time.Second
	if v := viper.GetInt("REQUEST_TIMEOUT_SECONDS"); v > 0 {
		requestTimeout = time.Duration(v) * time.Second
	}
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  requestTimeout,
		WriteTimeout: requestTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server listen error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Infof("shutdown signal received, shutting down gracefully")

	shutdownTimeout := 10 * time.Second
	if v := viper.GetInt("SERVER_SHUTDOWN_TIMEOUT"); v > 0 {
		shutdownTimeout = time.Duration(v) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown error: %v", err)
	}

	if sqlDB, err := database.GetDB().DB(); err == nil && sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("database close error: %v", err)
		}
	}

	logger.Infof("shutdown complete")
}
