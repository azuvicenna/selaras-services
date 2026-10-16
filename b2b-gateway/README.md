# Selaras SIMRS - B2B Partner API Gateway

Layanan API Gateway Enterprise (Business-to-Business) untuk integrasi rumah sakit dengan mitra asuransi kesehatan swasta (**Prudential, Allianz, AdMedika, AIA, Sinarmas, dll.**), penjamin korporat (Corporate Medical Check-Up), dan laboratorium / faskes rekanan RS Selaras Medika.

---

## 🏢 Arsitektur & Fitur Utama

1. **Partner Authentication & HMAC-SHA256 Cryptography**:
   - Autentikasi berbasis `X-Partner-Id` dan `X-API-Key`.
   - Verifikasi integritas pesan anti-tampering dan anti-replay attacks via HMAC-SHA256:
     $$\text{Signature} = \text{Hex}\Big(\text{HMAC-SHA256}\big(\text{SecretKey},\, \text{PartnerID} + \text{":"} + \text{Timestamp} + \text{":"} + \text{Method} + \text{":"} + \text{Path}\big)\Big)$$
   - Toleransi perbedaan waktu (timestamp skew) maksimal 300 detik (5 menit).
2. **Per-Partner Token Bucket Rate Limiting**:
   - Melindungi downstream services SIMRS dari lonjakan panggilan batch atau skrip mitra.
   - Mengembalikan `HTTP 429 Too Many Requests` dan header `Retry-After: 1` ketika kuota terlampaui.
3. **Asuransi Swasta & Pre-Authorization (Cashless)**:
   - `POST /api/v1/b2b/insurance/eligibility`: Pengecekan eligibilitas peserta, status kepesertaan, limit tahunan, deductible, dan hak kelas rawat.
   - `POST /api/v1/b2b/insurance/pre-authorization`: Pengajuan pre-approval tindakan medis / rawat inap terencana.
4. **Manajemen Klaim (Claims Processing)**:
   - `POST /api/v1/b2b/insurance/claims`: Pengajuan klaim tagihan pasien beserta nomor pre-auth dan diagnosa ICD-10.
   - `GET /api/v1/b2b/insurance/claims`: Pelacakan status verifikasi dan pembayaran klaim secara real-time.
5. **Corporate Medical Check-Up (MCU)**:
   - `POST /api/v1/b2b/corporate/mcu/register`: Pendaftaran massal (batch) karyawan korporat untuk paket pemeriksaan MCU.
   - `GET /api/v1/b2b/corporate/mcu/batches`: Monitoring progres pelaksanaan batch MCU.
6. **Webhooks Asinkron**:
   - `POST /api/v1/b2b/webhooks`: Pendaftaran callback URL mitra untuk notifikasi otomatis event klaim disetujui, pasien pulang (*discharge*), atau hasil MCU selesai.
   - `GET /api/v1/b2b/webhooks`: Daftar webhook aktif milik mitra.
7. **Observabilitas & Dokumentasi API**:
   - `GET /healthz` & `GET /health` - Health check gateway.
   - `GET /metrics` - Prometheus metrics exporter.
   - `GET /docs` - Swagger UI dokumentasi API interaktif.
   - `GET /openapi.json` - OpenAPI 3.0 specification file.

---

## 🔑 Header Autentikasi B2B

| Header | Wajib | Keterangan | Contoh |
| :--- | :---: | :--- | :--- |
| `X-Partner-Id` | Ya | Kode unik identitas mitra | `PRU-INDONESIA`, `ALLIANZ-ID`, `ADMEDIKA-TPA`, `TELKOM-CORP` |
| `X-API-Key` | Ya | Kunci API rahasia mitra | `pru-api-key-live-2026` |
| `X-Timestamp` | Opsional* | UNIX timestamp dalam detik | `1773665420` *(Wajib jika signature check aktif)* |
| `X-Signature` | Opsional* | Hex-encoded HMAC-SHA256 signature | `e3b0c44298fc1c149afbf4c8996fb924...` |

---

## 🚀 Panduan Menjalankan

### Konfigurasi Environment (`.env`)

```env
HTTP_PORT=8082
METRICS_PORT=9092
ENV=development
CORS_ALLOWED_ORIGINS=*
JWT_SECRET=selaras-hospital-jwt-secret-key-32bytes!
B2B_API_KEY_SECRET=selaras-b2b-secret-salt-key-32bytes!
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
ENABLE_SIGNATURE_CHECK=false

# Downstream gRPC Services
AUTH_SERVICE_ADDR=localhost:50057
PATIENT_SERVICE_ADDR=localhost:50051
PRACTITIONER_SERVICE_ADDR=localhost:50056
ADMISSION_SERVICE_ADDR=localhost:50055
BED_SERVICE_ADDR=localhost:50054
MASTER_SERVICE_ADDR=localhost:50053
EMR_SERVICE_ADDR=localhost:50052
BILLING_SERVICE_ADDR=localhost:50058
PHARMACY_SERVICE_ADDR=localhost:50059
LABORATORY_SERVICE_ADDR=localhost:50060
RADIOLOGY_SERVICE_ADDR=localhost:50061
DOCUMENT_SERVICE_ADDR=localhost:50062
```

### Menjalankan secara Mandiri (Local Go)

```bash
# 1. Jalankan unit test
go test -v -count=1 ./b2b-gateway/...

# 2. Build dan jalankan server
go run ./b2b-gateway/cmd/main.go
```

### Menjalankan via Docker Compose

```bash
# Dari root repositori
docker compose --profile apps up -d b2b-gateway
```

Akses portal dokumentasi API:
- Swagger UI: `http://localhost:8082/docs`
- Health check: `http://localhost:8082/healthz`
