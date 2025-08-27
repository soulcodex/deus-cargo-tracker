set shell := ["bash", "-cu"]

# Define colors for use inside recipes

ERROR_COLOR := '\033[0;31m'
NOTICE_COLOR := '\033[0;36m'
RESET_COLOR := '\033[0m'

# Docker env variables

export COMPOSE_BAKE := "true"

# Print the available recipes.
help:
    just --list
    echo "To execute a recipe, run: just <recipe-name> (e.g. just install)"

# Install the pre-commit Git hook.
install-git-hook:
    @echo -e "\n🔧 Installing pre-commit Git hook..."
    @cp scripts/git-pre-commit-hook .git/hooks/pre-commit
    @chmod +x .git/hooks/pre-commit
    @echo -e "\n✅ Pre-commit hook installed at .git/hooks/pre-commit"

# install the default environment file
install-env:
    @echo -e "\n🔧 Installing environment files..." && \
    [ ! -f .env ] && cp .env.default .env || echo "❌ .env already exists, skipping copy." && \
    [ ! -f test/.test.env ] && cp .env test/.test.env || echo "❌ test/.test.env already exists, skipping copy."

# Install required tools and deps.
install:
    #!/usr/bin/env bash
    echo -e "🔧 Installing tools...\n"

    echo -e "\n📦 Installing GolangCI-Lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.0

    echo -e "\n📦 Installing moq..."
    go install github.com/matryer/moq@latest

    just install-env && just install-git-hook

# Lint the code using gofmt, go vet, and golangci-lint.
lint:
    go fmt ./... && \
    go vet ./... && \
    golangci-lint run --timeout 2m

# Run tests suites
test:
    @echo -e "\n🔍 Running tests..."
    @go test --race ./... -count=1 --timeout=30s

# Run API service locally.
run-api:
    @echo -e "\n⚡ Running API..."
    @just install-env && echo -e "\n"
    @go run cmd/api/main.go

# Run consumer service locally.
run-consumer:
    @echo -e "\n⚡ Running Consumers..."
    @just install-env && echo -e "\n"
    @go run cmd/consumer/main.go

# Start the application components through docker compose.
up:
    docker compose \
    	-f ./deployments/docker-compose/docker-compose.yaml \
    	up -d --build --force-recreate

# Shutdown the application components through docker compose.
down:
    docker compose \
    	-f ./deployments/docker-compose/docker-compose.yaml \
    	down --remove-orphans --volumes

# Generate go code from go:generate directives.
go-generate:
    go generate ./...
