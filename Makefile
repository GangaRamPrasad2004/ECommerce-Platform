.PHONY: help run build dev lint format migrate-create migrate-up migrate-down docker-up docker-down docs-generate graph-generate

help:
	@echo "Available commands:"
	@echo "  make build          Build the application"
	@echo "  make run            Run the application"
	@echo "  make dev            Run the application in development mode"
	@echo "  make lint           Format and lint the code"
	@echo "  make format         Format the code"
	@echo "  make migrate-create name=<name>  Create a new SQL migration pair"
	@echo "  make migrate-up     Apply database migrations"
	@echo "  make migrate-down   Roll back database migrations"
	@echo "  make docker-up      Start Docker services"
	@echo "  make docker-down    Stop Docker services"

build:
	@echo "Building all binaries...."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
    		if [ -d "$$cmd" ]; then \
    			binary=$$(basename $$cmd); \
    			echo "Building $$binary..."; \
    			go build -o bin/$$binary ./$$cmd; \
    		fi \
    	done
run:
	@go run ./cmd/api

dev:
	@go run ./cmd/api

lint:format
	@golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .


docs-generate:
	@mkdir -p docs
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

migrate-create:
	@migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	@migrate -path db/migrations -database "postgresql://postgres:Renuka@b0315@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	@migrate -path db/migrations -database "postgresql://postgres:Renuka@b0315@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	@docker compose -f docker/docker-compose.yml up -d

docker-down:
	@docker compose -f docker/docker-compose.yml down

graph-generate:
	@go get github.com/99designs/gqlgen@v0.17.78
	@go run github.com/99designs/gqlgen generate