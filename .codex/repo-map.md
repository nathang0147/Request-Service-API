# Request-Service-API Repo Map

## Business Flow Summary

This repository is intended to host service-side request orchestration for VC-related frontend integrations.

Target business flow:

1. a frontend or operator-facing client sends a request to this service
2. the service validates and normalizes the request
3. the orchestration layer coordinates VC ecosystem dependencies
4. integrations call sibling or upstream systems when needed
5. the service returns a normalized result, error, or follow-up callback state

Current reality:

- the repository is bootstrap-only
- no runtime request flow is implemented yet
- the repo currently contains workflow docs and the local `.codex` SDLC control plane only

## Major Entrypoints

Current repo-native entrypoints:

- `README.md`: one-line repo purpose
- `AGENTS.md`: repo-local SDLC rules
- `.codex/repo-map.md`: fast mini-map
- `.codex/system-map.yaml`: machine-readable stage and lifecycle rules
- `.codex/ownership-map.yaml`: repo and sibling ownership map

Planned runtime entrypoints once code exists:

- `src/http/**`: request entrypoints and transport layer
- `src/orchestration/**`: request lifecycle coordinator
- `src/integrations/**`: outbound VC ecosystem clients and adapters
- `src/domain/**`: domain logic and invariants
- `src/config/**`: service config and environment loading

## Key Folders And Ownership

- `AGENTS.md`
  - repo-local workflow contract
- `.codex/`
  - repo-local SDLC engine
- `docs/plans/`
  - design and implementation planning record
- `src/http/**`
  - future local ownership for transport and validation
- `src/orchestration/**`
  - future local ownership for request lifecycle logic
- `src/integrations/**`
  - future local ownership for external adapters
- `src/domain/**`
  - future local ownership for pure domain rules
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
  - inbound service request from frontend or operator workflow
- `VerificationResult`
  - normalized success or failure result returned to the caller
- `CallbackEvent`
  - async callback or follow-up event emitted by an external dependency
- `IntegrationDependency`
  - sibling or upstream service required to complete the request

## Verification Commands

Commands that are valid today:

- `git diff --check`
- `git status --short`
- `find .codex -maxdepth 3 -type f | sort`
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
  - upstream protocol reference for wallet, verifier, OID4VCI, and related VC behavior

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
