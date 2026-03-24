# Request-Service-API Task Templates

These templates are for repo-local execution and for handoff from a parent controller in `17 Thesis`.

Use them to keep task routing explicit:

- `repo-entry-template.md`: when a portfolio-level controller routes a task into this repository
- `repo-implementation-template.md`: when the local parent agent assigns concrete engineering work
- `cross-repo-handoff-template.md`: when the requested change belongs in another repository first
- `verification-review-template.md`: when capturing what was verified, what was not, and why

Guidelines:

- keep repo ownership explicit
- name sibling-repo dependencies instead of hiding them
- define in-scope files before implementation starts
- keep verification proportional to the actual repository state
- if the repo is still bootstrap-only, say so instead of inventing runtime commands
