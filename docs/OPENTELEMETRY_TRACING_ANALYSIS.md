# Analisis Pengerjaan: OpenTelemetry Tracing

> Dokumen untuk diskusi sebelum eksekusi. Bukan implementasi.

**Kesimpulan singkat:** Untuk **satu service (monolith)**, **request_id saja sudah cukup**. Tidak wajib menambah trace_id/OpenTelemetry. OTel/trace_id baru relevan kalau nanti mau pakai Jaeger/Tempo atau ada banyak service (microservices). Implementasi di bawah bersifat opsional.

**Catatan penilaian:** Jika ada penilaian "tracing 6/10" karena belum ada OpenTelemetry, itu dari sudut pandang **distributed tracing**. Untuk konteks monolith, request_id + propagasi context + LogStart/LogFinish + FromContext + ARRIVED/RESPONSE SENT **sudah memadai** dan pantas dinilai lebih tinggi (mis. 8–9/10). Nilai 6 tidak seharusnya dipakai untuk menilai tracing monolith yang hanya pakai request_id.

---

## 1. Kondisi saat ini

| Komponen | Status |
|----------|--------|
| **Request ID** | Sudah: `RequestIDMiddleware` → UUID / `X-Request-ID`, disimpan di `gin.Context` + `c.Request.Context()` via `logger.RequestIDContextKey`. |
| **Logger** | `pkg/logger`: formatter sudah cetak `[request_id]`, `LogStart`/`LogFinish`, `FromContext(ctx)` baca request_id dari context. |
| **Urutan middleware** | `RequestID` → `RequestLog` → `Metrics`. Request log dan metrics bergantung pada request_id. |
| **Config** | Viper + `config.SetupConfig()`, validasi required keys. Belum ada env khusus tracing. |
| **Shutdown** | Graceful shutdown di `main.go` (server, DB). Belum ada shutdown tracer. |

Kesimpulan: fondasi request-scoped logging dan request_id sudah rapi. Yang ditambah nanti: **trace_id/span_id + export ke backend (OTLP/Jaeger)**, dan tetap **tetap memakai request_id** (bisa dipetakan ke trace_id atau hidup berdampingan).

---

## 1.1 Request ID vs Trace ID — kapan beda, kapan mirip?

Di **satu service** (monolith), request_id dan trace_id memang mirip: keduanya “identitas satu request”. Bedanya baru terasa kalau ada **banyak service** dan kamu mau pakai tools seperti Jaeger.

| | Request ID | Trace ID |
|---|------------|----------|
| **Format** | Bebas (biasanya UUID), e.g. `550e8400-e29b-41d4-a716-446655440000` | Standar OTel/W3C (biasanya 32 char hex), e.g. `4bf92f3577b34da6a3ce929d0e0e4736` |
| **Siapa yang paham** | Aplikasi kita + log (grep, Loki, ELK) | Aplikasi + **Jaeger/Tempo/backend tracing** (bisa gambar timeline, latency per langkah) |
| **Satu service** | Satu request = satu request_id. Cukup untuk korelasi log. | Satu request = satu trace_id (+ span_id). Bisa dipetakan 1:1 dengan request_id. **Fungsinya mirip.** |
| **Banyak service** | Bisa satu ID yang di-forward (X-Request-ID) ke semua service → korelasi log lintas service. | Satu **trace_id** untuk seluruh alur; tiap service punya **span_id** sendiri. Backend tracing bisa gambar pohon: Gateway → Auth → User → Payment. |

**Kesimpulan singkat:**  
- **Hanya satu service (monolith):** request_id saja sudah cukup untuk korelasi log. Trace_id “tambahan” kalau mau nanti integrasi Jaeger atau standar OTel; bisa juga trace_id di-set sama dengan request_id.  
- **Banyak service (microservices):** trace_id + span_id dipakai agar Jaeger/Tempo bisa mengelompokkan semua span dalam satu trace dan menampilkan timeline. Request_id tetap berguna untuk “ID yang manusia baca” (support, ticket).

Contoh konkret di bawah.

---

### Contoh 1: Satu service (API kita saja)

Satu HTTP request masuk → satu request_id, satu trace (satu span).

