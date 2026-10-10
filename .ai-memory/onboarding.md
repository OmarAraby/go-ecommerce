---
tags: [onboarding, entry-point]
type: onboarding
updated: 2026-06-12
---

# Onboarding — go-ecommerce

## What this is

E-commerce backend written in Go — a hands-on project for learning backend
engineering (REST, auth, databases, Clean Architecture, and later gRPC
microservices). Built slice-by-slice; each commit is a vertical "Slice N".

## Read order

`onboarding.md` → `context.md` → `architecture.md` → `decisions.md` →
`progress.md` → `active-tasks.md` → `session-log.md`

## Golden rules

- **Respect the Dependency Rule.** Dependencies point inward. `domain` imports
  nothing. `application` defines repository interfaces; `infrastructure`
  implements them. `api` is the outermost layer.
- **Layer-First Clean Architecture** — folders are layers, not features.
- Keep services, repository interfaces, and Postgres implementations in sync.
- DB access goes through `sqlc`-generated queries (`make sqlc`), not hand-rolled SQL in handlers.
- Update `.ai-memory/` when project facts change (see CLAUDE.md).

## Environment setup

- Go 1.26
- PostgreSQL running locally (DB `ecommerce` on `localhost:5432`)
- Env: `DB_PASSWORD` (used by Makefile `DB_URL`); see `config/` and `.env` (godotenv)
- Tools on `GOPATH/bin`: `migrate`, `sqlc`; `protoc` at `~/protoc/bin/protoc`

## How to run

```bash
make run          # go run ./cmd/api
make build        # build api.exe
make migrate-up   # apply DB migrations
make migrate-down # roll back one migration
make sqlc         # regenerate query code
make proto        # regenerate gRPC stubs into gen/
```

## How to test

```bash
make test         # go test ./...
```
