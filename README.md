# Selaras Hospital Information System (SIMRS) - Backend Services

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![gRPC](https://img.shields.io/badge/gRPC-v1.79-244c5a?style=flat&logo=grpc)](https://grpc.io)
[![Buf](https://img.shields.io/badge/Protobuf-Buf%20v2-blue?style=flat&logo=buf)](https://buf.build)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![MinIO](https://img.shields.io/badge/MinIO-Object%20Storage-c72c48?style=flat&logo=minio)](https://min.io)
[![Redis](https://img.shields.io/badge/Redis-7-dc382d?style=flat&logo=redis)](https://redis.io)

**Selaras** adalah platform backend Sistem Informasi Manajemen Rumah Sakit (SIMRS) modern berbasis arsitektur *microservices* terdistribusi, dirancang untuk skalabilitas tinggi, kepatuhan regulasi Rekam Medis Elektronik (RME Permenkes No. 24/2022), serta integrasi Kemenkes SATUSEHAT dan BPJS Kesehatan.

---

## 🏛️ Arsitektur Sistem

```
                      +-----------------------------+
                      |   Klien Web / Mobile Apps   |
                      +-----------------------------+
                                     |
                                     v [HTTP / REST / WebSocket]
            +-------------------------------------------------+
            |     BFF & Gateway Layer                         |
            |     • bff-web (SIMRS Portal Kasir, Poli, dll.)  |
            |     • bff-mobile (Aplikasi Pasien & Dokter)     |
            |     • b2b-gateway (Integrasi Asuransi / Faskes) |
            +-------------------------------------------------+
                                     |
                                     v [gRPC / Protobuf]
+-----------------------------------------------------------------------------+
|                       Core Clinical & Admin Microservices                   |
|                                                                             |
|  [auth-service]        [patient-service]     [practitioner-service]         |
|  :50057 (JWT/RBAC)     :50051 (Demografi)    :50056 (Dokter, Nakes, Bidan)  |
|                                                                             |
|  [admission-service]   [bed-service]         [master-service]               |
|  :50055 (Pendaftaran)  :50054 (Bangsal)      :50053 (ICD-10, ICD-9-CM, Wil) |
|                                                                             |
|  [emr-service]         [billing-service]     [pharmacy-service]             |
|  :50052 (RME & SOAP)   :50058 (Kasir & Tarif):50059 (Obat, Resep, Stok)    |
|                                                                             |
|  [laboratory-service]  [radiology-service]   [document-service]             |
|  :50060 (LIS, LOINC)   :50061 (RIS, DICOM)   :50062 (20 Typst PDF Templates)|
+-----------------------------------------------------------------------------+
             |                                       |
             v [Event Streaming / Broker]             v [Object Storage]
  +-------------------------+             +----------------------+
  | RabbitMQ Message Broker |             | MinIO Object Storage |
  +-------------------------+             +----------------------+
             |                               (Berkas Pasien & Dokumen)
             v
+---------------------------------------+
|      Asynchronous Workers             |
|  • bpjs-worker (VClaim, PCare, SEP)   |
|  • satusehat-worker (FHIR Sync)       |
|  • notification-worker (WA, Email)    |
+---------------------------------------+
```

---

## 📊 Matriks Layanan & Port Mapping

| No | Nama Service | Port gRPC | Port Metrik | Database (Postgres) | Bucket MinIO / External | Status |
|:--:|:-------------|:---------:|:-----------:|:-------------------:|:------------------------|:------:|
| 1  | [**`patient-service`**](./patient-service) | `50051` | `9092` | `patient_db` (`patient`) | `patient-documents` | ✅ Selesai |
| 2  | [**`emr-service`**](./emr-service) | `50052` | `9093` | `emr_db` (`emr`) | `emr-documents` | ✅ Selesai |
| 3  | [**`master-service`**](./master-service) | `50053` | `9094` | `master_db` (`master`) | - | ✅ Selesai |
| 4  | [**`bed-service`**](./bed-service) | `50054` | `9095` | `bed_db` (`bed`) | - | ✅ Selesai |
| 5  | [**`admission-service`**](./admission-service) | `50055` | `9096` | `admission_db` (`admission`) | - | ✅ Selesai |
| 6  | [**`practitioner-service`**](./practitioner-service) | `50056` | `9097` | `practitioner_db` (`practitioner`) | - | ✅ Selesai |
| 7  | [**`auth-service`**](./auth-service) | `50057` | `9098` | `auth_db` (`auth`) | Redis (JWT blacklist & Captcha) | ✅ Selesai |
| 8  | [**`billing-service`**](./billing-service) | `50058` | `9099` | `billing_db` (`billing`) | - | ✅ Selesai |
| 9  | [**`pharmacy-service`**](./pharmacy-service) | `50059` | `9100` | `pharmacy_db` (`pharmacy`) | - | ✅ Selesai |
| 10 | [**`laboratory-service`**](./laboratory-service) | `50060` | `9101` | `laboratory_db` (`laboratory`) | - | ✅ Selesai |
| 11 | [**`radiology-service`**](./radiology-service) | `50061` | `9102` | `radiology_db` (`radiology`) | DICOM / PACS URLs | ✅ Selesai |
| 12 | [**`document-service`**](./document-service) | `50062` | `9103` | `document_db` (`document`) | `medical-documents` + Typst CLI | ✅ Selesai |
| 13 | [**`bff-web`**](./bff-web) | HTTP `8080` | `9090` | - | Konsumsi gRPC 12 Services + Swagger UI | ✅ Selesai |
| 14 | [**`bff-mobile`**](./bff-mobile) | HTTP `8081` | `9091` | - | Mobile Backend-For-Frontend, Komposit & Swagger UI | ✅ Selesai |
| 15 | [**`b2b-gateway`**](./b2b-gateway) | HTTP `8082` | `9092` | - | Enterprise B2B Gateway (Asuransi Swasta & MCU) | ✅ Selesai |
| 16 | [**`satusehat-worker`**](./satusehat-worker) | HTTP `8085` | `9095` | - | SATUSEHAT FHIR R4 & RabbitMQ Consumer | ✅ Selesai |
| 17 | [**`bpjs-worker`**](./bpjs-worker) | HTTP `8086` | `9096` | - | BPJS VClaim 2.0, Antrean, Aplicares Consumer | ✅ Selesai |
| 18 | [**`notification-worker`**](./notification-worker) | HTTP `8087` | `9097` | - | WhatsApp & Email via RabbitMQ Consumer | ✅ Selesai |
| 19 | [**`prometheus`**](http://localhost:9089) | - | `9089` | `selaras-prometheus-data` | Time-Series Metrics Scraper | ✅ Selesai |
| 20 | [**`grafana`**](http://localhost:3000) | - | `3000` | `selaras-grafana-data` | Dashboard Visualisasi SIMRS Overview | ✅ Selesai |

---

## 🛠️ Persyaratan Lingkungan (Prerequisites)

- **Go**: Versi 1.26+ (menggunakan [Go Workspaces (`go.work`)](go.work))
- **Docker & Docker Compose**: Docker Engine 24+
- **Buf CLI**: v1.40+ (untuk kompilasi proto: `brew install bufbuild/buf/buf` atau via binary)
- **Typst CLI**: (opsional untuk local generator document-service: `typst --version`)

---

## 🚀 Panduan Memulai Cepat (Quick Start)

### 1. Inisialisasi Dependensi Infrastruktur
Jalankan PostgreSQL, Redis, MinIO (lengkap dengan bucket otomatis), dan RabbitMQ menggunakan Docker Compose:

```bash
# Menggunakan Makefile
make compose-infra

# Atau langsung via Docker Compose
docker compose up -d postgres redis minio minio-init rabbitmq
```

Akses Dashboard Infrastruktur:
- **MinIO Console**: http://localhost:9001 (User: `minioadmin`, Pass: `minioadmin`)
- **RabbitMQ Management**: http://localhost:15672 (User: `selaras`, Pass: `selaras123`)
- **PostgreSQL**: `localhost:5432` (Superuser: `postgres` / `postgres`, password default user service: `change-me`)
- **Prometheus Dashboard**: http://localhost:9089 (Scrape seluruh 18 microservices & RabbitMQ)
- **Grafana SIMRS Dashboard**: http://localhost:3000 (User: `admin` / Pass: `selaras`)

### 2. Sinkronisasi Go Workspaces & Modul
```bash
go work sync
```

### 3. Eksekusi Migrasi Database & Seeding Sekaligus
Sistem dilengkapi dengan **Unified Migration Runner** untuk mengeksekusi migrasi ke-12 database PostgreSQL secara otomatis:
```bash
# Migrasi seluruh 12 database ke versi terbaru
make migrate-all
# Atau via PowerShell (Windows):
.\scripts\migrate.ps1 up
# Atau via Bash (Linux/macOS):
./scripts/migrate.sh up

# Cek status migrasi seluruh 12 database
make migrate-status
# Atau: .\scripts\migrate.ps1 status

# Seeding data master awal (ICD-10, users, poli, bed, tarif, dll.)
make seed-all
# Atau: .\scripts\migrate.ps1 seed
```

### 4. Kompilasi Protobuf (Jika melakukan perubahan kontrak API)
```bash
make proto
# Atau: buf generate
```

### 5. Menjalankan Service di Komputer Lokal
Setiap service memiliki file `.env.example`. Salin menjadi `.env` lalu jalankan secara mandiri:
```bash
cd auth-service
cp .env.example .env
go run ./cmd
```

Atau jalankan seluruh service sekaligus dalam kontainer Docker:
```bash
make compose-apps
# Atau: docker compose --profile apps up -d --build
```

### 6. Menjalankan Pengujian (Testing) & Lint
```bash
# Menjalankan seluruh unit & integration tests di 18 services
make test

# Menjalankan linter golangci-lint
make lint

# Menjalankan End-to-End (E2E) Patient Journey Integration Test
# (Pastikan stack aplikasi aktif terlebih dahulu via: make compose-apps)
make test-e2e
# Atau via interactive CLI runner dengan visualisasi step-by-step:
make run-e2e
# Atau via PowerShell (Windows):
.\scripts\run-e2e.ps1
# Atau via Bash (Linux/macOS):
./scripts/run-e2e.sh
```

---

## 📄 Standar & Konvensi
- **Protokol Komunikasi Internal:** gRPC dengan Protobuf v2 (Strict backward compatibility via Buf breaking checker).
- **Format Respons Database:** PostgreSQL dengan Connection Pool tuning (`pgxpool`) untuk beban produksi rumah sakit.
- **Observability:** Prometheus metrics (`/metrics`) dan Health check endpoint (`/healthz` + gRPC standard health check).
- **Format Dokumen Medis:** Tipografi presisi tinggi berbasis [Typst](https://typst.app) untuk 20 template surat medis resmi (Rujukan, Kematian/Visum, Informed Consent, Resume Medis, Odontogram, Hasil Lab/Radiologi, dll.).
