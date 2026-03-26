# Request Service Bootstrap Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Bootstrap `Request-Service-API` into a runnable Go service that persists verification requests, creates `walt` verifier sessions, handles callbacks, exposes normalized HR-facing APIs, and is ready for Docker, CI, and Cloud Run.

**Architecture:** Use a contract-first vertical slice inside a modular monolith. Align the repo-control-plane docs first, then build the app shell, persistence, domain flow, provider boundary, callback path, and delivery surfaces in that order.

**Tech Stack:** Go, chi, PostgreSQL, pgx, sqlc, golang-migrate, zap, Docker, Docker Compose, GitHub Actions, Cloud Run

---

### Task 1: Align Repo Control-Plane Docs To The Go Architecture

**Files:**
- Modify: `AGENTS.md`
- Modify: `.codex/repo-map.md`
- Modify: `.codex/ownership-map.yaml`
- Modify: `.codex/tests/acceptance-checklist.md`
- Reference: `docs/plans/architecture_fixed.md`
- Reference: `docs/plans/2026-03-25-request-service-bootstrap-design.md`

**Step 1: Update intended runtime boundaries**

Replace generic future `src/**` placeholders with:

- `cmd/api/**`
- `internal/verification/**`
- `internal/callback/**`
- `internal/provider/**`
- `internal/provider/walt/**`
- `internal/persistence/postgres/**`
- `internal/platform/**`
- `internal/shared/**`
- `db/migrations/**`
- `db/query/**`
- `api/openapi.yaml`
- `deployments/**`
- `.github/workflows/**`

**Step 2: Add architecture-derived verification assertions**

Add checklist assertions that prove:

- only one provider is supported in v1
- HR API stays business-level
- callback path is normalized before domain update
- provider DTOs do not leak into domain boundaries

**Step 3: Run structural verification**

Run:

```bash
git diff --check
find .codex -maxdepth 3 -type f | sort
sed -n '1,260p' AGENTS.md
```

Expected:

- no whitespace errors
- repo docs reflect the Go package layout
- no reintroduction of legacy task orchestration

### Task 2: Bootstrap The Go Module And Entry Point

**Files:**
- Create: `go.mod`
- Create: `cmd/api/main.go`
- Create: `Makefile`
- Modify: `README.md`

**Step 1: Create module and base commands**

Add:

- `go.mod`
- basic `Makefile` targets for `fmt`, `vet`, `test`, `build`
- README bootstrap instructions

**Step 2: Add minimal main entry**

`cmd/api/main.go` should:

- load config
- initialize logger
- build router
- create `http.Server`
- listen on `$PORT`
- support graceful shutdown

**Step 3: Run bootstrap checks**

Run:

```bash
go test ./...
go build ./cmd/api
```

Expected:

- `go test ./...` passes or reports no test files
- `go build ./cmd/api` succeeds

### Task 3: Add Platform Shell

**Files:**
- Create: `internal/platform/config/config.go`
- Create: `internal/platform/logger/logger.go`
- Create: `internal/platform/router/router.go`
- Create: `internal/platform/middleware/request_id.go`
- Create: `internal/platform/middleware/logging.go`
- Create: `internal/shared/apierror/apierror.go`
- Create: `internal/shared/clock/clock.go`

**Step 1: Write failing smoke tests where useful**

Add focused tests for:

- config defaulting/required env validation
- router construction
- request ID middleware behavior

**Step 2: Run tests to verify failure**

Run:

```bash
go test ./internal/platform/... ./internal/shared/...
```

Expected:

- failures for missing implementations

**Step 3: Implement minimal platform code**

Implement:

- env-based config with `PORT`, `DATABASE_URL`, `WALT_BASE_URL`, `WALT_API_KEY`, `CALLBACK_BASE_URL`, `LOG_LEVEL`, `DEFAULT_PROVIDER`, `CALLBACK_AUTH_SECRET`
- zap logger builder
- chi router with health route placeholder
- request ID and request logging middleware
- normalized API error helpers

**Step 4: Re-run tests**

Run:

```bash
go test ./internal/platform/... ./internal/shared/...
```

Expected:

- platform/shared tests pass

### Task 4: Add Database Schema, Migrations, And sqlc

**Files:**
- Create: `db/migrations/0001_init.up.sql`
- Create: `db/migrations/0001_init.down.sql`
- Create: `db/query/verification_requests.sql`
- Create: `db/query/verification_sessions.sql`
- Create: `db/query/verification_events.sql`
- Create: `sqlc.yaml`

