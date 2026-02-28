# Next Task: MCP Server Gap Completion

## Quick Resume Instructions

Drop this file into your conversation to quickly resume work on this project.

## Project: 20260301.01.mcp-server-gap-completion

**Description**: Close all gaps between the MCP server (100+ tools, 22 domains, 2 MCP resources) and the full Planton Cloud API surface (~564 proto files, 30+ domains). Add missing bounded contexts: Connect (credentials), IAM, full ResourceManager CRUD, StackJob AI-native tools, CloudResource lifecycle completion, PromotionPolicy, FlowControlPolicy, and expanded MCP resources.

**Goal**: Add ~60-70 missing tools and 5+ MCP resources across 8+ new/expanded domains.

**Tech Stack**: Go/gRPC/MCP

## Task Inventory

### TIER 1 -- Critical Gaps

| Task | Description | Est. Tools | Status |
|------|-------------|------------|--------|
| T02 | Architecture Decision: Generic vs Per-Type Credential Tools | 0 (design) | COMPLETED |
| T03 | ResourceManager: Organization Full CRUD | 4 | NOT STARTED |
| T04 | ResourceManager: Environment Full CRUD | 4 | NOT STARTED |
| T05 | Connect Domain: Credential Management (depends on T02) | 25-30 | NOT STARTED |
| T06 | StackJob: AI-Native Tools (error recommendation, IaC resources) | 5 | NOT STARTED |
| T07 | CloudResource: Lifecycle Completion (purge + cleanup) | 2 | NOT STARTED |

### TIER 2 -- Important Gaps

| Task | Description | Est. Tools | Status |
|------|-------------|------------|--------|
| T08 | IAM Domain: Access Control (Team, Policy, Role, ApiKey, Identity) | 12-15 | NOT STARTED |
| T09 | InfraPipeline: Missing Trigger Variants | 2 | NOT STARTED |
| T10 | PromotionPolicy: Cross-Environment Deployment Governance | 4 | NOT STARTED |
| T11 | FlowControlPolicy: Change Approval Workflows | 3 | NOT STARTED |

### TIER 3 -- MCP Resources

| Task | Description | Est. Resources | Status |
|------|-------------|----------------|--------|
| T12 | Expand MCP Resources (api-resource-kinds, credential-types, presets, catalogs) | 5+ | NOT STARTED |

### TIER 4 -- Explore / Deferred

| Task | Description | Status |
|------|-------------|--------|
| T13 | Investigation: Runner Domain Accessibility | NOT STARTED |
| T14 | Investigation: Non-Streaming Log Retrieval | NOT STARTED |
| T15 | MCP Prompts (Exploratory) | NOT STARTED |

## Execution Order

T02 -> T03/T04/T06/T07 (parallel) -> T05 -> T08/T09/T10/T11 (parallel) -> T12 -> T13/T14/T15

## Essential Files to Review

### 1. Gap Analysis Reference
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/.cursor/plans/mcp_server_gap_analysis_1c322248.plan.md
```

### 2. Current Task Plan
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/tasks/
```

### 3. Project Documentation
- **README**: `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/README.md`

### 4. Existing Domain Patterns (reference implementations)
- **CloudResource**: `internal/domains/infrahub/cloudresource/`
- **StackJob**: `internal/domains/infrahub/stackjob/`
- **Service**: `internal/domains/servicehub/service/`
- **Organization (current)**: `internal/domains/resourcemanager/organization/`
- **Server registration**: `internal/server/server.go`

## Knowledge Folders to Check

### Design Decisions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/design-decisions/
```

### Coding Guidelines
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/coding-guidelines/
```

### Wrong Assumptions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/wrong-assumptions/
```

### Don't Dos
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/dont-dos/
```

## Planton API References (for implementation)

Key protobuf API paths in `/Users/suresh/scm/github.com/plantonhq/planton/apis/ai/planton/`:
- **ResourceManager**: `resourcemanager/organization/v1/`, `resourcemanager/environment/v1/`
- **Connect**: `connect/*credential/v1/` (20+ credential types)
- **IAM**: `iam/team/v1/`, `iam/iampolicy/v2/`, `iam/iamrole/v1/`, `iam/apikey/v1/`, `iam/identityaccount/v1/`
- **StackJob**: `infrahub/stackjob/v1/`
- **CloudResource**: `infrahub/cloudresource/v1/`
- **InfraPipeline**: `infrahub/infrapipeline/v1/`
- **PromotionPolicy**: `resourcemanager/promotionpolicy/v1/`
- **FlowControlPolicy**: `infrahub/flowcontrolpolicy/v1/` (verify path)

## Resume Checklist

When starting a new session:

1. [ ] Read the latest checkpoint (if any) from `checkpoints/`
2. [ ] Check task inventory above for current status
3. [ ] Review any design decisions in `design-decisions/`
4. [ ] Check coding guidelines in `coding-guidelines/`
5. [ ] Review wrong assumptions and dont-dos
6. [ ] Pick the next task and begin

## Current Status

**Created**: 2026-03-01
**Current Task**: T02 (COMPLETED -- Design Decision: DD-01)
**Next Task**: T03/T04/T06/T07 (parallel, independent of each other)
**Status**: Ready for implementation

---

*This file provides direct paths to all project resources for quick context loading.*
