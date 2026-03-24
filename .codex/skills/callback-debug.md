# callback-debug

## When To Use

Use when debugging callback handling, follow-up events, duplicate deliveries, timeout recovery, or stuck orchestration state.

## Inputs Required

- task objective
- callback payload or reproduction notes
- relevant request lifecycle artifact
- current task state and blockers

## Investigation Steps

1. Confirm the callback entrypoint or intended callback seam.
2. Check whether the issue is missing input, duplicate input, delayed input, or bad state transition.
3. Compare expected callback effect with actual task state.
4. Identify idempotency, ordering, and timeout assumptions.
5. Reduce the issue to the smallest reproducible state transition failure.

## Expected Outputs

- failure classification
- likely faulty transition or missing guard
- next repair or handoff recommendation

## Stop Conditions

- the callback failure mode is classified clearly
- the next actor can reproduce or fix the issue

## Escalation Conditions

- the failing callback source is controlled by another repo or external system
- no reliable reproduction input exists
