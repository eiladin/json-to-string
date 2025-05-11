.PHONY: build test clean

# Build the application
build:
	go build -o json-to-string ./cmd/json-to-string

# Run all tests
test:
	go test -v ./...

# Run unit tests only (skip integration tests)
test-unit:
	go test -v -short ./...

# Run the package tests only
test-pkg:
	go test -v ./pkg/...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean build artifacts
clean:
	rm -f json-to-string
	rm -f coverage.out coverage.html
	go clean

# Install the application
install:
	go install ./cmd/json-to-string

# Format Go code
fmt:
	go fmt ./...

# Vet Go code
vet:
	go vet ./...

# Run all checks
check: fmt vet test 