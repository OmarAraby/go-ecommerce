GOBIN  := $(shell go env GOPATH)/bin
DB_URL := postgres://postgres:$(DB_PASSWORD)@localhost:5432/ecommerce?sslmode=disable

.PHONY: run build migrate-up migrate-down sqlc

run:
	go run ./cmd/api

build:
	go build -o api.exe ./cmd/api

migrate-up:
	$(GOBIN)/migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	$(GOBIN)/migrate -path db/migrations -database "$(DB_URL)" down 1

sqlc:
	$(GOBIN)/sqlc generate

test:
	go test ./...
