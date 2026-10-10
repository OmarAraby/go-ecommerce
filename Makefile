GOBIN    := $(shell go env GOPATH)/bin
PROTOC   := $(HOME)/protoc/bin/protoc
DB_URL   := postgres://postgres:$(DB_PASSWORD)@localhost:5432/ecommerce?sslmode=disable

.PHONY: run build migrate-up migrate-down sqlc proto

run:
	go run ./cmd/api

build:
	go build -o api.exe ./cmd/api

migrate-up:
	$(GOBIN)/migrate -path internal/infrastructure/migrations -database "$(DB_URL)" up

migrate-down:
	$(GOBIN)/migrate -path internal/infrastructure/migrations -database "$(DB_URL)" down 1

sqlc:
	$(GOBIN)/sqlc generate

test:
	go test ./...

proto:
	$(PROTOC) \
		--proto_path=proto \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		notification/v1/notification.proto
