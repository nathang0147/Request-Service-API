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
make migrate-create name=add_reason_code
make migrate-up
make migrate-down
make compose-up
make compose-down
```

The bootstrap shell listens on `PORT`, defaulting to `8080` when unset.

## Local Database Flow

Schema changes live in `db/migrations/` and are applied with `golang-migrate`.

Typical local flow:

```sh
make migrate-up
make compose-up
```

To create a new migration template:

```sh
make migrate-create name=add_reason_code
```

`make compose-up` runs the `postgres`, `migrate`, and `request-service-api` services from `deployments/docker-compose.yml`. The API waits for both PostgreSQL health and successful migration completion before starting.
