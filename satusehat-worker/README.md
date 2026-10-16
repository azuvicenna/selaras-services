# SATUSEHAT Worker Microservice

Worker asinkron untuk menyinkronkan data rekam medis klinis (Encounter, Condition, Observation, Procedure, MedicationRequest, DiagnosticReport, AllergyIntolerance, Composition) ke platform **Kemenkes SATUSEHAT** berbasis standar internasional **HL7 FHIR R4**.

Setelah resource berhasil diposting ke SATUSEHAT dan mendapatkan FHIR Resource ID dari Kemenkes, worker secara otomatis memanggil callback gRPC ke `emr-service` (`SyncSatusehatID`) untuk menyimpan `satusehat_id` ke dalam database PostgreSQL.

---

## 🛠️ Fitur Utama

- **OAuth 2.0 Client Credentials Token Manager**:
  - Otomatis melakukan token exchange ke server autentikasi SATUSEHAT (`/accesstoken?grant_type=client_credentials`).
  - Thread-safe caching token di memori.
  - Auto-refresh sebelum token kedaluwarsa dan fallback invalidasi jika menerima HTTP 401.
- **HL7 FHIR R4 Resource Builder**:
  - **Encounter**: Triase, rawat jalan, rawat inap, status kedatangan, identifier RS.
  - **Condition**: Kodifikasi ICD-10 untuk diagnosa primer/sekunder.
  - **Observation**: Vital signs, LOINC coding, hasil lab dengan satuan UCUM.
  - **Procedure**: Tindakan klinis berstandar ICD-9-CM & SNOMED CT.
  - **MedicationRequest**: Resep obat berbasis KFA (Katalog Farmasi & Alat Kesehatan).
  - **DiagnosticReport**: Hasil pemeriksaan penunjang lab & radiologi berstandar LOINC.
  - **AllergyIntolerance**: Data alergi pasien berkodifikasi SNOMED CT.
  - **Composition**: Resume medis / catatan perkembangan pasien terintegrasi.
- **RabbitMQ Consumer Resilience**:
  - Prefetch QoS concurrency worker.
  - Topic exchange & queue binding (`satusehat.sync.#`).
  - Dead Letter Queue (`satusehat.sync.dlq`) jika retries melebihi batas maksimum (`RABBITMQ_MAX_RETRIES`).
  - Graceful stop & in-flight delivery draining.
- **EMR Callback Interoperability**:
  - Menghubungkan gRPC client ke `emr-service:50052` untuk memanggil `SyncSatusehatID`.
- **HTTP Operational & Observability Endpoints**:
  - `GET /healthz` - Health check status service.
  - `GET /metrics` - Prometheus metrics (`satusehat_sync_total`, `satusehat_sync_duration_seconds`, `satusehat_queue_messages_total`, `satusehat_worker_up`).
  - `POST /api/v1/sync` - Endpoint manual / synchronous trigger untuk debugging, backfill, dan webhook.

---

## 🚀 Menjalankan Service

### 1. Menjalankan Langsung (Local Go)

Pastikan RabbitMQ dan `emr-service` berjalan, lalu:

```bash
cd satusehat-worker
cp .env.example .env
go run cmd/main.go
```

### 2. Menjalankan via Docker Compose

```bash
# Dari root project:
docker compose --profile apps up -d satusehat-worker
```

---

## 📬 Format Pesan Antrean RabbitMQ (`satusehat.sync`)

Pesan dikirimkan dalam format JSON ke antrean `satusehat.sync`:

```json
{
  "task_id": "0199f7d0-1234-7000-8000-000000000001",
  "resource_type": "encounter",
  "resource_id": "0199f7d0-9999-7000-8000-000000000002",
  "action": "CREATE",
  "payload": {
    "local_id": "REG-20261016-0001",
    "patient_ihs": "10000000001",
    "patient_name": "Budi Santoso",
    "practitioner_ihs": "N10000001",
    "practitioner_name": "dr. Andi Wijaya, Sp.PD",
    "location_id": "LOC-POLI-01",
    "location_name": "Poli Penyakit Dalam",
    "status": "finished",
    "class_code": "AMB",
    "start_time": "2026-10-16T08:00:00+07:00",
    "end_time": "2026-10-16T08:30:00+07:00"
  }
}
```

Atau jika pengirim sudah menyusun payload FHIR R4 lengkap:

```json
{
  "task_id": "0199f7d0-1234-7000-8000-000000000001",
  "resource_type": "condition",
  "resource_id": "0199f7d0-9999-7000-8000-000000000003",
  "payload": {
    "resourceType": "Condition",
    "clinicalStatus": {
      "coding": [{"system": "http://terminology.hl7.org/CodeSystem/condition-clinical", "code": "active"}]
    },
    "category": [
      {"coding": [{"system": "http://terminology.hl7.org/CodeSystem/condition-category", "code": "encounter-diagnosis"}]}
    ],
    "code": {
      "coding": [{"system": "http://hl7.org/fhir/sid/icd-10", "code": "A09", "display": "Infectious gastroenteritis"}]
    },
    "subject": {"reference": "Patient/10000000001"},
    "encounter": {"reference": "Encounter/ENC-FHIR-987"}
  }
}
```

---

## 🧪 Menjalankan Unit Tests

```bash
cd satusehat-worker
go test -v -count=1 ./...
```
