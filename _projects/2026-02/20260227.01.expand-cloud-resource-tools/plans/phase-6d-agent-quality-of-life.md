---
name: Phase 6D Implementation
overview: Implement 3 Phase 6D tools (check_slug_availability, search_cloud_object_presets, get_cloud_object_preset) expanding the MCP server from 10 to 13 tools. This plan flags 4 design decisions and 1 proto surprise that need your input before implementation begins.
todos:
  - id: decisions
    content: Get user input on 4 design decisions before writing any code
    status: completed
  - id: slug-tool
    content: Implement check_slug_availability in cloudresource/slug.go + add to tools.go
    status: completed
  - id: preset-domain
    content: Create preset/ domain package with tools.go, search.go, get.go
    status: completed
  - id: server-registration
    content: Register 3 new tools in server.go (10 -> 13)
    status: completed
  - id: resolve-kind-refactor
    content: If Decision 1 = B, extract resolveKind to shared domains package
    status: completed
  - id: verify-build
    content: Run go build, go test, check linter errors
    status: completed
isProject: false
---

# Phase 6D: Agent Quality-of-Life — Implementation Plan

## Surprise Found During Proto Analysis

**The revised plan (T01_2) claims** `CloudResourceSlugAvailabilityCheckResponse` returns "if taken, the existing resource's ID." **The actual proto** only has a single field:

```proto
message CloudResourceSlugAvailabilityCheckResponse {
  bool is_available = 1;
}
```

There is no `existing_resource_id` field. The tool can only tell the agent "available" or "not available." This is still fully sufficient for the agent's use case (validate before `apply`), and we should not try to work around this. **Impact: None on implementation — just correcting the plan's documented output.**

---

## Design Decisions Requiring Input

### Decision 1: `resolveKind` duplication (3rd copy)

The `resolveKind` function (5 lines, maps PascalCase string to `CloudResourceKind` enum) currently exists in two places:

- `cloudresource/kind.go`
- `stackjob/enum.go`

Phase 6D adds a third copy in `preset/`. The Phase 6B decision was to duplicate intentionally to avoid cross-domain coupling. Three copies of a trivial function is still defensible, but we're approaching the threshold where extracting to the shared `domains` package (already imported by every domain) would be cleaner.

**Options:**

- **A) Duplicate again** — follow Phase 6B precedent, keep domains fully decoupled, accept 3 copies
- **B) Extract to `domains` package** — single source of truth, all domains already import it, but changes the precedent

**My recommendation: B** — `resolveKind` depends only on the `cloudresourcekind` proto package (a leaf dependency), not on any domain. Moving it to `domains` doesn't create coupling between domains. It's a proto-enum utility, not domain logic. The same argument would apply if we later need `resolveProvider`.

### Decision 2: `search_cloud_object_presets` — org-context flag behavior

The `SearchCloudObjectPresetsByOrgContextInput` has two boolean flags:

- `is_include_official` — include platform-level official presets
- `is_include_organization_presets` — include org-created custom presets

**The revised plan says:** set `is_include_official = true` when `org` is provided.
**But what about `is_include_organization_presets`?** The plan is silent on this.

When an agent searches presets in an org context, it almost certainly wants to see everything relevant: both official presets AND any org-customized ones. Exposing these two flags to the agent would add cognitive load with no real value.

**My recommendation:** When `org` is provided, set both flags to `true`. Don't expose them as MCP tool inputs. The agent always gets the most useful result set. If a use case arises for fine-grained filtering, we add it later.

### Decision 3: `providers` filter on preset search

Both search RPCs accept a `repeated CloudResourceProvider providers` field (aws, gcp, azure, etc.). Exposing this would require a `resolveProviders` function (new enum resolver).

**My recommendation: Do not expose in Phase 6D.** The `kind` filter already implies a provider (e.g., `AwsEksCluster` is obviously AWS). Adding `providers` increases tool complexity without clear agent value. We can add it in Hardening if needed.

### Decision 4: Pagination on preset search

Both search RPCs accept `PageInfo`. Unlike stack jobs (which can be numerous), preset datasets are small — typically 1-5 presets per kind.

**My recommendation: Do not expose pagination.** Send the request without `PageInfo` (or with a generous default like `page_size=50`). The result sets are inherently bounded. This keeps the tool surface minimal.

---

## Implementation Plan (after decisions are confirmed)

### Tool 1: `check_slug_availability` — in existing `cloudresource/` domain

**New file:** [internal/domains/cloudresource/slug.go](internal/domains/cloudresource/slug.go)

- Domain function `CheckSlugAvailability(ctx, serverAddress, org, env, kind CloudResourceKind, slug string) (string, error)`
- Uses `CloudResourceQueryControllerClient` (same client as `get.go` — already established)
- Builds `CloudResourceSlugAvailabilityCheckRequest{Org, Env, CloudResourceKind, Slug}`
- Returns `domains.MarshalJSON(resp)` — yields `{"isAvailable": true/false}`

**Modified file:** [internal/domains/cloudresource/tools.go](internal/domains/cloudresource/tools.go)

