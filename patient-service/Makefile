-include .env
export

# Default Database Configurations
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= postgres
DB_PASSWORD ?= postgres
DB_NAME     ?= patient_db
DB_URL       = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Protobuf Directories
PROTO_DIR = proto/patient/v1
OUT_DIR   = gen/patient/v1

.PHONY: proto run test build up down migrate-up migrate-down clean

# -----------------------------------------------------------------------------
# PROTOBUF GENERATION
# -----------------------------------------------------------------------------
proto:
	@mkdir -p $(OUT_DIR)
	protoc \
		--proto_path=proto \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "✓ Generated Protobuf & gRPC files successfully in $(OUT_DIR)"

# -----------------------------------------------------------------------------
# APPLICATION COMMANDS
# -----------------------------------------------------------------------------
run:
	go run ./cmd

test:
	go test ./... -cover

build: proto
	go build -o bin/patient-service ./cmd

# -----------------------------------------------------------------------------
# DOCKER COMMANDS
# -----------------------------------------------------------------------------
up:
	docker compose up -d --build

down:
	docker compose down

# -----------------------------------------------------------------------------
# DATABASE MIGRATIONS
# -----------------------------------------------------------------------------
migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
clean:
	rm -rf bin/ gen/
	@echo "✓ Cleaned build binaries and generated proto files"