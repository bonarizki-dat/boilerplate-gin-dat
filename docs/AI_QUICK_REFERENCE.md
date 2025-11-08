# AI Agent Quick Reference

> **Print this mentally before every code change!**

---

## ⚠️ FIRST TIME HERE?

**🚨 READ [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) FIRST!**

That file contains the absolute non-negotiable rules (100 lines).
This file is for quick templates and checklists.

---

## ⚡ THE 5 COMMANDMENTS

```
1. 📏 File >300 lines?        → STOP. Split it.
2. 📐 Function >100 lines?    → STOP. Extract functions.
3. 🧪 No tests?               → STOP. Write tests first.
4. ❌ Error ignored (_, _)?   → STOP. Handle it.
5. 📝 Exported without docs?  → STOP. Document it.
```

**VIOLATE = CODE REJECTED**

---

## 🎯 Before Writing ANY Code

```bash
# Ask yourself:
□ What layer am I in? (Controller/Service/Repository)
□ Am I following dependency direction?
□ Will this file exceed 300 lines? → Plan to split
□ Will this function exceed 100 lines? → Plan to extract
□ Do I need tests? → Yes, ALWAYS for services
□ Is this documented? → Required for exported items
```

---

## 📐 Size Limits (HARD LIMITS)

```
File:     MAX 300 lines  (warning at 250)
Function: MAX 100 lines  (warning at 80)
```

**Approaching limit?**
- Stop and refactor NOW
- Don't wait until you exceed
- Split proactively

---

## 🏗️ Architecture Cheat Sheet

```
Request Flow:
Router → Controller → Service → Repository → Database

Layers:
┌──────────────┐
│  Controller  │  ← HTTP only, <50 lines/function
├──────────────┤
│   Service    │  ← Business logic, <100 lines/function
├──────────────┤
│  Repository  │  ← CRUD only, return models
├──────────────┤
│    Model     │  ← Data structures only
└──────────────┘

Dependencies:
Controller  →  Service  →  Repository  →  Model
    ↓           ↓             ↓
   DTO        Utils        Database
```

**Forbidden:**
- ❌ Controller with business logic
- ❌ Service accessing database directly
- ❌ Repository with business logic
- ❌ Circular dependencies

---

## 🔥 Templates

### Controller (15-30 lines)
```go
func GetResource(c *gin.Context) {
    // 1. Parse input (5 lines)
    id := c.Param("id")

    // 2. Call service (3 lines)
    resource, err := service.GetByID(id)
    if err != nil {
        utils.HandleErrors(c, http.StatusNotFound, nil, err.Error())
        return
    }

    // 3. Return response (2 lines)
    c.JSON(http.StatusOK, resource)
}
```

### Service (30-80 lines)
```go
func (s *ResourceService) Create(dto dto.CreateRequest) (*models.Resource, error) {
    // 1. Validate (5 lines)
    if err := s.validate(dto); err != nil {
        return nil, fmt.Errorf("validation: %w", err)
    }

    // 2. Business logic (10-20 lines)
    existing, _ := s.repo.FindByName(dto.Name)
    if existing != nil {
        return nil, types.ErrDuplicateEntry
    }

    // 3. Transform (5-10 lines)
    resource := s.buildResource(dto)

    // 4. Persist (5 lines)
    if err := s.repo.Create(resource); err != nil {
        logger.Errorf("create failed: %v", err)
        return nil, fmt.Errorf("create: %w", err)
    }

    // 5. Log & return (3 lines)
    logger.Infof("Resource created: ID=%d", resource.ID)
    return resource, nil
}
```

### Repository (15-25 lines)
```go
func (r *resourceRepository) FindByID(id uint) (*models.Resource, error) {
    var resource models.Resource

    err := r.db.
        Preload("RelatedEntity").
        First(&resource, id).
        Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, types.ErrNotFound
        }
        return nil, fmt.Errorf("query failed: %w", err)
    }

    return &resource, nil
}
```

---

## ✅ Error Handling Pattern

```go
// ✅ ALWAYS do this:
result, err := someFunction()
if err != nil {
    logger.Errorf("context: %v", err)                    // Log
    return fmt.Errorf("operation failed: %w", err)       // Wrap with %w
}

// ❌ NEVER do this:
result, _ := someFunction()                              // Ignored!
result, err := someFunction()
if err != nil {
    panic(err)                                           // Panic!
}
result, err := someFunction()
return err                                               // Not wrapped!
```

---

## 📝 Documentation Pattern

```go
// ✅ CORRECT:
// CreateUser creates a new user account with validation.
//
// Returns ErrDuplicateEntry if email exists.
// Returns ErrValidation if input is invalid.
func CreateUser(dto CreateUserRequest) (*User, error) {
    // implementation
}

// ❌ WRONG:
// Create user
func CreateUser(dto CreateUserRequest) (*User, error) {

// ❌ WRONG:
func CreateUser(dto CreateUserRequest) (*User, error) {  // No comment
```

