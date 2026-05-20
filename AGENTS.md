# AGENTS.md

## Project Rules

- Work from the project root: `/home/bojan/develop/horisen/Krstenica/Krstenica-gui/krstenica`.
- Do not revert or overwrite user changes unless explicitly requested.
- Treat the current git worktree as potentially dirty; inspect `git status --short` before and after edits.
- Keep changes focused on the requested task.
- Prefer small, reviewable changes over broad refactors.
- For UI changes, preserve the existing Gin template/HTMX structure and current visual language unless the user asks for redesign.
- Do not remove untracked files unless the user explicitly asks.
- For local testing, remember that full `go test ./...` may require local Postgres on `127.0.0.1:5560`.

## Coding Style

- Use idiomatic Go and run `gofmt` after editing Go files.
- Keep handlers thin; put business rules in service code when behavior is domain logic.
- Keep DTO changes explicit and compatible with existing JSON/form tags.
- Prefer existing helper functions and patterns before adding new abstractions.
- Keep template changes simple and compatible with the existing HTMX flow.
- Use ASCII in code and docs unless the file already uses Serbian Cyrillic/Latin UI text.
- Avoid adding comments unless they explain non-obvious behavior.

## Documentation Requirements

- Update `SESSION.md` whenever work is completed, interrupted, or handed off.
- Document completed work, changed files, pending tasks, blockers, and next steps.
- Include verification commands and results when tests or checks are run.
- If tests cannot be run, document the exact blocker.
- Record untracked files that are noticed but not part of the task.
- Keep documentation concise and useful for resuming work.

## SESSION.md Maintenance

- Always update `SESSION.md` before ending work.
- `SESSION.md` is the authoritative handoff for the next agent/session.
- Keep `SESSION.md` current with the real worktree state.
- Do not leave completed implementation details only in chat; summarize them in `SESSION.md`.
