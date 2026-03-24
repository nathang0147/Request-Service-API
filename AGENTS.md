# AGENTS.md

## Repository Role

This repository is the local home for `Request-Service-API` engineering work.
It should own backend request orchestration, service-side integration logic, and API contracts that sit between IU frontend flows and the VC ecosystem.

Treat this repository as the execution boundary once work has been routed here by a parent controller.

## Current State

- This repository is currently bootstrap-only and contains only a minimal `README.md` and `LICENSE`.
- Do not invent build scripts, test runners, package managers, or source paths that are not actually present.
- If runtime or project tooling is missing, say so explicitly and either:
  - limit the task to docs/config/bootstrap work, or
  - identify the minimum missing scaffold that must be added first.

## Cross-Repo Map

- `IU-cert-university`: likely downstream frontend or operator-facing consumer for flows exposed by this service. UI and browser workflow changes belong there unless the task is explicitly about this API.
- `IU-VC-registry`: source of truth for registry contexts, credential schemas, DID documents, and public registry metadata. Do not invent or fork registry-owned changes here.
- `waltid-identity`: upstream protocol reference for wallet, verifier, OID4VCI, and related VC ecosystem behavior. Use it as a reference, not as an IU-owned implementation target.

If a task depends on authority in another repository, stop at the handoff boundary and report the dependency instead of guessing.

## Expected Code Boundaries

When code is added, keep the repository layered and predictable:

- `src/http/**`: HTTP routes, controllers, transport validation, response mapping
- `src/orchestration/**`: request lifecycle coordination, state transitions, retries, callback handling
- `src/integrations/**`: external VC ecosystem clients, adapters, and protocol-facing boundaries
- `src/domain/**`: pure domain logic, policies, models, and invariants
- `src/config/**`: environment loading and config wiring
- `tests/**`: contract, integration, and workflow verification

These paths are target boundaries, not a claim that they already exist.

## Hard Boundaries

- Keep transport-layer concerns out of pure domain logic.
- Keep external client code isolated from orchestration and domain rules.
- Do not couple this repository directly to sibling-repo internals.
- Do not make cross-repo edits from this checkout unless the user explicitly asks for work in the other repository.
- If this repo is still missing source scaffolding, prefer the smallest viable structure over speculative architecture.

## Commands

Safe baseline commands in the current repository state:

- inspect files: `rg`, `find`, `sed`
- inspect history: `git log --oneline -5`
- inspect working tree: `git status --short`

When package or runtime tooling appears:

- prefer repository scripts from `package.json`, `Makefile`, or equivalent
- prefer targeted checks over broad sweeps
- do not guess commands that are not defined locally

## Verification Expectations

- For docs or config-only changes, use structural verification such as:
  - `git diff --check`
  - `find .codex -maxdepth 3 -type f | sort`
  - `sed -n '1,220p' AGENTS.md`
- For code changes, run the smallest relevant real project command set available in the repository.
- If the repository lacks runnable tooling, do not claim runtime verification.
- Be explicit about what was verified, what was not verified, and why.

## Multiagent Use

- A thesis-level controller in the parent `17 Thesis` folder may decide whether work should enter this repository.
- Once work is routed here, this repository's `AGENTS.md` and `.codex/agents/` are the source of truth for local execution.
- The parent agent should own planning, integration, and final verification.
- Spawn read-only explorers first when the request is unclear or crosses boundaries.
- Use implementation workers only after file ownership is clear.
- Give workers disjoint ownership sets. Good future splits in this repo are:
  - `src/http/**`
  - `src/orchestration/**`
  - `src/integrations/**`
  - `src/domain/**`
  - `docs/**`
- Do not run parallel workers on the same subtree or the same feature slice.

## Project Subagents

This repository includes project-scoped subagents under `.codex/agents/` for:

- repository and path mapping
- request-flow and orchestration analysis
- API contract review
- cross-repo dependency scouting
- small, scoped implementation
- verification planning and evidence review

Use them as helpers, not as a substitute for parent-agent judgment.
