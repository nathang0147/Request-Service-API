# Request-Service-API

Request Service API for HR-facing verification orchestration across the VC ecosystem.

## Bootstrap Status

This repository is in active bootstrap. The current runnable shell provides:

- a Go module root
- an HTTP entrypoint at `cmd/api/main.go`
- a `GET /healthz` endpoint
- graceful shutdown support
- repo-local Go cache handling through `Makefile`

## Local Commands

```sh
make fmt
make vet
make test
make build
```

The bootstrap shell listens on `PORT`, defaulting to `8080` when unset.
