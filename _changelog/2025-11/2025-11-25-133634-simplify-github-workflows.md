# Simplify GitHub Workflows

**Date**: November 25, 2025  
**Type**: Infrastructure / DevOps Cleanup  
**Impact**: Medium - Reduces noise and maintenance overhead

## Summary

Removed excessive GitHub automation and checks that were creating noise without adding value. Streamlined CI/CD to focus only on essential testing and release workflows. This change aligns with a pragmatic approach to automation - only maintain what's necessary and avoid over-engineering for a project of this scale.

## Problem Statement

The repository had accumulated multiple GitHub workflows and automation tools that were creating more noise than value:

### Before State
```
.github/
├── workflows/
│   ├── ci.yml              # 2 jobs: tests + golangci-lint (failing)
│   ├── codeql.yml          # Security scanning + dependency review
│   ├── stale.yml           # Auto-close inactive issues/PRs
│   └── release.yml         # GoReleaser + Docker
└── dependabot.yml          # Up to 20 PRs per week
```

### Pain Points

1. **golangci-lint Constant Failures**: The CI workflow included golangci-lint but had no `.golangci.yml` config file, causing continuous failures

2. **Dependabot Noise**: Configured to create up to 20 PRs per week:
   - 10 PRs for Go module updates
   - 5 PRs for GitHub Actions updates
   - 5 PRs for Docker image updates

3. **Unnecessary Security Theater**: CodeQL scanning and dependency review were running on every push/PR, which is excessive for:
   - A small team project
   - Where GitHub's MCP server (the reference implementation) shows this level of automation isn't required for success

4. **Stale Bot Overhead**: Auto-closing issues and PRs after inactivity creates busy work and can frustrate contributors

5. **Over-Engineering**: The workflows were more defensive than necessary for a project of this scale and team size

## Solution

Adopted a minimalist approach to CI/CD: keep only what's essential, remove what creates noise.

### After State
```
.github/
└── workflows/
    ├── ci.yml              # 1 job: basic tests + build
    └── release.yml         # GoReleaser + Docker (unchanged)
```

### Guiding Principles

1. **Manual Over Automatic**: Prefer manual dependency updates over automated PR spam
2. **Essential Over Comprehensive**: Only maintain checks that are truly necessary
3. **Pragmatic Over Defensive**: Don't over-engineer for theoretical problems
4. **Focus Over Noise**: Reduce distractions to focus on actual development

## Changes Made

### 1. Removed CodeQL Security Scanning
- **File Deleted**: `.github/workflows/codeql.yml`
- **Rationale**: Security scanning every push/PR is excessive for this project size
- **Impact**: No more weekly scheduled security scans or PR dependency reviews

### 2. Removed Stale Issues Bot
- **File Deleted**: `.github/workflows/stale.yml`
- **Rationale**: Auto-closing creates noise; team will manually triage
- **Impact**: No more automated stale labels or issue closures

### 3. Removed Dependabot
- **File Deleted**: `.github/dependabot.yml`
- **Rationale**: Up to 20 PRs per week is too noisy; prefer manual updates
- **Impact**: Dependencies updated manually when needed

### 4. Simplified CI Workflow
- **File Modified**: `.github/workflows/ci.yml`
- **Change**: Removed entire `golangci-lint` job
- **Kept**: Essential checks only:
  - Go module verification
  - `go vet` for correctness
  - `go fmt` for formatting
  - Tests with race detection
  - Build verification
- **Rationale**: golangci-lint was failing without config; basic checks are sufficient

### 5. Release Workflow Unchanged
- **File**: `.github/workflows/release.yml`
- **Status**: No changes - this is essential for releases
- **Includes**: GoReleaser for binaries, Docker multi-arch builds

## Benefits

1. **Less Noise**: No more automated PRs or stale bot comments
2. **Cleaner CI**: Tests pass consistently without linter config issues
3. **Lower Maintenance**: Fewer workflows to maintain and debug
4. **Faster Feedback**: CI runs faster with fewer jobs
5. **More Focus**: Team focuses on actual features, not automated overhead

## Technical Details

### CI Workflow Structure (After)
```yaml
name: CI
on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Setup Go 1.22
      - Cache Go modules
      - Download dependencies
      - Verify dependencies (go mod verify)
      - Run go vet
      - Run go fmt check
      - Run tests with race detector
      - Build binary
```

### What We Still Maintain

**CI Pipeline** - Essential quality checks:
- Dependency verification
- Static analysis (go vet)
- Code formatting (go fmt)
- Unit tests with race detection
- Build verification

**Release Pipeline** - Automated releases:
- GoReleaser for GitHub releases
- Multi-architecture Docker images
- Semantic versioning via git tags

### What We Don't Maintain Anymore

**Security Scanning**:
- CodeQL analysis removed
- Dependency review removed
- Team will handle security updates manually

**Issue Management**:
- No stale bot
- Manual triage only

**Dependency Updates**:
- No Dependabot PRs
- Manual updates when needed

## Migration Notes

No migration needed - this is purely a removal of automation. Developers should:

1. **CI Failures**: If CI was failing before due to golangci-lint, it should now pass
2. **Dependencies**: Update Go modules manually using `go get -u` when needed
3. **Security**: Check GitHub's security tab manually for alerts
4. **Issues/PRs**: Team manually closes stale items as needed

## Rollback Plan

If automation is needed again, files are in git history:
```bash
# Restore individual workflows
git checkout HEAD^ -- .github/workflows/codeql.yml
git checkout HEAD^ -- .github/workflows/stale.yml
git checkout HEAD^ -- .github/dependabot.yml

# Restore golangci-lint job in ci.yml
git diff HEAD^ -- .github/workflows/ci.yml
```

## Conclusion

This change reflects a pragmatic approach to DevOps automation: maintain only what adds clear value, remove what creates noise. For a project of this scale, essential CI testing and automated releases are sufficient. This reduces maintenance overhead and allows the team to focus on building features rather than managing automation.

The project now has a clean, minimal CI/CD setup that actually works consistently and doesn't create distractions.
