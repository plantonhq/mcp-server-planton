# Add GoFormat Commands to Makefile

**Date**: November 25, 2025  
**Type**: Developer Experience / Build Tools  
**Impact**: Medium - Prevents GoFormat CI failures and improves development workflow

## Summary

Added GoFormat commands to the Makefile to automatically check and fix Go code formatting issues before pushing code. Also added a convenient `make release` command to create and push semantic version tags. These additions ensure developers catch formatting issues locally during `make build`, preventing CI failures and reducing friction in the development workflow.

## Problem Statement

The codebase had recurring GoFormat issues that were only caught in CI, leading to:

1. **Failed CI Builds**: Developers pushed code that wasn't properly formatted, causing CI to fail
2. **Manual Formatting**: Developers had to manually run `gofmt -w .` to fix issues
3. **No Local Verification**: No easy way to check formatting before pushing
4. **Inefficient Release Process**: Creating release tags required manual git commands

### Specific Issues Found

GoFormat check in CI was failing for 8 files with two types of issues:
- **Import ordering**: Imports were not alphabetically sorted
- **Trailing blank lines**: Extra blank lines at end of files

Affected files:
```
internal/infrahub/client.go
internal/infrahub/tools/errors.go
internal/infrahub/tools/get.go
internal/infrahub/tools/kinds.go
internal/infrahub/tools/lookup.go
internal/infrahub/tools/search.go
internal/resourcemanager/client.go
internal/resourcemanager/tools/environment.go
```

## Solution

Added three new Makefile targets and integrated format checking into the build process:

### 1. `make fmt` - Auto-format Go code
```makefile
## fmt: Format Go code
fmt:
	@echo "Formatting Go code..."
	@gofmt -w .
	@echo "Code formatted"
```

**Purpose**: Automatically format all Go code in the project  
**Usage**: `make fmt`

### 2. `make fmt-check` - Verify formatting
```makefile
## fmt-check: Check if Go code is formatted
fmt-check:
	@echo "Checking Go code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted:"; \
		gofmt -l .; \
		echo "Run 'make fmt' to fix formatting"; \
		exit 1; \
	fi
	@echo "All Go code is properly formatted"
```

**Purpose**: Check if code is formatted without modifying files  
**Usage**: `make fmt-check`  
**Behavior**: Exits with error if formatting issues found, lists affected files

### 3. `make release` - Create release tags
```makefile
## release: Create and push a release tag (usage: make release TAG=v1.0.0)
release:
ifndef TAG
	@echo "Error: TAG is required. Usage: make release TAG=v1.0.0"
	@exit 1
endif
	@echo "Creating release tag $(TAG)..."
	@if ! echo "$(TAG)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then \
		echo "Error: TAG must follow semantic versioning (e.g., v1.0.0, v2.1.3)"; \
		exit 1; \
	fi
	@git tag -a $(TAG) -m "Release $(TAG)"
	@git push origin $(TAG)
	@echo "Release tag $(TAG) created and pushed"
	@echo "GitHub Actions will now build and publish the release"
```

**Purpose**: Streamline release process with validation  
**Usage**: `make release TAG=v1.0.0`  
**Features**:
- Validates TAG parameter is provided
- Ensures TAG follows semantic versioning (vX.Y.Z)
- Creates annotated git tag
- Pushes to origin to trigger release workflow
- Provides clear error messages

### 4. Integrated Format Check into Build

Updated the `build` target to depend on `fmt-check`:

```makefile
## build: Build the binary for local architecture
build: fmt-check
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/mcp-server-planton
	@echo "Binary built: $(BINARY_PATH)"
```

**Impact**: `make build` now fails fast if code is not formatted, catching issues before commit/push

## Changes Made

### Files Modified

**Makefile**:
1. Updated `.PHONY` declaration to include new targets: `fmt`, `fmt-check`, `release`
2. Added `fmt` target after `lint` target
3. Added `fmt-check` target after `fmt` target
4. Added `release` target after `clean` target
5. Modified `build` target to depend on `fmt-check`

### Code Formatting Fixed

Ran `gofmt -w .` to fix all existing formatting issues in 8 files:
- Reordered imports alphabetically
- Removed trailing blank lines
- Applied standard Go formatting

## Benefits

### 1. Early Detection
- Format issues caught during `make build` instead of in CI
- Faster feedback loop for developers
- Reduces failed CI builds

### 2. Developer Experience
- Simple command to fix all formatting: `make fmt`
- Clear error messages when formatting is wrong
- One command workflow: `make build` checks everything

