# CLAUDE.md — go-ecommerce

Project-level instructions for AI assistants working in this repo.

## Persistent Memory (.ai-memory)

This project uses a persistent memory directory: **`.ai-memory/`**. Commit it to git.

**At session start**, read in this order:
1. `.ai-memory/onboarding.md`
2. `.ai-memory/context.md`
3. `.ai-memory/architecture.md`
4. `.ai-memory/decisions.md`
5. `.ai-memory/progress.md`
6. `.ai-memory/active-tasks.md`
7. `.ai-memory/session-log.md`

## Before making changes

- Confirm the change respects the **Dependency Rule** (dependencies point inward;
  `domain` imports nothing).
- Check `active-tasks.md` and `decisions.md` for relevant context/constraints.
- Use semantic navigation (LSP) over text search for Go code.
- For DB changes: edit SQL in `queries/`, run `make sqlc`; add migrations under
  `internal/infrastructure/migrations`.

## Memory-update rules

- When an architectural decision is made → add an ADR to `decisions.md` and update the index.
- When project facts/scope change → update `context.md` / `architecture.md`.
- Keep `progress.md` milestones current as slices land.
- Do **not** duplicate what the code or git history already records.

## Task-tracking rules

- Use `/e4:new-task` to open a detailed task in `tasks/` and register it in `active-tasks.md`.
- Move tasks through Current → In review → done; keep `active-tasks.md` honest.

## Session-summary rules

- At session end (`/e4:end-session`), prepend a dated entry to `session-log.md`
  with a TL;DR, and update `progress.md` / `active-tasks.md`.

## Coding rules

- Idiomatic Go; stdlib `net/http` (no heavy framework).
- Keep services, repository interfaces, and Postgres implementations in sync.
- After editing, ensure the build passes (`make build`) and tests run (`make test`).
- Fix compile errors, imports, and type issues immediately.
