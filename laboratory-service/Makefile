-include .env
export

# Default Configurations
GRPC_PORT   ?= 50060
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= laboratory
DB_PASSWORD ?= change-me
DB_NAME     ?= laboratory_db
DB_URL       = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Protobuf Directories
PROTO_DIR = proto/laboratory/v1
OUT_DIR   = gen/laboratory/v1

.DEFAULT_GOAL := help

.PHONY: help proto run test test-coverage build up down logs restart \
        migrate-up migrate-down migrate-force migrate-create migrate-version \
        seed lint fmt tidy clean

# -----------------------------------------------------------------------------
# HELP
# -----------------------------------------------------------------------------
help: ## Menampilkan daftar perintah yang tersedia
	@echo "Perintah yang tersedia di laboratory-service:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# -----------------------------------------------------------------------------
# PROTOBUF GENERATION
# -----------------------------------------------------------------------------
proto: ## Generate protobuf & gRPC Go stubs using buf
	cd .. && buf generate laboratory-service/proto
	@echo "✓ Generated Protobuf & gRPC files successfully in $(OUT_DIR)"

# -----------------------------------------------------------------------------
# APPLICATION COMMANDS
# -----------------------------------------------------------------------------
run: ## Menjalankan service secara lokal
	go run ./cmd

build: ## Build binary executable service
	go build -o bin/laboratory-service ./cmd

tidy: ## Merapikan dan memverifikasi go.mod dependencies
	go mod tidy
	go mod verify

seed: ## Menjalankan seeder data awal master uji lab
	go run ./db/seeders/main.go -addr localhost:$(GRPC_PORT) -file db/seeders/tests.json

# -----------------------------------------------------------------------------
# CODE QUALITY & TESTING
# -----------------------------------------------------------------------------
test: ## Menjalankan unit & integration tests
	go test ./... -cover

test-coverage: ## Menjalankan test dan membuat laporan coverage HTML
	@mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "✓ Laporan coverage dibuat di coverage/coverage.html"

lint: ## Menjalankan buf lint dan go vet
	cd .. && buf lint laboratory-service/proto
	go vet ./...

fmt: ## Memformat kode Go
	go fmt ./...

# -----------------------------------------------------------------------------
# DOCKER & INFRASTRUCTURE
# -----------------------------------------------------------------------------
up: ## Menjalankan seluruh stack via docker-compose
	docker compose up -d --build

down: ## Menghentikan seluruh stack docker-compose
	docker compose down

logs: ## Melihat log running service di docker
	docker compose logs -f app

restart: down up ## Restart docker stack

# -----------------------------------------------------------------------------
# DATABASE MIGRATIONS (golang-migrate)
# -----------------------------------------------------------------------------
migrate-up: ## Menjalankan semua migrasi database yang belum dieksekusi
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down: ## Rollback 1 langkah migrasi terakhir
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

migrate-force: ## Memaksa versi migrasi jika terjadi dirty state (e.g. make migrate-force v=1)
	migrate -path db/migrations -database "$(DB_URL)" force $(v)

migrate-version: ## Melihat versi migrasi database saat ini
	migrate -path db/migrations -database "$(DB_URL)" version

migrate-create: ## Membuat file migrasi baru (e.g. make migrate-create name=add_lab_indexes)
	migrate create -ext sql -dir db/migrations -seq $(name)

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
clean: ## Menghapus artifact build dan coverage
	rm -rf bin/ coverage/
