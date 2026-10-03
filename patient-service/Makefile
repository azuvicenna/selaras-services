-include .env
export

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_URL   = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: run test build up down migrate-up migrate-down

run:
	go run ./cmd

test:
	go test ./... -cover

build:
	go build -o bin/patient-service ./cmd

up:
	docker compose up -d --build

down:
	docker compose down

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1