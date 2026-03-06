---
name: MCP Prompts v1
overview: Add 5 MCP prompts to the Planton MCP server — the first implementation of the third MCP primitive (alongside the existing 172 tools and 5 resources). Each prompt encodes a cross-domain workflow that an LLM couldn't easily discover from tool descriptions alone.
todos:
  - id: helpers
    content: Add PromptResult and UserMessage helpers to internal/domains/toolresult.go
    status: completed
  - id: package
    content: Create internal/domains/prompts/doc.go with package documentation
    status: completed
  - id: debug
    content: Implement debug_failed_deployment prompt (debug_deployment.go)
    status: completed
  - id: impact
    content: Implement assess_change_impact prompt (assess_impact.go)
    status: completed
  - id: explore
    content: Implement explore_infrastructure prompt (explore_infrastructure.go)
    status: completed
  - id: provision
    content: Implement provision_cloud_resource prompt (provision_resource.go)
    status: completed
  - id: access
    content: Implement manage_access prompt (manage_access.go)
    status: completed
  - id: register
    content: Create register.go and wire prompts into server.go
    status: completed
  - id: verify
    content: Build and verify compilation, check lints
    status: completed
isProject: false
---

# T15: MCP Prompts — Cross-Domain Workflow Templates

## What Are MCP Prompts?

MCP prompts are the third primitive in the Model Context Protocol (alongside tools and resources). They are **pre-built conversation starters** — when a user selects a prompt in their MCP client (Cursor, Claude Desktop, etc.), the server returns a list of messages that get injected into the conversation. The LLM then processes these messages and uses the available tools to fulfill the intent.

Prompts add value when they encode **non-obvious multi-step workflows** that cross domain boundaries. A user could always type "help me debug my deployment," but the prompt encodes the *recommended diagnostic sequence* — which tools to use, in what order, and why — that an LLM wouldn't discover from tool descriptions alone.

## Architecture

### Placement: Cross-Cutting Domain

Prompts span multiple bounded contexts (a debug workflow touches stackjob, cloudresource, and graph domains). They belong in a cross-cutting package, following the `discovery` precedent:

```
internal/domains/prompts/
    doc.go                      -- Package documentation
    register.go                 -- Register(srv *mcp.Server)
    debug_deployment.go         -- debug_failed_deployment prompt
    assess_impact.go            -- assess_change_impact prompt
    explore_infrastructure.go   -- explore_infrastructure prompt
    provision_resource.go       -- provision_cloud_resource prompt
    manage_access.go            -- manage_access prompt
```

### Per-File Pattern (mirrors tool pattern)

Each prompt file contains three things:

- **Prompt definition**: `func XxxPrompt() *mcp.Prompt` — name, description, arguments
- **Handler**: `func XxxHandler() mcp.PromptHandler` — builds and returns messages
- **Text builder**: `func buildXxxText(args...) string` — pure function, independently testable

### Shared Helpers

Add two helpers to [internal/domains/toolresult.go](internal/domains/toolresult.go) (alongside existing `TextResult` and `ResourceResult`):

```go
func PromptResult(description string, messages ...*mcp.PromptMessage) *mcp.GetPromptResult
func UserMessage(text string) *mcp.PromptMessage
```

### Server Registration

Add `registerPrompts(srv)` to [internal/server/server.go](internal/server/server.go):

```go
func New(cfg *config.Config) *Server {
    srv := mcp.NewServer(...)
    registerTools(srv, cfg.ServerAddress)
    registerResources(srv)
    registerPrompts(srv)        // NEW
    return &Server{...}
}

func registerPrompts(srv *mcp.Server) {
    prompts.Register(srv)
    slog.Info("prompts registered")
}
```

### Content Style: Hybrid

Each prompt states the **goal** clearly, then provides a **recommended tool sequence** as guidance (not a rigid script). The LLM adapts to the user's specific situation. This balances discoverability of non-obvious tools with the LLM's ability to reason about what applies.

### No gRPC Calls

Prompt handlers are **purely static** — they interpolate string arguments into pre-written text. No gRPC calls, no dynamic data, no failure modes. This keeps them fast, testable, and reliable.

---

## The 5 Prompts

### 1. `debug_failed_deployment`

**Why it matters**: Most common pain point. Requires 5+ tools across 2 domains (stackjob + cloudresource + graph) in a non-obvious sequence. Most users don't know `get_error_resolution_recommendation` exists.

