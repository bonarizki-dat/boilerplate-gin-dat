# Configuration Management Guide

This document explains how to configure the application for different environments.

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Environment Detection (APP_ENV)](#environment-detection-app_env)
3. [Environment Variables](#environment-variables)
4. [Configuration Files](#configuration-files)
5. [Validation](#validation)
6. [Security Best Practices](#security-best-practices)
7. [Environment-Specific Setup](#environment-specific-setup)
8. [Environment-Based Features](#environment-based-features)
9. [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Start

### 1. Copy Example Configuration

```bash
cp .env.example .env
```

### 2. Generate Secure Secrets

```bash
# Generate SECRET
openssl rand -base64 32

# Generate JWT_SECRET
openssl rand -base64 32
```

### 3. Update .env File

```env
# Replace with generated secrets
SECRET=<your-generated-secret-here>
JWT_SECRET=<your-generated-jwt-secret-here>

# Update database credentials
MASTER_DB_NAME=your_actual_db_name
MASTER_DB_USER=your_actual_db_user
MASTER_DB_PASSWORD=your_actual_db_password
MASTER_DB_HOST=localhost
```

### 4. Run the Application

```bash
go run main.go
```

The application will validate all required configuration on startup and **fail fast** if any required values are missing or using insecure defaults.

---

## 🌍 Environment Detection (APP_ENV)

The application uses `APP_ENV` to determine which environment it's running in. This controls feature flags, logging behavior, and other environment-specific configurations.

### Supported Environments

| Environment | APP_ENV Value | Use Case |
|-------------|---------------|----------|
| **Development** | `development` | Local development, debugging |
| **Staging** | `staging` | Pre-production testing |
| **Production** | `production` | Live production environment |

### Setting Environment

```env
# .env file
APP_ENV=development  # or staging, or production
```

If `APP_ENV` is not set, it defaults to `development`.

### Environment Helper Functions

Use these helper functions in your code to conditionally execute logic:

```go
import "github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"

// Get current environment
env := config.GetEnvironment()  // Returns: "development", "staging", or "production"

// Check specific environment
if config.IsDevelopment() {
    // Only runs in development
    logger.Debugf("Detailed debug info here")
}

if config.IsProduction() {
    // Only runs in production
    initMetrics()
    enableStrictSecurity()
}

if config.IsStaging() {
    // Only runs in staging
    enableTestingFeatures()
}

// Check debug mode (auto-enabled in development)
if config.IsDebugEnabled() {
    // Debug mode active
}
```

### Real-World Examples

**Example 1: Enable SQL Logging Only in Development**
```go
// internal/adapters/database/database.go
func DbConnection(masterDSN, replicaDSN string) error {
    logMode := config.IsDevelopment() // true in dev, false in prod

    masterDB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        // Only log SQL in development
        DisableForeignKeyConstraintWhenMigrating: config.IsDevelopment(),
    })

    if config.IsDevelopment() {
        masterDB = masterDB.Debug() // Enable SQL logging
    }

    return nil
}
```

**Example 2: Seed Test Data Only in Development**
```go
// migrations/migrate.go
func Migrate() {
    // Run migrations
    database.MasterDB.AutoMigrate(&models.User{})

    // Seed test data only in development
    if config.IsDevelopment() {
        seedTestData()
    }
}

func seedTestData() {
    logger.Infof("Seeding test data for development...")
    // Create test users, etc.
}
```

**Example 3: Different Rate Limits Per Environment**
```go
// internal/app/routers/index.go
func RegisterAuthRoutes(router *gin.Engine) {
    var rateLimit int

    if config.IsProduction() {
        rateLimit = 10 // Strict limit in production
    } else {
        rateLimit = 1000 // Relaxed for development/testing
    }

    authRoutes := router.Group("/auth")
    authRoutes.Use(middlewares.RateLimitMiddlewareWithConfig(rateLimit, rateLimit*2))
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
    }
}
```

**Example 4: Enable Profiling in Non-Production**
```go
// cmd/main.go or middleware
if !config.IsProduction() {
    // Enable pprof profiling endpoints
    router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
}
```

**Example 5: Environment-Specific Error Messages**
```go
// pkg/utils/response.go
func InternalServerError(c *gin.Context, err error, message string) {
    if config.IsDevelopment() {
        // Show detailed error in development
        HandleErrors(c, http.StatusInternalServerError, err, fmt.Sprintf("%s: %v", message, err))
    } else {
        // Hide details in production
        HandleErrors(c, http.StatusInternalServerError, nil, message)
    }
}
```

### Environment Constants

Available constants for comparison:

```go
import "github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"

config.EnvDevelopment  // "development"
config.EnvStaging      // "staging"
config.EnvProduction   // "production"

// Usage
if config.GetEnvironment() == config.EnvProduction {
    // Production-specific code
}
```

---

## 🔑 Environment Variables

### Required Variables

The following environment variables **MUST** be set. The application will not start if any are missing:

| Variable | Description | Example | Notes |
|----------|-------------|---------|-------|
| `SECRET` | Application secret key | Generated via openssl | Min 32 chars, MUST be changed from example |
| `JWT_SECRET` | JWT signing secret | Generated via openssl | Min 32 chars, MUST be changed from example |
| `SERVER_HOST` | Server bind address | `0.0.0.0` | Use `0.0.0.0` to bind all interfaces |
| `SERVER_PORT` | Server port | `8000` | Any available port |
| `MASTER_DB_NAME` | Master database name | `my_app_db` | Primary database for writes |
| `MASTER_DB_USER` | Master database user | `postgres` | User with write permissions |
| `MASTER_DB_PASSWORD` | Master database password | `secure_password` | Strong password required |
| `MASTER_DB_HOST` | Master database host | `localhost` | Database server address |
| `MASTER_DB_PORT` | Master database port | `5432` | PostgreSQL default: 5432 |

### Optional Variables

| Variable | Description | Default | Notes |
|----------|-------------|---------|-------|
| `APP_ENV` | Application environment | `development` | Values: `development`, `staging`, `production` |
| `DEBUG` | Debug mode | Auto (true in dev) | Set to `True` only in development |
| `ALLOWED_HOSTS` | Allowed hosts | `0.0.0.0` | Comma-separated list |
| `SERVER_TIMEZONE` | Server timezone | `Asia/Dhaka` | IANA timezone name |
| `MASTER_DB_LOG_MODE` | Enable DB query logging | `True` | Set to `False` in production |
| `MASTER_SSL_MODE` | Database SSL mode | `disable` | Use `require` in production |

### Replica Database (Optional)

For read scaling with database replicas:

| Variable | Description | Notes |
|----------|-------------|-------|
| `REPLICA_DB_NAME` | Replica database name | Can be same as master in development |
| `REPLICA_DB_USER` | Replica database user | Read-only user recommended |
| `REPLICA_DB_PASSWORD` | Replica database password | |
| `REPLICA_DB_HOST` | Replica database host | |
| `REPLICA_DB_PORT` | Replica database port | |
| `REPLICA_SSL_MODE` | Replica SSL mode | Use `require` in production |

---

## 📁 Configuration Files

### .env (Main Configuration)

- **Location:** Root directory
- **Purpose:** Main configuration file loaded on startup
- **Git:** `.gitignore`d - NEVER commit this file
- **Contains:** Actual secrets and environment-specific values

### .env.example (Template)

- **Location:** Root directory
- **Purpose:** Template showing all required variables
- **Git:** Committed to repository
- **Contains:** Placeholder values and documentation

### config/ Package

Configuration is managed through the `pkg/config` package:

```go
// config.go - Main config loader with validation
func SetupConfig() error

// db.go - Database configuration
func DbConfiguration() (masterDSN, replicaDSN string)

// server.go - Server configuration
func ServerConfig() string
```

---

## ✅ Validation

### Automatic Validation on Startup

The application automatically validates configuration when starting:

```go
// main.go
if err := config.SetupConfig(); err != nil {
    logger.Fatalf("config SetupConfig() error: %s", err)
}
```

### What Gets Validated

1. **Required Keys:** All required environment variables must be present
2. **Secret Security:** Secrets cannot use example/default values
3. **Secret Length:** Secrets must be at least 32 characters
4. **Value Presence:** No empty values for required fields

### Validation Error Examples

```bash
# Missing required key
Error: missing required config keys: JWT_SECRET, MASTER_DB_PASSWORD

# Using example value
Error: SECRET must be changed from example value and be at least 32 characters long.
Generate with: openssl rand -base64 32

# Empty value
Error: SERVER_PORT cannot be empty
```

---

## 🔒 Security Best Practices

### 1. Generate Strong Secrets

**❌ DON'T:**
```env
SECRET=12345
JWT_SECRET=mysecret
```

**✅ DO:**
```bash
# Generate cryptographically secure secrets
openssl rand -base64 32
# Output: K7gNU3sdo+OL0wNhqoVWhr3g6s1xYv72ol/pe/Unols=
```

### 2. Never Commit Secrets

**❌ DON'T:**
- Commit `.env` file to git
- Put real secrets in `.env.example`
- Hardcode secrets in source code

**✅ DO:**
- Keep `.env` in `.gitignore`
- Use placeholder values in `.env.example`
- Load secrets from environment variables

### 3. Use Different Secrets Per Environment

**❌ DON'T:**
```env
# Same secret everywhere
Development: SECRET=abc123
Production:  SECRET=abc123  # BAD!
```

**✅ DO:**
```env
# Different secrets per environment
Development: SECRET=dev_secret_K7gNU3sdo+OL0wNhqoVWhr3g
Production:  SECRET=prod_secret_X9mPQ7tzo+RL2xOisqXYks4h
```

### 4. Rotate Secrets Regularly

**Best Practice:**
- Rotate secrets every 90 days
- Rotate immediately if compromised
- Keep old secrets for brief transition period

### 5. Secure Secret Storage

**For Production:**
- **AWS:** Use AWS Secrets Manager or Parameter Store
- **GCP:** Use Google Secret Manager
- **Azure:** Use Azure Key Vault
- **Kubernetes:** Use Kubernetes Secrets
- **HashiCorp:** Use Vault

**Example with AWS Secrets Manager:**
```bash
# Store secret
aws secretsmanager create-secret \
    --name prod/myapp/jwt-secret \
    --secret-string "your-secret-here"

# Retrieve in application
# (requires AWS SDK integration)
```

---

## 🌍 Environment-Specific Setup

### Development Environment

```env
# .env.development (if implementing multi-env support)
DEBUG=True
SERVER_HOST=localhost
SERVER_PORT=8000

MASTER_DB_HOST=localhost
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_dev
MASTER_DB_LOG_MODE=True
MASTER_SSL_MODE=disable

# Use same DB for replica in dev
REPLICA_DB_HOST=localhost
REPLICA_DB_NAME=myapp_dev
```

### Staging Environment

```env
# .env.staging
DEBUG=False
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

MASTER_DB_HOST=staging-db.internal.company.com
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_staging
MASTER_DB_LOG_MODE=False
MASTER_SSL_MODE=require

REPLICA_DB_HOST=staging-db-replica.internal.company.com
```

### Production Environment

```env
# .env.production
DEBUG=False
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

MASTER_DB_HOST=prod-db.internal.company.com
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_production
MASTER_DB_LOG_MODE=False
MASTER_SSL_MODE=require

REPLICA_DB_HOST=prod-db-replica.internal.company.com
```

### Loading Environment-Specific Config (Future Enhancement)

```go
// Example: Load based on APP_ENV
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

configFile := fmt.Sprintf(".env.%s", env)
viper.SetConfigFile(configFile)
```

---

## 🐛 Troubleshooting

### Application Won't Start

**Error:** `missing required config keys: JWT_SECRET`

**Solution:**
1. Ensure `.env` file exists
2. Check all required variables are present
3. Verify no empty values

```bash
# Check .env file
cat .env | grep JWT_SECRET
```

---

### "SECRET must be changed from example value"

**Error:** `SECRET must be changed from example value and be at least 32 characters long`

**Solution:**
```bash
# Generate new secret
openssl rand -base64 32

# Update .env
SECRET=<paste-generated-secret-here>
```

---

### Database Connection Failed

**Error:** `database DbConnection error: connection refused`

**Possible Causes:**
1. Wrong `MASTER_DB_HOST` or `MASTER_DB_PORT`
2. Database not running
3. Incorrect credentials
4. Firewall blocking connection

**Solution:**
```bash
# Test database connection
psql -h localhost -U your_user -d your_db

# Check if PostgreSQL is running
pg_isready -h localhost -p 5432
```

---

### "JWT secret not configured"

**Error:** `JWT secret not configured`

**Solution:**
Ensure `JWT_SECRET` is set in `.env`:
```env
JWT_SECRET=your-generated-jwt-secret-min-32-chars
```

---

## 8️⃣ Environment-Based Features

Use environment detection to conditionally enable/disable features and adjust behavior.

### Basic Usage

```go
import "github.com/bonarizki-dat/boilerplate-gin-dat/pkg/config"

// Get current environment
env := config.GetEnvironment()
fmt.Println("Running in:", env) // development, staging, or production

// Check specific environments
if config.IsDevelopment() {
    fmt.Println("Development mode active")
}

if config.IsProduction() {
    fmt.Println("Production mode active")
}
```

### Real-World Examples

**1. Conditional Database Logging**

```go
// internal/adapters/database/database.go
func DbConnection(masterDSN, replicaDSN string) error {
    if config.IsDevelopment() {
        masterDB = masterDB.Debug() // Enable SQL logging only in dev
    }
    return nil
}
```

**2. Environment-Specific Seeding**

```go
// migrations/migrate.go
func Migrate() {
    database.MasterDB.AutoMigrate(&models.User{})

    if config.IsDevelopment() {
        seedTestData() // Only seed in development
    }
}
```

**3. Different Rate Limits Per Environment**

```go
var limit int
if config.IsProduction() {
    limit = 10 // Strict in production
} else {
    limit = 1000 // Relaxed in dev/staging
}
authRoutes.Use(middlewares.RateLimitMiddlewareWithConfig(limit, limit*2))
```

**4. Development-Only Debug Endpoints**

```go
if config.IsDevelopment() {
    router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
}
```

**5. Environment-Specific Error Messages**

```go
func InternalServerError(c *gin.Context, err error, message string) {
    if config.IsDevelopment() {
        // Show details in dev
        HandleErrors(c, 500, err, fmt.Sprintf("%s: %v", message, err))
    } else {
        // Hide details in production
        HandleErrors(c, 500, nil, message)
    }
}
```

**6. Feature Flags**

```go
type FeatureFlags struct {
    EmailVerificationRequired bool
    TwoFactorAuthEnabled      bool
    BetaFeatures              bool
}

func GetFeatureFlags() FeatureFlags {
    if config.IsProduction() {
        return FeatureFlags{
            EmailVerificationRequired: true,
            TwoFactorAuthEnabled:      true,
            BetaFeatures:              false,
        }
    }
    // Relaxed for development
    return FeatureFlags{
        EmailVerificationRequired: false,
        TwoFactorAuthEnabled:      false,
        BetaFeatures:              true,
    }
}
```

**7. Environment-Specific Logging**

```go
func InitLogger() {
    if config.IsProduction() {
        logrus.SetLevel(logrus.WarnLevel) // Only warnings in prod
        logrus.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logrus.SetLevel(logrus.DebugLevel) // Everything in dev
        logrus.SetFormatter(&logrus.TextFormatter{})
    }
}
```

**8. Conditional Middleware**

```go
// Only enable metrics in staging and production
if !config.IsDevelopment() {
    router.Use(middlewares.MetricsMiddleware())
}

// Security headers only in production
if config.IsProduction() {
    router.Use(middlewares.SecurityHeadersMiddleware())
}
```

### Security Best Practices

**Don't expose environment in responses:**

```go
// ❌ BAD
router.GET("/status", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "environment": config.GetEnvironment(), // ❌ Security risk!
    })
})

// ✅ GOOD
router.GET("/status", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})
```

**Environment-based security:**

```go
if config.IsProduction() {
    minPasswordLength = 12 // Strict in production
} else {
    minPasswordLength = 6  // Relaxed for testing
}
```

---

## 🔗 Related Documentation

- [README.md](../README.md) - Project overview and quick start
- [CODING_STANDARDS.md](CODING_STANDARDS.md) - Section 12: Configuration
- [Security Guidelines](CODING_STANDARDS.md#9-security-guidelines)

---

## 📝 Configuration Checklist

Before deploying to any environment:

```
□ Copied .env.example to .env
□ Generated unique SECRET (min 32 chars)
□ Generated unique JWT_SECRET (min 32 chars)
□ Updated all database credentials
□ Set DEBUG=False for production
□ Enabled SSL for database (MASTER_SSL_MODE=require)
□ Verified all required variables are set
□ Tested application starts without errors
□ Never committed .env file to git
□ Documented where production secrets are stored
```

---

**Last Updated:** 2025-11-09
**Maintainer:** Boilerplate Team
