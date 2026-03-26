# AGENTS.md — Request Service API (Go) Bootstrap Guide

## 1. Project Identity

This repository implements a **Request Service API** for verifiable credential verification orchestration.

Its role is to act as the **front door** for an HR web application. The HR page must remain simple and must **not** need to know VC, verifier, protocol, or wallet details.

The Request Service API is responsible for:

- receiving a business-level verification initiation request from the HR page
- creating a verifier session through an external verifier ecosystem
- storing internal state and audit trail
- returning frontend-friendly data such as QR/deep link/session status
- receiving callback/final result from the verifier flow
- exposing a normalized verification status back to HR

For **v1**, only **one verifier ecosystem is supported: walt.id**.

This service is **not** the verifier itself and is **not** the policy engine.
The actual verification flow is:

- HR page -> Request Service API
- Request Service API -> walt.id verifier
- wallet/user interacts with verifier session
- verifier -> external trusted policy layer directly
- verifier result -> Request Service API
- HR page polls Request Service API for final normalized status

## 2. Strict Scope for v1

### In scope

- One Go service
- One PostgreSQL database
- One provider adapter: `walt`
- One northbound HR-facing REST API
- One callback endpoint for verifier result
- Persistent storage of requests, sessions, events
- Docker support
- Docker Compose for local dev
- GitHub Actions CI
- Cloud Run-ready containerization

### Out of scope

Do **not** build:

- microservices
- message queues
- Kafka / RabbitMQ
- Redis
- event sourcing
- CQRS
- runtime plugin systems
- multi-provider marketplace
- HR-configurable protocol/provider selection
- issuance flow
- wallet app
- policy engine
- complex admin UI
- generic framework abstractions with no immediate need

This is a **modular monolith**.

## 3. Architecture Decision

Use a **domain/feature-based package layout** with clear boundaries.

Primary organization should be by **business capability**, not by global technical layers.

Use these design patterns only where helpful:

- **Adapter Pattern** for provider integration
- **Strategy/Resolver Pattern** for provider selection
- **Repository Pattern** for persistence abstraction
- **Mapper Pattern** for DTO/domain/provider translation
- **Application Service Pattern** for orchestration logic

Do not create excessive abstraction or “clean architecture theater”.

## 4. Tech Stack

Use the following stack unless there is a strong reason not to:

- Go
- chi router
- PostgreSQL
- pgx
- sqlc
- golang-migrate
- zap logger
- Docker
- Docker Compose
- GitHub Actions
- Cloud Run compatible container

### Notes

- Prefer explicit SQL via `sqlc` over ORM
- Prefer simple `net/http` client wrappers for provider calls
- Use env-based configuration
- Keep startup simple and fast
- Make the app listen on `$PORT` for Cloud Run

## 5. Final Project Layout

Use this exact structure unless there is a strong implementation reason to adjust slightly:

```text
request-service/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── verification/
│   │   ├── domain.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── handler.go
│   │   ├── dto.go
│   │   ├── mapper.go
│   │   └── errors.go
│   ├── callback/
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── dto.go
│   │   └── mapper.go
│   ├── provider/
│   │   ├── provider.go
│   │   ├── resolver.go
│   │   └── walt/
│   │       ├── client.go
│   │       ├── provider.go
│   │       ├── dto.go
│   │       ├── mapper.go
│   │       └── auth.go
│   ├── persistence/
│   │   └── postgres/
│   │       ├── db.go
│   │       ├── verification_repository.go
│   │       ├── callback_repository.go
│   │       └── sqlc/
│   ├── platform/
│   │   ├── config/
│   │   ├── logger/
│   │   ├── middleware/
│   │   └── router/
│   └── shared/
│       ├── apierror/
│       ├── util/
│       └── clock/
├── db/
│   ├── migrations/
│   └── query/
├── api/
│   └── openapi.yaml
├── deployments/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── cloudrun.yaml
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
├── Makefile
├── go.mod
└── README.md
```

## 6. Domain Model

Define a normalized internal model. Provider-specific payloads must not leak into the whole codebase.

### `VerificationRequest`

Represents the business request from HR.

**Fields:**

- `ID`
- `BusinessRef`
- `CandidateRef`
- `Provider`
- `Status`
- `Verified`
- `ReasonCode`
- `CreatedAt`
- `UpdatedAt`

### `VerificationSession`

Represents the provider session associated with a request.

**Fields:**

- `ID`
- `VerificationRequestID`
- `Provider`
- `ProviderSessionID`
- `QRCodeURL`
- `DeepLink`
- `OfferURL`
- `ExpiresAt`
- `RawCreateResponse`
- `CreatedAt`