- **Arguments**: `resource_id` (optional), `stack_job_id` (optional)
- **Workflow encoded**:
  - Find the failed stack job (`get_latest_stack_job` or `get_stack_job`)
  - Get error details, then call `get_error_resolution_recommendation` for AI-suggested fix
  - Review IaC resource state (`find_iac_resources_by_stack_job`)
  - Inspect the (credential-free) IaC input (`get_stack_job_input`) for config errors
  - Check upstream dependencies (`get_dependencies`) for cascading failures
  - Summarize findings and recommended next steps

### 2. `assess_change_impact`

**Why it matters**: Safety-critical workflow before destructive operations. The platform has rich impact analysis tools (graph domain) that most users don't discover. This prompt is an architectural guardrail.

- **Arguments**: `resource_id` (required), `change_type` (optional, default: "delete")
- **Workflow encoded**:
  - Run `get_impact_analysis` to quantify blast radius (directly + transitively affected resources)
  - Use `get_dependents` to enumerate downstream resources
  - Use `get_dependencies` to understand upstream context
  - Check `list_resource_versions` for recent change history
  - Present risk assessment with concrete numbers and recommendation

### 3. `explore_infrastructure`

**Why it matters**: Discovery/onboarding workflow. With 172 tools, new users don't know where to start. The graph tools are powerful but their relationships are non-obvious. This prompt provides a top-down exploration path.

- **Arguments**: `org_id` (optional)
- **Workflow encoded**:
  - Start from `list_organizations` if no org specified
  - Use `get_organization_graph` for the full topology
  - Summarize: environments, cloud resources, services, credentials
  - Highlight resources with recent failures
  - Reference `api-resource-kinds://catalog` for platform context

### 4. `provision_cloud_resource`

**Why it matters**: Longest multi-step workflow in the platform, crossing 4+ domains (connect, cloudresource, stackjob, discovery). Ensures prerequisites are checked before the user hits an error 5 steps in.

- **Arguments**: `kind` (optional), `org_id` (optional), `env_id` (optional)
- **Workflow encoded**:
  - Discover available types via `cloud-resource-kinds://catalog` and `cloud-resource-schema://{kind}`
  - Verify credentials exist (`search_credentials`)
  - Confirm provider connection (`resolve_default_provider_connection`)
  - Check for pre-approved configurations (`search_cloud_object_presets`)
  - Create the resource (`apply_cloud_resource`)
  - Monitor deployment (`get_latest_stack_job`)
  - Handle failure with `get_error_resolution_recommendation`

### 5. `manage_access`

**Why it matters**: IAM is consistently the most confusing domain for platform users. The Planton IAM model (principals, policies, roles, teams, API keys) has a learning curve. This prompt guides users through discovery, audit, and action.

- **Arguments**: `org_id` (optional), `resource_id` (optional)
- **Workflow encoded**:
  - Discover who has access (`list_principals`)
  - Check what a specific principal can access (`list_resource_access`)
  - Verify specific permissions (`check_authorization`)
  - Understand available roles (`list_iam_roles_for_resource_kind`)
  - Grant or revoke access (`create_iam_policy` / `delete_iam_policy`)
  - Warn about nuclear options (`revoke_org_access`)

---

## Key Design Properties

- **No domain logic in prompts**: Prompts contain only text — they reference tool names but don't call tools themselves. They're a UX layer, not a business logic layer.
- **Prompt text references tool names as stable identifiers**: Tool names are part of the MCP server's public API surface and are unlikely to change. If they do, prompts are updated alongside.
- **Each prompt is self-contained**: No prompt depends on another prompt. Each encodes a complete workflow.
- **Arguments are always optional where possible**: The prompt text handles "if not provided" cases, letting the LLM discover the needed information via tools.

## Files to Create (7)

- `internal/domains/prompts/doc.go`
- `internal/domains/prompts/register.go`
- `internal/domains/prompts/debug_deployment.go`
- `internal/domains/prompts/assess_impact.go`
- `internal/domains/prompts/explore_infrastructure.go`
- `internal/domains/prompts/provision_resource.go`
- `internal/domains/prompts/manage_access.go`

## Files to Modify (2)

- `internal/domains/toolresult.go` — add `PromptResult` and `UserMessage` helpers
- `internal/server/server.go` — add `registerPrompts(srv)` and import

