# Bootstrap Acceptance Checklist

Use these structural checks until runtime code exists. Each assertion must be
verified by repo inspection, not by free-form interpretation.

## Architecture Assertions

- [ ] Only one provider is supported in v1: `walt`
  Method: inspect `AGENTS.md` and `.codex/repo-map.md`
  Expected result: both files state that only the local v1 `walt` provider is supported
- [ ] The HR API stays business-level and does not expose VC, verifier, or wallet details
  Method: inspect `.codex/repo-map.md`
  Expected result: business flow and bootstrap notes describe a business-level API only
- [ ] The callback path is normalized before any domain update
  Method: inspect `AGENTS.md` and `.codex/repo-map.md`
  Expected result: both files require normalization before request-state mutation
- [ ] Provider DTOs do not leak into domain boundaries
  Method: inspect `.codex/repo-map.md`
  Expected result: provider-specific DTOs are explicitly confined to provider adapters

## SDLC Anti-Regression Assertions

- [ ] Only the 4 core agents exist
  Method: `find .codex/agents -maxdepth 1 -type f | sort`
  Expected result: only `designer_guard.toml`, `implementer.toml`, `intake_mapper.toml`, and `verifier.toml`
- [ ] No legacy `.codex/tasks/*` namespace exists
  Method: `! find .codex -maxdepth 3 -print | sort | rg -q '\.codex/tasks'`
  Expected result: command exits successfully because no legacy namespace exists
- [ ] Stage outputs are template-driven, not reply-chain driven
  Method: `rg -n 'stage_exit_artifacts:|blocked_stage_exits_require_stage_message: true|blocked stage exits must attach one|scope-report.md|task-envelope.md|design-review.md|implementation-report.md|verification-report.md|cross-repo-handoff.md|stage-message.md' AGENTS.md .codex/system-map.yaml .codex/agents`
  Expected result: active workflow surfaces map every stage to a template artifact and require `stage-message.md` on blocked exits
- [ ] Task, message, and handoff records stay distinct
  Method: inspect `AGENTS.md` and `.codex/system-map.yaml`
  Expected result: separate models and required fields exist for each record type
- [ ] No persisted or resumable subagent execution surface exists
  Method: `rg -n 'resume_allowed|terminated_executions_are_never_resumed|subagent_executions_do_not_persist|free_form_reply_chains_are_not_workflow_state' AGENTS.md .codex/system-map.yaml`
  Expected result: executions are non-durable and never resumed
- [ ] `AGENTS.md` and `.codex/system-map.yaml` reflect the same lifecycle model
  Method: inspect both files together
  Expected result: same stages, same owners, same clean-exit status mapping, same blocked and handoff rules, and same execution lifecycle rules

## Tracking

- [x] `.codex/tests/acceptance-checklist.md` exists in the worktree and is not ignored by git
