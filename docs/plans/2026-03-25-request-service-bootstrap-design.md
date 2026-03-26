# Request Service Bootstrap Design

## Goal

Turn the architecture in `docs/plans/architecture_fixed.md` into a repo-local bootstrap design for `Request-Service-API` without overbuilding beyond the current v1 scope.

The runtime target is a single Go service that sits in front of HR, owns verification request/session/event persistence, integrates only with `walt.id` in v1, and returns normalized status to HR.

## Chosen Approach

Use a contract-first vertical-slice bootstrap.

Why this approach:

- it matches the architecture's modular-monolith constraint
- it avoids infra-first drift
- it keeps the first working path focused on the real business flow
- it lets repo-local SDLC docs stay aligned with the runtime structure before code appears

Rejected alternatives:

- layer-first foundation bootstrap: too much upfront scaffolding before core flow exists
- demo-first mock provider bootstrap: high rework risk and weak audit/persistence fidelity

## Architecture Import Strategy

The architecture should be imported into this repository in two layers.

1. Repo-control-plane alignment
   - update `AGENTS.md`
   - update `.codex/repo-map.md`
   - update `.codex/ownership-map.yaml`
   - keep `.codex/system-map.yaml` focused on SDLC lifecycle, not runtime package details

2. Runtime bootstrap
   - add the Go service structure
   - add database, provider, callback, and deployment surfaces
   - add only the minimum abstractions required for the v1 flow

## Runtime Target Shape

The runtime codebase should move toward this structure:

```text
cmd/api/main.go
internal/verification/*
internal/callback/*
internal/provider/provider.go
internal/provider/resolver.go
internal/provider/walt/*
internal/persistence/postgres/*
internal/platform/config/*
internal/platform/logger/*
internal/platform/middleware/*
internal/platform/router/*
internal/shared/apierror/*
internal/shared/util/*
internal/shared/clock/*
db/migrations/*
db/query/*
api/openapi.yaml
deployments/Dockerfile
deployments/docker-compose.yml
deployments/cloudrun.yaml
.github/workflows/ci.yml
.github/workflows/deploy.yml
Makefile
go.mod
README.md
```

## Repo Alignment Changes

Before code bootstrap, align the repo-local maps to the architecture:

- replace future `src/**` placeholders with the Go package layout above
- keep repo ownership centered on request orchestration, provider integration, callback handling, and API contracts
- preserve cross-repo boundaries:
  - `IU-cert-university` owns frontend behavior
  - `IU-VC-registry` owns schemas, DID docs, contexts, and registry metadata
  - `waltid-identity` remains a protocol reference, not the default implementation target

## Core Runtime Capabilities

The bootstrap must produce these capabilities in order:

1. service startup with env-based config and structured logging
2. PostgreSQL persistence for requests, sessions, and events
3. normalized verification domain and state machine
4. HR-facing create and status endpoints
5. provider interface and `walt` adapter
6. end-to-end create-session flow
7. callback ingestion and status update flow
8. audit trail persistence
9. local Docker execution
10. CI and Cloud Run readiness

## Domain Model

The normalized internal model should include:

- `VerificationRequest`
  - business request from HR
  - fields: `ID`, `BusinessRef`, `CandidateRef`, `Provider`, `Status`, `Verified`, `ReasonCode`, `CreatedAt`, `UpdatedAt`
- `VerificationSession`
  - provider-facing session metadata
  - fields: `ID`, `VerificationRequestID`, `Provider`, `ProviderSessionID`, `QRCodeURL`, `DeepLink`, `OfferURL`, `ExpiresAt`, `RawCreateResponse`, `CreatedAt`
- `VerificationEvent`
  - audit trail entry
  - fields: `ID`, `VerificationRequestID`, `Source`, `EventType`, `Payload`, `CreatedAt`

## State Model

Public normalized lifecycle:

- `CREATED`
- `SESSION_CREATED`
- `PENDING`
- `VERIFIED`
- `FAILED`
- `EXPIRED`

Internal code may use extra intermediate markers, but northbound API responses must stay simple.

## API Boundaries

Northbound API:

- `POST /api/v1/verification-requests`
- `GET /api/v1/verification-requests/{id}`

Provider callback:

- `POST /api/v1/callbacks/walt`

Rules:

- HR sends only business data
- provider choice is internal and defaults to `walt`
- provider DTOs do not leak into public API or core domain
- callback logic normalizes provider payloads before domain updates

## Persistence Shape

Required tables:

- `verification_requests`
- `verification_sessions`
- `verification_events`

Rules:

- every verification request is persisted
- every provider session is persisted
- every significant callback/update creates an audit event
- raw provider payloads may be stored for audit/debugging

## Verification Strategy

Early bootstrap verification should progress in layers:

1. structural verification while only docs exist
2. build and startup checks once `go.mod` and `cmd/api` exist
3. focused unit tests for domain, mappers, handlers, and provider boundaries
4. repository/integration tests once DB wiring exists
5. end-to-end request/session/callback happy-path tests once core flow exists
6. Docker build and CI verification after packaging is added

## Constraints To Preserve

Do not add during bootstrap:

- microservices
- queues or event buses
- Redis
- CQRS or event sourcing
- provider marketplaces
- HR-configurable provider/protocol choice
- generic plugin systems
- ORM-first persistence

Keep the codebase explicit, small, and reviewable.

## Open Questions

These inputs are still needed before implementation is fully closed:

- exact `walt.id` session-creation request/response contract
- exact callback auth/signature contract
- whether HR-facing auth is mandatory in v1 or environment-specific
- whether any registry-owned metadata must be provided from another repo

## Recommended Execution Order

1. align repo docs to the architecture
2. bootstrap Go module and app shell
3. add DB schema, migrations, and `sqlc`
4. implement verification domain/service/repository
5. add HR endpoints
6. add provider interface and `walt` scaffold
7. make create-session flow work end-to-end
8. add callback flow and audit events
9. add Docker, Compose, CI, and Cloud Run assets
10. finalize OpenAPI, tests, and repo docs
