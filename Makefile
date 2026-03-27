GO ?= go
GO_ENV := GOCACHE="$(CURDIR)/.gocache" GOMODCACHE="$(CURDIR)/.gomodcache"
BINARY := .bin/request-service-api
COMPOSE := docker compose -f deployments/docker-compose.yml
MIGRATE_DATABASE_URL ?= postgres://postgres:postgres@postgres:5432/request_service?sslmode=disable
MIGRATE_PATH := /migrations

.PHONY: fmt vet test build compose-up compose-down migrate-create migrate-up migrate-down

fmt:
	$(GO_ENV) $(GO) fmt ./...

vet:
	$(GO_ENV) $(GO) vet ./...

test:
	$(GO_ENV) $(GO) test ./...

build:
	mkdir -p .bin
	$(GO_ENV) $(GO) build -o $(BINARY) ./cmd/api

compose-up:
	$(COMPOSE) up --build

compose-down:
	$(COMPOSE) down

migrate-create:
	test -n "$(name)"
	$(COMPOSE) run --rm migrate create -ext sql -dir=$(MIGRATE_PATH) -seq $(name)

migrate-up:
	$(COMPOSE) up -d postgres
	$(COMPOSE) run --rm migrate -path=$(MIGRATE_PATH) -database='$(MIGRATE_DATABASE_URL)' up

migrate-down:
	$(COMPOSE) up -d postgres
	$(COMPOSE) run --rm migrate -path=$(MIGRATE_PATH) -database='$(MIGRATE_DATABASE_URL)' down 1
