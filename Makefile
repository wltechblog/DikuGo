# Declare phony targets to always run regardless of file timestamps
.PHONY: all dikugo clean test

# Default target
all: dikugo

# Always compile clean - no file dependencies
dikugo:
	@echo "Building DikuGo (clean build)..."
	go clean
	go build -o dikugo cmd/dikugo/main.go
	@echo "Build complete: dikugo"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f dikugo
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests complete"