- Add `CheckSlugAvailabilityInput` struct with 4 required fields: `org`, `env`, `kind`, `slug`
- Add `CheckSlugAvailabilityTool()` and `CheckSlugAvailabilityHandler(serverAddress)`
- Handler calls `resolveKind(input.Kind)` (already in `kind.go`), then delegates to `CheckSlugAvailability()`
- Update package doc comment (5 tools -> 6 tools)

**Complexity: Low** — thin RPC wrapper, reuses existing `resolveKind` and `CloudResourceQueryControllerClient`.

### Tool 2: `search_cloud_object_presets` — new `preset/` domain

**New file:** [internal/domains/preset/search.go](internal/domains/preset/search.go)

- Domain function `Search(ctx, serverAddress string, input SearchInput) (string, error)`
- `SearchInput` struct: `Org string`, `Kind string`, `SearchText string`
- Branching logic:
  - If `Org != ""`: call `SearchCloudObjectPresetsByOrgContext` with `IsIncludeOfficial=true`, `IsIncludeOrganizationPresets=true`
  - If `Org == ""`: call `SearchOfficialCloudObjectPresets`
- If `Kind != ""`: resolve via `resolveKind` and set `CloudResourceKind` on the request
- Returns `domains.MarshalJSON(resp)` — yields list of `ApiResourceSearchRecord`

**Import paths:**

- `infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"` (for `InfraHubSearchQueryControllerClient`)
- `cloudresourcekind` from openmcf (for kind resolution)

**Note:** The Go package name for the search stubs will be `infrahub` (from the path `search/v1/infrahub/`). Since we're creating a new `preset` domain package, there's no naming collision — but we should alias to `infrahubsearch` for clarity, consistent with how `cloudresourcesearch` was aliased in `cloudresource/list.go`.

### Tool 3: `get_cloud_object_preset` — same `preset/` domain

**New file:** [internal/domains/preset/get.go](internal/domains/preset/get.go)

- Domain function `Get(ctx, serverAddress, presetID string) (string, error)`
- Uses `CloudObjectPresetQueryControllerClient.Get(ctx, &apiresource.ApiResourceId{Value: presetID})`
- Returns `domains.MarshalJSON(resp)` — full `CloudObjectPreset` including YAML content

**Import path:**

- `cloudobjectpresetv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudobjectpreset/v1"`
- `apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"`

### Tool definitions: `preset/tools.go`

**New file:** [internal/domains/preset/tools.go](internal/domains/preset/tools.go)

- Package doc comment documenting proto service provenance (both `InfraHubSearchQueryController` and `CloudObjectPresetQueryController`)
- `SearchCloudObjectPresetsInput` — `kind` (optional), `org` (optional), `search_text` (optional)
- `SearchTool()` / `SearchHandler()` — validates nothing required, delegates to `Search()`
- `GetCloudObjectPresetInput` — `id` (required)
- `GetTool()` / `GetHandler()` — validates `id != ""`, delegates to `Get()`

### Server registration

**Modified file:** [internal/server/server.go](internal/server/server.go)

- Import `preset` package
- Add 3 new `mcp.AddTool` calls (1 for slug check, 2 for preset)
- Update tool count log: 10 -> 13
- Update tool name list

### If Decision 1 = B (extract resolveKind)

**Modified file:** [internal/domains/helpers.go](internal/domains/helpers.go) (or new `internal/domains/kind.go`)

- Move `resolveKind` to `domains` package as `ResolveKind` (exported)
- Update all 3 callers: `cloudresource/kind.go`, `stackjob/enum.go`, `preset/search.go`
- Remove duplicate in `stackjob/enum.go`

---

## File Change Summary


| Action      | File                                      | What                                  |
| ----------- | ----------------------------------------- | ------------------------------------- |
| Create      | `internal/domains/preset/tools.go`        | Tool defs + handlers for search + get |
| Create      | `internal/domains/preset/search.go`       | Search domain function with branching |
| Create      | `internal/domains/preset/get.go`          | Get domain function                   |
| Create      | `internal/domains/cloudresource/slug.go`  | CheckSlugAvailability domain function |
| Modify      | `internal/domains/cloudresource/tools.go` | Add slug check tool def + handler     |
| Modify      | `internal/server/server.go`               | Register 3 new tools                  |
| Conditional | `internal/domains/kind.go` or similar     | If Decision 1 = B                     |


## Testing Strategy

- `check_slug_availability`: No unit tests needed — thin RPC wrapper with no domain logic. `resolveKind` is already tested.
- `search_cloud_object_presets`: If Decision 1 = B and `resolveKind` moves to `domains`, existing tests still cover it. The branching logic (org vs no-org) is straightforward conditional — no enum resolution or complex mapping.
- `get_cloud_object_preset`: No unit tests — trivial ID-to-RPC proxy.

Following the Phase 6C precedent: no unit tests for thin wrappers with no domain logic.

## Verification

- `go build ./...` must pass
- `go test ./...` must pass (including existing tests)
- Zero linter errors

