# request-lifecycle-review

## When To Use

Use when the task involves request orchestration, retries, callbacks, timeouts, state transitions, or coordination across external dependencies.

## Inputs Required

- task scope
- lifecycle design or implementation artifacts
- repo map and ownership map

## Investigation Steps

1. Identify the request entrypoint and terminal outcomes.
2. Map state transitions from intake to completion or failure.
3. Check retry, timeout, and callback assumptions.
4. Look for points where external failures can leave local state inconsistent.
5. Record required invariants and open risks.

## Expected Outputs

- lifecycle map
- invariant checklist
- state-risk findings

## Stop Conditions

- the request path and failure modes are explicit
- blocking lifecycle gaps are documented

## Escalation Conditions

- state ownership actually belongs in another repo
- callback behavior depends on undocumented external systems
