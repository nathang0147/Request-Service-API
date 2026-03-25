GO ?= go
GO_ENV := GOCACHE=$(CURDIR)/.gocache GOMODCACHE=$(CURDIR)/.gomodcache
BINARY := .bin/request-service-api

.PHONY: fmt vet test build

fmt:
	$(GO_ENV) $(GO) fmt ./...

vet:
	$(GO_ENV) $(GO) vet ./...

test:
	$(GO_ENV) $(GO) test ./...

build:
	mkdir -p .bin
	$(GO_ENV) $(GO) build -o $(BINARY) ./cmd/api
