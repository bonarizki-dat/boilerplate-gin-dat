# Daftar Improvement Boilerplate

Dokumen ini berisi usulan perbaikan yang bisa diterapkan bertahap. Prioritas: **Quick wins** → **Medium** → **Larger**.

**Rencana implementasi (langkah per langkah, file yang diubah, checklist):** [IMPROVEMENT_PLAN.md](IMPROVEMENT_PLAN.md)

---

## Quick wins (low effort, high impact)

### 1. Seragamkan nama env DB log mode
- **Masalah:** `internal/adapters/database/database.go` pakai `viper.GetBool("DB_LOG_MODE")`, sementara `.env.example` dan docs pakai `MASTER_DB_LOG_MODE`.
- **Tindakan:** Ganti di `database.go` jadi `MASTER_DB_LOG_MODE` (atau sebaliknya satu nama di semua tempat). Disarankan pakai `MASTER_DB_LOG_MODE` agar konsisten dengan env DB lain.

### 2. Default timezone selaras dengan .env.example
- **Masalah:** `main.go` set default `SERVER_TIMEZONE` = `"Asia/Dhaka"`, `.env.example` pakai `SERVER_TIMEZONE=UTC`.
- **Tindakan:** Default di code pakai `"UTC"` (atau hapus default dan wajibkan env), supaya dev lokal konsisten dengan contoh env.

### 3. Graceful shutdown
- **Masalah:** `main.go` hanya `router.Run()`; saat SIGTERM/SIGINT proses exit tanpa menutup DB dan server dengan rapi.
- **Tindakan:** Di `main.go`: listen `os.Interrupt` / `syscall.SIGTERM`, panggil `server.Shutdown(ctx)` (pakai `http.Server` + `ListenAndServe`), lalu tutup koneksi DB (mis. `sqlDB.Close()` dari `database.DB.DB()`). Dokumentasi di README/CONFIGURATION.

### 4. Perbaikan kecil di database adapter
- **Masalah:** `database.go` punya `var err error` di package level; pengecekan `if err != nil` ada setelah blok yang bisa mengubah `err`.
- **Tindakan:** Pakai `err :=` lokal di dalam fungsi (mis. `db, err := gorm.Open(...)`), dan pastikan pengecekan `err` tepat setelah operasi yang mengisi `err`.

---

## Medium (beberapa jam–hari)

### 5. Integration test contoh
- **Masalah:** `tests/integration/api/` dan `tests/integration/database/` hanya berisi `.gitkeep`; tidak ada contoh integration test.
- **Tindakan:** Tambah 1–2 contoh: e.g. integration test `GET /health` (tanpa DB mock) dan/atau test koneksi DB dengan testcontainer atau DB sementara. Dokumentasi singkat di `tests/README.md`.

### 6. Repository interfaces (optional untuk testability)
- **Masalah:** Repository berupa function-based concrete (e.g. `repositories.GetUserByEmail`); untuk unit test service yang ketat perlu mock.
- **Tindakan:** (Opsional) Per fitur yang butuh mock, definisikan interface (e.g. `UserRepository`) di `domain/repositories` dan inject ke service. Implementasi tetap di `internal/domain/repositories`. Bisa bertahap (mulai auth/user saja).

### 7. Pemakaian pkg/types errors (konsistensi)
- **Masalah:** `pkg/types/errors.go` mendefinisikan `APIError`, `ErrNotFound`, dll., tapi controller saat ini pakai `utils.BadRequest`/`auth.Err*` dan tidak memetakan ke types tersebut.
- **Tindakan:** Pilih salah satu: (a) pakai `types.APIError` / `types.Err*` di layer service/controller dan map ke response utils, atau (b) dokumentasikan bahwa `pkg/types/errors` dipakai untuk kasus tertentu (e.g. client SDK) dan response HTTP tetep lewat `pkg/utils`. Hindari duplikasi konsep error.

### 8. Request timeout
- **Masalah:** Tidak ada timeout global per request; request yang hang bisa menahan worker.
- **Tindakan:** Tambah middleware timeout (e.g. `context.WithTimeout` 30s) untuk request, atau set `ReadTimeout`/`WriteTimeout` di `http.Server` saat implement graceful shutdown.

---

## Larger (refactor / fitur baru)

### 9. Structured config (pakai struct yang di-Unmarshal)
- **Masalah:** `config.SetupConfig()` unmarshal ke `Configuration` struct tapi kebanyakan akses tetap lewat `viper.GetString` di banyak tempat.
- **Tindakan:** Ekspos config lewat struct global atau getter yang baca dari struct (e.g. `config.Get().Server.Port`), dan kurangi ketergantungan langsung ke viper di app code. Validasi tetap di `ValidateConfig()`.

### 10. Tracing (OpenTelemetry atau X-Request-ID lengkap)
- **Masalah:** Sudah ada request ID; belum ada distributed tracing (trace/span) untuk observability di production.
- **Tindakan:** Tambah opsi integrasi OpenTelemetry (trace ID + span), atau dokumentasikan pola propagate request ID ke log dan ke downstream service. Bisa optional (feature flag / build tag).

### 11. Health check dependency detail
- **Masalah:** Health hanya cek DB; tidak ada endpoint atau field untuk dependency lain (e.g. cache, external API) di masa depan.
- **Tindakan:** Rancang `HealthService` agar bisa register checker per dependency (interface `HealthChecker`) dan agregasi status; endpoint `/health` tetap satu, isi bisa diperluas.

### 12. Rate limit per user (selain per IP)
- **Masalah:** Rate limit saat ini per IP; untuk API key atau JWT mungkin ingin limit per user.
- **Tindakan:** Tambah opsi middleware rate limit by user ID (dari JWT) atau by API key, dengan fallback ke per-IP. Bisa konfigurasi via env.

---

## Ringkasan prioritas

| # | Item                         | Effort | Impact |
|---|------------------------------|--------|--------|
| 1 | DB_LOG_MODE → MASTER_DB_LOG_MODE | Kecil  | Konsistensi |
| 2 | Default timezone UTC         | Kecil  | Konsistensi |
| 3 | Graceful shutdown            | Kecil  | Production-ready |
| 4 | Perbaikan err di database.go | Kecil  | Correctness |
| 5 | Contoh integration test      | Medium | DX + CI |
| 6 | Repository interfaces        | Medium | Testability |
| 7 | Pemakaian types errors       | Medium | Konsistensi |
| 8 | Request timeout              | Medium | Stability |
| 9 | Structured config            | Besar  | Maintainability |
| 10| Tracing                      | Besar  | Observability |
| 11| Health check extensible      | Medium | Scalability |
| 12| Rate limit per user          | Medium | Fitur |

---

**Cara pakai:** Pilih nomor yang mau dikerjakan, lalu buat task/PR per item. Untuk quick wins (1–4), bisa digabung dalam satu PR.
