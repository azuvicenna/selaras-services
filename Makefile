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