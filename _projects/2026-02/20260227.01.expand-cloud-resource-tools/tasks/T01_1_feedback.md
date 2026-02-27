# T01 Feedback — Developer Review Responses

**Date**: 2026-02-27
**Reviewer**: Developer

## Responses to Open Questions

### Q1: `find` API auth level

**Decision**: Do NOT use `find` — it is a platform-operator RPC meant for re-indexing the search index. Do NOT use `streamByOrg` either — it is a raw database stream not intended for user-facing listing.

Instead, use `CloudResourceSearchQueryController.getCloudResourcesCanvasView` from the search domain (`ai.planton.search.v1.infrahub.cloudresource`). This is the correct user-facing API for discovering cloud resources — it queries the search index with server-side filtering (envs, kinds, text search) and requires only `get` permission on the organization.

### Q2: Destroy confirmation

**Decision**: No confirmation mechanism at the MCP level. The agent (LLM) is responsible for confirming destructive intent with the user before calling the tool. MCP protocol does not have a built-in confirmation mechanism — tools are fire-and-forget from the protocol's perspective.

The tool description should include a clear warning about the destructive nature.

### Q3: Stack job response size

**Decision**: Return the full stack job. Stack jobs do not contain secrets (sensitive data is behind the operator-only `getCloudResourceStackExecuteInput` RPC). Full response is appropriate.

### Q4: Preset YAML content

**Decision**: Split into search + get pattern:
- `search_cloud_object_presets` → returns lightweight search records (metadata only)
- `get_cloud_object_preset` → returns the full `CloudObjectPreset` including YAML content

This pattern can be applied to other resources in the future where search returns metadata and get returns full content.

This adds 1 additional tool, bringing the total to **14 new tools (17 total)**.

### Q5: Phase ordering

**Decision**: Confirmed as proposed: 6A → 6B → 6C → 6D → 6E → Hardening.
