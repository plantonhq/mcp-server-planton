---
name: Private module CI auth
overview: "Configure GitHub PAT-based authentication for private Go module downloads across all three build surfaces: CI workflow, GoReleaser release workflow, and Docker image build."
todos:
  - id: ci-yml
    content: "Refine ci.yml: move GOPRIVATE to job-level env, remove it from step-level"
    status: completed
  - id: release-goreleaser
    content: "Update release.yml goreleaser job: add git config step and GOPRIVATE env"
    status: completed
  - id: release-docker
    content: "Update release.yml docker job: pass GH_PAT as BuildKit secret"
    status: completed
  - id: dockerfile
    content: "Update Dockerfile: use --mount=type=secret for go mod download"
    status: completed
isProject: false
---

# Fix Private Go Module Access in CI/CD

## Problem

`go mod download` fails in CI because two dependencies are in private GitHub repos:

- `github.com/plantonhq/planton/apis`
- `github.com/plantonhq/openmcf`

The Go toolchain cannot authenticate against GitHub, producing: `fatal: could not read Username for 'https://github.com'`.

## Approach

Use a GitHub Personal Access Token (stored as `GH_PAT` repo secret) to authenticate Git HTTPS requests, and set `GOPRIVATE` to bypass the public module proxy for all `plantonhq` modules.

There are **three separate build surfaces** that need this, each requiring a slightly different mechanism.

## Prerequisites (manual, outside of code)

- Create a GitHub PAT (classic) with `repo` scope — or a fine-grained token with read access to `plantonhq/planton` and `plantonhq/openmcf`
- Add it as a repository secret named `GH_PAT` under **Settings > Secrets and variables > Actions** in the `mcp-server-planton` repo

---

## Surface 1: CI Workflow — `[.github/workflows/ci.yml](.github/workflows/ci.yml)`

Already partially patched in an earlier edit. Refinements needed:

- Move `GOPRIVATE` from step-level to **job-level** `env` so it applies uniformly to every Go command (`go mod download`, `go vet`, `go test`, `go build`) rather than just the download step. This avoids a subtle failure if any step triggers an implicit module fetch.
- Keep the existing `Configure Git for private Go modules` step as-is.

Final shape of the relevant section:

```yaml
jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    env:
      GOPRIVATE: github.com/plantonhq/*
    steps:
      # ... checkout, setup-go ...
      - name: Configure Git for private Go modules
        run: git config --global url."https://${{ secrets.GH_PAT }}@github.com/".insteadOf "https://github.com/"
      # ... remaining steps unchanged ...
```

---

## Surface 2: Release Workflow — `[.github/workflows/release.yml](.github/workflows/release.yml)`

Two independent jobs need attention:

### 2a. `goreleaser` job

GoReleaser internally runs `go build`, which triggers module downloads. Add the same Git credential config and `GOPRIVATE` env:

```yaml
goreleaser:
  name: GoReleaser
  runs-on: ubuntu-latest
  env:
    GOPRIVATE: github.com/plantonhq/*
  steps:
    - name: Checkout
      # ...
    - name: Set up Go
      # ...
    - name: Configure Git for private Go modules
      run: git config --global url."https://${{ secrets.GH_PAT }}@github.com/".insteadOf "https://github.com/"
    - name: Install syft
      # ...
    - name: Run GoReleaser
      # ... (keep GITHUB_TOKEN as-is for release publishing)
```

### 2b. `docker` job

The Docker build calls `go mod download` inside the container. We need to pass the PAT as a **BuildKit secret** (never baked into an image layer). The `docker/build-push-action` supports this via the `secrets` input:

```yaml
- name: Build and push Docker image
  uses: docker/build-push-action@v6
  with:
    context: .
    platforms: linux/amd64,linux/arm64
    push: true
    tags: ${{ steps.meta.outputs.tags }}
    labels: ${{ steps.meta.outputs.labels }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
    secrets: |
      gh_pat=${{ secrets.GH_PAT }}
```

---

## Surface 3: Dockerfile — `[Dockerfile](Dockerfile)`

Use a BuildKit secret mount to configure Git credentials only during the `go mod download` step, ensuring the token never persists in any image layer:

```dockerfile
COPY go.mod go.sum ./
RUN --mount=type=secret,id=gh_pat \
    git config --global url."https://$(cat /run/secrets/gh_pat)@github.com/".insteadOf "https://github.com/" && \
    GOPRIVATE=github.com/plantonhq/* go mod download
```

Key points:

- The secret mount is ephemeral — it exists only for the duration of this `RUN` instruction
- The `git config` written here lives in the builder stage only and is discarded in the final runtime stage (multi-stage build)
- No token leaks into the final image

---

## What does NOT change

- `[.goreleaser.yaml](.goreleaser.yaml)` — No changes needed. GoReleaser inherits the Git config and env from the runner.
- `[go.mod](go.mod)` — No changes. Import paths stay as `github.com/plantonhq/*`.
- Local development — Developers already have Git credentials configured locally.

## Security considerations

- The PAT is never committed to the repo or baked into a Docker image layer.
- BuildKit secret mounts are the Docker-recommended way to handle build-time credentials.
- If the org later adopts GitHub App installation tokens (more secure, auto-rotating), the only change would be how `GH_PAT` is populated in the workflow — the consuming code stays identical.

