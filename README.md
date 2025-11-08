<div align="center">

# 🚀 Go Gin Enterprise Boilerplate

**Production-Ready Starter Kit for Building Scalable RESTful APIs**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.11-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](docs/CODING_STANDARDS.md)

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Architecture](#-architecture) • [Testing](#-testing)

</div>

---

## 📖 What is This?

**This is not just another boilerplate** - This is a **company-standard starter kit** designed to accelerate API development while maintaining high code quality, consistency, and best practices across all projects.

### 🎯 Purpose

This boilerplate serves as the **foundation for all Go API projects** in our organization. It eliminates the need to set up authentication, database connections, testing infrastructure, and project structure from scratch for every new project.

### 💡 Why This Boilerplate?

| Problem | Solution |
|---------|----------|
| ❌ Starting from zero for each project | ✅ Production-ready foundation with auth, DB, testing |
| ❌ Inconsistent code across projects | ✅ Enforced coding standards and design patterns |
| ❌ Poor documentation | ✅ 5000+ lines of comprehensive documentation |
| ❌ No testing examples | ✅ Complete test suite with examples |
| ❌ Security vulnerabilities | ✅ Built-in security best practices |
| ❌ AI agents breaking conventions | ✅ AI-friendly docs with critical rules |

### 🎁 What You Get

Start your next API project with:
- ✅ **JWT Authentication** - Login, register, protected routes ready to go
- ✅ **Clean Architecture** - Proven layered structure (Controllers → Services → Repositories)
- ✅ **Complete Documentation** - 5000+ lines covering every aspect
- ✅ **Testing Infrastructure** - Unit tests, integration tests, examples included
- ✅ **Security Built-in** - SQL injection prevention, password hashing, token security
- ✅ **Database Ready** - PostgreSQL with GORM, migrations, master-replica support
- ✅ **Docker Support** - Development and production configurations
- ✅ **AI-Ready** - Comprehensive guides for AI-assisted development

**Time to First API:** ~5 minutes instead of days 🚀

---

## ✨ Features

### Core Features

- 🔐 **JWT Authentication** - Secure login/register endpoints with token-based auth
- 🏗️ **Clean Architecture** - Layered design with clear separation of concerns
- 🗄️ **GORM Integration** - PostgreSQL with master-replica configuration
- ✅ **Request Validation** - Built-in validation using go-playground/validator
- 📝 **Structured Logging** - Logrus integration with custom formatting
- 🔌 **Middleware Support** - CORS, Auth, Rate Limiting middleware
- 🛡️ **Rate Limiting** - Per-IP rate limiting to prevent abuse and brute force attacks
- 🗃️ **Database Migrations** - SQL migrations and AutoMigrate support
- 🐳 **Docker Support** - Dev and prod Docker configurations with live reload
- 🧪 **Comprehensive Testing** - Service and controller test examples
- 📊 **DataTables Integration** - Server-side pagination, search, and sorting

### What Makes This Different

- 📚 **World-Class Documentation** - 5000+ lines covering standards, patterns, and AI guides
- 🤖 **AI-Friendly** - Specialized docs for AI-assisted development ([docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md))
- 🛡️ **Security First** - SQL injection prevention, secure password handling, token security
- 🎯 **Battle-Tested Patterns** - Proven in production environments
- 📏 **Enforced Standards** - File size limits, function limits, test coverage requirements
- 🧩 **Modular Design** - Easy to extend, hard to break

---

## 🚀 Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/your-org/go-gin-boilerplate.git
cd go-gin-boilerplate

# 2. Copy environment file
cp .env.example .env

# 3. Install dependencies
go get .

# 4. Run the application
go run main.go

# ✅ Visit: http://localhost:8000/health
```

**Using Docker:**
```bash
# Development with live reload
make dev

# Production build
make production
```

**First API Call:**
```bash
# Register a user
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'

# ✅ You now have a working API with authentication!
```

---

## 📚 Documentation

> **🎯 For New Developers:** Start with the documentation to understand the architecture and coding standards.

### 📖 Reading Guide

**For AI Agents - START HERE ⚠️**

If you're an AI agent or using AI-assisted development:

1. **[docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md)** ⚠️ **READ FIRST** (100 lines)
   - Non-negotiable patterns (struct-based, response utilities)
   - Absolute rules that MUST be followed
   - **Skip this = Code rejected**

2. **[docs/AI_QUICK_REFERENCE.md](docs/AI_QUICK_REFERENCE.md)** (405 lines)
   - Quick templates for controllers, services, repositories
   - The 5 Commandments (size limits)
   - Testing checklist

3. **Use as Reference:**
   - [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md) - Complete coding standards
   - [docs/DESIGN_PATTERNS.md](docs/DESIGN_PATTERNS.md) - Architecture patterns

**For Human Developers**

| Document | Size | Purpose |
|----------|------|---------|
| [docs/README.md](docs/README.md) | Quick | Documentation navigation guide |
| [docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md) | 100 lines | Critical rules summary |
| [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md) | 1955 lines | Comprehensive coding standards |
| [docs/DESIGN_PATTERNS.md](docs/DESIGN_PATTERNS.md) | 2479 lines | Architecture and design patterns |
| [TESTING.md](TESTING.md) | Full | Complete testing guide |
| [tests/README.md](tests/README.md) | Quick | Test organization |

**Quick Links:**
- 🏗️ [Architecture Overview](#-architecture)
- 🔐 [Authentication Guide](#-authentication-endpoints)
- 🧪 [Testing Guide](TESTING.md)
- 🗄️ [Database Migrations](internal/adapters/database/migrations/sql/README.md)
- 🐳 [Docker Setup](#-docker-development)

---

## 🏗️ Architecture

This boilerplate follows **Clean Architecture** principles with a layered approach:

```
┌─────────────────────────────────────────┐
│           HTTP Request                  │
└──────────────┬──────────────────────────┘
               │
        ┌──────▼──────┐
        │  Router     │  Route definitions
        └──────┬──────┘
               │
        ┌──────▼──────────┐
        │  Middleware     │  Auth, CORS, etc.
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │  Controllers    │  HTTP handlers (thin layer)
        │  - Validate     │  Max 50 lines per function
        │  - Call Service │
        │  - Return JSON  │
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │  Services       │  Business logic (fat layer)
        │  - Validation   │  Max 100 lines per function
        │  - Processing   │  Struct-based with DI
        │  - Orchestrate  │
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │  Repositories   │  Data access (CRUD only)
        │  - Create       │  Function-based
        │  - Read         │  Return models
        │  - Update       │
        │  - Delete       │
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │  Database       │  PostgreSQL + GORM
        └─────────────────┘
```

### Directory Structure

```
project/
├── cmd/                        # Commands (migrate, seeder)
├── internal/                   # Private application code
│   ├── adapters/              # External adapters (DB, cache)
│   ├── app/
│   │   ├── controllers/       # HTTP handlers (struct-based)
│   │   ├── dto/               # Data Transfer Objects
│   │   ├── middlewares/       # Gin middlewares
│   │   ├── routers/           # Route definitions
│   │   └── services/          # Business logic (struct-based)
│   └── domain/
│       ├── models/            # Database entities (GORM)
│       └── repositories/      # Data access layer
├── pkg/                       # Public reusable packages
│   ├── config/               # Configuration management
│   ├── logger/               # Logging infrastructure
│   ├── types/                # Shared types
│   └── utils/                # Utility functions
│       └── response.go       # MUST use for all responses
├── tests/                    # ALL tests go here
│   ├── unit/                 # Unit tests
│   │   ├── controllers/
│   │   ├── services/
│   │   └── repositories/
│   ├── integration/          # Integration tests
│   └── fixtures/             # Test data
├── docs/                     # Documentation
└── main.go                   # Application entry point
```

**Key Principles:**
- ✅ Controllers are thin (validation + call service)
- ✅ Services contain all business logic
- ✅ Repositories only do CRUD operations
- ✅ No circular dependencies
- ✅ Dependency injection via constructors
- ❌ Controllers never access database directly
- ❌ Repositories never contain business logic

---

## 🔐 API Endpoints

### Authentication Endpoints

#### Register a New User

```bash
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}

