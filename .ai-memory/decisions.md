---
tags: [decisions, adr]
type: decisions
updated: 2026-06-12
---

# Decision log (ADR)

| #   | Date       | Decision                                          | Status   |
|-----|------------|---------------------------------------------------|----------|
| 001 | 2026-06-12 | Use Layer-First Clean Architecture                | Accepted |
| 002 | 2026-06-12 | stdlib `net/http` instead of a web framework      | Accepted |
| 003 | 2026-06-12 | PostgreSQL via `pgx` + `sqlc` (generated queries) | Accepted |
| 004 | 2026-06-12 | JWT + bcrypt for authentication                   | Accepted |
| 005 | 2026-06-12 | gRPC for the notification microservice            | Accepted |

---

## ADR-001 — Layer-First Clean Architecture
**Context:** Learning project; want clear boundaries and the Dependency Rule.
**Decision:** Organize `internal/` by layer (domain/application/infrastructure/api).
**Consequences:** Strong inward-only dependencies; folder = layer, not feature.

## ADR-002 — stdlib net/http
**Context:** Goal is to learn Go fundamentals, not a framework.
**Decision:** Build routing/middleware on `net/http`.
**Consequences:** More manual wiring, deeper understanding, fewer dependencies.

## ADR-003 — PostgreSQL + pgx + sqlc
**Decision:** Type-safe generated queries via `sqlc`; `pgx` driver; `golang-migrate` for migrations.
**Consequences:** SQL lives in `queries/`; run `make sqlc` after query changes.

## ADR-004 — JWT + bcrypt auth
**Decision:** Stateless JWT auth; bcrypt password hashing.

## ADR-005 — gRPC notification service
**Decision:** Introduce a gRPC notification service (`proto/` + generated `gen/`) as the first microservice step.

> Add new ADRs as decisions are made; update the index table.
