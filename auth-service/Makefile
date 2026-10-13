-include .env
export

# Default Configurations
GRPC_PORT   ?= 50057
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= auth
DB_PASSWORD ?= change-me
DB_NAME     ?= auth_db
DB_URL       = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Protobuf Directories
PROTO_DIR = proto/auth/v1
OUT_DIR   = gen/auth/v1

.DEFAULT_GOAL := help

.PHONY: help proto run test test-coverage build up down logs restart \
        migrate-up migrate-down migrate-force migrate-create migrate-version \
        seed lint fmt tidy clean

# -----------------------------------------------------------------------------
# HELP
# -----------------------------------------------------------------------------
help: ## Menampilkan daftar perintah yang tersedia
	@echo "Perintah yang tersedia di auth-service:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# -----------------------------------------------------------------------------
# PROTOBUF GENERATION
# -----------------------------------------------------------------------------
proto: ## Generate protobuf & gRPC Go stubs using buf
	cd .. && buf generate auth-service/proto
	@echo "✓ Generated Protobuf & gRPC files successfully in $(OUT_DIR)"

# -----------------------------------------------------------------------------
# APPLICATION COMMANDS
# -----------------------------------------------------------------------------
run: ## Menjalankan service secara lokal
	go run ./cmd

build: ## Build binary executable service
	go build -o bin/auth-service ./cmd

tidy: ## Merapikan dan memverifikasi go.mod dependencies
	go mod tidy
	go mod verify

seed: ## Menjalankan seeder data awal user
	go run ./db/seeders/main.go -addr localhost:$(GRPC_PORT) -file db/seeders/users.json

# -----------------------------------------------------------------------------
# CODE QUALITY & TESTING
# -----------------------------------------------------------------------------
test: ## Menjalankan unit & integration tests
	go test ./... -cover

test-coverage: ## Menjalankan test dan membuat laporan coverage HTML
	@mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "✓ Coverage report generated at coverage/coverage.html"

lint: ## Menjalankan golangci-lint
	golangci-lint run ./...

fmt: ## Format code menggunakan go fmt
	go fmt ./...

# -----------------------------------------------------------------------------
# DOCKER & CONTAINER COMMANDS
# -----------------------------------------------------------------------------
up: ## Menjalankan seluruh environment (database + app) via docker-compose
	docker compose up -d

down: ## Menghentikan seluruh container docker-compose
	docker compose down

logs: ## Melihat log output dari app container
	docker compose logs -f app

restart: down up ## Restart docker compose environment

# -----------------------------------------------------------------------------
# DATABASE MIGRATIONS
# -----------------------------------------------------------------------------
migrate-up: ## Menjalankan migrasi database ke versi terbaru
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down: ## Rollback 1 langkah migrasi database
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-force: ## Memaksa state migrasi ke versi tertentu (misal: make migrate-force V=1)
	migrate -path db/migrations -database "$(DB_URL)" force $(V)

migrate-version: ## Melihat versi migrasi saat ini
	migrate -path db/migrations -database "$(DB_URL)" version

migrate-create: ## Membuat file migrasi baru (misal: make migrate-create NAME=add_fields)
	migrate create -ext sql -dir db/migrations -seq $(NAME)

clean: ## Membersihkan binary dan temporary files
	rm -rf bin/ coverage/
