# Request-Service-API Repo Map

## Business Flow Summary

This repository is intended to host service-side request orchestration for HR-facing VC verification requests.

Target business flow:

1. an HR-facing client sends a business-level verification request to this service
2. the service validates and normalizes the request without exposing VC, verifier, or wallet details
3. the orchestration layer coordinates the local `walt` provider adapter implemented in this repo
4. callback payloads are normalized before any domain update
5. the service returns a normalized result, error, or follow-up callback state

Current reality:

- the repository is bootstrap-only
- no runtime request flow is implemented yet
- the repo currently contains workflow docs and the local `.codex` SDLC control plane only
- `.codex/tests/acceptance-checklist.md` holds architecture-derived bootstrap assertions

## Major Entrypoints

Current repo-native entrypoints:

- `README.md`: one-line repo purpose
- `AGENTS.md`: repo-local SDLC rules
- `.codex/repo-map.md`: fast mini-map
- `.codex/system-map.yaml`: machine-readable stage and lifecycle rules
- `.codex/ownership-map.yaml`: repo and sibling ownership map
- `.codex/tests/acceptance-checklist.md`: bootstrap acceptance assertions

Planned runtime entrypoints once code exists:

- `cmd/api/main.go`: service entrypoint
- `internal/verification/**`: business-level verification request lifecycle
- `internal/callback/**`: callback ingestion and normalization
- `internal/provider/**`: provider resolution and adapter boundary
- `internal/provider/walt/**`: local v1 `walt` provider adapter implementation in this repo
- `internal/persistence/postgres/**`: PostgreSQL repositories and storage
- `internal/platform/config/**`: service config and environment loading
- `internal/platform/logger/**`: structured logging
- `internal/platform/middleware/**`: HTTP middleware
- `internal/platform/router/**`: route wiring
- `internal/shared/apierror/**`: API error mapping
- `internal/shared/util/**`: shared helpers
- `internal/shared/clock/**`: clock abstraction
- `db/migrations/**`: schema migrations
- `db/query/**`: SQL definitions
- `api/openapi.yaml`: northbound contract
- `deployments/**`: container, compose, and Cloud Run assets
- `.github/workflows/**`: CI and deployment automation
- `go.mod`: Go module definition
- `Makefile`: developer task entrypoints if later added
- `tests/**`: contract, integration, and workflow verification

## Key Folders And Ownership

- `AGENTS.md`
  - repo-local workflow contract
- `.codex/`
  - repo-local SDLC engine
- `.codex/tests/`
  - tracked bootstrap acceptance assertions
- `docs/plans/`
  - design and implementation planning record
- `internal/verification/**`
  - future local ownership for request lifecycle and status handling
- `internal/callback/**`
  - future local ownership for callback normalization
- `internal/provider/**`
  - future local ownership for provider selection and adapter boundaries
- `internal/persistence/postgres/**`
  - future local ownership for persisted request, session, and event storage
- `internal/platform/**`
  - future local ownership for platform concerns like config, logging, middleware, and routing
- `internal/shared/**`
  - future local ownership for shared helpers and abstractions
- `db/**`
  - future local ownership for SQL migrations and query definitions
- `api/openapi.yaml`
  - future local ownership for the northbound HTTP contract
- `deployments/**`
  - future local ownership for container and deployment assets
- `.github/workflows/**`
  - future local ownership for CI and deployment automation
- `go.mod`
  - future local ownership for the Go module and dependency graph
- `Makefile`
  - future local ownership for developer task entrypoints if later introduced
- `tests/**`
  - future local ownership for verification

## Important Domain Objects

Current SDLC workflow objects:

- `Task`
  - durable work unit across all stages
- `Message`
  - communication event attached to a task
- `HandoffArtifact`
  - formal stage-transfer packet
- `SubagentExecution`
  - ephemeral stage invocation record, not persisted

Planned service-domain objects once implementation exists:

- `VerificationRequest`
  - inbound business-level verification request from the HR-facing client
- `VerificationSession`
  - provider session state, QR/deep link data, expiry, and raw create-session response
- `VerificationEvent`
  - append-only audit trail of request and callback transitions

Provider-specific DTOs stay inside provider adapter boundaries and must not become shared domain objects.

Normalized public lifecycle once implementation exists:

- `CREATED`
- `SESSION_CREATED`
- `PENDING`
- `VERIFIED`
- `FAILED`
- `EXPIRED`

## Verification Commands

Commands that are valid today:

- `git diff --check`
- `git status --short`
- `find .codex -maxdepth 3 -type f | sort`
- `find .codex/agents -maxdepth 1 -type f | sort`
- `! find .codex -maxdepth 3 -print | sort | rg -q '\.codex/tasks'`
- `rg -n 'stage_exit_artifacts:|blocked_stage_exits_require_stage_message: true|blocked stage exits must attach one|scope-report.md|task-envelope.md|design-review.md|implementation-report.md|verification-report.md|cross-repo-handoff.md|stage-message.md' AGENTS.md .codex/system-map.yaml .codex/agents`
- `rg -n 'resume_allowed|terminated_executions_are_never_resumed|subagent_executions_do_not_persist|free_form_reply_chains_are_not_workflow_state' AGENTS.md .codex/system-map.yaml`
- `sed -n '1,260p' AGENTS.md`
- `sed -n '1,260p' .codex/system-map.yaml`

Commands that are not valid yet:

- package-manager scripts
- build commands
- runtime tests

Use structural verification until real tooling exists.

## Sibling-Repo Dependencies

- `../IU-cert-university`
  - frontend and operator-facing consumer
  - owns browser-visible UX and frontend workflow behavior
- `../IU-VC-registry`
  - authority for registry contexts, credential schemas, DID documents, and public metadata
- `../waltid-identity`
  - upstream protocol reference only for wallet, verifier, OID4VCI, and related VC behavior
  - not the IU implementation repo for the local v1 `walt` adapter

## Bootstrap Notes

- repo implementation status: bootstrap-only
- codebase density: minimal
- fastest starting point for any new task:
  1. read `AGENTS.md`
  2. read this file
  3. read `.codex/system-map.yaml`
  4. inspect the actual file tree with `find` or `rg --files`
- do not pretend target folders already exist
- do not over-design the service before a real request path is introduced
- for v1, only the `walt` provider is supported
- keep the HR-facing API business-level and free of provider-specific details
- normalize callback paths before updating domain state
