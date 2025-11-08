# 📚 Documentation Guide

Welcome to the project documentation. This guide helps you navigate all documentation files.

---

## 🚨 FOR AI AGENTS - START HERE

### Reading Order (MANDATORY)

1. **[`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md)** ⚠️ **START HERE** (100 lines)
   - Absolute non-negotiable rules
   - Struct-based patterns (MANDATORY)
   - Response utilities (MANDATORY)
   - Test location rules (MANDATORY)
   - **READ THIS FIRST OR YOUR CODE WILL BE REJECTED**

2. **[`AI_QUICK_REFERENCE.md`](./AI_QUICK_REFERENCE.md)** (405 lines)
   - Quick templates for controllers, services, repositories
   - The 5 Commandments (file/function size limits)
   - Testing checklist
   - Common patterns

3. **Critical Sections from Long Docs:**
   - **CODING_STANDARDS.md:**
     - Lines 900-1100: Struct-based patterns
     - Lines 1429-1475: Response format
     - Lines 1479-1584: Response utilities
     - Lines 840-1009: Testing requirements

   - **DESIGN_PATTERNS.md:**
     - Lines 900-948: Controller pattern
     - Lines 984-1016: Service pattern
     - Lines 1070-1100: Response pattern
     - Lines 439-492: Dependency injection

4. **Use as Reference:**
   - Full CODING_STANDARDS.md when needed
   - Full DESIGN_PATTERNS.md for detailed patterns

---

## 📖 Documentation Files

### For AI Agents

| File | Size | Purpose | When to Read |
|------|------|---------|--------------|
| **00_AI_CRITICAL_RULES.md** | 100 lines | Non-negotiable rules | **FIRST - ALWAYS** |
| **AI_QUICK_REFERENCE.md** | 405 lines | Templates & checklists | Before writing code |
| **CODING_STANDARDS.md** | 1955 lines | Complete coding standards | Reference (focus on critical sections) |
| **DESIGN_PATTERNS.md** | 2479 lines | Architecture patterns | Reference (focus on critical sections) |

### For Developers

| File | Purpose |
|------|---------|
| **CODING_STANDARDS.md** | Comprehensive coding standards, naming conventions, best practices |
| **DESIGN_PATTERNS.md** | Architecture patterns, layer responsibilities, implementation guides |
| **AI_QUICK_REFERENCE.md** | Quick templates and decision trees |
| **00_AI_CRITICAL_RULES.md** | Quick reference for critical rules |

---

## 🎯 Quick Navigation

### I Want To...

**Write a new controller:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1)
2. Template: `AI_QUICK_REFERENCE.md` → Templates section
3. Details: `DESIGN_PATTERNS.md` lines 900-948

**Write a new service:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1)
2. Template: `AI_QUICK_REFERENCE.md` → Templates section
3. Details: `DESIGN_PATTERNS.md` lines 984-1016

**Return a response:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 2)
2. Utils: `pkg/utils/response.go`
3. Details: `CODING_STANDARDS.md` lines 1479-1584

**Write tests:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 3)
2. Guide: `tests/README.md`
3. Details: `CODING_STANDARDS.md` lines 840-1009

**Handle errors:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 2)
2. Details: `CODING_STANDARDS.md` → Error Handling section

**Use dependency injection:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 4)
2. Details: `DESIGN_PATTERNS.md` lines 439-492

---

## 🚫 Common Mistakes (Don't Do This)

### Mistake #1: Not Reading Critical Rules
```
❌ Skipping 00_AI_CRITICAL_RULES.md
✅ Reading it first (takes 3 minutes)
```

**Result of skipping:** Code rejected, need to refactor everything.

### Mistake #2: Using Standalone Functions
```go
❌ func Register(ctx *gin.Context) { }  // Rejected
✅ func (ctrl *AuthController) Register(c *gin.Context) { }
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 1

### Mistake #3: Direct c.JSON() Calls
```go
❌ c.JSON(200, gin.H{"data": user})  // Rejected
✅ utils.Ok(c, user, "Success")
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 2

### Mistake #4: Co-located Tests
```
❌ internal/app/services/auth_service_test.go  // Rejected
✅ tests/unit/services/auth_service_test.go
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 3

### Mistake #5: Exceeding File Size
```
❌ File with 400 lines  // Rejected
✅ Split into multiple focused files (max 300 lines)
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 1

---

## 📊 Documentation Statistics

- **Total Documentation:** ~5,000 lines
- **Critical Rules:** 100 lines (2%)
- **Must Read Before Coding:** 505 lines (10%)
- **Reference Material:** 4,434 lines (90%)

**Efficiency Tip:** Read 10% (critical + quick ref) first, use 90% as reference.

---

## 🔄 Document Updates

**Last Updated:** 2025-11-09

**Recent Changes:**
- Added `00_AI_CRITICAL_RULES.md` (NEW - 2025-11-09)
- Updated all docs with critical section pointers (2025-11-09)
- Added tests/ directory structure requirements (2025-11-09)

---

## ✅ Checklist Before First Code Contribution

```
□ Read 00_AI_CRITICAL_RULES.md (100 lines)
□ Read AI_QUICK_REFERENCE.md (405 lines)
□ Understand struct-based pattern requirement
□ Understand response utilities requirement
□ Understand tests/ directory requirement
□ Know file size limits (300 lines max)
□ Know function size limits (100 lines max)
```

**Time Required:** 15-20 minutes
**Time Saved:** Hours of refactoring

---

## 📞 Questions?

If documentation is unclear:
1. Check `00_AI_CRITICAL_RULES.md` first
2. Search in `CODING_STANDARDS.md` or `DESIGN_PATTERNS.md`
3. Look for examples in existing code
4. Ask the team

**For AI Agents:** If you're unsure, ASK. Don't guess and violate critical rules.
