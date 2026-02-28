---
name: Clean Makefile auto-release
overview: Replace the manual `version=` release target with Stigmer's auto-incrementing semver pattern (defaulting to patch bump), and restructure the entire Makefile for clarity using section headers and auto-generated help.
todos:
  - id: rewrite-makefile
    content: "Rewrite the mcp-server-planton Makefile: add `bump` variable, replace release target with auto-increment logic, convert to auto-help with section headers"
    status: completed
  - id: verify-release
    content: Verify the new release target by running `make help` and `make release --dry-run` style check
    status: completed
isProject: false
---

# Clean Makefile with Auto-Increment Release

## Target file

`[Makefile](Makefile)` in `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/`

## What changes

### 1. Add `bump` variable (default: patch)

Add `bump ?= patch` at the top, matching Stigmer's convention. Usage becomes:

```makefile
bump ?= patch
```

- `make release` -- auto-increments patch (e.g. v0.5.2 -> v0.5.3)
- `make release bump=minor` -- increments minor (e.g. v0.5.2 -> v0.6.0)
- `make release bump=major` -- increments major (e.g. v0.5.2 -> v1.0.0)

### 2. Replace the `release` target

Replace the current `release` target (lines 60-85) with Stigmer's auto-increment logic, adapted for this repo (single tag, no `mcp-server/` prefix tag):

```makefile
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
```

Key differences from Stigmer: keeps `build` as a prerequisite (existing behavior), and only creates a single tag (no `mcp-server/` secondary tag).

### 3. Replace hand-written `help` with auto-generated help

Replace the 20+ line manual `help` target with Stigmer's `awk`-based auto-extraction:

```makefile
help: ## Show available targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
```

### 4. Add `##`  descriptions to every target

Convert existing comment-above-target style to inline `##`  style so `help` auto-discovers them:

- `build: ## Build binary to bin/$(BINARY)`
- `install: ## Install to GOPATH/bin`
- `test: ## Run tests with race detection`
- `lint: ## Run golangci-lint (or go vet)`
- `fmt: ## Format Go source files`
- `vet: ## Run go vet (excludes gen/)`
- `tidy: ## Run go mod tidy`
- `docker-build: ## Build Docker image`
- `docker-run: ## Run Docker image`
- `clean: ## Remove build artifacts`
- `codegen-schemas: ## Stage 1: proto to JSON schemas`
- `codegen-types: ## Stage 2: JSON schemas to Go types`
- `codegen: ## Full codegen pipeline (Stage 1 + 2)`

### 5. Add section separators

Organize targets into visual sections matching Stigmer's style:

```
# --- Help ---
# --- Build ---
# --- Test and Lint ---
# --- Codegen ---
# --- Docker ---
# --- Release ---
# --- Clean ---
```

### 6. Remove the `force` flag

The `force=true` escape hatch for overwriting tags is removed. If a computed tag already exists, the release fails with a clear error -- the developer should investigate rather than blindly overwrite.

## What does NOT change

- `BINARY`, `CMD`, `IMAGE`, `VERSION` variables at the top
- Build flags and ldflags
- All existing target behaviors (build, test, lint, etc.)
- `.PHONY` declarations
- The `build` prerequisite on `release`

## Summary of removed lines

- The entire manual `help` target (lines 99-126) replaced by 2-line auto-help
- The `ifndef version` / manual version logic in `release` (lines 60-85)
- Standalone comment lines above each target (replaced by inline `##`)

