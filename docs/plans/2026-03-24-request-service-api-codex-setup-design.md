# Request-Service-API Codex Setup Design

## Goal

Create a repo-local agent setup for `Request-Service-API` that is safe for engineering work inside this repository and easy for a parent controller in `17 Thesis` to route into without overriding repo-local rules.

## Context

This repository is currently at bootstrap state with only `README.md` and `LICENSE` committed. Sibling repositories already use a consistent local-agent pattern:

- root `AGENTS.md`
- `.codex/config.toml`
- `.codex/agents/*.toml`

The user also wants this repo to participate in a broader thesis-level orchestration model.

## Decision

Use a hybrid federated setup:

- keep repository-specific instructions and subagents inside this repository
- make those instructions controller-friendly so a parent orchestration layer can hand work into this repo cleanly
- do not attempt to centralize repo-specific behavior in the parent folder

## Why Not Central-Only

A controller one directory above the repo can decide which repository should own a task, but it should not be the source of truth for:

- repository boundaries
- local verification commands
- file ownership conventions
- cross-repo handoff rules once execution has entered the repo

Those concerns drift over time and belong in the repo that owns the code.

## Planned Local Structure

```text
Request-Service-API/
  AGENTS.md
  docs/
    plans/
      2026-03-24-request-service-api-codex-setup-design.md
      2026-03-24-request-service-api-codex-setup.md
  .codex/
    config.toml
    agents/
      api_contract_guard.toml
      cross_repo_dependency_scout.toml
      repo_mapper.toml
      request_flow_specialist.toml
      targeted_implementer.toml
      verification_guard.toml
    tasks/
      README.md
      cross-repo-handoff-template.md
      repo-entry-template.md
      repo-implementation-template.md
      verification-review-template.md
```

## AGENTS.md Role

`AGENTS.md` will define:

- the repository role as verifier-orchestration API work for the VC ecosystem
- the current repository state as bootstrap / low-scaffold
- future code-boundary guidance without pretending the code already exists
- cross-repo relationships with `IU-cert-university`, `IU-VC-registry`, and `waltid-identity`
- multiagent rules for controller handoff, local execution, and ownership splits
- verification behavior for a repo that may not yet have runtime scripts

## Subagent Set

### `repo_mapper`

Read-only path finder for exact files, docs, and verification entry points before any edit starts.

### `request_flow_specialist`

Read-only analyst for request lifecycle, orchestration boundaries, state transitions, retries, callback handling, and VC-ecosystem touch points.

### `api_contract_guard`

Read-only reviewer focused on HTTP surface, payload contracts, auth assumptions, idempotency, error mapping, and backward-compatibility risk.

### `cross_repo_dependency_scout`

Read-only scout for deciding whether a requested change belongs here or in a sibling repo first, with explicit handoff notes for a parent controller.

### `targeted_implementer`

Workspace-write worker for small, clearly assigned edits after the parent agent has set file ownership.

### `verification_guard`

Read-only verifier for choosing the smallest relevant validation commands and for flagging when the repo lacks enough tooling to claim confidence.

## Task Template Set

Task templates under `.codex/tasks/` will support the controller pattern without depending on controller-owned files:

- `repo-entry-template.md`: intake from the thesis-level controller into this repo
- `repo-implementation-template.md`: local engineering execution contract
- `cross-repo-handoff-template.md`: explicit dependency/handoff record
- `verification-review-template.md`: verification scope and evidence record

## Controller Contract

The parent `17 Thesis` orchestration layer should do only four things before routing work here:

1. decide whether `Request-Service-API` is the owning repo
2. summarize external dependencies and sibling-repo context
3. assign a narrow objective
4. let this repo's `AGENTS.md` and `.codex/agents/` take over once execution begins

This repo should not assume the parent controller can safely override local boundaries.

## Validation

Because the repository is currently bootstrap-only, validation for this setup is structural rather than runtime-heavy:

- `find .codex -maxdepth 3 -type f | sort`
- `sed -n '1,220p' AGENTS.md`
- `git diff --check`

If package-level tooling appears later, the repo-local instructions can be tightened around real scripts and test commands.
