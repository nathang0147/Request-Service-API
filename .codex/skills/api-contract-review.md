# api-contract-review

## When To Use

Use when the task changes or evaluates request and response shapes, HTTP behavior, validation, error mapping, or consumer-facing compatibility.

## Inputs Required

- task scope
- request and response contract artifacts
- relevant docs or route files
- sibling consumer assumptions if known

## Investigation Steps

1. Identify the contract boundary being changed or reviewed.
2. Check request shape, response shape, and error shape.
3. Look for ambiguous status-code or validation behavior.
4. Check idempotency and duplicate-request handling assumptions.
5. Record compatibility risks for downstream consumers.

## Expected Outputs

- contract findings
- compatibility risks
- recommended fixes or explicit non-goals

## Stop Conditions

- the contract surface is narrow and understood
- compatibility risk is stated clearly

## Escalation Conditions

- the canonical contract lives in another repo
- there is no durable contract artifact to review
