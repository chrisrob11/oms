# Makefile for basic Go project tasks

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOCMD=$(GOBASE)/cmd
GO_ENV_PATH=$(shell go env GOPATH)

# Compile and place binaries into the bin directory.
build:
	@echo "  >  Building binaries..."
	@GOBIN=$(GOBIN) CGO_ENABLED=0 go install  $(GOCMD)/...

# Generate Go code based on source changes.
generate:
	@echo "  >  Generating dependency files..."
	@go generate ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t operationwave:latest -f ./deployments/Dockerfile .

# Run linters using golangci-lint.
lint:
	@echo "  >  Linting code..."
	${GO_ENV_PATH}/bin/golangci-lint --color always run

test:
	@echo ">> running tests"
	@go test -race ./...

coverage:
	@echo ">> running tests with coverage on"
	@go test -coverprofile out.coverage -v -race ./...
	@go tool cover -html out.coverage
# Execute all steps: generate code, lint, and build.
all: generate lint build ./...

tools:
	@echo "  >  Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

# devtools are run on dev boxes but not CI
devtools:
	@go install github.com/rubenv/sql-migrate/...@latest
	@go install github.com/jmattheis/goverter/cmd/goverter@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

migrate-up:
	@echo "  >  Migrating db..."
	sql-migrate up --config=./migration_config.yml

# Clean up generated files and binaries.
clean:
	@echo "  >  Cleaning build cache"
	@go clean ./...
	@echo "  >  Removing binaries..."
	@rm -rf $(GOBIN)/*

.PHONY: build generate lint all clean docker-build migrate-up
