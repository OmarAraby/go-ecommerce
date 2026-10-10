---
tags: [context, scope]
type: context
updated: 2026-06-12
---

# Context

## Purpose

A learning-focused e-commerce backend in Go. Primary goal: master backend
engineering with Go — REST APIs, authentication, relational data, Clean
Architecture, and eventually gRPC-based microservices.

## Goals

- Build a working e-commerce API (users, auth, products, orders).
- Practice idiomatic Go and Clean Architecture boundaries.
- Grow toward microservices (notification service over gRPC is being scaffolded).

## Scope

- **In scope:** auth (JWT + bcrypt), product CRUD with image upload, orders,
  pagination/filtering/sorting, rate limiting, graceful shutdown, gRPC
  notification service.
- **Out of scope (for now):** front-end, payment gateways, deployment infra.

## Constraints

- Standard library `net/http` for HTTP (no heavy web framework).
- PostgreSQL via `pgx` + `sqlc` (generated, type-safe queries).
- Keep the Dependency Rule intact at all times.

## Domain glossary

- **Slice** — one vertical feature increment; commits are tagged "Slice N".
- **Composition Root** — `cmd/api/main.go`, where dependencies are wired.
- **Repository interface** — defined in `application`, implemented in `infrastructure/postgres`.

> Fill in further domain terms (order lifecycle, product image rules, etc.) as they solidify.
