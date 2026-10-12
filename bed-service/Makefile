-include .env
export

# Default Configurations
GRPC_PORT   ?= 50054
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= bed
DB_PASSWORD ?= change-me
DB_NAME     ?= bed_db
DB_URL       = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Protobuf Directories
PROTO_DIR = proto/bed/v1
OUT_DIR   = gen/bed/v1

.DEFAULT_GOAL := help

.PHONY: help proto run test test-coverage build up down logs restart \
        migrate-up migrate-down migrate-force migrate-create migrate-version \
        seed lint fmt tidy clean

# -----------------------------------------------------------------------------
# HELP
# -----------------------------------------------------------------------------
help: ## Menampilkan daftar perintah yang tersedia
	@echo "Perintah yang tersedia di bed-service:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# -----------------------------------------------------------------------------
# PROTOBUF GENERATION
# -----------------------------------------------------------------------------
proto: ## Generate protobuf & gRPC Go stubs using buf
	cd .. && buf generate bed-service/proto
	@echo "✓ Generated Protobuf & gRPC files successfully in $(OUT_DIR)"

# -----------------------------------------------------------------------------
# APPLICATION COMMANDS
# -----------------------------------------------------------------------------
run: ## Menjalankan service secara lokal
	go run ./cmd

build: ## Build binary executable service
	go build -o bin/bed-service ./cmd

tidy: ## Merapikan dan memverifikasi go.mod dependencies
	go mod tidy
	go mod verify

seed: ## Menjalankan seeder data awal tempat tidur
	go run ./db/seeders/main.go -addr localhost:$(GRPC_PORT) -file db/seeders/bed_data.json

# -----------------------------------------------------------------------------
# CODE QUALITY & TESTING
# -----------------------------------------------------------------------------
test: ## Menjalankan unit & integration tests
	go test ./... -cover

test-coverage: ## Menjalankan test dan membuat laporan coverage HTML
	@mkdir -p coverage
	go test "-coverprofile=coverage/coverage.out" ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "✓ Coverage report generated at coverage/coverage.html"

lint: ## Menjalankan golangci-lint
	golangci-lint run ./...

fmt: ## Format kode menggunakan gofmt
	gofmt -s -w .

# -----------------------------------------------------------------------------
# DOCKER COMMANDS
# -----------------------------------------------------------------------------
up: ## Menyalakan seluruh service Docker (Postgres, App)
	docker compose up -d --build

down: ## Mematikan container Docker
	docker compose down

restart: ## Restart container aplikasi
	docker compose restart app

logs: ## Melihat streaming logs dari container aplikasi
	docker compose logs -f app

# -----------------------------------------------------------------------------
# DATABASE MIGRATIONS
# -----------------------------------------------------------------------------
migrate-up: ## Menjalankan semua migrasi database yang belum dieksekusi
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down: ## Rollback 1 step migrasi database
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-version: ## Melihat status versi migrasi database saat ini
	migrate -path db/migrations -database "$(DB_URL)" version

migrate-force: ## Memperbaiki status dirty migrasi (penggunaan: make migrate-force v=1)
	@if [ -z "$(v)" ]; then echo "Error: tentukan versi, contoh: make migrate-force v=1"; exit 1; fi
	migrate -path db/migrations -database "$(DB_URL)" force $(v)

migrate-create: ## Membuat file migrasi baru (penggunaan: make migrate-create name=add_column)
	@if [ -z "$(name)" ]; then echo "Error: tentukan nama, contoh: make migrate-create name=add_column"; exit 1; fi
	migrate create -ext sql -dir db/migrations -seq $(name)

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
clean: ## Membersihkan hasil build dan coverage
	rm -rf bin/ coverage/
	@echo "✓ Cleaned build binaries and coverage reports"
