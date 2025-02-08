.PHONY: build run test dev clean install-deps install-air

# Build the application
build:
	go build -o bin/app cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install development dependencies
install-deps:
	go mod download
	go mod tidy

# Install Air for hot reloading
install-air:
	go install github.com/cosmtrek/air@latest

# Run with hot reload using Air
dev:
	air

# Default target
all: build 