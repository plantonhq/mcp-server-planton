# Project: 20260307.01.proto-contract-sync

## Overview
Migrate MCP server tools to match restructured protobuf contracts. The connect domain renamed all *Credential types to *ProviderConnection/*Connection, new resources were added across configmanager, IAM, and infrahub, and 7 entirely new domains appeared in the proto definitions.

**Created**: 2026-03-07
**Status**: Active 🟢

## Project Information

### Primary Goal
Get the MCP server build green by migrating broken connect/credential and connect/github imports, then add tool implementations for new resources (secretbackend, variablegroup, serviceaccount, iacprovisionermapping) and evaluate new domains (agentfleet, billing, copilot, search, reporting, integration, runner) for tool coverage.

### Timeline
**Target Completion**: 1-2 weeks

### Technology Stack
Go/gRPC/MCP

### Project Type
Migration

### Affected Components
internal/domains/connect/credential, internal/domains/connect/github, internal/domains/configmanager, internal/domains/iam, internal/domains/infrahub, schemas/credentials, gen/go

## Project Context

### Dependencies
Proto code generation (gen/go/) must be complete before migration. New proto types from planton repo already generated locally.

### Success Criteria
- 1. go build ./... succeeds with zero errors. 2. All 19 credential types migrated to new *ProviderConnection/*Connection contracts. 3. Tool implementations exist for secretbackend
- variablegroup
- serviceaccount
- iacprovisionermapping. 4. New domains evaluated and prioritized tools implemented. 5. Credential schemas and resources updated to match new type names.

### Known Risks & Mitigations
1. Spec field changes beyond renames could break apply/get tool handlers. 2. Security model change (plaintext secrets -> secret slug references) may require rethinking redaction logic. 3. Large surface area - 24 connection types in connect domain alone. 4. New domains may have incomplete or unstable proto contracts.

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