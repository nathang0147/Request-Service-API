# Request-Service-API

Request Service API for HR-facing verification orchestration across the VC ecosystem.

## Docs

- [Architecture overview](docs/architecture-overview.md)
- [Bootstrap architecture source](docs/plans/architecture_fixed.md)
- [Bootstrap design](docs/plans/2026-03-25-request-service-bootstrap-design.md)
- [Bootstrap implementation plan](docs/plans/2026-03-25-request-service-bootstrap-implementation-plan.md)

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

When started through Docker Compose, the API is exposed on `http://localhost:8081` and PostgreSQL is published on `localhost:5433`. Internal container ports remain `8080` for the API and `5432` for PostgreSQL.

For local verifier integration, the request-service container defaults to `WALT_VERIFIER_MODE=legacy` with `WALT_VERIFIER_BASE_URL=http://host.docker.internal:7003`, which matches the older `waltid-identity` verifier surface that still emits `presentation_definition` flows for the current wallet UI. `WALT_BEARER_TOKEN` is optional. `WALT_VC_POLICY_WEBHOOK_URL` points the legacy verifier at the IU trust-policy webhook, defaulting locally to `http://host.docker.internal:8787/api/verifier/policies/vc`. If you later migrate the wallet/UI to DCQL, switch to `WALT_VERIFIER_MODE=verifier2` and point `WALT_VERIFIER_BASE_URL` at the verifier2 surface on port `7004`.

Use separate machine and browser redirect settings when the verifier runs in Docker and the browser runs on your host:

- `CALLBACK_BASE_URL` -> machine callback base for `/api/v1/callbacks/walt`
- `PUBLIC_REDIRECT_URL_TEMPLATE` -> browser-visible success/error redirect target, with `$id` replaced by the verifier session id
- `PUBLIC_BASE_URL` -> fallback browser base for `/api/v1/verification-requests/{requestId}` if no redirect template is configured

Local Docker defaults are:

- `CALLBACK_BASE_URL=http://host.docker.internal:8081`
- `PUBLIC_REDIRECT_URL_TEMPLATE=http://localhost:7102/success/$$id` in `docker-compose.yml` so Compose passes a literal `$id`
- `PUBLIC_BASE_URL=http://localhost:8081`

This split matters because `localhost` means different things in different execution contexts:

- inside the verifier container, `localhost` is the verifier container itself, so callbacks to `http://localhost:8081` fail
- in the host browser wallet UI, `http://localhost:7102/success/$id` gives the normal `waltid-identity` success page

The `waltid-identity` verifier portal itself uses this same pattern, sending `successRedirectUri: ${window.location.origin}/success/$id` in [pages/verify/index.tsx](/Users/thanganguyen/Documents/10-19%20Repo/17%20Thesis/waltid-identity/waltid-applications/waltid-web-portal/pages/verify/index.tsx#L89).

When you set `PUBLIC_REDIRECT_URL_TEMPLATE` through Docker Compose, escape the dollar sign as `$$id`. Otherwise Compose treats `$id` as an environment variable and strips it.

`response_uri` in the OpenID4VP request is not the success redirect target. It is only the verifier endpoint that receives the wallet submission.

## Local Database Flow

Schema changes live in `db/migrations/` and are applied with `golang-migrate`.

Typical local flow:

```sh
make compose-up
```

To create a new migration template:

```sh
make migrate-create name=add_reason_code
```

`make compose-up` runs the `postgres`, `migrate`, and `request-service-api` services from `deployments/docker-compose.yml`. The API waits for both PostgreSQL health and successful migration completion before starting, so you do not need to run `make migrate-up` first for normal local startup.