```
Request:  POST /auth/login
request_id: 550e8400-e29b-41d4-a716-446655440000   ← sudah ada sekarang
trace_id:  4bf92f3577b34da6a3ce929d0e0e4736        ← dari OTel (bisa dipetakan ke request_id)
span_id:  00f067aa0ba902b7                           ← satu span “HTTP POST /auth/login”
```

**Log (dengan request_id saja — seperti sekarang):**
```
INFO ... [550e8400-e29b-41d4-a716-446655440000] ARRIVED REQUEST ...
INFO ... [550e8400-e29b-41d4-a716-446655440000] START AuthController.Login
INFO ... [550e8400-e29b-41d4-a716-446655440000] FINISH AuthService.Login (SUCCESS) duration=12ms
INFO ... [550e8400-e29b-41d4-a716-446655440000] RESPONSE SENT status=200 ...
```

**Kalau ditambah trace_id (opsional):**  
Bisa tampil `[request_id=...] [trace_id=...]`. Di monolith, trace_id tidak menambah info baru untuk “satu request” — fungsinya mirip request_id. Manfaatnya: nanti kalau export ke Jaeger, satu “trace” itu = satu request; atau kalau nanti ada service lain yang dipanggil, trace_id bisa dipakai sebagai “benang merah” yang sama.

**Sample penggunaan trace_id di sini:**  
- Tetap pakai **request_id** untuk grep log dan response header `X-Request-ID`.  
- **Trace_id** dipakai saat kamu mau lihat request yang sama di Jaeger (satu baris = satu trace = satu request). Bisa juga di log ditampilkan trace_id supaya satu nilai bisa dipakai baik di log maupun di Jaeger.

---

### Contoh 2: Banyak service (Gateway → Auth → User)

Satu request dari user melewati 3 service. Request_id bisa satu yang di-forward; trace_id satu untuk seluruh perjalanan; tiap service punya span_id sendiri.

```
User → Gateway → Auth Service → User Service
         (span_id=S1)   (span_id=S2)   (span_id=S3)
         trace_id = T1 (sama di semua)
         request_id = R1 (sama di semua, dari header X-Request-ID)
```

**Log di Gateway:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] calling Auth ...
```
**Log di Auth:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] validating token ...
INFO ... [request_id=R1] [trace_id=T1] calling User Service ...
```
**Log di User Service:**
```
INFO ... [request_id=R1] [trace_id=T1] ARRIVED REQUEST ...
INFO ... [request_id=R1] [trace_id=T1] fetch user from DB ...
```

Di **log**: request_id (atau trace_id) dipakai untuk grep satu alur.  
Di **Jaeger**: trace_id T1 dipakai untuk mengelompokkan tiga span (S1, S2, S3) jadi satu “trace” dan menampilkan timeline + latency per service (mis. Gateway 2ms, Auth 50ms, User 10ms). Itu yang request_id saja tidak bisa berikan (bukan format yang dipahami Jaeger).

**Sample penggunaan trace_id di sini:**  
- **Request_id R1:** untuk manusia (support: “kirim request ID R1”), dan korelasi log sederhana (grep R1).  
- **Trace_id T1:** untuk tool (Jaeger/Tempo): buka trace T1 → lihat semua span S1, S2, S3 dalam satu view dan lihat di service mana waktu habis.

---

### Ringkasan sample penggunaan

| Situasi | Request ID | Trace ID |
|--------|------------|----------|
| Satu service, cuma mau korelasi log | Cukup. Grep request_id, tampilkan di X-Request-ID. | **Tidak wajib.** Boleh dipakai kalau mau siap Jaeger atau standar OTel; bisa di-set = request_id. |
| Satu service, mau pakai Jaeger | Tetap dipakai untuk log & response. | Dipakai oleh Jaeger sebagai “satu trace = satu request”. |
| Banyak service, korelasi log | Satu ID di-forward (X-Request-ID) ke semua service. | Satu ID (trace_id) + span per service; Jaeger pakai untuk timeline & latency per service. |

Jadi: trace_id **tidak menggantikan** request_id. Request_id tetap untuk “satu request” di log dan response. Trace_id tambahan untuk **ecosystem tracing** (standar, format, dan tools seperti Jaeger), dan baru “wajib” kalau kamu mau visualisasi lintas service atau integrasi dengan backend tracing.

