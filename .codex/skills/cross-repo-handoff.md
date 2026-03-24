# cross-repo-handoff

## When To Use

Use when the local task cannot continue because the canonical change belongs in another repository or an upstream reference must be consulted first.

## Inputs Required

- current task state
- ownership map
- findings from the current stage
- exact blocking dependency

## Investigation Steps

1. Identify why the local repo cannot proceed honestly.
2. Name the owning repo or external dependency.
3. Capture exact scope that belongs elsewhere.
4. Capture what this repo should do later, if anything.
5. Package the result as a formal handoff artifact.

## Expected Outputs

- owning repo
- handoff reason
- exact required inputs for the next repo or actor
- downstream follow-up notes for this repo

## Stop Conditions

- the owning repo is unambiguous
- the next actor has enough input to continue without reconstructing chat history

## Escalation Conditions

- ownership is disputed or ambiguous
- more than one sibling repo appears to own the canonical change
