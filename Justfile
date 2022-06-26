_default:
  @just --list --unsorted

# Run all unit tests
test:
  go test -race -cover ./...

# Lint all code
lint:
  golangci-lint run
