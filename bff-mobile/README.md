# Selaras SIMRS - Mobile Backend-For-Frontend (BFF) Gateway

Layanan API Gateway Backend-For-Frontend (BFF) khusus untuk aplikasi mobile pasien dan tenaga medis (**Flutter & React Native**) di lingkungan RS Selaras Medika. Layanan ini mengoptimalkan komunikasi jaringan nirkabel seluler dengan menyediakan agregasi komposit (composite endpoints), verifikasi versi aplikasi seluler, dan konsolidasi data dari 12 downstream gRPC microservices.

---

## 📱 Keunggulan & Arsitektur Mobile-First

1. **Composite Dashboard (`/api/v1/home`)**:
   - Menghindari *waterfall requests* dan konsumsi baterai/data berlebih pada perangkat seluler.
   - Menggabungkan data profil pasien, jadwal janji temu terdekat, status antrean aktif, dan ketersediaan tempat tidur dalam satu panggilan REST tunggal secara asinkron (`sync.WaitGroup`).
2. **Mobile Device & App Version Enforcement (`DeviceMiddleware`)**:
   - Membaca header `X-Device-Id`, `X-Platform` (Android/iOS), dan `X-App-Version`.
   - Melakukan evaluasi semantik terhadap versi minimal aplikasi (`MIN_APP_VERSION`). Jika versi aplikasi di bawah batas, gateway mengembalikan response `HTTP 426 Upgrade Required` beserta `update_url` Play Store / App Store.
3. **Stateless Security**:
   - Integrasi token JWT pada header `Authorization: Bearer <token>`.
   - Menginjeksi identitas pengguna (`user_id`, `role`, `patient_id`) ke dalam context request.
4. **Resilient Downstream Integration**:
   - Koneksi gRPC terpusat ke 12 microservices inti (Auth, Patient, Practitioner, Admission, Bed, Master, EMR, Billing, Pharmacy, Laboratory, Radiology, Document).
   - Penanganan timeout konteks dan pemetaan status error gRPC ke JSON HTTP response standar.

---

## 🛠️ Daftar Endpoint API

### 1. Sistem & Observabilitas
- `GET /healthz` - Health check status gateway.
- `GET /health` - Alias health check.
- `GET /metrics` - Prometheus metrics exporter.
- `GET /docs` - Swagger UI dokumentasi API interaktif.
- `GET /openapi.json` - OpenAPI 3.0 specification file.

### 2. Autentikasi Mobile
- `POST /api/v1/auth/login` - Login pengguna aplikasi mobile (Pasien/Nakes).
- `POST /api/v1/auth/refresh` - Refresh access token yang telah kedaluwarsa.
- `GET /api/v1/auth/me` *(Protected)* - Mengambil profil akun pengguna aktif.

### 3. Komposit Beranda
- `GET /api/v1/home` *(Protected)* - Dashboard komposit mobile pasien.

### 4. Pendaftaran & Reservasi Janji Temu (Booking)
- `GET /api/v1/booking/polyclinics` *(Protected)* - Daftar poliklinik / departemen yang tersedia.
- `GET /api/v1/booking/doctors` *(Protected)* - Daftar dokter dan spesialisasi.
- `POST /api/v1/booking/reservations` *(Protected)* - Membuat reservasi janji temu poliklinik.
- `GET /api/v1/booking/my-reservations` *(Protected)* - Riwayat janji temu pasien.

### 5. Antrean Rumah Sakit (Live Queue)
- `GET /api/v1/queue/live?poli_id=...` *(Protected)* - Status antrean langsung, nomor yang dipanggil, sisa antrean di depan, dan estimasi waktu tunggu.

### 6. Rekam Medis Pasien (EMR)
- `GET /api/v1/records/encounters` *(Protected)* - Riwayat kunjungan dan rawat jalan/inap.
- `GET /api/v1/records/encounters/detail?encounter_id=...` *(Protected)* - Detail diagnosis, catatan klinis, dan tanda vital.

### 7. Hasil Penunjang Medis (Laboratorium & Radiologi)
- `GET /api/v1/lab/results` *(Protected)* - Riwayat dan hasil pemeriksaan laboratorium pasien.
- `GET /api/v1/radiology/results` *(Protected)* - Riwayat dan hasil pemeriksaan radiologi pasien.

### 8. Tagihan & Pembayaran (Billing)
- `GET /api/v1/billing/invoices` *(Protected)* - Daftar faktur tagihan pasien.
- `GET /api/v1/billing/invoices/detail?invoice_id=...` *(Protected)* - Rincian item tindakan, obat, dan pembayaran.

### 9. Ketersediaan Tempat Tidur Publik
- `GET /api/v1/beds/availability` *(Public)* - Ringkasan ketersediaan kamar dan ranjang per kelas perawatan (VVIP, VIP, Kelas 1, 2, 3, ICU).

---

## 🚀 Panduan Menjalankan

### Konfigurasi Environment (`.env`)

```env
HTTP_PORT=8081
METRICS_PORT=9091
ENV=development
CORS_ALLOWED_ORIGINS=*
JWT_SECRET=selaras-hospital-jwt-secret-key-32bytes!
MIN_APP_VERSION=1.0.0
APP_UPDATE_URL=https://play.google.com/store/apps/details?id=com.selaras.mobile

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

### Menjalankan secara Lokal

```bash
# 1. Jalankan unit test
go test -v -count=1 ./bff-mobile/...

# 2. Build dan jalankan server
go run ./bff-mobile/cmd/main.go
```

### Menjalankan via Docker Compose

```bash
# Dari root repositori
docker compose --profile apps up -d bff-mobile
```

Akses gateway di browser atau Postman:
- Dokumentasi API Swagger UI: `http://localhost:8081/docs`
- Health check: `http://localhost:8081/healthz`
