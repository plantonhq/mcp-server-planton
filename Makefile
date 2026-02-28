bump ?= patch

.PHONY: build install test lint fmt vet tidy docker-build docker-run clean release codegen-schemas codegen-types codegen help

BINARY  := mcp-server-planton
CMD     := ./cmd/mcp-server-planton
IMAGE   := ghcr.io/plantonhq/mcp-server-planton
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.DEFAULT_GOAL := help

# ─── Help ─────────────────────────────────────

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ─── Build ────────────────────────────────────

build: ## Build binary to bin/$(BINARY)
	go build -ldflags="-s -w -X github.com/plantonhq/mcp-server-planton/internal/server.buildVersion=$(VERSION)" -o bin/$(BINARY) $(CMD)

install: ## Install to GOPATH/bin
	go install $(CMD)

# ─── Test & Lint ──────────────────────────────

test: ## Run tests with race detection
	go test -v -race -timeout 30s ./...

lint: ## Run golangci-lint (or go vet)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi

fmt: ## Format Go source files
	go fmt ./...

vet: ## Run go vet (excludes gen/)
	go vet $$(go list ./... | grep -v '/gen/')

tidy: ## Run go mod tidy
	go mod tidy

# ─── Codegen ──────────────────────────────────

codegen-schemas: ## Stage 1: proto to JSON schemas
	go run ./tools/codegen/proto2schema/ --all

codegen-types: ## Stage 2: JSON schemas to Go types
	rm -rf gen/infrahub/cloudresource/
	go run ./tools/codegen/generator/ --schemas-dir=schemas --output-dir=gen/infrahub/cloudresource

codegen: codegen-schemas codegen-types ## Full codegen pipeline (Stage 1 + 2)

# ─── Docker ───────────────────────────────────

docker-build: ## Build Docker image
	docker build -t $(IMAGE):latest .

docker-run: ## Run Docker image
	docker run -i --rm \
		-e PLANTON_API_KEY=$(PLANTON_API_KEY) \
		-e PLANTON_APIS_GRPC_ENDPOINT=$(PLANTON_APIS_GRPC_ENDPOINT) \
		$(IMAGE):latest

# ─── Release ──────────────────────────────────

release: build ## Tag and push a release (usage: make release [bump=patch|minor|major])
	@LATEST_TAG=$$(git tag -l "v*" | sort -V | tail -n1); \
	[ -z "$$LATEST_TAG" ] && LATEST_TAG="v0.0.0"; \
	VERSION=$$(echo $$LATEST_TAG | sed 's/^v//'); \
	MAJOR=$$(echo $$VERSION | cut -d. -f1); \
	MINOR=$$(echo $$VERSION | cut -d. -f2); \
	PATCH=$$(echo $$VERSION | cut -d. -f3); \
	case $(bump) in \
		major) MAJOR=$$((MAJOR + 1)); MINOR=0; PATCH=0 ;; \
		minor) MINOR=$$((MINOR + 1)); PATCH=0 ;; \
		patch) PATCH=$$((PATCH + 1)) ;; \
		*) echo "error: invalid bump '$(bump)' (use patch|minor|major)" && exit 1 ;; \
	esac; \
	NEW_TAG="v$$MAJOR.$$MINOR.$$PATCH"; \
	if git rev-parse "$$NEW_TAG" >/dev/null 2>&1; then \
		echo "error: tag $$NEW_TAG already exists" && exit 1; \
	fi; \
	echo "$$LATEST_TAG -> $$NEW_TAG"; \
	git tag -a "$$NEW_TAG" -m "Release $$NEW_TAG"; \
	git push origin "$$NEW_TAG"

# ─── Clean ────────────────────────────────────

clean: ## Remove build artifacts
	rm -rf bin/ dist/
