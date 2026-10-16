SERVICES := $(patsubst %/go.mod,%,$(wildcard */go.mod))
PKGS := $(addsuffix /...,$(addprefix ./,$(SERVICES)))

.PHONY: proto lint test build

proto:
	buf generate

lint:
	buf lint
	@for s in $(SERVICES); do (cd $$s && golangci-lint run ./...) || exit 1; done

test:
	go test $(PKGS)

build:
	go build $(PKGS)

compose-infra:
	docker compose up -d postgres redis minio minio-init rabbitmq

compose-apps:
	docker compose --profile apps up -d --build

compose-down:
	docker compose down

migrate-all:
	go run ./scripts/migrator -action=up

migrate-down:
	go run ./scripts/migrator -action=down

migrate-status:
	go run ./scripts/migrator -action=status

migrate-reset:
	go run ./scripts/migrator -action=reset

seed-all:
	go run ./scripts/migrator -action=seed

compose-migrate:
	docker compose --profile tools run --rm migration-runner