**Step 1: Encode the schema from the architecture**

Create tables:

- `verification_requests`
- `verification_sessions`
- `verification_events`

Include:

- UUID primary keys
- FK from sessions/events to requests
- unique provider session ID
- JSONB storage for raw provider payloads

**Step 2: Add SQL queries for repository operations**

At minimum:

- create request
- get request by ID
- update request status
- create session
- get session by provider session ID
- create event
- list events for a request if useful

**Step 3: Generate sqlc code**

Run:

```bash
sqlc generate
```

Expected:

- generated code under the configured output directory

### Task 5: Add Postgres Integration

**Files:**
- Create: `internal/persistence/postgres/db.go`
- Create: `internal/persistence/postgres/verification_repository.go`
- Create: `internal/persistence/postgres/callback_repository.go`
- Create: `internal/persistence/postgres/sqlc/` generated files
- Test: `internal/persistence/postgres/*_test.go`

**Step 1: Write failing repository tests**

Cover:

- create/get verification request
- create session and look up by provider session ID
- append audit event
- update normalized status and reason code

**Step 2: Run repository tests to verify failure**

Run:

```bash
go test ./internal/persistence/postgres/... -v
```

Expected:

- repository tests fail before implementation

**Step 3: Implement DB wiring and repository methods**

Use:

- `pgx`
- generated `sqlc` queries
- explicit mapping between DB rows and domain models

**Step 4: Re-run repository tests**

Run:

```bash
go test ./internal/persistence/postgres/... -v
```

Expected:

- repository tests pass

### Task 6: Implement Verification Domain And Service

**Files:**
- Create: `internal/verification/domain.go`
- Create: `internal/verification/repository.go`
- Create: `internal/verification/service.go`
- Create: `internal/verification/dto.go`
- Create: `internal/verification/mapper.go`
- Create: `internal/verification/errors.go`
- Test: `internal/verification/service_test.go`

**Step 1: Write failing service tests**

Cover:

- create request initializes normalized request state
- create request persists request and session
- service returns normalized session details to HR
- request lookup returns normalized status/verified/reason
- provider failure maps to `PROVIDER_SESSION_CREATE_FAILED`

**Step 2: Run service tests to verify failure**

Run:

```bash
go test ./internal/verification/... -v
```

Expected:

- service tests fail before service implementation exists

**Step 3: Implement the domain and orchestration**

Add:

- normalized domain types
- repository interfaces near the consumer
- orchestration service for create request and get status
- mapper boundary between DTOs, domain, and provider input/output

**Step 4: Re-run service tests**

Run:

```bash
go test ./internal/verification/... -v
```

Expected:

- verification tests pass

### Task 7: Add HR-Facing HTTP Endpoints

**Files:**
- Create: `internal/verification/handler.go`
- Test: `internal/verification/handler_test.go`
- Modify: `internal/platform/router/router.go`

**Step 1: Write failing handler tests**

Cover:

- `POST /api/v1/verification-requests`
- `GET /api/v1/verification-requests/{id}`
- 4xx on invalid JSON/input
- 404 for unknown request
- normalized 5xx on internal failure

**Step 2: Run handler tests to verify failure**

Run:

```bash
go test ./internal/verification/... -run Handler -v
```

Expected:

- handler tests fail before routes/handlers exist

**Step 3: Implement minimal handlers**

Rules:

- handlers remain thin
- validation is strict
- response payload stays HR-friendly
- do not leak provider payloads or internal stack traces

**Step 4: Re-run handler tests**

Run:

```bash
go test ./internal/verification/... -run Handler -v
```

Expected:

- handler tests pass

### Task 8: Add Provider Boundary And Walt Adapter

**Files:**
- Create: `internal/provider/provider.go`
- Create: `internal/provider/resolver.go`
- Create: `internal/provider/walt/client.go`
- Create: `internal/provider/walt/provider.go`
- Create: `internal/provider/walt/dto.go`
- Create: `internal/provider/walt/mapper.go`
- Create: `internal/provider/walt/auth.go`
- Test: `internal/provider/walt/provider_test.go`

**Step 1: Write failing provider tests**

Cover:

- resolver always returns `walt` in v1
- create session maps internal request to walt payload
- provider response maps to normalized session output
- callback parse normalizes provider event shape

**Step 2: Run provider tests to verify failure**

Run:

```bash
go test ./internal/provider/... -v
```

Expected:

- provider tests fail before adapter implementation

**Step 3: Implement the interface and adapter**

Rules:

