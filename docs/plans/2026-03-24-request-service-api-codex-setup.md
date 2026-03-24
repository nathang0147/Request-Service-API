# Request-Service-API Codex SDLC Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refactor the repository's Codex setup into a lightweight stage-based SDLC workflow with 4 core agents, persistent repo maps, reusable skills, explicit task/message/handoff lifecycle, and ephemeral subagent execution rules.

**Architecture:** The repository keeps a small repo-local control plane under `.codex/`. Durable workflow state lives in task and artifact templates, stage ownership lives in `AGENTS.md` plus `.codex/system-map.yaml`, and specialized procedures move into reusable skill documents. Subagents are stage-local executors only and are never resumed.

**Tech Stack:** Markdown, TOML, YAML, repo-local orchestration docs

---

### Task 1: Update The Design Record

**Files:**
- Modify: `docs/plans/2026-03-24-request-service-api-codex-setup-design.md`
- Modify: `docs/plans/2026-03-24-request-service-api-codex-setup.md`

**Step 1: Rewrite the design doc**

Replace the old specialized-agent design with:

- 4 stage agents only
- explicit `task`, `message`, and `handoff artifact` definitions
- explicit ephemeral `subagent execution` definition
- persistent repo maps and reusable skills

**Step 2: Rewrite the implementation plan**

Write a plan that covers:

- `AGENTS.md`
- `.codex/config.toml`
- `.codex/repo-map.md`
- `.codex/system-map.yaml`
- `.codex/ownership-map.yaml`
- 4 agent TOMLs
- 6 skill documents
- 7 template files
- deletion of old specialized agent/task files

**Step 3: Validate docs**

Run: `sed -n '1,260p' docs/plans/2026-03-24-request-service-api-codex-setup-design.md`
Expected: explicit task/message/handoff/execution lifecycle model

**Step 4: Commit**

```bash
git add docs/plans/2026-03-24-request-service-api-codex-setup-design.md docs/plans/2026-03-24-request-service-api-codex-setup.md
git commit -m "docs: update codex SDLC migration design"
```

### Task 2: Replace Root Workflow Instructions

**Files:**
- Modify: `AGENTS.md`

**Step 1: Rewrite repository workflow rules**

Make `AGENTS.md` define:

- the repo-local SDLC purpose
- the stage model
- the durable task model
- the message model
- the handoff artifact model
- the ephemeral subagent execution model

**Step 2: Add explicit anti-pattern rules**

State clearly:

- free-form reply chains are not the primary workflow record
- subagent executions terminate after one stage outcome
- terminated executions are never resumed

**Step 3: Validate**

Run: `sed -n '1,260p' AGENTS.md`
Expected: clear lifecycle and stage rules

**Step 4: Commit**

```bash
git add AGENTS.md
git commit -m "docs: rewrite Request-Service-API agent workflow"
```

### Task 3: Replace The `.codex` Control Plane

**Files:**
- Modify: `.codex/config.toml`
- Create: `.codex/repo-map.md`
- Create: `.codex/system-map.yaml`
- Create: `.codex/ownership-map.yaml`
- Delete: `.codex/agents/api_contract_guard.toml`
- Delete: `.codex/agents/cross_repo_dependency_scout.toml`
- Delete: `.codex/agents/repo_mapper.toml`
- Delete: `.codex/agents/request_flow_specialist.toml`
- Delete: `.codex/agents/targeted_implementer.toml`
- Delete: `.codex/agents/verification_guard.toml`
- Create: `.codex/agents/intake_mapper.toml`
- Create: `.codex/agents/designer_guard.toml`
- Create: `.codex/agents/implementer.toml`
- Create: `.codex/agents/verifier.toml`
- Delete: `.codex/tasks/README.md`
- Delete: `.codex/tasks/cross-repo-handoff-template.md`
- Delete: `.codex/tasks/repo-entry-template.md`
- Delete: `.codex/tasks/repo-implementation-template.md`
- Delete: `.codex/tasks/verification-review-template.md`
- Create: `.codex/templates/task-envelope.md`
- Create: `.codex/templates/scope-report.md`
- Create: `.codex/templates/design-review.md`
- Create: `.codex/templates/implementation-report.md`
- Create: `.codex/templates/verification-report.md`
- Create: `.codex/templates/cross-repo-handoff.md`
- Create: `.codex/templates/stage-message.md`
- Create: `.codex/skills/repo-bootstrap-scan.md`
- Create: `.codex/skills/api-contract-review.md`
- Create: `.codex/skills/request-lifecycle-review.md`
- Create: `.codex/skills/callback-debug.md`
- Create: `.codex/skills/cross-repo-handoff.md`
- Create: `.codex/skills/minimal-verification.md`

**Step 1: Replace maps and agent definitions**

Create the persistent maps and the 4 stage agent definitions.

**Step 2: Replace task templates with lifecycle templates**

Create templates for:

- durable task record
- intake scope report
- stage message
- design review output
- implementation report
- verification report
- cross-repo handoff artifact

**Step 3: Add reusable skills**

Add short deterministic skill documents with:

- when to use
- inputs required
- investigation steps
- expected outputs
- stop conditions
- escalation conditions

**Step 4: Remove old files**

Delete the specialized agent TOMLs and the old `.codex/tasks/*` layout.

**Step 5: Validate**

Run: `find .codex -maxdepth 3 -type f | sort`
Expected: only the new control-plane files remain

**Step 6: Commit**

```bash
git add .codex
git commit -m "chore: simplify codex SDLC control plane"
```

### Task 4: Structural Verification

**Files:**
- Verify: `AGENTS.md`
- Verify: `.codex/**`
- Verify: `docs/plans/**`

**Step 1: Check patch correctness**

Run: `git diff --check`
Expected: no output

**Step 2: Check working tree**

Run: `git status --short`
Expected: clean after commit

**Step 3: Review the final file set**

Run: `find .codex -maxdepth 3 -type f | sort`
Expected: 4 agent files, 6 skill files, 7 template files, 3 map/config files

**Step 4: Review recent commits**

Run: `git log --oneline -3`
Expected: design-update commit and migration commit visible

**Step 5: Commit**

```bash
git add AGENTS.md .codex docs/plans
git commit -m "docs: finalize codex SDLC migration"
```