---

## 🧪 Testing Checklist

```go
⚠️  CRITICAL: ALL tests MUST be in tests/ directory
□ Create tests/unit/{layer}/{filename}_test.go
□ Use package {layer}_test (external test)
□ Import the package being tested
□ Test happy path
□ Test 2+ error cases
□ Use table-driven tests if >3 scenarios
□ Mock dependencies
□ Assert expectations
□ Run: go test ./tests/...
□ Coverage >70% for services
```

**Example test location:**
```
internal/app/services/auth_service.go
→ tests/unit/services/auth_service_test.go (package services_test)

internal/app/controllers/user_controller.go
→ tests/unit/controllers/user_controller_test.go (package controllers_test)
```

---

## 🚨 Forbidden Patterns

```go
❌ panic() in business logic
❌ _, _ = someFunc()              // Ignored error
❌ "SELECT * FROM " + table       // SQL injection
❌ if x { if y { if z { } } }     // Too nested (>3 levels)
❌ password := "hardcoded"        // Hardcoded secrets
❌ log.Printf()                   // Use logger.Infof()
❌ file size >300 lines
❌ function >100 lines
❌ No tests for services
❌ Exported function without docs
```

---

## 🎨 Naming Conventions

```go
// Files
✅ user_service.go
❌ UserService.go, user-service.go

// Packages
✅ package services
❌ package user_services, package Services

// Variables
✅ user, userID, httpClient
❌ u, usrID, http_client

// Functions
✅ GetUserByID, CreateTransaction
❌ get_user, GetUser (too generic)

// Constants
✅ const MaxRetryAttempts = 3
✅ const StatusPending Status = "pending"
❌ const MAX_RETRY = 3
```

---

## 🔍 Pre-Commit Checklist

```bash
□ All functions <100 lines?
□ All files <300 lines?
□ All errors handled?
□ All exported items documented?
□ Tests written and passing?
□ No hardcoded secrets?
□ No panic() in business logic?
□ No SQL string concatenation?
□ No ignored errors (_, _)?
□ gofmt applied?

# Run these:
gofmt -w .
go vet ./...
go test ./...
```

---

## 🚀 When Refactoring Large Files

**If file >300 lines:**

1. **Identify boundaries**
   - Group related functions
   - Find logical separations

2. **Create new files**
   ```
   service.go           → service.go (main)
                        → service_helpers.go
                        → service_validators.go
                        → service_transformers.go
   ```

3. **Move code**
   - Keep related functions together
   - Maintain package cohesion

4. **Update imports**

5. **Run tests**
   ```bash
   go test ./...
   ```

**If function >100 lines:**

1. **Extract logical blocks**
   ```go
   // Before: 200 lines
   func Process() { ... }

   // After: Multiple focused functions
   func Process() {           // 20 lines - orchestration
       data := parse()
       validated := validate(data)
       transformed := transform(validated)
       save(transformed)
   }

   func parse() { }           // 30 lines
   func validate() { }        // 25 lines
   func transform() { }       // 40 lines
   func save() { }            // 20 lines
   ```

---

## 💡 Common Patterns

### Pagination
```go
func List(page, pageSize int) ([]*Model, int64, error) {
    if page < 1 { page = 1 }
    if pageSize < 1 || pageSize > 100 { pageSize = 20 }

    var items []*Model
    var total int64

    db := r.db.Model(&Model{})
    db.Count(&total)

    err := db.
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&items).
        Error

    return items, total, err
}
```

### Transactions
```go
func Process(data Data) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := step1(tx, data); err != nil {
            return err  // Auto rollback
        }
        if err := step2(tx, data); err != nil {
            return err  // Auto rollback
        }
        return nil  // Auto commit
    })
}
```

### Validation
```go
type Request struct {
    Name  string `validate:"required,min=3,max=255"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=150"`
}

func validate(req Request) error {
    return validator.New().Struct(req)
}
```

---

## 📚 Quick Links

- **Full Standards:** `CODING_STANDARDS.md`
- **AI Rules:** `AI_AGENT_RULES.md`
- **Refactoring Plan:** `REFACTORING_PLAN.md`
- **Quality Checks:** `make -f Makefile.quality help`

---

## 🎯 Remember

```
Small files    = Easy to understand
Small functions = Easy to test
Good tests     = Confident refactoring
Good docs      = Happy developers

✅ Quality > Quantity
✅ Simple > Complex
✅ Clear > Clever
```

---

**Print this before every commit!**
**Follow the rules strictly!**
**Your future self will thank you!**
