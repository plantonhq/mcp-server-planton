# Project: 20260227.02.expand-infrahub-mcp-tools

## Overview
Add InfraChart, InfraProject, InfraPipeline, Graph, ConfigManager, Audit, StackJob commands, Deployment Component catalog tools to the MCP server, and restructure generated code under domain-scoped directories.

**Created**: 2026-02-27
**Status**: Active 🟢

## Project Information

### Primary Goal
Expand the MCP server from 18 tools (cloud resource CRUD only) to ~55+ tools covering the full InfraHub composition, pipeline observability, dependency intelligence, configuration lifecycle, audit trail, and operational control surface.

### Timeline
**Target Completion**: 2-3 weeks

### Technology Stack
Go/gRPC/MCP

### Project Type
Feature Development

### Affected Components
internal/domains/infrahub/, internal/domains/graph/, internal/domains/configmanager/, internal/domains/audit/, gen/, cmd/server

## Project Context

### Dependencies
Planton backend APIs must be accessible for all new RPCs (infrachart, infraproject, infrapipeline, graph, configmanager, audit, stackjob essentials)

### Success Criteria
- InfraChart search/get/build tools working
- InfraProject full CRUD with search
- InfraPipeline list/get/run/cancel
- Graph dependency and impact analysis tools
- ConfigManager variable and secret CRUD
- Audit version history tools
- StackJob rerun/cancel/preflight
- Deployment component and IaC module catalog
- Generated code restructured under domain-scoped directories

### Known Risks & Mitigations
Large scope requires careful phased rollout, Graph domain may need Neo4j connectivity, Generated code restructuring requires updating all import paths

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