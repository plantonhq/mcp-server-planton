# Project: 20260227.01.expand-cloud-resource-tools

## Overview
Expand the MCP server's cloud resource tool surface from the current 3 tools (apply, get, delete) to 16 tools covering the full lifecycle: listing/search, infrastructure destroy, stack job observability, org/env discovery, slug validation, presets, locks, rename, env var maps, and cross-resource reference resolution.

**Created**: 2026-02-27
**Status**: Active 🟢

## Project Information

### Primary Goal
Give AI agents full autonomous capability over cloud resource lifecycle — from discovering their operating context (orgs, environments) through CRUD operations to observing provisioning outcomes (stack jobs) and managing operational concerns (locks, presets, references).

### Timeline
**Target Completion**: Flexible / no hard deadline

### Technology Stack
Go/gRPC/MCP

### Project Type
Feature Development

### Affected Components
internal/domains/cloudresource/, internal/domains/stackjob/ (new), internal/domains/organization/ (new), internal/domains/environment/ (new), internal/domains/preset/ (new), internal/server/server.go

## Project Context

### Dependencies
Planton gRPC APIs (plantonhq/planton/apis), existing cloud resource domain patterns from Phase 1-5 refactor

### Success Criteria
- All 13 new tools registered and functional: list_cloud_resources
- destroy_cloud_resource
- get_stack_job_status
- list_stack_jobs
- list_organizations
- list_environments
- check_slug_availability
- search_cloud_object_presets
- list_cloud_resource_locks
- remove_cloud_resource_locks
- rename_cloud_resource
- get_env_var_map
- resolve_value_references. Unit tests for all pure domain logic. README updated.

### Known Risks & Mitigations
None beyond standard API integration risks — some query APIs may have pagination patterns that need careful handling, and streaming RPCs (streamByOrg) are not usable via standard MCP tool responses.

## Project Structure

This project follows the **Next Project Framework** for structured multi-day development:

- **`tasks/`** - Detailed task planning and execution logs (update freely)
- **`checkpoints/`** - Major milestone summaries (⚠️ ASK before creating)
- **`design-decisions/`** - Significant architectural choices (⚠️ ASK before creating)
- **`coding-guidelines/`** - Project-wide code standards (⚠️ ASK before creating)
- **`wrong-assumptions/`** - Important misconceptions (⚠️ ASK before creating)
- **`dont-dos/`** - Critical anti-patterns (⚠️ ASK before creating)

**📌 IMPORTANT**: Knowledge folders require developer permission. See [coding-guidelines/documentation-discipline.md](coding-guidelines/documentation-discipline.md)

## Current Status

### Active Task
See [tasks/](tasks/) for the current task being worked on.

### Latest Checkpoint
See [checkpoints/](checkpoints/) for the most recent project state.

### Progress Tracking
- [x] Project initialized
- [ ] Initial analysis complete
- [ ] Core implementation
- [ ] Testing and validation
- [ ] Documentation finalized
- [ ] Project completed

## How to Resume Work

**Quick Resume**: Simply drag and drop the `next-task.md` file into your AI conversation.

The `next-task.md` file contains:
- Direct paths to all project folders
- Current status information
- Resume checklist
- Quick commands

## Quick Links

- [Next Task](next-task.md) - **Drag this into chat to resume**
- [Current Task](tasks/)
- [Latest Checkpoint](checkpoints/)
- [Design Decisions](design-decisions/)
- [Coding Guidelines](coding-guidelines/)

## Documentation Discipline

**CRITICAL**: AI assistants must ASK for permission before creating:
- Checkpoints
- Design decisions
- Guidelines
- Wrong assumptions
- Don't dos

Only task logs (T##_1_feedback.md, T##_2_execution.md) can be updated without permission.

## Notes

_Add any additional notes, links, or context here as the project evolves._