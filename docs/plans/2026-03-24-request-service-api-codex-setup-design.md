# Request-Service-API Codex SDLC Migration Design

## Goal

Replace the current specialized `.codex` setup with a simpler repo-local SDLC engine that is fast to load, explicit about lifecycle state, and easy to route through in multi-repo work.

## Problems In The Current Setup

1. Repo understanding is slower than it should be because there is no persistent repo-native mini-map.
2. Agent handoff is implied by reply flow instead of formalized through durable artifacts.
3. Reusable debugging and review procedures live as prompts/templates instead of reusable skills.
4. The current six-agent layout is too granular for a bootstrap repository and creates coordination overhead.
5. The current model does not clearly separate `task`, `message`, and `handoff artifact`.
6. The current model does not clearly separate persistent work state from short-lived subagent execution state.

## Decision

Adopt a stage-based SDLC model with:

- 4 core repo-local subagents only
- persistent repo maps
- reusable repo-local skills
- deterministic templates for durable task state, stage messages, and handoff artifacts
- ephemeral subagent executions that terminate after a single stage responsibility

## SDLC Stages

The repository workflow is intentionally linear and low-complexity:

1. Intake
2. Repo Mapping
3. Design Review
4. Implementation
5. Verification
6. Cross-Repo Handoff

## Core Concepts

### Task

A `task` is the durable SDLC work unit. It persists across stages and is the source of truth for workflow state.

Required fields:

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

Task statuses:

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

A `message` is a communication event attached to a task. A message is not the task and does not own task state.

Required fields:

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

A `handoff artifact` is the formal transfer packet between stages.

Required fields:

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

A subagent execution is distinct from a task lifecycle. Each invocation is a short-lived stage executor.

Required fields:

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
- later work always starts as a fresh invocation from persisted task state and artifacts

## Agent Model

Keep only four stage executors:

1. `intake_mapper`
2. `designer_guard`
3. `implementer`
4. `verifier`

These agents are not long-lived workers. They are ephemeral stage executors that:

- load current task state and relevant artifacts
- execute one assigned stage responsibility
- emit deterministic outputs
- update task or handoff state
- terminate immediately after completion, block, or handoff emission

## Stage Ownership

- `intake_mapper`
  - owns `Intake`
  - owns `Repo Mapping`
  - owns `Cross-Repo Handoff` packaging
- `designer_guard`
  - owns `Design Review`
- `implementer`
  - owns `Implementation`
- `verifier`
  - owns `Verification`

## Stage Exit Artifacts

- Intake -> `scope-report.md`
- Repo Mapping -> `task-envelope.md`
- Design Review -> `design-review.md`
- Implementation -> `implementation-report.md`
- Verification -> `verification-report.md`
- Cross-Repo Handoff -> `cross-repo-handoff.md`

Messages can be attached at any stage using `stage-message.md`, but messages are never a substitute for task state or handoff artifacts.

## Persistent Repo Understanding Layer

Add three always-available repo-local map files:

- `.codex/repo-map.md`
  - fast repo mini-map for humans and agents
- `.codex/system-map.yaml`
  - machine-readable workflow, state transitions, stage ownership, and artifact emission rules
- `.codex/ownership-map.yaml`
  - machine-readable local vs sibling-repo responsibility map

This is the mechanism for faster repo understanding without specialized permanent agents.

## Skill Model

Move reusable specialized procedures into repo-local skills:

- `repo-bootstrap-scan.md`
- `api-contract-review.md`
- `request-lifecycle-review.md`
- `callback-debug.md`
- `cross-repo-handoff.md`
- `minimal-verification.md`

Each skill must define:

- when to use
- inputs required
- investigation steps
- expected outputs
- stop conditions
- escalation conditions

## Planned Repository Layout

```text
Request-Service-API/
  AGENTS.md
  docs/plans/
    2026-03-24-request-service-api-codex-setup-design.md
    2026-03-24-request-service-api-codex-setup.md
  .codex/
    config.toml
    repo-map.md
    system-map.yaml
    ownership-map.yaml
    agents/
      intake_mapper.toml
      designer_guard.toml
      implementer.toml
      verifier.toml
    skills/
      repo-bootstrap-scan.md
      api-contract-review.md
      request-lifecycle-review.md
      callback-debug.md
      cross-repo-handoff.md
      minimal-verification.md
    templates/
      task-envelope.md
      scope-report.md
      design-review.md
      implementation-report.md
      verification-report.md
      cross-repo-handoff.md
      stage-message.md
```

## Repo-Local Operation

The parent or controller role stays thin:

- decide high-level repo ownership
- create or route the task
- pass in relevant upstream context
- let repo-local stage execution own the rest

This repository becomes a local SDLC execution engine. It should not depend on external orchestration frameworks, long-lived agent pools, or reply-chain reconstruction.

## Bootstrap-Specific Constraints

This repository is still bootstrap-stage with only minimal content committed. The SDLC model must therefore:

- never invent runtime scripts that do not exist
- prefer structural verification until real project tooling exists
- document target folder boundaries as intended future structure, not as existing implementation

## Validation

Structural validation for the migration:

- `git diff --check`
- `find .codex -maxdepth 3 -type f | sort`
- `sed -n '1,260p' AGENTS.md`
- `sed -n '1,260p' .codex/system-map.yaml`

Because the repository is bootstrap-only, there is no broader runtime verification claim.
