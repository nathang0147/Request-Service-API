# minimal-verification

## When To Use

Use at Verification stage or whenever a task is about to be claimed complete.

## Inputs Required

- current task state
- changed files
- available local commands
- implementation report

## Investigation Steps

1. Discover what commands actually exist in the repo.
2. Choose the smallest command set that honestly tests the touched surface.
3. Prefer structural checks when no runtime tooling exists.
4. Run the commands and record the exact outcome.
5. State residual risk for anything not verified.

## Expected Outputs

- exact commands run
- what each command proves
- verification status
- residual risk statement

## Stop Conditions

- verification evidence matches the actual change
- unverified areas are explicit

## Escalation Conditions

- the repo lacks any command that can verify the change directly
- verification depends on another repo or missing scaffold