- use `net/http` wrappers, not a heavy client framework
- keep DTOs and mappers inside `internal/provider/walt`
- do not leak `walt` structs into `internal/verification`

**Step 4: Re-run provider tests**

Run:

```bash
go test ./internal/provider/... -v
```

Expected:

- provider tests pass

### Task 9: Wire The Create-Session Flow End-To-End

**Files:**
- Modify: `internal/verification/service.go`
- Modify: `internal/provider/resolver.go`
- Test: `internal/verification/service_e2e_test.go`

**Step 1: Write failing flow test**

Cover:

- create request
- resolve `walt`
- create provider session
- persist request/session
- return normalized `PENDING` response with QR/deep link/session expiry

**Step 2: Run flow test to verify failure**

Run:

```bash
go test ./internal/verification/... -run EndToEnd -v
```

Expected:

- end-to-end flow test fails before final wiring

**Step 3: Implement minimal end-to-end wiring**

Re-use the existing service, repository, and provider boundaries. Avoid new abstractions unless a real duplication forces them.

**Step 4: Re-run flow test**

Run:

```bash
go test ./internal/verification/... -run EndToEnd -v
```

Expected:

- end-to-end flow test passes

### Task 10: Add Callback Flow And Audit Events

**Files:**
- Create: `internal/callback/service.go`
- Create: `internal/callback/handler.go`
- Create: `internal/callback/dto.go`
- Create: `internal/callback/mapper.go`
- Test: `internal/callback/service_test.go`
- Test: `internal/callback/handler_test.go`
- Modify: `internal/platform/router/router.go`

**Step 1: Write failing callback tests**

Cover:

- callback auth boundary
- payload parse through provider adapter
- request/session state update
- audit event persisted
- normalized HTTP response to callback sender

**Step 2: Run callback tests to verify failure**

Run:

```bash
go test ./internal/callback/... -v
```

Expected:

- callback tests fail before callback implementation

**Step 3: Implement callback flow**

Rules:

- validate auth/signature if contract is known
- keep auth verification behind a clear boundary if final contract is still TBD
- normalize callback semantics before touching domain state

**Step 4: Re-run callback tests**

Run:

```bash
go test ./internal/callback/... -v
```

Expected:

- callback tests pass

### Task 11: Add Contracts, Packaging, And Delivery Surfaces

**Files:**
- Create: `api/openapi.yaml`
- Create: `deployments/Dockerfile`
- Create: `deployments/docker-compose.yml`
- Create: `deployments/cloudrun.yaml`
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/deploy.yml`
- Modify: `README.md`

**Step 1: Add OpenAPI contract**

Describe:

- create request endpoint
- get request status endpoint
- callback endpoint
- normalized status and error shapes

**Step 2: Add local packaging**

Create:

- Dockerfile
- Docker Compose with app + postgres

**Step 3: Add CI/CD workflows**

CI must run:

```bash
go fmt ./...
go vet ./...
go test ./...
docker build -f deployments/Dockerfile .
```

CD can remain a scaffold, but it must match the Cloud Run deployment target.

**Step 4: Verify packaging**

Run:

```bash
go test ./...
docker build -f deployments/Dockerfile .
```

Expected:

- tests pass
- Docker image builds successfully

### Task 12: Final Verification And Documentation Sweep

**Files:**
- Modify: `README.md`
- Modify: `.codex/repo-map.md`
- Modify: `.codex/ownership-map.yaml`
- Modify: `.codex/tests/acceptance-checklist.md`

**Step 1: Reconcile docs with the landed runtime structure**

Update repo docs so they describe the actual Go layout and active verification commands.

**Step 2: Run full verification set**

Run:

```bash
git diff --check
go fmt ./...
go vet ./...
go test ./...
docker build -f deployments/Dockerfile .
```

Expected:

- whitespace clean
- formatting clean
- vet clean
- tests pass
- Docker build succeeds

**Step 3: Commit in small logical batches**

Recommended commit sequence:

```bash
git add AGENTS.md .codex/repo-map.md .codex/ownership-map.yaml .codex/tests/acceptance-checklist.md
git commit -m "docs: align repo control plane with bootstrap architecture"

git add go.mod cmd/api/main.go internal/platform internal/shared Makefile README.md
git commit -m "feat: bootstrap go service shell"

git add db sqlc.yaml internal/persistence
git commit -m "feat: add postgres persistence foundation"

git add internal/verification internal/provider internal/callback api deployments .github
git commit -m "feat: add verification flow and delivery surfaces"
```
