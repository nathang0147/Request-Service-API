# Request-Service-API Codex Setup Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a repo-local Codex agent setup that defines repository boundaries, controller-friendly subagents, and reusable task templates for future engineering work.

**Architecture:** The repository keeps its own `AGENTS.md` and `.codex/` directory as the source of truth for local execution. A parent controller above this repo may route work into the repository, but once inside, repo-local instructions own boundaries, agent roles, and verification expectations.

**Tech Stack:** Markdown, TOML, repo-local agent/task scaffolding

---

### Task 1: Add Repository-Level Agent Instructions

**Files:**
- Create: `AGENTS.md`
- Reference: `README.md`
- Reference: `docs/plans/2026-03-24-request-service-api-codex-setup-design.md`

**Step 1: Draft the repository role and current-state guidance**

Write `AGENTS.md` so it:

- defines the repo as the local owner for Request-Service-API engineering work
- states the repo is currently bootstrap state and should not invent missing tooling
- explains cross-repo relationships with `IU-cert-university`, `IU-VC-registry`, and `waltid-identity`

**Step 2: Add execution and boundary rules**

Add rules that:

- keep future HTTP, orchestration, integration, and domain layers separated
- forbid guessing cross-repo changes from this checkout
- require repo-local instructions to take precedence once work is routed here

**Step 3: Validate the document**

Run: `sed -n '1,220p' AGENTS.md`
Expected: clear sections for role, current state, cross-repo map, boundaries, verification, and multiagent use

**Step 4: Commit**

```bash
git add AGENTS.md docs/plans/2026-03-24-request-service-api-codex-setup-design.md docs/plans/2026-03-24-request-service-api-codex-setup.md
git commit -m "docs: add Request-Service-API agent design"
```

### Task 2: Add Repo-Local Codex Agent Definitions

**Files:**
- Create: `.codex/config.toml`
- Create: `.codex/agents/repo_mapper.toml`
- Create: `.codex/agents/request_flow_specialist.toml`
- Create: `.codex/agents/api_contract_guard.toml`
- Create: `.codex/agents/cross_repo_dependency_scout.toml`
- Create: `.codex/agents/targeted_implementer.toml`
- Create: `.codex/agents/verification_guard.toml`

**Step 1: Add the base `.codex` config**

Write `.codex/config.toml` with conservative thread, depth, and runtime settings aligned with sibling repos.

**Step 2: Add read-only specialist agents**

Create the read-only TOML files with:

- a single narrow responsibility each
- repo-specific guidance for bootstrap conditions
- explicit instruction not to edit files

**Step 3: Add the implementation worker**

Create `targeted_implementer.toml` for small, assigned edits only, with minimal verification rules and a rule to report missing tooling instead of inventing commands.

**Step 4: Validate the agent tree**

Run: `find .codex -maxdepth 3 -type f | sort`
Expected: config file, six agent TOMLs, and the task-template files from Task 3

**Step 5: Commit**

```bash
git add .codex
git commit -m "chore: add Request-Service-API codex agents"
```

### Task 3: Add Controller-Friendly Task Templates

**Files:**
- Create: `.codex/tasks/README.md`
- Create: `.codex/tasks/repo-entry-template.md`
- Create: `.codex/tasks/repo-implementation-template.md`
- Create: `.codex/tasks/cross-repo-handoff-template.md`
- Create: `.codex/tasks/verification-review-template.md`

**Step 1: Add a task-template overview**

Document when each template should be used and how a parent controller should hand work into this repo.

**Step 2: Add the repo-entry template**

Include:

- why this repo owns the task
- in-scope and out-of-scope files
- sibling-repo dependencies
- local success criteria

**Step 3: Add implementation and handoff templates**

Create templates for:

- local implementation execution
- explicit cross-repo dependency handoff
- verification evidence capture

**Step 4: Validate the templates**

Run: `sed -n '1,200p' .codex/tasks/README.md`
Expected: clear guidance for controller intake, local execution, handoff, and verification usage

**Step 5: Commit**

```bash
git add .codex/tasks
git commit -m "docs: add Request-Service-API task templates"
```

### Task 4: Structural Verification

**Files:**
- Verify: `AGENTS.md`
- Verify: `.codex/config.toml`
- Verify: `.codex/agents/*.toml`
- Verify: `.codex/tasks/*.md`

**Step 1: Check for whitespace and patch errors**

Run: `git diff --check`
Expected: no output

**Step 2: Review the created structure**

Run: `find .codex -maxdepth 3 -type f | sort`
Expected: only the intended config, agent, and task files appear

**Step 3: Review the root instructions**

Run: `sed -n '1,260p' AGENTS.md`
Expected: repo-local boundaries and controller handoff guidance are explicit

**Step 4: Commit**

```bash
git add AGENTS.md .codex docs/plans/2026-03-24-request-service-api-codex-setup-design.md docs/plans/2026-03-24-request-service-api-codex-setup.md
git commit -m "chore: bootstrap Request-Service-API agent workspace"
```