---

## 2. Yang perlu dikerjakan (scope pengerjaan)

### 2.1 Dependensi

- `go.opentelemetry.io/otel`
- `go.opentelemetry.io/otel/sdk` (trace)
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp` (atau `otlptracegrpc`) — export ke collector
- **Opsional:** `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin` — middleware Gin resmi OTel, atau kita bikin middleware sendiri yang tipis.

Kalau mau **tanpa backend dulu** (hanya inject trace_id ke context & log): cukup SDK trace + `trace.NewNoopTracerProvider()` atau provider dengan `NewBatchSpanProcessor(NewNoopExporter())` sehingga tidak ada jaringan. Nanti saat ada Jaeger/OTLP collector, ganti ke exporter sungguhan.

### 2.2 Inisialisasi tracer (main.go)

- **Feature flag:** mis. `ENABLE_TRACING` (default `false`).
- Jika `true`:
  - Buat `TracerProvider` (resource: service name, version).
  - Set global `otel.SetTracerProvider(...)`.
  - Opsional: set propagator `otel.SetTextMapPropagator(propagation.TraceContext{})` agar trace_id/span_id ikut ke HTTP header ke service lain.
- Jika `false`: pakai noop provider (atau tidak set apa-apa; SDK biasanya noop by default).
- **Shutdown:** saat graceful shutdown, panggil `TracerProvider.Shutdown(ctx)` agar span yang belum ter-export tetap dikirim (dengan timeout).

Ini menambah sedikit kode di `main.go` dan satu blok if/else atau fungsi `initTracing() bool`.

### 2.3 Middleware tracing

- **Lokasi:** middleware baru, mis. `internal/app/middlewares/tracing.go` (atau `otelgin` dari contrib).
- **Urutan:** tracing harus bisa akses request_id dan context. Opsi:
  - **Opsi A:** RequestID → **Tracing** → RequestLog → Metrics (trace_id/span_id tersedia untuk log).
  - **Opsi B:** Tracing → RequestID → RequestLog → Metrics (trace_id dihasilkan di awal, request_id tetap untuk backward compatibility).
- **Perilaku:**
  - Start span per request (nama mis. `HTTP GET /api/...`).
  - Inject span context ke `c.Request.Context()` (dan ke gin context jika perlu).
  - Set response header `X-Trace-ID` (opsional) agar client bisa referensi.
  - Defer `span.End()`.
- **Sampling:** bisa pakai `ParentBased(root=AlwaysSample)` atau `TraceIDRatioBased(0.1)` dan baca ratio dari env (mis. `TRACE_SAMPLE_RATIO=1.0` atau `0.1`). Saat `ENABLE_TRACING=false`, middleware bisa no-op (langsung `c.Next()` tanpa membuat span).

Ini tidak mengubah kontrak controller/service; mereka hanya dapat context yang sudah berisi span (dan trace_id bisa diambil dari context untuk log).

### 2.4 Integrasi dengan logger (trace_id / span_id di log)

- **Tujuan:** agar log aggregator (Loki, ELK) bisa korelasi dengan trace backend (Jaeger) via trace_id.
- **Opsi:**
  - **Minimal:** middleware set `trace_id` (dan opsional `span_id`) di gin context; di `pkg/logger` formatter atau helper baca dari context (mirip `request_id`) dan cetak `[trace_id=...]` atau gabung dengan request_id.
  - **Lanjutan:** `logger.FromContext(ctx)` bisa baca span dari `otelspan.SpanFromContext(ctx)` dan ambil `SpanContext().TraceID().String()` lalu tambah ke `logrus.Fields`. Dengan begitu setiap log yang pakai `FromContext(ctx)` otomatis dapat trace_id/span_id tanpa ubah tiap handler.
- **Backward compatibility:** request_id tetap ada dan tetap dicetak; trace_id tambahan. Saat tracing dimatikan, trace_id kosong/tidak dicetak.

Perubahan terbatas di `pkg/logger` (context key opsional, formatter, dan `FromContext`) plus memastikan controller/service selalu pakai `ctx` dari request (sudah begitu).

### 2.5 Konfigurasi (env)

- `ENABLE_TRACING` — boolean, default `false`.
- `OTEL_SERVICE_NAME` — nama service (default dari binary/app name).
- `OTEL_EXPORTER_OTLP_ENDPOINT` — opsional; kalau kosong bisa pakai noop exporter atau tidak export (hanya in-process trace_id di log).
- `TRACE_SAMPLE_RATIO` — opsional, 0.0–1.0 (default 1.0 jika tracing on).

Tidak perlu masuk ke `ValidateConfig()` required keys; semua optional. Cukup baca di `main.go` atau di helper init tracing.

### 2.6 Outbound HTTP (propagasi ke service lain)

- Kalau nanti ada client HTTP ke service lain (auth, payment, dll.), inject trace context ke header dengan `otelgin` atau `propagation.Inject(ctx, carrier)` agar satu trace_id melewati banyak service. Ini bisa fase kedua setelah middleware in-process jalan.

---

## 3. Opsi scope (untuk diskusi)

| Opsi | Isi | Effort | Cocok kalau |
|------|-----|--------|-------------|
| **A. Minimal (trace_id di log saja)** | Feature flag, noop provider atau provider tanpa exporter, middleware inject span ke context, logger baca trace_id dari context dan cetak di log. Tidak ada OTLP/Jaeger. | Kecil | Mau trace_id di log dulu, deploy backend tracing belakangan. |
| **B. Lengkap in-process + export** | Seperti A + OTLP HTTP exporter, env endpoint, sampling, shutdown. | Sedang | Sudah ada/mau pasang collector (Jaeger/OTLP). |
| **C. B + propagasi outbound** | Seperti B + dokumentasi/contoh inject ke HTTP client ke service lain. | Sedang+ | Multi-service dan mau satu trace across services. |

Rekomendasi untuk diskusi: mulai **A** atau **B**. Kalau belum ada Jaeger/collector, **A** cukup; nanti tinggal ganti ke exporter sungguhan.

---

## 4. Risiko & mitigasi

| Risiko | Mitigasi |
|--------|----------|
| **Overhead** | ENABLE_TRACING=false default; sampling (TRACE_SAMPLE_RATIO); export async (batch) by default di SDK. |
| **Tracing nyala tapi endpoint kosong** | Cek env; kalau endpoint kosong pakai noop exporter atau log warning dan disable export. |
| **Breaking change** | Tidak ubah signature controller/service; hanya context jadi berisi span. Request_id tetap. |
| **Ketergantungan baru** | Pin versi OTel di go.mod; tetap gunakan go 1.25. |

---

## 5. File yang kemungkinan tersentuh

| File / area | Perubahan |
|-------------|-----------|
| `go.mod` / `go.sum` | Tambah dependency OTel. |
| `main.go` | Init tracer (conditional), shutdown `TracerProvider`. |
| `internal/app/middlewares/tracing.go` | Baru: middleware span + inject context. |
| `internal/app/routers/router.go` | Register tracing middleware (setelah RequestID). |
| `pkg/logger/logger.go` | Opsional: baca trace_id/span_id dari context, tambah ke formatter / FromContext. |
| `docs/OBSERVABILITY.md` | Update: langkah konfigurasi, env, contoh log dengan trace_id. |
| `.env.example` | Tambah ENABLE_TRACING, OTEL_SERVICE_NAME, OTEL_EXPORTER_OTLP_ENDPOINT, TRACE_SAMPLE_RATIO (optional). |

Tidak wajib ubah controller/service kalau kita hanya inject context di middleware dan logger baca dari context.

---

## 6. Checklist diskusi (putuskan dulu)

Ini lima keputusan yang perlu diambil sebelum coding. Penjelasan singkat per poin:

---

### 6.1 Scope: A, B, atau C?

**Maksud:** Sejauh apa fitur tracing yang mau kita pasang sekarang?

| Pilihan | Arti | Contoh |
|--------|------|--------|
| **A** | Trace ID cuma muncul di **log** (file/stdout). Tidak kirim data ke mana-mana. | Nanti di log ada `[request_id=...] [trace_id=abc123]`. Cari error pakai trace_id. Tidak perlu pasang Jaeger. |
| **B** | Seperti A + **kirim span ke backend** (Jaeger, Grafana Tempo, dll.) lewat OTLP. | Bisa lihat timeline request di UI Jaeger. Perlu jalankan collector/backend. |
| **C** | Seperti B + kalau API kita **panggil service lain** (HTTP), trace_id ikut di header sehingga satu trace bisa lintas service. | Untuk arsitektur microservices; satu request dari gateway → auth → payment tetap satu trace. |

**Saran:** Belum punya Jaeger/collector → pilih **A**. Sudah ada atau mau pasang → **B**. Multi-service dan mau satu trace → **C**.

---

### 6.2 Middleware: otelgin vs custom?

**Maksud:** Pakai library jadi atau bikin middleware sendiri?

| Pilihan | Arti |
|--------|------|
| **otelgin** | Pakai `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`. Sudah siap: start span, inject context, propagasi. Kurang kontrol halus (mis. nama span). |
| **Custom tipis** | Kita sendiri bikin middleware: start span, inject ke `c.Request.Context()`, defer End. Kontrol penuh, kode sedikit lebih banyak. |

**Saran:** Mau cepat dan standar → **otelgin**. Mau minimal dependency atau nama span spesifik → **custom**.

---

### 6.3 Logger: trace_id di setiap baris atau cuma ARRIVED/RESPONSE?

**Maksud:** Di mana saja trace_id (dan span_id) muncul di log?

| Pilihan | Arti |
|--------|------|
| **Setiap baris** | Semua log (START, FINISH, FromContext, dll.) otomatis dapat field `trace_id` (mirip `request_id`). Satu trace_id bisa dipakai grep di seluruh log. |
| **Hanya ARRIVED / RESPONSE SENT** | Trace_id cuma di baris yang dikeluarkan middleware (ARRIVED REQUEST, RESPONSE SENT). Baris lain tetap pakai request_id saja. |

**Saran:** Supaya konsisten dengan request_id dan enak korelasi log–trace → **setiap baris**.

---

### 6.4 Default ENABLE_TRACING: false atau true (development)?

**Maksud:** Kalau env `ENABLE_TRACING` tidak di-set, tracing nyala atau mati?

| Pilihan | Arti |
|--------|------|
| **false (default)** | Tracing mati sampai kamu set `ENABLE_TRACING=true`. Aman untuk production yang belum siap. |
| **true untuk dev** | Bisa bikin: di development default true, di production default false (baca dari APP_ENV). Development langsung dapat trace_id tanpa ubah .env. |

**Saran:** Tetap **false** global biar tidak kaget; yang mau nyalakan set explicit `ENABLE_TRACING=true`.

---

### 6.5 Nama env: setuju atau mau diganti?

**Maksud:** Nama variabel environment yang dipakai.

| Env yang diusulkan | Fungsi |
|--------------------|--------|
| `ENABLE_TRACING` | Nyala/mati tracing (true/false). |
| `OTEL_SERVICE_NAME` | Nama service di trace (mis. `boilerplate-gin-api`). Kalau kosong, pakai nama binary. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | URL collector (mis. `http://jaeger:4318`). Kalau kosong dan scope A → tidak kirim ke mana-mana. |
| `TRACE_SAMPLE_RATIO` | Berapa persen request yang di-trace (0.0–1.0). Mis. 0.1 = 10%. Untuk kurangi beban. |

**Saran:** Nama itu standar (OTEL_*) dan cocok dengan ecosystem OpenTelemetry. Kalau mau konsisten dengan konvensi project (mis. prefix `TRACING_`), bisa dipakai nama lain asal didokumentasikan.

---

### Ringkasan: isi checklist

- [ ] **Scope:** A / B / C
- [ ] **Middleware:** otelgin / custom
- [ ] **Logger:** trace_id setiap baris / hanya ARRIVED & RESPONSE SENT
- [ ] **Default ENABLE_TRACING:** false / true (dev only)
- [ ] **Nama env:** pakai usulan di atas / mau ganti (sebutkan)

Kalau poin-poin di atas sudah disepakati, langkah berikutnya bisa dirinci jadi task per file atau langsung ke implementasi bertahap.

---

**Last updated:** 2025-02-03 (analisis untuk diskusi, belum implementasi)
