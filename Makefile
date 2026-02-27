.PHONY: build install test lint fmt vet tidy docker-build docker-run clean release codegen-schemas codegen-types codegen help

BINARY  := mcp-server-planton
CMD     := ./cmd/mcp-server-planton
IMAGE   := ghcr.io/plantonhq/mcp-server-planton
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.DEFAULT_GOAL := help

# Build the server binary into bin/.
build:
	go build -ldflags="-s -w -X github.com/plantonhq/mcp-server-planton/internal/server.buildVersion=$(VERSION)" -o bin/$(BINARY) $(CMD)

# Install the binary to GOPATH/bin.
install:
	go install $(CMD)

# Run unit tests with race detection.
test:
	go test -v -race -timeout 30s ./...

# Run golangci-lint (falls back to go vet).
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi

# Format all Go source files.
fmt:
	go fmt ./...

# Run go vet on hand-written packages (generated code in gen/ is excluded
# because jsonschema-go's escaped-comma tag convention triggers false positives).
vet:
	go vet $$(go list ./... | grep -v '/gen/')

# Tidy Go modules.
tidy:
	go mod tidy

# Build Docker image.
docker-build:
	docker build -t $(IMAGE):latest .

# Run Docker image with environment variables.
docker-run:
	docker run -i --rm \
		-e PLANTON_API_KEY=$(PLANTON_API_KEY) \
		-e PLANTON_APIS_GRPC_ENDPOINT=$(PLANTON_APIS_GRPC_ENDPOINT) \
		$(IMAGE):latest

# Remove build artifacts.
clean:
	rm -rf bin/ dist/

# Create and push a release tag (usage: make release version=v1.0.0 [force=true]).
release: build
ifndef version
	$(error version is required. Usage: make release version=v1.0.0)
endif
	@if ! echo "$(version)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then \
		echo "Error: version must follow semantic versioning (e.g., v1.0.0)"; \
		exit 1; \
	fi
	@if git rev-parse $(version) >/dev/null 2>&1; then \
		if [ "$(force)" = "true" ]; then \
			git tag -d $(version); \
		else \
			echo "Error: Tag $(version) already exists locally. Use force=true to recreate."; \
			exit 1; \
		fi \
	fi
	@if git ls-remote --tags origin | grep -q "refs/tags/$(version)$$"; then \
		if [ "$(force)" = "true" ]; then \
			git push origin :refs/tags/$(version); \
		else \
			echo "Error: Tag $(version) already exists remotely. Use force=true to recreate."; \
			exit 1; \
		fi \
	fi
	git tag -a $(version) -m "Release $(version)"
	git push origin $(version)

# Stage 1: Generate JSON schemas from OpenMCF provider protos.
codegen-schemas:
	go run ./tools/codegen/proto2schema/ --all

# Stage 2: Generate Go input types from JSON schemas.
codegen-types:
	rm -rf gen/infrahub/cloudresource/
	go run ./tools/codegen/generator/ --schemas-dir=schemas --output-dir=gen/infrahub/cloudresource

# Full codegen pipeline (Stage 1 + Stage 2).
codegen: codegen-schemas codegen-types

# Show available targets.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build:"
	@echo "  build           Build binary to bin/$(BINARY)"
	@echo "  install         Install to GOPATH/bin"
	@echo "  docker-build    Build Docker image"
	@echo "  docker-run      Run Docker image"
	@echo ""
	@echo "Test & Lint:"
	@echo "  test            Run tests with race detection"
	@echo "  lint            Run golangci-lint (or go vet)"
	@echo "  vet             Run go vet (excludes gen/)"
	@echo "  fmt             Format Go source files"
	@echo ""
	@echo "Codegen:"
	@echo "  codegen-schemas Stage 1: proto -> JSON schemas"
	@echo "  codegen-types   Stage 2: JSON schemas -> Go types"
	@echo "  codegen         Full pipeline (Stage 1 + 2)"
	@echo ""
	@echo "Release:"
	@echo "  release         Create and push a release tag"
	@echo ""
	@echo "Misc:"
	@echo "  tidy            Run go mod tidy"
	@echo "  clean           Remove build artifacts"
