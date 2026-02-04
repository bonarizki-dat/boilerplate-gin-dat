# Rencana Implementasi Improvement

Dokumen ini mendampingi [IMPROVEMENTS.md](IMPROVEMENTS.md). Berisi **langkah konkret**, **file yang diubah**, dan **urutan eksekusi** per fase.

---

## Fase 1: Quick wins (1–2 jam)

**Tujuan:** Konsistensi config, production-ready shutdown, perbaikan bug kecil.  
**Output:** Satu PR (atau satu commit) berisi item 1–4.

| Step | Item | File yang diubah | Tindakan |
|------|------|------------------|----------|
| 1.1 | DB log mode | `internal/adapters/database/database.go` | Ganti `viper.GetBool("DB_LOG_MODE")` → `viper.GetBool("MASTER_DB_LOG_MODE")`. Cek tidak ada referensi lain ke `DB_LOG_MODE`. |
| 1.2 | Default timezone | `main.go` | Ganti `viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")` → `viper.SetDefault("SERVER_TIMEZONE", "UTC")`. |
| 1.3 | Graceful shutdown | `main.go` | (1) Buat `http.Server{Addr, Handler: router, ReadTimeout, WriteTimeout}`. (2) Start server di goroutine dengan `server.ListenAndServe()`. (3) Channel `sigChan` listen `os.Interrupt`, `syscall.SIGTERM`. (4) On signal: `server.Shutdown(context.WithTimeout(..., 10*time.Second))`, lalu `database.GetDB().DB()` → `Close()`. (5) `os.Exit(0)` atau return dari main. |
| 1.4 | Bug err di database | `internal/adapters/database/database.go` | Hapus `var err error` package-level. Di `DbConnection`: pakai `db, err := gorm.Open(...)`; pindahkan pengecekan `if err != nil` langsung setelah `gorm.Open`; assign `DB = db` hanya jika tidak error. |

**Verifikasi Fase 1:** `go build ./...`, `go test ./tests/...`, jalankan app lalu Ctrl+C dan pastikan log shutdown rapi.

**Dokumentasi:** Tambah paragraf singkat di README (Deployment) atau CONFIGURATION.md: graceful shutdown (SIGTERM/SIGINT), optional `SERVER_SHUTDOWN_TIMEOUT`.

---

## Fase 2: Medium – Stabilitas & testing (setengah–1 hari)

**Tujuan:** Request timeout, contoh integration test, dokumentasi error strategy.

| Step | Item | File yang diubah | Tindakan |
|------|------|------------------|----------|
| 2.1 | Request timeout | `internal/app/routers/router.go` atau `internal/app/middlewares/` | Tambah middleware: `context.WithTimeout(c.Request.Context(), 30*time.Second)`, replace request context, defer cancel. Atau set `ReadTimeout`/`WriteTimeout` di `http.Server` (Fase 1). Config via env `REQUEST_TIMEOUT_SECONDS` (default 30). |
| 2.2 | Integration test – health | `tests/integration/api/health_test.go` (baru) | Test: start app (atau httptest), `GET /health`, assert status 200 dan body berisi `"status":"healthy"` (atau unhealthy jika DB mati). Baca `tests/README.md` untuk cara run (e.g. `go test ./tests/integration/... -tags=integration` atau env). |
| 2.3 | Integration test – DB (opsional) | `tests/integration/database/connection_test.go` (baru) | Test: panggil `database.DbConnection(masterDSN, replicaDSN)` dengan DSN dari env/testcontainer, lalu ping. Skip jika tidak ada DB. |
| 2.4 | Dokumen error strategy | `docs/CODING_STANDARDS.md` atau `pkg/types/errors.go` | Tambah subbagian: “Error strategy – response HTTP selalu lewat `pkg/utils`; domain errors (e.g. `auth.Err*`) untuk business logic; `pkg/types/errors` untuk [kasus X].” Supaya tidak duplikasi dan konsisten. |

**Verifikasi Fase 2:** Integration test jalan (dengan/skip DB); timeout bisa diuji dengan handler yang sleep.

---

## Fase 3: Medium – Testability & fitur (1–2 hari)

**Tujuan:** Repository interfaces (auth/user), health extensible, rate limit per user.

