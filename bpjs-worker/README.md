# BPJS Worker Microservice

Worker asinkron untuk menjembatani sistem rumah sakit (SIMRS) dengan layanan terintegrasi **BPJS Kesehatan**, mencakup **VClaim v2.0**, **Antrean RS (Mobile JKN)**, dan **Aplicares (Ketersediaan Tempat Tidur)**.

---

## 🛠️ Fitur Utama

- **BPJS TrustMark & Web Service Cryptography**:
  - Penghitungan signature HMAC-SHA256 sesuai standar TrustMark Kemenkes/BPJS:
    $$\text{Signature} = \text{Base64}\Big(\text{HMAC-SHA256}\big(\text{consID} + \text{"&"} + \text{timestamp},\, \text{secretKey}\big)\Big)$$
  - Auto-derivation AES Key (SHA-256 dari `consID + secretKey + timestamp`) dan 16-byte IV.
  - Dekripsi otomatis payload AES-256-CBC response BPJS VClaim 2.0.
  - Enkripsi AES-256-CBC payload untuk pengiriman data terlindungi.
- **VClaim v2.0 Services**:
  - **Peserta**: Pengecekan data kepesertaan dan eligibilitas pasien berdasarkan NIK atau Nomor Kartu BPJS.
  - **Rujukan**: Pencarian data rujukan faskes 1 / faskes 2 (No. Kunjungan, diagnosa, poli rujukan, PPK asal).
  - **SEP (Surat Eligibilitas Peserta)**:
    - Pembuatan SEP Rawat Jalan & Rawat Inap (`InsertSEP`).
    - Hapus SEP (`DeleteSEP`).
    - Update Tanggal Pulang SEP (`UpdateTglPulang`).
- **Aplicares (Ketersediaan Tempat Tidur)**:
  - Pembaruan kapasitas dan ketersediaan tempat tidur RS secara real-time (`AplicaresUpdateBed`).
- **Antrean RS (Mobile JKN Task ID)**:
  - Pembaruan waktu tunggu pasien (Task ID 1 s.d. 7) untuk integrasi antrean online Mobile JKN.
- **RabbitMQ Consumer Resilience**:
  - Antrean `bpjs.sync` dengan topic exchange `bpjs.exchange` dan dead-letter queue `bpjs.sync.dlq`.
  - Dukungan retries otomatis dengan backoff sebelum pengalihan ke DLQ.
- **HTTP Operational & Observability**:
  - `GET /healthz` - Health check status service.
  - `GET /metrics` - Prometheus metrics (`bpjs_request_total`, `bpjs_request_duration_seconds`, `bpjs_queue_messages_total`, `bpjs_worker_up`).
  - `POST /api/v1/vclaim/peserta` - Pengecekan peserta langsung via REST.
  - `POST /api/v1/vclaim/sep` - Pembuatan SEP langsung.
  - `POST /api/v1/aplicares/bed` - Update ketersediaan bed ke Aplicares.
  - `POST /api/v1/sync` - Eksekusi tugas antrean secara manual/sinkron untuk debugging dan testing.

---

## 🚀 Menjalankan Service

### 1. Menjalankan Langsung (Local Go)

```bash
cd bpjs-worker
cp .env.example .env
go run cmd/main.go
```

### 2. Menjalankan via Docker Compose

```bash
# Dari root project:
docker compose --profile apps up -d bpjs-worker
```

---

## 📬 Format Pesan Antrean RabbitMQ (`bpjs.sync`)

### 1. Pengecekan Peserta (`CHECK_PARTICIPANT`)
```json
{
  "task_id": "0199f7d0-0001-7000-8000-000000000001",
  "action": "CHECK_PARTICIPANT",
  "payload": {
    "nik": "3171010000000001",
    "tgl_sep": "2026-10-16"
  }
}
```

### 2. Pembuatan SEP (`CREATE_SEP`)
```json
{
  "task_id": "0199f7d0-0002-7000-8000-000000000002",
  "action": "CREATE_SEP",
  "resource_id": "local-encounter-uuid",
  "payload": {
    "request": {
      "t_sep": {
        "noKartu": "0001234567890",
        "tglSep": "2026-10-16",
        "ppkPelayanan": "0115R001",
        "jnsPelayanan": "2",
        "klsRawat": "1",
        "noMR": "MR-00102",
        "rujukan": {
          "asalRujukan": "1",
          "tglRujukan": "2026-10-15",
          "noRujukan": "0115B0011026P000001",
          "ppkRujukan": "0115B001"
        },
        "catatan": "Kunjungan poli umum",
        "diagAwal": "A09",
        "poli": {
          "tujuan": "INT",
          "eksekutif": "0"
        },
        "cob": { "cob": "0" },
        "katarak": { "katarak": "0" },
        "jaminan": { "lakaLantas": "0" },
        "tujuanKunjungan": "0",
        "noTelp": "081234567890",
        "user": "petugas_admisi"
      }
    }
  }
}
```

### 3. Update Ketersediaan Bed Aplicares (`APLICARES_SYNC_BED`)
```json
{
  "task_id": "0199f7d0-0003-7000-8000-000000000003",
  "action": "APLICARES_SYNC_BED",
  "payload": {
    "kodekelas": "VIP",
    "koderuang": "R001",
    "namaruang": "Ruang Anggrek VIP 1",
    "kapasitas": 10,
    "tersedia": 3,
    "tersediapria": 1,
    "tersediawanita": 2,
    "tersediapriawanita": 0
  }
}
```

---

## 🧪 Menjalankan Unit Tests

```bash
cd bpjs-worker
go test -v -count=1 ./...
```
