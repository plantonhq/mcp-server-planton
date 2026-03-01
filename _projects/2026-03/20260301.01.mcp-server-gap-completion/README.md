# Project: 20260301.01.mcp-server-gap-completion

## Overview
Close all gaps between the MCP server (100+ tools) and the Planton Cloud API surface (~564 proto files). Add missing bounded contexts: Connect (credentials), IAM, full ResourceManager CRUD, StackJob AI-native tools, CloudResource lifecycle completion, PromotionPolicy, FlowControlPolicy, and expanded MCP resources.

**Created**: 2026-03-01
**Status**: Active 🟢

## Project Information

### Primary Goal
Add ~60-70 missing tools and 5+ MCP resources across 8+ new/expanded domains to achieve comprehensive coverage of the Planton Cloud platform API surface.

### Timeline
**Target Completion**: 2-3 weeks

### Technology Stack
Go/gRPC/MCP

### Project Type
Feature Development

### Affected Components
internal/domains/connect/*, internal/domains/resourcemanager/*, internal/domains/iam/*, internal/domains/infrahub/stackjob/, internal/domains/infrahub/cloudresource/, internal/domains/infrahub/infrapipeline/, internal/server/server.go, MCP resources

## Project Context

### Dependencies
Planton gRPC API server must expose all referenced RPC endpoints for Connect, IAM, and policy domains

### Success Criteria
- All Tier 1 gaps closed (Connect credentials + ResourceManager CRUD + StackJob AI tools + CloudResource lifecycle)
- All Tier 2 gaps closed (IAM + InfraPipeline triggers + Promotion/FlowControl policies)
- MCP resources expanded from 2 to 7+
- All new tools follow established domain patterns
- Comprehensive unit tests for pure domain logic

### Known Risks & Mitigations
Connect domain has 20+ credential types requiring architecture decision (generic vs per-type tools), Runner domain accessibility via same gRPC server is unknown, Non-streaming log retrieval endpoints may not exist server-side requiring backend work

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