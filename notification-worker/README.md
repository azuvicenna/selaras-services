# Notification Worker Microservice (Apache Kafka)

Worker asinkron untuk pengiriman notifikasi rumah sakit berskala tinggi (*high-throughput event streaming*) menggunakan **Apache Kafka**. Mendukung multi-channel: **WhatsApp** dan **Email (HTML templating)**.

---

## 🛠️ Fitur Utama

- **Apache Kafka Event-Driven Architecture**:
  - Topik utama: `hospital.notifications`
  - Consumer Group: `notification-worker-group`
  - Dead Letter Queue (DLQ): `hospital.notifications.dlq`
  - Automatic offset management, concurrency worker pool, dan retry backoff policy.
- **Multi-Channel Notification Providers**:
  - **WhatsApp**: Mendukung Mock Provider (untuk development/testing) dan HTTP Gateway Provider (kompatibel dengan Meta WhatsApp Business Cloud API, Fonnte, Wablas, dan generic webhooks).
  - **Email**: Mendukung Mock Provider dan SMTP Provider dengan output HTML responsive dan styled card.
- **Dynamic Template Engine**:
  - `APPOINTMENT_REMINDER`: Pengingat jadwal konsultasi poliklinik dan dokter.
  - `REGISTRATION_WELCOME`: Notifikasi registrasi pasien baru dan Nomor Rekam Medis (No. RM).
  - `QUEUE_TICKET_UPDATE`: Pemanggilan antrean poliklinik, admisi, dan farmasi.
  - `LAB_RESULT_READY`: Pemberitahuan hasil pemeriksaan laboratorium selesai beserta tautan unduh.
  - `RADIOLOGY_RESULT_READY`: Pemberitahuan hasil pemeriksaan radiologi selesai.
  - `BILLING_INVOICE`: Konfirmasi pembayaran tagihan / invoice rumah sakit.
  - `OTP_CODE`: Kode verifikasi 2FA aman dengan kedaluwarsa 5 menit.
  - `CUSTOM`: Pesan teks bebas kustom.
- **Observability & REST APIs**:
  - `GET /healthz` - Health check status service.
  - `GET /metrics` - Prometheus metrics (`notification_sent_total`, `notification_duration_seconds`, `kafka_messages_processed_total`, `notification_worker_up`).
  - `POST /api/v1/notify/whatsapp` - Pengiriman WhatsApp langsung.
  - `POST /api/v1/notify/email` - Pengiriman Email langsung.
  - `POST /api/v1/notify` - Publish event ke Kafka topic (`hospital.notifications`), atau kirim langsung dengan parameter `?direct=true`.

---

## 🚀 Menjalankan Service

### 1. Menjalankan Langsung (Local Go)

Pastikan Kafka berjalan (misal via Docker Compose), lalu:

```bash
cd notification-worker
cp .env.example .env
go run cmd/main.go
```

### 2. Menjalankan via Docker Compose

```bash
# Dari root project:
docker compose --profile apps up -d notification-worker
```

---

## 📬 Format Event Message Kafka (`hospital.notifications`)

Pesan dikirimkan dalam format JSON ke Kafka topic `hospital.notifications`:

### 1. Pengingat Janji Temu (WhatsApp)
```json
{
  "id": "notif-001",
  "channel": "WHATSAPP",
  "recipient": "+6281234567890",
  "type": "APPOINTMENT_REMINDER",
  "data": {
    "patient_name": "Budi Santoso",
    "practitioner_name": "dr. Andi Wijaya, Sp.PD",
    "poli_name": "Poli Penyakit Dalam",
    "schedule_time": "Senin, 19 Oktober 2026, 09:00 WIB",
    "queue_number": "A-012"
  }
}
```

### 2. Bukti Pembayaran Tagihan (Email HTML)
```json
{
  "id": "notif-002",
  "channel": "EMAIL",
  "recipient": "pasien@example.com",
  "type": "BILLING_INVOICE",
  "data": {
    "patient_name": "Citra Dewi",
    "invoice_number": "INV-20261016-0045",
    "amount": "750.000"
  }
}
```

---

## 🧪 Menjalankan Unit Tests

```bash
cd notification-worker
go test -v -count=1 ./...
```
