.PHONY: test test-coverage lint fmt check build clean help

# Run all tests
test:
	go test -v -race ./...

# Generate coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linters
lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Run all quality checks (format, lint, test)
check: fmt lint test
	@echo "All checks passed!"

# Build project
build:
	go build -v ./...

# Clean generated files
clean:
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Install development tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Display help
help:
	@echo "Available targets:"
	@echo "  test            - Run all tests with race detection"
	@echo "  test-coverage   - Generate test coverage report"
	@echo "  lint            - Run linters (go vet, golangci-lint)"
	@echo "  fmt             - Format code (gofmt, goimports)"
	@echo "  check           - Run all quality checks (fmt, lint, test)"
	@echo "  build           - Build all packages"
	@echo "  clean           - Clean generated files and caches"
	@echo "  install-tools   - Install development tools"
	@echo "  help            - Display this help message"