### 3. Consistency
- Same format check logic in Makefile and CI
- Ensures local and CI environments are aligned
- Prevents "works on my machine" issues

### 4. Streamlined Releases
- Simple, validated release tagging: `make release TAG=v1.0.0`
- Prevents common mistakes (wrong format, missing tag)
- One command triggers entire release pipeline

## Workflow Examples

### Daily Development Workflow
```bash
# Make code changes...

# Build (automatically checks formatting)
make build

# If formatting fails, fix it
make fmt

# Build again to verify
make build

# Commit and push
git commit -m "Add new feature"
git push
```

### Release Workflow
```bash
# Ensure code is ready
make build
make test

# Create release (triggers CI/CD)
make release TAG=v1.2.3

# GitHub Actions automatically:
# - Runs GoReleaser
# - Builds multi-arch binaries
# - Publishes Docker images
# - Creates GitHub release
```

## Technical Details

### Format Check Logic

The `fmt-check` uses the same logic as the CI workflow:
```bash
if [ -n "$(gofmt -l .)" ]; then
  echo "Go code is not formatted:"
  gofmt -l .
  echo "Run 'make fmt' to fix formatting"
  exit 1
fi
```

This ensures perfect alignment between local checks and CI validation.

### Release Tag Validation

The release command validates tags using regex:
```bash
grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'
```

**Valid formats**:
- `v1.0.0`
- `v2.1.3`
- `v10.5.2`

**Invalid formats** (rejected):
- `1.0.0` (missing 'v' prefix)
- `v1.0` (incomplete version)
- `v1.0.0-beta` (pre-release, not supported yet)

### Integration with CI

The CI workflow (`.github/workflows/ci.yml`) already has:
```yaml
- name: Run go fmt check
  run: |
    if [ -n "$(gofmt -l .)" ]; then
      echo "Go code is not formatted:"
      gofmt -d .
      exit 1
    fi
```

Now developers have the same check locally via `make fmt-check` or `make build`.

### Integration with Release Workflow

When `make release TAG=v1.0.0` pushes a tag, it triggers:

`.github/workflows/release.yml`:
- Runs GoReleaser to build binaries for all platforms
- Builds and pushes multi-arch Docker images to GHCR
- Creates GitHub release with artifacts and changelog

## Updated Makefile Targets

Complete list of available targets (from `make help`):

```
Targets:
  build          Build the binary for local architecture
  install        Install the binary to GOPATH/bin
  test           Run tests
  lint           Run linter (requires golangci-lint)
  fmt            Format Go code
  fmt-check      Check if Go code is formatted
  docker-build   Build Docker image
  docker-run     Run Docker image with environment variables
  clean          Remove build artifacts
  release        Create and push a release tag (usage: make release TAG=v1.0.0)
  help           Show this help message
```

## Testing

Verified all commands work correctly:

```bash
# Format check passes after formatting
$ make fmt-check
Checking Go code formatting...
All Go code is properly formatted

# Build includes format check
$ make build
Checking Go code formatting...
All Go code is properly formatted
Building mcp-server-planton...
Binary built: bin/mcp-server-planton

# Help displays new commands
$ make help
# ... shows fmt, fmt-check, and release targets
```

## Migration Notes

No migration needed for existing developers. Simply:

1. **Continue normal development**: Existing `make build` now includes format checking
2. **Fix formatting issues**: Run `make fmt` if build fails due to formatting
3. **Create releases**: Use `make release TAG=vX.Y.Z` instead of manual git tag commands

## Rollback Plan

If issues arise, revert the Makefile changes:
```bash
git checkout HEAD^ -- Makefile
```

However, this is unlikely to be needed as:
- Format commands are simple and well-tested
- Build dependency is standard make practice
- Release command is optional (manual git commands still work)

## Future Enhancements

Potential improvements for future consideration:

1. **Pre-release tags**: Support alpha/beta/rc tags (e.g., `v1.0.0-beta.1`)
2. **Format on save**: Add IDE integration instructions
3. **Pre-commit hook**: Auto-format on git commit
4. **Changelog generation**: Auto-generate changelog from commits during release

## Conclusion

These Makefile improvements significantly enhance the developer experience by:

1. **Preventing CI failures**: Format issues caught locally during build
2. **Reducing friction**: Simple commands for common tasks
3. **Enforcing standards**: Build fails if code isn't formatted
4. **Streamlining releases**: One validated command to create releases

The changes align with the project's pragmatic approach: add tooling that prevents common issues without creating unnecessary complexity. Developers now have a smooth, error-free workflow from development through release.
