# Project: 20260226.01.refactor-mcp-server-stigmer-patterns

## Overview
Complete refactoring of mcp-server-planton to follow the Stigmer MCP server architecture: domain-driven tool organization, two-stage codegen pipeline (proto → schema → Go input types), and consistent apply/delete/get tool patterns.

**Created**: 2026-02-26
**Status**: Active 🟢

## Project Information

### Primary Goal
Replace the entire current tool implementation with a clean, codegen-driven architecture modeled after stigmer/mcp-server. Start with infrahub/cloudresource apply and delete tools, then expand to other domains.

### Timeline
**Target Completion**: 2 weeks

### Technology Stack
Go/gRPC/MCP/Protobuf/CodeGen

### Project Type
Refactoring

### Affected Components
cmd/mcp-server-planton, internal/domains, internal/mcp, tools/codegen (new), gen (new)

## Project Context

### Dependencies
stigmer/mcp-server as reference architecture, stigmer/tools/codegen as codegen reference, planton/apis protobuf definitions, openmcf protobuf definitions

### Success Criteria
- 1) All current tools removed and replaced with codegen-driven tools 2) apply_cloud_resource and delete_cloud_resource tools working against Planton gRPC APIs 3) Codegen pipeline produces correct input types from proto definitions 4) Architecture matches Stigmer MCP server patterns

### Known Risks & Mitigations
1) CloudResource cloud_object (google.protobuf.Struct) wrapping adds complexity vs Stigmer flat resources 2) OpenMCF provider-specific schemas need special handling in codegen 3) Existing users may rely on current tool names/signatures

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