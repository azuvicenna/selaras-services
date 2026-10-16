# Notification Worker Microservice (RabbitMQ)

Worker asinkron untuk pengiriman notifikasi rumah sakit berskala tinggi menggunakan **RabbitMQ**. Mendukung multi-channel: **WhatsApp** dan **Email (HTML templating)**.

---

## 🛠️ Fitur Utama

- **RabbitMQ Event-Driven & Task Queue Architecture**:
  - Exchange: `notification.exchange` (topic)
  - Queue utama: `notification.send`
  - Routing Key: `notification.#`
  - Dead Letter Queue (DLQ): `notification.send.dlq` (native via `x-dead-letter-exchange`)
  - Concurrency worker pool dan automatic acknowledgement (ACK/NACK).
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
  - `GET /metrics` - Prometheus metrics (`notification_sent_total`, `notification_duration_seconds`, `notification_queue_messages_total`, `notification_worker_up`).
  - `POST /api/v1/notify/whatsapp` - Pengiriman WhatsApp langsung.
  - `POST /api/v1/notify/email` - Pengiriman Email langsung.
  - `POST /api/v1/notify` - Publish event ke antrean RabbitMQ (`notification.send`), atau kirim langsung dengan parameter `?direct=true`.

---

## 🚀 Menjalankan Service

### 1. Menjalankan Langsung (Local Go)

Pastikan RabbitMQ berjalan (misal via Docker Compose), lalu:

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

## 📬 Format Pesan Antrean RabbitMQ (`notification.send`)

### Contoh Payload JSON

```json
{
  "id": "notif-001",
  "channel": "WHATSAPP",
  "recipient": "+628123456789",
  "type": "APPOINTMENT_REMINDER",
  "subject": "Pengingat Janji Temu Poliklinik",
  "data": {
    "patient_name": "Budi Santoso",
    "doctor_name": "dr. Siti Rahma, Sp.A",
    "clinic_name": "Poli Anak",
    "appointment_date": "18 Oktober 2026",
    "appointment_time": "09:00 WIB",
    "queue_number": "A-012"
  }
}
```
