# go-ecommerce

E-commerce backend built in Go — a hands-on project to learn backend engineering with Go (REST, auth, databases, clean architecture, and later gRPC microservices).

## Tech Stack

- **Language:** Go 1.26
- **HTTP:** stdlib `net/http`
- **Database:** PostgreSQL + `pgx` / `sqlc`
- **Auth:** JWT + bcrypt
- **Architecture:** Clean Architecture (Layer-First)

## Architecture (Layer-First Clean Architecture)

```
cmd/api/main.go          → Composition Root (wires dependencies)
internal/
├── domain/              → Entities + domain errors
├── application/         → Services + repository interfaces
├── infrastructure/      → Postgres repositories + DB connection
└── api/                 → HTTP handlers + router + middleware
```

The Dependency Rule: dependencies point inward. `domain` imports nothing.

## Running

```bash
go run ./cmd/api
```

Server starts on `:8080`.

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```
