---
tags: [architecture, design]
type: architecture
updated: 2026-06-12
---

# Architecture

## System overview

Layer-First Clean Architecture monolith with the beginnings of a gRPC
notification microservice. Dependencies point inward; `domain` imports nothing.

```
cmd/api/main.go          → Composition Root (wires dependencies)
internal/
├── domain/              → Entities + domain errors
│   └── entities/
├── application/         → Services + repository interfaces
│   └── services/{user,product,order}/
├── infrastructure/      → Postgres repos, DB conn, migrations, auth, storage
│   ├── db/  postgres/  queries/  migrations/  auth/  storage/
└── api/                 → HTTP handlers, router, middleware
    └── {auth,users,products,orders,health,middleware,response,validate}/
proto/  + gen/           → gRPC notification service (protobuf + generated stubs)
config/                  → configuration loading
uploads/products/        → uploaded product images
```

## Tech stack

| Concern        | Choice                                  |
|----------------|-----------------------------------------|
| Language       | Go 1.26                                 |
| HTTP           | stdlib `net/http`                       |
| Database       | PostgreSQL + `pgx/v5`                   |
| Queries        | `sqlc` (generated, type-safe)           |
| Migrations     | `golang-migrate`                        |
| Auth           | JWT (`golang-jwt/v5`) + bcrypt (`x/crypto`) |
| Validation     | `go-playground/validator/v10`           |
| Config / env   | `joho/godotenv`                         |
| RPC            | gRPC + protobuf (notification service)  |

## Components

- **api** — HTTP routing, request validation, JSON responses, middleware
  (rate limiting per-IP, auth).
- **application/services** — use-case logic for user, product, order; owns
  repository interfaces.
- **infrastructure** — Postgres repository implementations, DB connection,
  migrations, sqlc queries, JWT/auth helpers, file storage for images.
- **domain/entities** — core types + domain errors.
- **notification (gRPC)** — `proto/notification/v1` + generated `gen/`.

## Data model

> TODO: document tables (users, products, product_images, orders, …) and
> relationships once schema is reviewed in `internal/infrastructure/migrations`.

## Design principles

- The Dependency Rule: dependencies point inward only.
- Interfaces defined by the consumer (`application`), implemented by `infrastructure`.
- Generated DB access (`sqlc`) over hand-written SQL in handlers.
- Vertical slices: ship one end-to-end feature per increment.
- Operational hygiene: graceful shutdown, per-IP rate limiting.
