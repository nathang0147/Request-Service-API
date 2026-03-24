# repo-bootstrap-scan

## When To Use

Use when the repository is sparse, newly bootstrapped, or when repo ownership and actual implementation surface are unclear.

## Inputs Required

- task objective
- current repo tree
- `AGENTS.md`
- `.codex/repo-map.md`
- `.codex/system-map.yaml`
- `.codex/ownership-map.yaml`

## Investigation Steps

1. Read the repo-native maps and root workflow contract.
2. Inspect the actual file tree and recent commits.
3. Separate current reality from intended future structure.
4. Identify the smallest valid local scope for the task.
5. List missing scaffolding or missing tooling explicitly.

## Expected Outputs

- current-state summary
- valid local scope
- bootstrap-specific blockers
- next recommended stage or owning repo

## Stop Conditions

- repo scope is clear enough for Intake or Repo Mapping
- the task is clearly outside this repo

## Escalation Conditions

- the objective depends on files or tooling that do not exist yet
- sibling repo ownership is stronger than local ownership
