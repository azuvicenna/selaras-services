# patient-service

Service gRPC untuk data pasien di aplikasi rumah sakit Selaras. Menyimpan data pasien di PostgreSQL, dengan nomor rekam medis (RM) dan ID dibuat otomatis oleh database.

## Arsitektur

```
handler (gRPC)  ->  usecase (aturan bisnis)  ->  domain  <-  repository (PostgreSQL)
```

| Folder | Isi |
|---|---|
| `cmd/` | `main.go`, merangkai semua bagian lalu menjalankan server |
| `internal/domain/` | Entity `Patient`, error domain, interface `PatientRepository` |
| `internal/usecase/` | Validasi dan logika bisnis, plus unit test |
| `internal/repository/` | Implementasi `PatientRepository` memakai pgx |
| `internal/handler/` | Implementasi gRPC, pemetaan proto dan error domain ke status code |
| `internal/config/` | Membaca environment variable |
| `proto/patient/v1/` | `patient.proto` dan kode hasil generate (`*.pb.go`, jangan diedit manual) |
| `db/migrations/` | Migrasi skema (golang-migrate) |
| `db/seeders/` | Data contoh |
| `tests/` | Disiapkan untuk integration test |

## Prasyarat

- Go 1.26.6
- Docker Desktop
- Tool berikut, dipasang sekali saja:

```powershell
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Folder `C:\Users\hp\go\bin` harus ada di PATH.

## Konfigurasi

Salin contoh environment lalu sesuaikan. File `.env` tidak boleh di-commit.

```powershell
Copy-Item .env.example .env
```

| Variabel | Default | Keterangan |
|---|---|---|
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL di laptop |
| `DB_USER` | - | User database |
| `DB_PASSWORD` | - | Password database |
| `DB_NAME` | - | Nama database |
| `GRPC_PORT` | `50051` | Port server gRPC |

## Menjalankan lewat Docker (database dan aplikasi)

Jalankan dari dalam `patient-service`:

```powershell
docker compose up -d --build
docker compose ps          # tunggu db berstatus healthy
docker compose logs app    # harus ada "patient-service listening"
```

Migrasi tidak berjalan otomatis, jalankan sekali lewat laptop (lihat bagian Migrasi).

## Menjalankan lokal (go run)

Hanya database yang dijalankan di Docker, aplikasinya lewat `go run`:

```powershell
docker compose up -d db
Get-Content .env | ForEach-Object { if ($_ -match '^\s*([^#=]+?)\s*=\s*(.*)$') { Set-Item "env:$($Matches[1])" $Matches[2] } }
go run ./cmd
```

Go tidak membaca `.env` sendiri, jadi baris `Get-Content` itu memuat isinya ke sesi PowerShell. Jangan jalankan container `app` bersamaan, karena port 50051 akan bentrok.

## Migrasi dan seeder

Dari dalam `patient-service`, sesuaikan user, password, dan nama database dengan `.env`:

```powershell
# terapkan semua migrasi
migrate -path db/migrations -database "postgres://patient:change-me@localhost:5432/patient_db?sslmode=disable" up

# batalkan 1 migrasi terakhir
migrate -path db/migrations -database "postgres://patient:change-me@localhost:5432/patient_db?sslmode=disable" down 1

# buat migrasi baru (kerangka file kosong, isi SQL-nya manual)
migrate create -ext sql -dir db/migrations -seq nama_migrasi

# isi data contoh (aman dijalankan berulang)
Get-Content db/seeders/001_patients.sql | docker compose exec -T db psql -U patient -d patient_db
```

Cek isi tabel:

```powershell
docker compose exec db psql -U patient -d patient_db -c "SELECT medical_record_no, name FROM patients;"
```

## Mengubah kontrak API

Dari root `selaras-services`:

```powershell
buf lint
buf generate
cd patient-service
go mod tidy
cd ..
go build ./patient-service/...
```

Urutannya: ubah `patient.proto`, `buf generate`, `go mod tidy`, lalu `go build`.

## Test dan lint

Dari root `selaras-services`:

```powershell
go test ./patient-service/internal/usecase/ -v -cover
go build ./patient-service/...
```

Dari dalam `patient-service`:

```powershell
golangci-lint run ./...
```

Unit test memakai repository palsu, jadi tidak butuh database.

## API

Paket proto: `patient.v1`, service: `PatientService`.

| RPC | Fungsi |
|---|---|
| `RegisterPatient` | Daftarkan pasien baru, ID dan nomor RM dibuat database |
| `GetPatient` | Ambil pasien berdasarkan ID |
| `ListPatients` | Daftar pasien (`limit` default 20, maksimal 100, urut terbaru) |
| `UpdatePatient` | Ganti seluruh data pasien (field kosong ikut mengosongkan data) |

Contoh dengan grpcurl (server mendukung reflection, jadi tidak perlu file `.proto`):

```powershell
'{"limit": 10}' | grpcurl -plaintext -d '@' localhost:50051 patient.v1.PatientService/ListPatients

'{"nik": "3301010404950004", "name": "Dewi Lestari", "birthDate": "1995-04-04T00:00:00Z", "gender": "GENDER_FEMALE"}' | grpcurl -plaintext -d '@' localhost:50051 patient.v1.PatientService/RegisterPatient

'{"id": "<uuid-pasien>"}' | grpcurl -plaintext -d '@' localhost:50051 patient.v1.PatientService/GetPatient
```

### Aturan validasi

- Nama wajib diisi (spasi saja ditolak).
- NIK harus 16 digit angka dan unik.
- Tanggal lahir wajib diisi dan tidak boleh di masa depan.
- Gender harus `GENDER_MALE` atau `GENDER_FEMALE`.

### Pemetaan error

| Error domain | Status gRPC |
|---|---|
| `ErrInvalidInput` | `InvalidArgument` |
| `ErrNotFound` (termasuk ID bukan UUID) | `NotFound` |
| `ErrAlreadyExists` (NIK duplikat) | `AlreadyExists` |
| lainnya | `Internal` (detail hanya di log server) |

## Troubleshooting

- **`go mod edit -module=...` error di PowerShell**: bungkus argumen dengan tanda kutip, misalnya `"-module=github.com/..."`. PowerShell memecah argumen yang mengandung titik.
- **`missing module declaration` saat `go work use`**: `go.mod` masih kosong. Isi dengan `go mod edit`, bukan `go mod init`.
- **Container db langsung berhenti**: PostgreSQL 18 mengharuskan volume dipasang ke `/var/lib/postgresql`, bukan `/var/lib/postgresql/data`.
- **`volumes must be a mapping`**: format `docker-compose.yaml` salah. Cek dengan `docker compose config`. Blok `volumes:` di paling bawah hanya berisi nama volume tanpa tanda `-`.
- **`connection refused` saat migrate**: container belum `healthy` atau Docker Desktop belum jalan. Cek `docker compose ps`.
- **Port 5432 atau 50051 bentrok**: ganti `DB_PORT` atau `GRPC_PORT` di `.env`, lalu sesuaikan angka di URL migrate. Cek pemakai port dengan `netstat -ano | findstr :5432`.
- **Reset database dari awal**: `docker compose down -v` menghapus semua data, lalu jalankan migrasi dan seeder lagi.
- **Ubah `.env` tidak berpengaruh ke database lama**: user dan password hanya dipakai saat volume pertama dibuat. Gunakan `docker compose down -v`.