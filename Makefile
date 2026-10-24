-include .env
export

.PHONY: db-up db-down run test up down

db-up:
	docker compose up -d --wait db

db-down:
	docker compose down

run:
	go run ./cmd/api

test:
	go vet ./...
	go test ./...

up:
	docker compose up -d --build --wait

down:
	docker compose down