### `VerificationEvent`

Represents audit trail entries.

**Fields:**

- `ID`
- `VerificationRequestID`
- `Source`
- `EventType`
- `Payload`
- `CreatedAt`

## 7. Verification State Machine

Use this normalized lifecycle:

- `CREATED`
- `SESSION_CREATED`
- `PENDING`
- `VERIFIED`
- `FAILED`
- `EXPIRED`

Optional internal event markers are acceptable, but the public status model should stay simple.

## 8. Public API Design

### `POST /api/v1/verification-requests`

**Purpose:**
Create a new verification request and initialize verifier session.

**Request body:**

```json
{
  "businessRef": "job-123",
  "candidateRef": "cand-456"
}
```

**Response body:**

```json
{
  "requestId": "vr_001",
  "status": "PENDING",
  "verified": false,
  "session": {
    "qrCodeUrl": "https://example.com/qr",
    "deepLink": "openid://...",
    "expiresAt": "2026-03-24T12:00:00Z"
  }
}
```

**Notes:**

- HR sends only business-level data
- HR must not send protocol/provider configuration
- Provider selection is internal and defaults to `walt`

### `GET /api/v1/verification-requests/{id}`

**Purpose:**
Return normalized status for HR polling.

**Response body:**

```json
{
  "requestId": "vr_001",
  "status": "VERIFIED",
  "verified": true,
  "reasonCode": null
}
```

Keep the response simple.

## 9. Callback API

### `POST /api/v1/callbacks/walt`

**Purpose:**
Receive callback/final result from the verifier flow.

**Responsibilities:**

- validate callback auth/signature if available
- parse provider payload
- map callback to normalized internal event/result
- update verification request/session state
- write audit event
- return appropriate HTTP response

Do not expose raw provider semantics through this endpoint.

## 10. Provider Abstraction

Define a provider interface under `internal/provider/provider.go`.

**Expected behavior:**

```go
type VerifierProvider interface {
    CreateSession(ctx context.Context, input CreateSessionInput) (*CreateSessionOutput, error)
    GetSessionStatus(ctx context.Context, providerSessionID string) (*ProviderStatusOutput, error)
    ParseCallback(ctx context.Context, body []byte, headers http.Header) (*ProviderCallbackEvent, error)
}
```

### v1 implementation

Only implement:

- `walt`

### Resolver

Create a resolver under `internal/provider/resolver.go`.

For v1:

- always resolve to `walt`

But keep the boundary so future providers can be added later without redesign.

## 11. Persistence Strategy

Use PostgreSQL with explicit SQL and `sqlc`.

Do not use an ORM unless there is a very strong implementation blocker.

### Required tables

#### `verification_requests`

**Columns:**

- `id` UUID primary key
- `business_ref` TEXT
- `candidate_ref` TEXT
- `provider` TEXT
- `status` TEXT
- `verified` BOOLEAN
- `reason_code` TEXT NULL
- `created_at` TIMESTAMPTZ
- `updated_at` TIMESTAMPTZ

#### `verification_sessions`

**Columns:**

- `id` UUID primary key
- `verification_request_id` UUID foreign key
- `provider` TEXT
- `provider_session_id` TEXT UNIQUE
- `qr_code_url` TEXT NULL
- `deep_link` TEXT NULL
- `offer_url` TEXT NULL
- `expires_at` TIMESTAMPTZ NULL
- `raw_create_response` JSONB
- `created_at` TIMESTAMPTZ

#### `verification_events`

**Columns:**

- `id` UUID primary key
- `verification_request_id` UUID foreign key
- `source` TEXT
- `event_type` TEXT
- `payload` JSONB
- `created_at` TIMESTAMPTZ

### Persistence rules

- every created verification request must be stored
- every provider session must be stored
- every significant callback/update must create an audit event
- raw provider payloads may be stored in JSONB for audit/debugging

## 12. Mapping Rules

Keep translation explicit.

You will need at least these mapping directions:

- HTTP request DTO -> internal `VerificationRequest`
- internal request -> provider `walt` create session payload
- provider response -> internal `VerificationSession`
- callback payload -> normalized status update/event
- internal domain -> HTTP response DTO

Do not mix transport structs, provider structs, and domain structs casually.

## 13. Error Handling Rules

Use structured application errors.

### Principles

- HR-facing API must return normalized errors
- do not leak raw `walt` API internals to HR
- persistence errors and provider errors must be logged with enough context
- invalid input must return proper 4xx responses
- unexpected internal/provider/db failures return 5xx with normalized message/code

### Suggested internal reason codes

