# Project: 20260228.01.servicehub-mcp-tools

## Overview
Add MCP tools for the ServiceHub domain — Service, Pipeline, VariablesGroup, SecretsGroup, DnsDomain, TektonPipeline, and TektonTask API resources.

**Created**: 2026-02-28
**Status**: Active 🟢

## Project Information

### Primary Goal
Implement 35 MCP tools across 7 ServiceHub bounded contexts (Service, Pipeline, VariablesGroup, SecretsGroup, DnsDomain, TektonPipeline, TektonTask), following the existing infrahub tool patterns.

### Timeline
**Target Completion**: 1 week

### Technology Stack
Go/gRPC/MCP

### Project Type
Feature Development

### Affected Components
internal/domains/servicehub/, internal/server/server.go

## Project Context

### Dependencies
Planton backend ServiceHub gRPC APIs must be available at the configured server address

### Success Criteria
- All 35 ServiceHub tools registered and functional
- consistent with existing infrahub tool patterns
- all tests passing

### Known Risks & Mitigations
ServiceHub gRPC API availability for integration testing, potential proto import issues for ServiceHub types

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