# Makefile for ZOGTest-Golang with OpenTelemetry

.PHONY: help build run test clean docker-up docker-down monitoring-start monitoring-stop monitoring-test

# Default target
help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run tests"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make docker-up          - Start all Docker services"
	@echo "  make docker-down        - Stop all Docker services"
	@echo "  make monitoring-start   - Start monitoring stack"
	@echo "  make monitoring-stop    - Stop monitoring stack"
	@echo "  make monitoring-test    - Test monitoring setup"
	@echo "  make deps               - Download dependencies"
	@echo "  make tidy               - Tidy dependencies"

# Build the application
build:
	@echo "Building application..."
	go build -v -o bin/zogtest-golang.exe .

# Run the application
run:
	@echo "Running application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Start all Docker services
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

# Stop all Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Start monitoring stack
monitoring-start:
	@echo "Starting monitoring stack..."
	powershell -ExecutionPolicy Bypass -File start-monitoring.ps1

# Stop monitoring stack
monitoring-stop:
	@echo "Stopping monitoring stack..."
	docker-compose down

# Test monitoring setup
monitoring-test:
	@echo "Testing monitoring setup..."
	powershell -ExecutionPolicy Bypass -File test-monitoring.ps1

# Development workflow
dev: docker-up
	@echo "Starting development environment..."
	@sleep 5
	@make run

# Full setup
setup: deps tidy monitoring-start
	@echo "Setup complete!"
	@echo "Run 'make run' to start the application"