- `PROVIDER_SESSION_CREATE_FAILED`
- `CALLBACK_INVALID`
- `VERIFICATION_REJECTED`
- `VERIFICATION_EXPIRED`
- `REQUEST_NOT_FOUND`

## 14. Logging and Observability

Use structured logging.

### Requirements

- include request ID / correlation ID where possible
- log incoming request path/method/status
- log provider session creation attempts
- log callback reception and status transition
- log DB failures with context
- avoid logging secrets/tokens/raw credentials

Do not overbuild observability for v1. No tracing system is required.

## 15. Security Baseline

Implement a practical v1 baseline.

### HR-facing API

Use simple API key or bearer token authentication if needed by environment.

### Callback endpoint

Validate callback authentication/signature if supported.
If the real callback auth contract is not fully known yet, scaffold the verification boundary clearly with TODOs and clean interface points.

### General

- strict JSON binding/validation
- do not trust external callback payload blindly
- redact secrets in logs
- do not return internal stack traces to clients

## 16. Configuration

Use environment variables only.

### Example variables

- `APP_ENV`
- `PORT`
- `DATABASE_URL`
- `WALT_BASE_URL`
- `WALT_API_KEY`
- `CALLBACK_BASE_URL`
- `LOG_LEVEL`
- `DEFAULT_PROVIDER`
- `CALLBACK_AUTH_SECRET`

Keep config explicit and minimal.

## 17. Local Development

Provide local development via Docker Compose.

`docker-compose.yml` should include:

- app
- postgres

Keep it simple. Do not add Redis, queues, or extra services unless strictly needed.

The app must be runnable both:

- locally on host
- in Docker Compose

## 18. CI/CD

### CI (`.github/workflows/ci.yml`)

Must at least:

- run `go fmt` or formatting check
- run `go vet`
- run `go test ./...`
- build the Docker image

### CD (`.github/workflows/deploy.yml`)

Can later:

- build image
- push to Artifact Registry or target registry
- deploy to Cloud Run

Do not overbuild CI/CD in the first pass.

## 19. Cloud Run Constraints

The service must be Cloud Run friendly.

### Requirements

- listen on `$PORT`
- stateless app process
- fast startup
- graceful shutdown
- env-config driven
- container-first mindset

Do not hardcode local assumptions.

## 20. Code Quality Rules

### General style

- keep packages cohesive and small
- avoid unnecessary interfaces
- place interfaces near the consumer
- prefer explicitness over cleverness
- keep handler thin
- keep service logic in service layer
- keep provider-specific code inside provider adapter
- keep SQL explicit and reviewable

### Avoid

- premature abstraction
- global god-packages
- giant utility folders
- leaking `walt` DTOs into domain/service code
- building fake extensibility without use
- spreading DB queries across handlers

## 21. What to Build First

Implementation order should be:

1. bootstrap Go module
2. config + logger + router
3. database connection + migrations
4. SQL schema + sqlc setup
5. verification domain/service/repository
6. public HR endpoints
7. provider interface + walt adapter scaffold
8. create session flow end-to-end
9. callback flow scaffold
10. audit event persistence
11. Docker + Compose
12. CI workflow
13. Cloud Run readiness

Do not start with over-polishing infra before the core flow works.

## 22. Deliverable Expectations

The repository should be left in a state where:

- the project builds cleanly
- migrations exist
- DB integration exists
- HTTP endpoints exist
- provider abstraction exists
- walt provider scaffold/implementation exists
- callback flow exists at least in a sane scaffolded form
- Docker and Compose run locally
- CI runs basic checks
- code is understandable by a human reviewer

## 23. Important Behavioral Constraints for the Agent

When generating code for this repository:

- do not invent extra domains
- do not build microservices
- do not add event buses
- do not add Redis/queues unless explicitly requested later
- do not generalize beyond current need
- do not create generic plugin systems
- do not move HR provider/protocol choice into the API
- do not couple domain code directly to `walt`
- do not use an ORM by default
- do not produce architecture that is larger than the product

## 24. Current Business Assumptions

These assumptions are intentional and should be preserved unless explicitly changed later:

- HR is treated as a simple client and should not need VC knowledge
- Request Service creates verifier session
- only one verifier ecosystem is supported now: `walt.id`
- provider-specific differences are internal
- request service is the front door, not the main verifier
- trusted policy layer is external and verifier talks to it directly
- final HR output is just a simple verification flag plus basic status/reason
- persistence is required for audit trail and session tracking
- future extensibility matters, but v1 must stay simple

## 25. Final Instruction

Build the repository from scratch according to this document.
Prefer safe, explicit, maintainable code over excessive abstraction.
Whenever there is a choice, choose the option that keeps the codebase simpler while preserving the architecture described above.
