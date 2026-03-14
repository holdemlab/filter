.PHONY: all build test lint fmt vet cover clean

all: fmt lint test

## Build — compile all packages
build:
	go build ./...

## Test — run all tests
test:
	go test ./... -count=1

## Lint — run golangci-lint
lint:
	golangci-lint run ./...

## Fmt — format code with gofmt and goimports
fmt:
	gofmt -w .
	goimports -w .

## Vet — run go vet
vet:
	go vet ./...

## Cover — run tests with coverage report
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out

## Clean — remove build artifacts
clean:
	go clean ./...
	rm -f coverage.out

## Help — show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
