package config

import (
	"fmt"
	"strings"

	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
	"github.com/spf13/viper"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

var appConfig *Configuration

// Get returns the loaded configuration. Call after SetupConfig().
// Returns nil if config has not been loaded.
func Get() *Configuration {
	return appConfig
}

// SetupConfig loads and validates configuration from .env file
// and stores it in the global config (accessible via Get()).
func SetupConfig() error {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error reading config file, %s", err)
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := ValidateConfig(); err != nil {
		logger.Errorf("Config validation failed, %v", err)
		return fmt.Errorf("config validation failed: %w", err)
	}

	appConfig = &Configuration{
		Server: ServerConfiguration{
			Port:   viper.GetString("SERVER_PORT"),
			Secret: viper.GetString("SECRET"),
		},
		Database: DatabaseConfiguration{
			Dbname:   viper.GetString("MASTER_DB_NAME"),
			Username: viper.GetString("MASTER_DB_USER"),
			Password: viper.GetString("MASTER_DB_PASSWORD"),
			Host:     viper.GetString("MASTER_DB_HOST"),
			Port:     viper.GetString("MASTER_DB_PORT"),
			LogMode:  viper.GetBool("MASTER_DB_LOG_MODE"),
		},
	}
	if appConfig.Server.Port == "" {
		appConfig.Server.Port = "8000"
	}
	host := viper.GetString("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	// ServerConfig() still needs host for addr; store in struct for future use
	appConfig.Server.Host = host

	logger.Infof("Configuration loaded and validated successfully")
	return nil
}

// ValidateConfig validates that all required configuration values are present
// and not using insecure default values
func ValidateConfig() error {
	// Required configuration keys
	requiredKeys := []string{
		"SECRET",
		"JWT_SECRET",
		"SERVER_HOST",
		"SERVER_PORT",
		"MASTER_DB_NAME",
		"MASTER_DB_USER",
		"MASTER_DB_PASSWORD",
		"MASTER_DB_HOST",
		"MASTER_DB_PORT",
	}

	// Check all required keys are present
	var missingKeys []string
	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required config keys: %s", strings.Join(missingKeys, ", "))
	}

	// Validate SECRET is not using example value
	secret := viper.GetString("SECRET")
	if strings.Contains(strings.ToUpper(secret), "CHANGE-THIS") ||
		strings.Contains(strings.ToLower(secret), "your-secret-key") ||
		len(secret) < 32 {
		return fmt.Errorf("SECRET must be changed from example value and be at least 32 characters long. Generate with: openssl rand -base64 32")
	}

	// Validate JWT_SECRET is not using example value
	jwtSecret := viper.GetString("JWT_SECRET")
	if strings.Contains(strings.ToUpper(jwtSecret), "CHANGE-THIS") ||
		strings.Contains(strings.ToLower(jwtSecret), "your-jwt-secret") ||
		len(jwtSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be changed from example value and be at least 32 characters long. Generate with: openssl rand -base64 32")
	}

	// Validate port is numeric
	port := viper.GetString("SERVER_PORT")
	if port == "" {
		return fmt.Errorf("SERVER_PORT cannot be empty")
	}

	logger.Debugf("Config validation passed for all required keys")
	return nil
}

// GetString returns a string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an integer configuration value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// Environment helper functions

// GetEnvironment returns the current application environment
// Defaults to "development" if APP_ENV is not set
func GetEnvironment() string {
	env := viper.GetString("APP_ENV")
	if env == "" {
		return EnvDevelopment
	}
	return env
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return GetEnvironment() == EnvDevelopment
}

// IsStaging returns true if running in staging environment
func IsStaging() bool {
	return GetEnvironment() == EnvStaging
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	return GetEnvironment() == EnvProduction
}

// IsDebugEnabled returns true if debug mode is enabled
// Automatically returns true for development environment unless explicitly set to false
func IsDebugEnabled() bool {
	// If DEBUG is explicitly set, use that value
	if viper.IsSet("DEBUG") {
		return viper.GetBool("DEBUG")
	}

	// Otherwise, default to true for development, false for others
	return IsDevelopment()
}
