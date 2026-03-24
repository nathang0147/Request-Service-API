# AGENTS.md

## Repository Role

`Request-Service-API` is the repo-local SDLC execution engine for backend request orchestration, service-side integration logic, and API contracts that sit between IU frontend flows and the VC ecosystem.

The parent controller above this repo stays thin:

- decide high-level repo ownership
- route work into this repository
- pass upstream context when needed

Once work enters this repository, repo-local workflow files under `.codex/` are authoritative.

## Current State

- This repository is still bootstrap-only.
- Current committed project content is minimal and does not yet include runtime source code or package tooling.
- Do not invent scripts, build systems, source folders, or test commands that do not exist.
- Prefer structural verification until real code and tooling exist.

## SDLC Stages

This repository uses a small stage-based workflow:

1. Intake
2. Repo Mapping
3. Design Review
4. Implementation
5. Verification
6. Cross-Repo Handoff

## Workflow Record Types

### Task

A `task` is the durable SDLC work unit. It is the source of truth for workflow state.

Required task fields:

- `task_id`
- `repo_name`
- `title`
- `objective`
- `stage`
- `status`
- `owner`
- `scope`
- `constraints`
- `inputs`
- `artifacts`
- `changed_files`
- `verification_status`
- `blockers`
- `next_expected_actor`

Allowed task statuses:

- `new`
- `mapped`
- `designed`
- `implementing`
- `implemented`
- `verifying`
- `done`
- `blocked`
- `handoff_required`
- `closed`

### Message

A `message` is a communication event attached to a task. A message is not the task and must not become the primary workflow record.

Required message fields:

- `message_id`
- `task_id`
- `from_actor`
- `to_actor`
- `message_type`
- `summary`
- `body`
- `related_artifacts`
- `timestamp`

### Handoff Artifact

A `handoff artifact` is the formal stage-transfer packet. It moves work between stages or between repos without forcing the next actor to reconstruct a reply chain.

Required handoff fields:

- `task_id`
- `repo_name`
- `from_stage`
- `to_stage`
- `objective`
- `scope`
- `assumptions`
- `findings`
- `changed_files`
- `open_risks`
- `blockers`
- `required_inputs`
- `verification_status`
- `next_expected_actor`

### Subagent Execution

A subagent execution is a short-lived stage executor, not a persistent worker.

Required execution fields:

- `execution_id`
- `task_id`
- `agent_name`
- `owned_stage`
- `started_at`
- `ended_at`
- `execution_outcome`
- `produced_artifacts`
- `next_expected_actor`

Allowed execution outcomes:

- `completed`
- `blocked`
- `handoff_emitted`
- `aborted`

Hard rules:

- tasks persist
- messages persist
- artifacts persist
- subagent executions do not persist
- terminated subagent executions are never resumed
- if later work is needed, start a fresh invocation from persisted task state and artifacts

## Core Subagents

Only these repo-local subagents exist:

- `intake_mapper`
- `designer_guard`
- `implementer`
- `verifier`

These subagents are ephemeral stage executors. Each invocation must:

- load the current task state and relevant persisted artifacts
- execute only its assigned stage responsibility
- emit deterministic outputs
- update task state or emit handoff state
- terminate immediately after completion, block, or handoff emission

Do not keep dormant subagents alive.

## Stage Ownership

- `intake_mapper`
  - owns `Intake`
  - owns `Repo Mapping`
  - owns `Cross-Repo Handoff`
- `designer_guard`
  - owns `Design Review`
- `implementer`
  - owns `Implementation`
- `verifier`
  - owns `Verification`

## Workflow Files

Start with the persistent repo-native maps:

- `.codex/repo-map.md`
- `.codex/system-map.yaml`
- `.codex/ownership-map.yaml`

Use repo-local reusable skills under `.codex/skills/`.
Use structured templates under `.codex/templates/`.

Agents should not rely on free-form reply chains as the primary workflow.

## Stage Exit Artifacts

- Intake -> `scope-report.md`
- Repo Mapping -> `task-envelope.md`
- Design Review -> `design-review.md`
- Implementation -> `implementation-report.md`
- Verification -> `verification-report.md`
- Cross-Repo Handoff -> `cross-repo-handoff.md`

Messages may be attached with `stage-message.md`, but messages do not replace task state or stage exit artifacts.

## Cross-Repo Boundaries

- `IU-cert-university` owns frontend and operator-facing application behavior.
- `IU-VC-registry` owns registry contexts, schemas, DID documents, and public metadata.
- `waltid-identity` is an upstream protocol reference, not the default implementation target for IU work.

If this repo is blocked on another repo's authority, set the task state to `handoff_required`, package the handoff formally, and stop.

## Intended Local Boundaries

When code is added later, keep responsibilities separated:

- `src/http/**`: transport, validation, response mapping
- `src/orchestration/**`: request lifecycle, state transitions, retries, callback handling
- `src/integrations/**`: external adapters and VC ecosystem clients
- `src/domain/**`: pure domain logic and invariants
- `src/config/**`: config and environment loading
- `tests/**`: contract, integration, and workflow verification

These are target boundaries, not a claim that they already exist.

## Verification Expectations

Current honest verification in this repository is structural:

- `git diff --check`
- `git status --short`
- `find .codex -maxdepth 3 -type f | sort`
- `sed -n '1,260p' AGENTS.md`

If later code adds real tooling, use the smallest relevant real command set and report exactly what was and was not verified.