# Response (201 Created)
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer"
  },
  "errors": null
}
```

#### Login

```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}

# Response (200 OK)
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer"
  },
  "errors": null
}
```

### Protected Routes

Routes under `/api/*` require JWT authentication:

```bash
GET /api/profile
Authorization: Bearer <your-jwt-token>

# Response (200 OK)
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "user_id": 1
  },
  "errors": null
}
```

### Public Endpoints

- `GET /health` - Health check endpoint
- `GET /datatables` - DataTables example with pagination/search

**Standard Response Format:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "errors": null
}
```

---

## 🧪 Testing

Run the complete test suite:

```bash
# Run all tests
go test ./tests/...

# Run with coverage
go test -cover ./tests/...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out

# Run specific test package
go test ./tests/unit/services/...

# Run with race detector
go test -race ./tests/...

# Run with verbose output
go test -v ./tests/...
```

**Test Coverage Goals:**
- Services: 70% minimum (85% target)
- Repositories: 70% minimum
- Controllers: 60% minimum
- Utils: 80% minimum

**Test Structure:**
```
tests/
├── unit/
│   ├── controllers/    # Controller HTTP tests
│   ├── services/       # Business logic tests
│   ├── repositories/   # Data access tests
│   └── utils/          # Utility tests
├── integration/
│   ├── api/           # End-to-end API tests
│   └── database/      # Database integration tests
└── fixtures/          # Test data (JSON, CSV)
```

**Example Test:**
```go
// tests/unit/services/auth_service_test.go
package services_test

import (
    "testing"
    "github.com/your-org/project/internal/app/services"
)

func TestAuthService_ValidateToken(t *testing.T) {
    service := services.NewAuthService()

    tests := []struct {
        name    string
        token   string
        wantErr bool
    }{
        {"Valid token", "valid.jwt.token", false},
        {"Invalid token", "invalid", true},
        {"Empty token", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := service.ValidateToken(tt.token)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

See [TESTING.md](TESTING.md) for comprehensive testing guide.

---

## ⚙️ Configuration

### Environment Variables

Create `.env` file from `.env.example`:

```bash
cp .env.example .env
```

**Essential Configuration:**

```env
# Server
SECRET=your-super-secret-jwt-key-change-this
DEBUG=True                    # Set False in production
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

# Database (Master)
MASTER_DB_NAME=your_database
MASTER_DB_USER=your_user
MASTER_DB_PASSWORD=your_password
MASTER_DB_HOST=localhost      # Use postgres_db for Docker
MASTER_DB_PORT=5432
MASTER_DB_LOG_MODE=True       # Set False in production

# Database (Replica) - Optional
REPLICA_DB_NAME=your_database
REPLICA_DB_USER=your_user
REPLICA_DB_PASSWORD=your_password
REPLICA_DB_HOST=localhost
REPLICA_DB_PORT=5432
```

**Security Notes:**
- ⚠️ Change `SECRET` in production
- ⚠️ Set `DEBUG=False` in production
- ⚠️ Set `MASTER_DB_LOG_MODE=False` in production
- ⚠️ Never commit `.env` to version control

### Database Configuration

**Local Development:**
```env
MASTER_DB_HOST=localhost
```

**Docker Development:**
```env
MASTER_DB_HOST=postgres_db
```

**Master-Replica Setup:**
- Master for writes (INSERT, UPDATE, DELETE)
- Replica for reads (SELECT)
- Automatic failover support

### JWT Configuration

```env
SECRET=your-jwt-secret-key-min-32-characters
```

- Token expiry: 24 hours (configurable in `auth_service.go`)
- Algorithm: HS256
- Claims: user_id, email, exp, iat

---

## 🐳 Docker Development

### Development with Live Reload

```bash
# Start development environment
make dev

# This starts:
# - PostgreSQL database (port 5432)
# - PG Admin (port 5050)
# - Go API with live reload (port 8000)
```

**Access Services:**
- API: [http://localhost:8000](http://localhost:8000)
- PG Admin: [http://localhost:5050](http://localhost:5050)
  - Email: `admin@admin.com`
  - Password: `root`
  - DB Host: `postgres_db`

### Production Build

```bash
# Build and run production containers
make production

# Build only
make build

# Clean up
make clean
```

**Docker Commands:**
- `make dev` - Development with live reload (Air)
- `make build` - Build production container
- `make production` - Build and run production
- `make clean` - Remove all containers and images

---

## 🗄️ Database Migrations

### Development Approach

Uses GORM AutoMigrate for quick iteration:

```go
// internal/adapters/database/migrations/migration.go
func Migrate() {
    models := []interface{}{
        &models.User{},
        // Add your models here
    }
    database.DB.AutoMigrate(models...)
}
```

### Production Approach

Uses SQL migration files with `golang-migrate`:

```bash
# Install golang-migrate
brew install golang-migrate  # macOS
# or
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql up

# Rollback last migration
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql down 1

# Create new migration
migrate create -ext sql -dir internal/adapters/database/migrations/sql -seq create_new_table
```

**Migration Files:**
```
internal/adapters/database/migrations/sql/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_posts_table.up.sql
└── 000002_create_posts_table.down.sql
```

See [migrations README](internal/adapters/database/migrations/sql/README.md) for detailed guide.

---

## 🚀 Deployment

### Prerequisites

- Go 1.24+
- PostgreSQL 13+
- Docker (optional)

### Production Checklist

```bash
□ Update .env with production values
□ Set DEBUG=False
□ Set strong SECRET key (min 32 chars)
□ Set MASTER_DB_LOG_MODE=False
□ Configure SSL for database
□ Run database migrations
□ Set up monitoring and logging
□ Configure reverse proxy (nginx)
□ Set up SSL/TLS certificates
□ Configure CORS for your domain
```

### Build for Production

```bash
# Build binary
go build -o api main.go

# Run binary
./api

# Or use Docker
make production
```

### Environment Variables (Production)

```env
DEBUG=False
SECRET=super-long-random-secret-key-min-32-chars
MASTER_DB_LOG_MODE=False
MASTER_SSL_MODE=require
```

---

## 📦 Tech Stack

### Core Framework
- **[Gin](https://github.com/gin-gonic/gin)** - High-performance HTTP web framework
- **[GORM](https://github.com/go-gorm/gorm)** - Fantastic ORM library for Golang
- **[Viper](https://github.com/spf13/viper)** - Configuration management

### Authentication & Security
- **[JWT-Go](https://github.com/golang-jwt/jwt)** - JSON Web Token implementation
- **[Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password hashing

### Validation & Logging
- **[Validator](https://github.com/go-playground/validator)** - Struct validation
- **[Logrus](https://github.com/sirupsen/logrus)** - Structured logging

### Testing
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit with assertions

### Development Tools
- **[Air](https://github.com/cosmtrek/air)** - Live reload for Go apps
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migrations

### Custom Libraries
- **[Datatables-Gin](https://github.com/bonarizki-dat/Datatables-Gin)** - DataTables integration

---

## 🎯 Project Standards

### The 5 Commandments

```
1. 📏 File >300 lines?        → STOP. Split it.
2. 📐 Function >100 lines?    → STOP. Extract functions.
3. 🧪 No tests?               → STOP. Write tests first.
4. ❌ Error ignored (_, _)?   → STOP. Handle it.
5. 📝 Exported without docs?  → STOP. Document it.
```

**Violate = Code Rejected**

### Code Quality Standards

- ✅ All controllers MUST be struct-based (NOT standalone functions)
- ✅ All services MUST be struct-based (NOT standalone functions)
- ✅ All responses MUST use `utils.Ok/Created/etc` (NOT `c.JSON`)
- ✅ All tests MUST be in `tests/` directory (NOT co-located)
- ✅ File size MAX 300 lines
- ✅ Function size MAX 100 lines
- ✅ Test coverage MIN 70% for services

See [docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md) for complete rules.

---

## 🤝 Contributing

We welcome contributions! Please follow our standards:

**Before Contributing:**
1. Read [docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md)
2. Read [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md)
3. Check existing issues and PRs

**Contribution Process:**
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Follow coding standards (struct-based, response utilities, etc.)
4. Write tests (minimum 70% coverage)
5. Commit with clear messages
6. Push to branch
7. Open Pull Request

**Pull Request Requirements:**
- ✅ All tests passing
- ✅ Code follows standards
- ✅ Documentation updated
- ✅ No linter errors
- ✅ Test coverage maintained

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💬 Support

**Issues & Questions:**
- 🐛 [Report Bug](https://github.com/your-org/go-gin-boilerplate/issues)
- 💡 [Request Feature](https://github.com/your-org/go-gin-boilerplate/issues)
- 📖 [Read Documentation](docs/README.md)

**Resources:**
- [Documentation Guide](docs/README.md)
- [Critical Rules](docs/00_AI_CRITICAL_RULES.md)
- [Coding Standards](docs/CODING_STANDARDS.md)
- [Design Patterns](docs/DESIGN_PATTERNS.md)
- [Testing Guide](TESTING.md)

---

## ⭐ Acknowledgments

- Gin Framework team for the excellent HTTP framework
- GORM team for the powerful ORM
- All contributors to the open-source packages used

---

<div align="center">

**Built with ❤️ by the team**

**[⬆ Back to Top](#-go-gin-enterprise-boilerplate)**

</div>
