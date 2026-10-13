-include .env
export

.PHONY: db-up db-down run test

db-up:
	docker compose up -d --wait db

db-down:
	docker compose down

run:
	go run ./cmd/api

test:
	go vet ./...
	go test ./...
