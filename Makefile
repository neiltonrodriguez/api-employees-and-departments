.PHONY: help test test-verbose test-coverage test-unit clean build run \
		docker-build docker-up docker-down docker-restart docker-test \
		docker-logs docker-logs-all prometheus-logs db-logs redis-logs \
		docker-clean docker-clean-volumes migrations-status check-ports \
		docker-ps docker-stop-all

DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifndef DOCKER_COMPOSE
	DOCKER_COMPOSE := docker compose
endif

help:
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "  API Employees and Departments - Makefile Commands"
	@echo "═══════════════════════════════════════════════════════════════"
	@echo ""
	@echo "📦 Local Development:"
	@echo "  make test              - Run all unit tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report (generates HTML)"
	@echo "  make test-unit         - Run only domain unit tests"
	@echo "  make clean             - Clean build artifacts and coverage files"
	@echo "  make build             - Build the application binary"
	@echo "  make run               - Run the application locally (requires .env)"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build      - Build all Docker images"
	@echo "  make docker-up         - Start all services (runs tests first)"
	@echo "  make docker-down       - Stop all services"
	@echo "  make docker-restart    - Restart all services"
	@echo "  make docker-test       - Run tests in Docker container"
	@echo "  make docker-clean      - Stop containers and remove images"
	@echo "  make docker-clean-volumes - Stop containers and remove volumes (⚠️  deletes data)"
	@echo ""
	@echo "📊 Logs:"
	@echo "  make docker-logs       - Show API logs (follow mode)"
	@echo "  make docker-logs-all   - Show all services logs"
	@echo "  make prometheus-logs   - Show Prometheus logs"
	@echo "  make db-logs           - Show PostgreSQL logs"
	@echo "  make redis-logs        - Show Redis logs"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  make migrations-status - Check migrations status"
	@echo ""
	@echo "🔧 Troubleshooting:"
	@echo "  make check-ports       - Check if required ports are available"
	@echo "  make docker-ps         - List all running containers"
	@echo "  make docker-stop-all   - Stop ALL Docker containers (use with caution)"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"

test:
	@echo "🧪 Running unit tests..."
	@go test ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/...

test-verbose:
	@echo "🧪 Running unit tests (verbose)..."
	@go test -v ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/...

test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo ""
	@echo "📈 Coverage Summary:"
	@go tool cover -func=coverage.out
	@echo ""
	@echo "✅ HTML coverage report generated: coverage.html"

test-unit:
	@echo "🧪 Running all domain unit tests..."
	@go test -v ./internal/domain/...

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f main
	@go clean
	@echo "✅ Clean complete"

build:
	@echo "🔨 Building application..."
	@go build -o main ./cmd/main.go
	@echo "✅ Build complete: ./main"

run: build
	@echo "🚀 Starting application..."
	@./main

docker-build:
	@echo "🐳 Building Docker images..."
	@$(DOCKER_COMPOSE) build
	@echo "✅ Docker images built"

docker-up:
	@echo "🚀 Starting all services..."
	@$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "✅ Services started successfully!"
	@echo ""
	@echo "📍 Available endpoints:"
	@echo "  • API:        http://localhost:8080"
	@echo "  • Swagger:    http://localhost:8080/docs/index.html"
	@echo "  • Metrics:    http://localhost:8080/metrics"
	@echo "  • Prometheus: http://localhost:9090"
	@echo "  • PostgreSQL: localhost:5432"
	@echo "  • Redis:      localhost:6380"

docker-down:
	@echo "🛑 Stopping all services..."
	@$(DOCKER_COMPOSE) down
	@echo "✅ All services stopped"

docker-restart:
	@echo "🔄 Restarting all services..."
	@$(DOCKER_COMPOSE) restart
	@echo "✅ Services restarted"

docker-test:
	@echo "🧪 Running tests in Docker..."
	@$(DOCKER_COMPOSE) build test
	@$(DOCKER_COMPOSE) run --rm test

docker-logs:
	@echo "📋 Showing API logs (Ctrl+C to exit)..."
	@$(DOCKER_COMPOSE) logs -f app

docker-logs-all:
	@echo "📋 Showing all services logs (Ctrl+C to exit)..."
	@$(DOCKER_COMPOSE) logs -f

prometheus-logs:
	@echo "📋 Showing Prometheus logs (Ctrl+C to exit)..."
	@$(DOCKER_COMPOSE) logs -f prometheus

db-logs:
	@echo "📋 Showing PostgreSQL logs (Ctrl+C to exit)..."
	@$(DOCKER_COMPOSE) logs -f db

redis-logs:
	@echo "📋 Showing Redis logs (Ctrl+C to exit)..."
	@$(DOCKER_COMPOSE) logs -f redis

docker-clean:
	@echo "🧹 Stopping containers and removing images..."
	@$(DOCKER_COMPOSE) down --rmi all
	@echo "✅ Cleanup complete"

docker-clean-volumes:
	@echo "⚠️  WARNING: This will delete all data (PostgreSQL, Redis, Prometheus)"
	@echo "Press Ctrl+C to cancel or Enter to continue..."
	@read -r
	@echo "🧹 Stopping containers and removing volumes..."
	@$(DOCKER_COMPOSE) down -v
	@echo "✅ All volumes removed"

migrations-status:
	@echo "📊 Checking migrations status..."
	@$(DOCKER_COMPOSE) exec db psql -U postgres -d companydb -c "SELECT version, description, installed_on FROM flyway_schema_history ORDER BY installed_rank;"

check-ports:
	@echo "🔍 Checking required ports..."
	@echo ""
	@echo "Port 8080 (API):"
	@lsof -i :8080 || echo "  ✅ Available"
	@echo ""
	@echo "Port 5432 (PostgreSQL):"
	@lsof -i :5432 || echo "  ✅ Available"
	@echo ""
	@echo "Port 6380 (Redis):"
	@lsof -i :6380 || echo "  ✅ Available"
	@echo ""
	@echo "Port 9090 (Prometheus):"
	@lsof -i :9090 || echo "  ✅ Available"
	@echo ""
	@echo "💡 If any port is in use, you can:"
	@echo "   1. Stop the service using that port"
	@echo "   2. Use 'make docker-stop-all' to stop all Docker containers"
	@echo "   3. Change the port in docker-compose.yml"

docker-ps:
	@echo "📦 Running Docker containers:"
	@docker ps -a

docker-stop-all:
	@echo "⚠️  WARNING: This will stop ALL Docker containers on your system"
	@echo "Press Ctrl+C to cancel or Enter to continue..."
	@read -r
	@echo "🛑 Stopping all Docker containers..."
	@docker stop $$(docker ps -aq) 2>/dev/null || echo "No containers to stop"
	@echo "✅ All containers stopped"
