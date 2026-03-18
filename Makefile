# DIVA Server Makefile
# Build, test, and development tools for the server

# Variables
BINARY_DIR = ./bin
SERVER_BINARY = $(BINARY_DIR)/diva-server
SERVER_MAIN = ./cmd/server/main.go

.PHONY: help build run test itest clean watch sqlc dev-build dev-up dev-down dev-logs dev-shell prod-build prod-up prod-down prod-logs prod-shell db-shell rebuild ps

# Default target
help:
	@echo "DIVA Server Commands:"
	@echo ""
	@echo "Build:"
	@echo "  build        Build server binary"
	@echo "  all          Build and test"
	@echo ""
	@echo "Run:"
	@echo "  run          Run server locally"
	@echo ""
	@echo "Development:"
	@echo "  watch        Live reload with air"
	@echo "  sqlc         Generate SQL code (run when SQL files change)"
	@echo "  dev-build    Build development image"
	@echo "  dev-up       Start development environment"
	@echo "  dev-down     Stop development environment"
	@echo "  dev-logs     View development logs"
	@echo "  dev-shell    Access development container shell"
	@echo ""
	@echo "Production:"
	@echo "  prod-build   Build production image"
	@echo "  prod-up      Start production environment"
	@echo "  prod-down    Stop production environment"
	@echo "  prod-logs    View production logs"
	@echo "  prod-shell   Access production container shell"
	@echo ""
	@echo "Utilities:"
	@echo "  clean        Clean up containers, volumes, and images"
	@echo "  ps           Show running containers"
	@echo "  db-shell     Access database shell"
	@echo "  rebuild      Force rebuild without cache"

# Build targets
all: build test

build:
	@echo "Building server..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(SERVER_BINARY) $(SERVER_MAIN)
	@echo "Server built: $(SERVER_BINARY)"

# Run targets
run: build
	@echo "Running server..."
	$(SERVER_BINARY)

# Development tools
watch:
	@if command -v air > /dev/null; then \
		echo "Starting live reload..."; \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
		echo "Starting live reload..."; \
		air; \
	fi

sqlc:
	@echo "Generating SQL code..."
	@sqlc generate
	@echo "SQL code generated!"

# Testing
test:
	@echo "Running tests..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Utilities
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@echo "Cleaning Docker resources..."
	@docker compose -f docker-compose.dev.yml down -v --remove-orphans 2>/dev/null || true
	@docker compose -f docker-compose.yml down -v --remove-orphans 2>/dev/null || true
	@docker system prune -f
	@echo "Clean completed!"

ps:
	@echo "Running containers:"
	@docker ps -a

db-shell:
	@echo "Accessing database shell..."
	@docker compose -f docker-compose.dev.yml exec diva_db psql -U ${DB_USER} -d ${DB_NAME}

rebuild:
	@echo "Force rebuilding without cache..."
	@docker compose -f docker-compose.dev.yml build --no-cache
	@docker compose -f docker-compose.yml build --no-cache

# Development Commands
dev-build:
	@echo "Building development image..."
	@docker compose -f docker-compose.dev.yml build

dev-up:
	@echo "Starting development environment..."
	@sqlc generate
	@docker compose -f docker-compose.dev.yml --env-file .env.dev up -d --build
	@echo "Development environment started!"

dev-down:
	@echo "Stopping development environment..."
	@docker compose -f docker-compose.dev.yml down

dev-logs:
	@echo "Following development logs..."
	@docker compose -f docker-compose.dev.yml logs -f

dev-shell:
	@echo "Accessing development container shell..."
	@docker compose -f docker-compose.dev.yml exec diva_server sh

# Production Commands
prod-build:
	@echo "Building production image..."
	@docker compose -f docker-compose.yml build

prod-up:
	@echo "Starting production environment..."
	@sqlc generate
	@docker compose -f docker-compose.yml --env-file .env up -d --build
	@echo "Production environment started!"

prod-down:
	@echo "Stopping production environment..."
	@docker compose -f docker-compose.yml down

prod-logs:
	@echo "Following production logs..."
	@docker compose -f docker-compose.yml logs -f

prod-shell:
	@echo "Accessing production container shell..."
	@docker compose -f docker-compose.yml exec diva_server sh
