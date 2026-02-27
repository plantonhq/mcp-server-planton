# Hardening: Documentation and Tool Description Review

## Context

All 5 feature phases (6A-6E) are complete. 18 tools are implemented across 5 domain packages. The codebase compiles and tests pass. But the documentation is frozen at the 3-tool era:

- [README.md](README.md) says "Three tools cover the full lifecycle" and lists only 3 tools
- [docs/tools.md](docs/tools.md) documents only 3 tools (229 lines) -- this is the detailed reference agents rely on for parameter tables, workflows, and error handling
- [docs/development.md](docs/development.md) lists only `cloudresource/` in the project structure and test table

## Surprise Discovery

**docs/tools.md is the most important update**, yet it wasn't explicitly called out in the original Hardening plan (the plan said "README update" and "docs/development.md update"). The tools.md file is the deep reference that agents and developers use for parameter tables, agent workflows, and decision guides. It currently documents only 3 tools and needs 15 new tool sections. Added to scope during planning.

## Execution Order

Tool description review (Go code) must happen FIRST because any wording changes there affect what we write in the docs.

---

## Task 1: Tool Description Review (Go code)

Review all 18 tool descriptions in Go `tools.go` files for consistency. Issues spotted during exploration:

- **"Planton platform" vs "Planton Cloud"**: Tool descriptions say "on the Planton platform" (e.g. `get_cloud_resource`, `delete_cloud_resource`), but README/docs use "Planton Cloud". Normalize to one term.
- **Dual-path identification phrasing**: Tools that accept ID-or-coordinates describe the pattern in slightly different wording. Standardize.
- **Cross-tool references**: Some tools mention related tools (e.g. `check_slug_availability` says "Use this before apply_cloud_resource"), others that should don't. Audit and add where it helps agent workflow.

Files touched:

- `internal/domains/cloudresource/tools.go` (11 tools)
- `internal/domains/stackjob/tools.go` (3 tools)

Scope: description strings and jsonschema tag text only. No logic changes.

---

## Task 2: Update README.md

The "Tools & Resources" section expanded from 3 tools to 18 with grouped tables organized by functional area.

---

## Task 3: Update docs/tools.md

Expanded from 229 lines / 3 tools to ~470 lines / 18 tools. Added Resource Identification Pattern section, 15 new tool sections, and expanded Agent Cheat Sheet with 6 decision guides.

---

## Task 4: Update docs/development.md

Updated Project Structure and Test Files table with new domain packages.

---

## What This Plan Does NOT Include

- **H4 (get.go refactor)**: Deferred. The refactor to use `resolveResource` in get.go is a code-quality improvement, not documentation. Can be a separate task.
- **H1 (new unit tests)**: Not in scope for this pass. Existing test coverage is adequate for the current domain logic.
- **docs/configuration.md**: No changes needed (no new env vars introduced).

## Verification

After all changes: `go build ./...` and `go test ./...` both pass (only description strings and documentation files changed).
