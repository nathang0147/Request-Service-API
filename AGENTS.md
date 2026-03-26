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
- blocked stage exits must attach a `stage-message.md` message before terminating

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

For a local implementation task, the expected status progression is:

- `new`
- `mapped`
- `designed`
- `implementing`
- `implemented`
- `verifying`
- `done`

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
Do not reintroduce legacy `.codex/tasks/` workflow state.

Agents should not rely on free-form reply chains as the primary workflow.

## Stage Exit Artifacts

- Intake -> `scope-report.md`
- Repo Mapping -> `task-envelope.md`
- Design Review -> `design-review.md`
- Implementation -> `implementation-report.md`
- Verification -> `verification-report.md`
- Cross-Repo Handoff -> `cross-repo-handoff.md`

Messages use `stage-message.md`; blocked stage exits must attach one, but messages do not replace task state or stage exit artifacts.

## Lifecycle Rules

Clean stage exits move work as follows:

- Intake -> Repo Mapping, task status remains `new`
- Repo Mapping -> Design Review, task status becomes `mapped`
- Design Review -> Implementation, task status becomes `designed`
- Implementation -> Verification, task status becomes `implemented`
- Verification -> workflow completion, task status becomes `done`
- Cross-Repo Handoff -> local workflow completion, task status becomes `closed`

Lifecycle transition constraints:

- only `intake_mapper` may advance work through Intake, Repo Mapping, and Cross-Repo Handoff
- only `designer_guard` may advance work through Design Review
- only `implementer` may advance work through Implementation
- only `verifier` may advance work through Verification
- active stage work may end early only as `blocked`, `handoff_required`, or `closed`
- entering Implementation uses task status `implementing`
- entering Verification uses task status `verifying`

## Cross-Repo Boundaries

- `IU-cert-university` owns frontend and operator-facing application behavior.
- `IU-VC-registry` owns registry contexts, schemas, DID documents, and public metadata.
- `waltid-identity` is an upstream protocol reference only; it is not the IU implementation repo for the local v1 `walt` adapter.

If this repo is blocked on another repo's authority, set the task state to `handoff_required`, package the handoff formally, and stop.

## Intended Local Boundaries

When code is added later, keep responsibilities separated:

- `cmd/api/main.go`: service entrypoint
- `internal/verification/**`: business-level verification request lifecycle and status handling
- `internal/callback/**`: verifier callback ingestion and normalization before domain updates
- `internal/provider/**`: provider resolution and adapter boundary
- `internal/provider/walt/**`: local v1 `walt` provider adapter implementation in this repo
- `internal/persistence/postgres/**`: PostgreSQL repositories and SQL-backed storage
- `internal/platform/config/**`: config and environment loading
- `internal/platform/logger/**`: structured logging setup
- `internal/platform/middleware/**`: HTTP middleware
- `internal/platform/router/**`: route wiring
- `internal/shared/apierror/**`: API error mapping
- `internal/shared/util/**`: shared helpers
- `internal/shared/clock/**`: clock abstraction
- `db/migrations/**`: schema migrations
- `db/query/**`: SQL query definitions
- `api/openapi.yaml`: northbound contract
- `deployments/**`: container, compose, and Cloud Run assets
- `.github/workflows/**`: CI and deploy automation
- `go.mod`: Go module definition
- `Makefile`: developer task entrypoints if later added
- `tests/**`: contract, integration, and workflow verification

These are target boundaries, not a claim that they already exist.

## Verification Expectations

Current honest verification in this repository is structural:

- `git diff --check`
- `git status --short`
- `find .codex -maxdepth 3 -type f | sort`
- `find .codex/agents -maxdepth 1 -type f | sort`
- `! find .codex -maxdepth 3 -print | sort | rg -q '\.codex/tasks'`
- `rg -n 'stage_exit_artifacts:|blocked_stage_exits_require_stage_message: true|blocked stage exits must attach one|scope-report.md|task-envelope.md|design-review.md|implementation-report.md|verification-report.md|cross-repo-handoff.md|stage-message.md' AGENTS.md .codex/system-map.yaml .codex/agents`
- `rg -n 'resume_allowed|terminated_executions_are_never_resumed|subagent_executions_do_not_persist|free_form_reply_chains_are_not_workflow_state' AGENTS.md .codex/system-map.yaml`
- `sed -n '1,260p' AGENTS.md`
- `sed -n '1,260p' .codex/system-map.yaml`
- `find .codex/agents -maxdepth 1 -type f | sort`
- `find .codex -maxdepth 3 -print | sort`
- `rg -n "\.codex/tasks|targeted_implementer|specialized agent|specialized_agents" .codex AGENTS.md`
- `rg -n "task-envelope.md|scope-report.md|design-review.md|implementation-report.md|verification-report.md|cross-repo-handoff.md|stage-message.md|free-form reply chains|reply chain" .codex AGENTS.md`
- `rg -n "^  task:|^  message:|^  handoff_artifact:|resume_allowed: false|terminated_executions_are_never_resumed: true|You do not resume|do not persist" .codex AGENTS.md`

If later code adds real tooling, use the smallest relevant real command set and report exactly what was and was not verified.
