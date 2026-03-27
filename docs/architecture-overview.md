# Architecture Overview

## Purpose

`Request-Service-API` is the backend entrypoint for HR-driven verification flows across the VC ecosystem.

Its job is to keep the HR application simple. The frontend should not need to know verifier-specific protocols, wallet interaction details, or provider-specific payload shapes.

This service is responsible for:

- receiving a business-level verification request from HR
- creating and tracking a verifier session
- persisting request, session, and audit state
- normalizing verifier outcomes into a simple status model
- exposing status back to the HR-facing application

For v1, the only supported verifier ecosystem is `walt.id`.

## System Context

The intended end-to-end flow is:

1. HR frontend sends a verification request to this service.
2. This service creates a provider session through `walt.id`.
3. The wallet holder completes the verifier flow outside the HR app.
4. The verifier sends a callback result to this service.
5. This service updates internal state and audit history.
6. The HR frontend polls this service for normalized status.

This repository is not the verifier itself, not the wallet, and not the policy engine.

## Repository Boundaries

This repository owns:

- request orchestration
- provider session lifecycle
- callback handling
- API contracts exposed to HR
- persistence of requests, sessions, and audit events

Related repositories stay responsible for their own domains:

- `IU-cert-university`: frontend and operator-facing behavior
- `IU-VC-registry`: registry metadata, schemas, contexts, DID documents
- `waltid-identity`: upstream protocol reference and ecosystem dependency

If a task depends on another repo's authority, the work should hand off instead of leaking that responsibility into this codebase.

## v1 Scope

In scope:

- one Go service
- one PostgreSQL database
- one provider adapter: `walt`
- one HR-facing REST API
- one verifier callback endpoint
- Docker and Docker Compose
- GitHub Actions CI
- Cloud Run-ready deployment shape

Out of scope:

- microservices
- message queues
- Redis
- CQRS or event sourcing
- provider marketplaces
- HR-configurable provider selection
- issuance workflows
- wallet implementation
- policy engine logic

The intended shape is a modular monolith, not a distributed system.

## Core Components

The runtime design groups code by capability:

- `cmd/api`
  - application entrypoint and dependency wiring
- `internal/verification`
  - verification request lifecycle, business orchestration, northbound API handlers
- `internal/callback`
  - verifier callback ingestion and normalized status update flow
- `internal/provider`
  - provider interface and provider resolution
- `internal/provider/walt`
  - `walt.id` client, auth, DTOs, and mapping
- `internal/persistence/postgres`
  - PostgreSQL access and repository implementations
- `internal/platform`
  - config, logging, middleware, router, and startup infrastructure
- `internal/shared`
  - small shared utilities such as API error handling and clock helpers
- `db`
  - migrations and `sqlc` query definitions
- `api/openapi.yaml`
  - public API contract
- `deployments`
  - local and cloud deployment assets

## Domain Model

The architecture normalizes provider-specific behavior into three core records.

### VerificationRequest

The business-level request initiated by HR.

Representative fields:

- `ID`
- `BusinessRef`
- `CandidateRef`
- `Provider`
- `Status`
- `Verified`
- `ReasonCode`
- `CreatedAt`
- `UpdatedAt`

### VerificationSession

The provider session linked to a verification request.

Representative fields:

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

### VerificationEvent

The audit trail for significant state changes and provider interactions.

Representative fields:

- `ID`
- `VerificationRequestID`
- `Source`
- `EventType`
- `Payload`
- `CreatedAt`

## Verification Lifecycle

Public normalized status should stay simple:

- `CREATED`
- `SESSION_CREATED`
- `PENDING`
- `VERIFIED`
- `FAILED`
- `EXPIRED`

Provider-specific states may exist internally, but they should be mapped into this normalized lifecycle before they reach the northbound API.

## API Surface

The intended northbound API is:

- `POST /api/v1/verification-requests`
  - create a verification request and initialize provider session
- `GET /api/v1/verification-requests/{id}`
  - read normalized verification status and session details
- `POST /api/v1/callbacks/walt`
  - receive verifier callback updates from `walt.id`

Rules:

- HR sends only business-level request data
- provider selection stays internal
- provider DTOs must not leak into domain or public API responses

## Persistence Model

The minimum persistent model uses three tables:

- `verification_requests`
- `verification_sessions`
- `verification_events`

Expected persistence behavior:

- every request is stored
- every provider session is stored
- every meaningful callback or status transition creates an audit event
- raw provider payloads may be kept when useful for audit and debugging

## Deployment Model

The intended production shape is:

- Cloud Run for the stateless Go service
- Cloud SQL for PostgreSQL
- Secret Manager for sensitive configuration

Local development should use Docker Compose for the service and PostgreSQL.

PostgreSQL should not be deployed inside Cloud Run for production.

## Current Status

This repository is still in bootstrap. The architecture is defined, and the implementation work should follow the package boundaries and lifecycle described above.

For deeper detail, read:

- `docs/plans/architecture_fixed.md`
- `docs/plans/2026-03-25-request-service-bootstrap-design.md`
- `docs/plans/2026-03-25-request-service-bootstrap-implementation-plan.md`