| Step | Item | File yang diubah | Tindakan |
|------|------|------------------|----------|
| 3.1 | Repository interface – User | `internal/domain/repositories/user_repo.go`, `internal/app/services/auth/` | Definisikan interface `UserRepository` (GetUserByEmail, CreateUser, UpdateUser, GetUserByRefreshToken, GetUserByPasswordResetToken). Implementasi concrete di `user_repo.go`. Inject ke `AuthService` via constructor; di router/wire, pass concrete. |
| 3.2 | Unit test auth dengan mock | `tests/unit/services/auth/` (atau existing auth_service_test) | Gunakan mock `UserRepository` (manual atau mockery) untuk test Register/Login tanpa DB. |
| 3.3 | Health check extensible | `internal/app/services/health_service.go`, `internal/app/dto/health_dto.go` | Definisikan interface `HealthChecker` (e.g. `Check(ctx) (string, error)`). HealthService punya `[]HealthChecker`; register DB checker. `CheckHealth` iterate checker, isi map. Nanti bisa tambah Redis/HTTP checker. |
| 3.4 | Rate limit per user | `internal/app/middlewares/rate_limit.go`, config | Tambah opsi: jika ada user_id di context (setelah auth), pakai key `user:{id}` untuk limiter; else tetap per IP. Env `RATE_LIMIT_USE_USER=true`. |

**Verifikasi Fase 3:** Unit test auth dengan mock lulus; health tetap 200/503; rate limit by user bisa diuji manual/curl.

---

## Fase 4: Larger – Config & observability (2+ hari)

**Tujuan:** Config terpusat, tracing (opsional).

| Step | Item | File yang diubah | Tindakan |
|------|------|------------------|----------|
| 4.1 | Structured config | `pkg/config/config.go`, `pkg/config/server.go`, `pkg/config/db.go` | Setelah `Unmarshal`, simpan ke global struct atau getter `Config() *Configuration`. Ganti panggilan `viper.GetString("SERVER_PORT")` dll. di app menjadi `config.Get().Server.Port` (atau getter per section). Lakukan bertahap per package. |
| 4.2 | Tracing (opsional) | `pkg/logger` atau pkg baru `pkg/trace` | Integrasi OpenTelemetry: init tracer, middleware inject span, log trace_id. Feature flag `ENABLE_TRACING=false`. Dokumen di OBSERVABILITY.md. |

**Verifikasi Fase 4:** App jalan dengan config struct; jika tracing on, trace_id muncul di log.

---

## Checklist eksekusi

Gunakan checklist ini saat mengerjakan (copy ke issue/PR description).

**Fase 1 – Quick wins**
- [ ] 1.1 Ganti DB_LOG_MODE → MASTER_DB_LOG_MODE di database.go
- [ ] 1.2 Default SERVER_TIMEZONE = UTC di main.go
- [ ] 1.3 Implement graceful shutdown di main.go + doc
- [ ] 1.4 Perbaiki err handling di database.go
- [ ] Build & test lulus

**Fase 2 – Stabilitas & testing**
- [ ] 2.1 Request timeout (middleware atau Server)
- [ ] 2.2 Integration test GET /health
- [ ] 2.3 (Opsional) Integration test DB
- [ ] 2.4 Dokumen error strategy

**Fase 3 – Testability & fitur**
- [ ] 3.1 UserRepository interface + inject ke AuthService
- [ ] 3.2 Unit test auth dengan mock
- [ ] 3.3 HealthChecker interface + register DB
- [ ] 3.4 Rate limit per user (opsional)

**Fase 4 – Config & observability**
- [ ] 4.1 Config struct/getter, kurangi viper di app
- [ ] 4.2 (Opsional) OpenTelemetry + doc

---

## Urutan rekomendasi

1. **Sekarang:** Fase 1 (quick wins) – risiko rendah, manfaat langsung.
2. **Berikutnya:** Fase 2 (timeout + integration test + error doc).
3. **Sesuai kebutuhan:** Fase 3 untuk testability dan fitur (auth mock, health extensible, rate limit user).
4. **Jangka panjang:** Fase 4 (structured config, tracing).

Link ke detail tiap improvement: [IMPROVEMENTS.md](IMPROVEMENTS.md).
