.PHONY: help test test-verbose test-coverage test-unit clean build run docker-build docker-up docker-down docker-test docker-logs

help:
	@echo "Available commands:"
	@echo "  make test              - Run all unit tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-unit         - Run only domain unit tests"
	@echo "  make clean             - Clean build artifacts and coverage files"
	@echo "  make build             - Build the application binary"
	@echo "  make run               - Run the application locally"
	@echo "  make docker-build      - Build Docker images"
	@echo "  make docker-up         - Start all Docker services (runs tests first)"
	@echo "  make docker-down       - Stop all Docker services"
	@echo "  make docker-test       - Run tests in Docker"
	@echo "  make docker-logs       - Show Docker logs"

test:
	go test ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/...

test-verbose:
	go test -v ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/...

test-coverage:
	go test -v ./internal/domain/validators/... ./internal/domain/employee/... ./internal/domain/department/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

test-unit:
	go test -v ./internal/domain/...

clean:
	rm -f coverage.out coverage.html
	rm -f main
	go clean

build:
	go build -o main ./cmd/main.go

run: build
	./main

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-test:
	docker-compose build test
	docker-compose run --rm test

docker-logs:
	docker-compose logs -f app
