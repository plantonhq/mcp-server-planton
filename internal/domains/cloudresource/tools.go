// Package cloudresource provides the MCP tools and resource templates for the
// CloudResource domain, backed by the CloudResourceCommandController,
// CloudResourceQueryController, and CloudResourceSearchQueryController RPCs
// on the Planton backend.
//
// Eleven tools are exposed:
//   - apply_cloud_resource: create or update (accepts opaque cloud_object map;
//     typed validation via generated parsers in gen/cloudresource/)
//   - get_cloud_resource: retrieve by ID or by (kind, org, env, slug)
//   - delete_cloud_resource: remove by ID or by (kind, org, env, slug)
//   - list_cloud_resources: query the search index for resources in an org
//   - destroy_cloud_resource: tear down cloud infrastructure (keeps record)
//   - check_slug_availability: verify slug uniqueness within (org, env, kind)
//   - list_cloud_resource_locks: show lock status, holder, and queue
//   - remove_cloud_resource_locks: force-remove all locks on a resource
//   - rename_cloud_resource: change a resource's display name
//   - get_env_var_map: extract env vars and secrets from a manifest YAML
//   - resolve_value_references: resolve all valueFrom refs in a resource
package cloudresource

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/gen/cloudresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_cloud_resource
// ---------------------------------------------------------------------------

// ApplyCloudResourceInput defines the parameters for the apply_cloud_resource tool.
type ApplyCloudResourceInput struct {
	CloudObject map[string]any `json:"cloud_object" jsonschema:"required,The full OpenMCF cloud resource object. Must contain api_version, kind, metadata (with name, org, env), and spec. Read cloud-resource-kinds://catalog for available kinds, then cloud-resource-schema://{kind} for the spec format."`
}

// ApplyTool returns the MCP tool definition for apply_cloud_resource.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_cloud_resource",
		Description: "Create or update a cloud resource on the Planton platform (idempotent). " +
			"The cloud_object must follow the OpenMCF format: " +
			"{ api_version, kind, metadata: { name, org, env }, spec: { ... } }. " +
			"Step 1: Read cloud-resource-kinds://catalog to discover supported kinds and api_versions. " +
			"Step 2: Read cloud-resource-schema://{kind} to get the full spec definition. " +
			"Step 3: Call this tool with the assembled cloud_object.",
	}
}

// ApplyHandler returns the typed tool handler for apply_cloud_resource.
// serverAddress is captured at registration time; the API key is read from
// context at call time.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyCloudResourceInput) (*mcp.CallToolResult, any, error) {
		co := input.CloudObject
		if co == nil {
			return nil, nil, fmt.Errorf("cloud_object is required")
		}

		kindStr, err := extractKindFromCloudObject(co)
		if err != nil {
			return nil, nil, err
		}

		parseFn, ok := cloudresource.GetParser(kindStr)
		if !ok {
			return nil, nil, fmt.Errorf(
				"unsupported cloud resource kind %q — read cloud-resource-kinds://catalog for all valid kinds",
				kindStr,
			)
		}

		normalizedObject, err := parseFn(co)
		if err != nil {
			return nil, nil, fmt.Errorf("cloud_object validation failed: %w", err)
		}

		cr, err := buildCloudResource(co, kindStr, normalizedObject)
		if err != nil {
			return nil, nil, err
		}

		text, err := Apply(ctx, serverAddress, cr)
		if err != nil {
			return nil, nil, err
		}

		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_cloud_resource
// ---------------------------------------------------------------------------

// GetCloudResourceInput defines the parameters for the get_cloud_resource tool.
// Exactly one identification path must be provided: either id alone, or all of
// kind + org + env + slug.
type GetCloudResourceInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}

// GetTool returns the MCP tool definition for get_cloud_resource.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_cloud_resource",
		Description: "Get a cloud resource from the Planton platform. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together. " +
			"Returns the full resource including metadata, spec, and status.",
	}
}

// GetHandler returns the typed tool handler for get_cloud_resource.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetCloudResourceInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}

		text, err := Get(ctx, serverAddress, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_cloud_resource
// ---------------------------------------------------------------------------

// DeleteCloudResourceInput defines the parameters for the delete_cloud_resource
// tool. Exactly one identification path must be provided: either id alone, or
// all of kind + org + env + slug.
type DeleteCloudResourceInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}

// DeleteTool returns the MCP tool definition for delete_cloud_resource.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_cloud_resource",
		Description: "Delete a cloud resource from the Planton platform. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together. " +
			"Returns the deleted resource.",
	}
}

// DeleteHandler returns the typed tool handler for delete_cloud_resource.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteCloudResourceInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}

		text, err := Delete(ctx, serverAddress, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_cloud_resources
// ---------------------------------------------------------------------------

// ListCloudResourcesInput defines the parameters for the list_cloud_resources
// tool. The org field is required; all other fields are optional filters.
type ListCloudResourcesInput struct {
	Org        string   `json:"org"                   jsonschema:"required,Organization identifier."`
	Envs       []string `json:"envs,omitempty"        jsonschema:"Environment slugs to filter by."`
	SearchText string   `json:"search_text,omitempty" jsonschema:"Free-text search query."`
	Kinds      []string `json:"kinds,omitempty"       jsonschema:"PascalCase cloud resource kinds to filter by (e.g. AwsVpc). Read cloud-resource-kinds://catalog for valid kinds."`
}

// ListTool returns the MCP tool definition for list_cloud_resources.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_cloud_resources",
		Description: "List cloud resources in an organization from the Planton platform. " +
			"Returns resources grouped by environment and kind. " +
			"Optionally filter by environment slugs, resource kinds, or free-text search. " +
			"Read cloud-resource-kinds://catalog for valid kind values.",
	}
}

// ListHandler returns the typed tool handler for list_cloud_resources.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListCloudResourcesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListCloudResourcesInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}

		kinds, err := resolveKinds(input.Kinds)
		if err != nil {
			return nil, nil, err
		}

		text, err := List(ctx, serverAddress, input.Org, input.Envs, input.SearchText, kinds)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// destroy_cloud_resource
// ---------------------------------------------------------------------------

// DestroyCloudResourceInput defines the parameters for the
// destroy_cloud_resource tool. Exactly one identification path must be
// provided: either id alone, or all of kind + org + env + slug.
type DestroyCloudResourceInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}

// DestroyTool returns the MCP tool definition for destroy_cloud_resource.
func DestroyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "destroy_cloud_resource",
		Description: "Destroy the cloud infrastructure (Terraform/Pulumi destroy) for a resource " +
			"while keeping the resource record on the Planton platform. " +
			"This tears down the actual cloud resources (VPCs, clusters, databases, etc.). " +
			"Use delete_cloud_resource to remove the record itself. " +
			"WARNING: This is a destructive operation that will destroy real cloud infrastructure. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together.",
	}
}

// DestroyHandler returns the typed tool handler for destroy_cloud_resource.
func DestroyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DestroyCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DestroyCloudResourceInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}

		text, err := Destroy(ctx, serverAddress, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// check_slug_availability
// ---------------------------------------------------------------------------

// CheckSlugAvailabilityInput defines the parameters for the
// check_slug_availability tool. All four fields are required — slugs are
// unique within the composite key (org, env, kind).
type CheckSlugAvailabilityInput struct {
	Org  string `json:"org"  jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Env  string `json:"env"  jsonschema:"required,Environment identifier. Use list_environments to discover available environments."`
	Kind string `json:"kind" jsonschema:"required,PascalCase cloud resource kind (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds."`
	Slug string `json:"slug" jsonschema:"required,The slug to check for availability."`
}

// CheckSlugAvailabilityTool returns the MCP tool definition for check_slug_availability.
func CheckSlugAvailabilityTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_slug_availability",
		Description: "Check whether a cloud resource slug is available within the scoped " +
			"composite key (org, env, kind). Slugs must be unique within this scope. " +
			"Use this before apply_cloud_resource to verify that the desired slug is not already taken.",
	}
}

// CheckSlugAvailabilityHandler returns the typed tool handler for check_slug_availability.
func CheckSlugAvailabilityHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckSlugAvailabilityInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CheckSlugAvailabilityInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Env == "" {
			return nil, nil, fmt.Errorf("'env' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}

		kind, err := domains.ResolveKind(input.Kind)
		if err != nil {
			return nil, nil, err
		}

		text, err := CheckSlugAvailability(ctx, serverAddress, input.Org, input.Env, kind, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_cloud_resource_locks
// ---------------------------------------------------------------------------

// ListCloudResourceLocksInput defines the parameters for the
// list_cloud_resource_locks tool. Exactly one identification path must be
// provided: either id alone, or all of kind + org + env + slug.
type ListCloudResourceLocksInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}

// ListLocksTool returns the MCP tool definition for list_cloud_resource_locks.
func ListLocksTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_cloud_resource_locks",
		Description: "List lock information for a cloud resource on the Planton platform. " +
			"Returns whether the resource is locked, current lock holder details " +
			"(workflow ID, acquired timestamp, TTL remaining), and any workflows " +
			"waiting in the lock queue. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together.",
	}
}

// ListLocksHandler returns the typed tool handler for list_cloud_resource_locks.
func ListLocksHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListCloudResourceLocksInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListCloudResourceLocksInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}

		text, err := ListLocks(ctx, serverAddress, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// remove_cloud_resource_locks
// ---------------------------------------------------------------------------

// RemoveCloudResourceLocksInput defines the parameters for the
// remove_cloud_resource_locks tool. Exactly one identification path must be
// provided: either id alone, or all of kind + org + env + slug.
type RemoveCloudResourceLocksInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}

// RemoveLocksTool returns the MCP tool definition for remove_cloud_resource_locks.
func RemoveLocksTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "remove_cloud_resource_locks",
		Description: "Remove all locks (active lock and wait queue) for a cloud resource " +
			"on the Planton platform. Returns details about what was removed " +
			"(active lock removed, queue entries cleared). " +
			"WARNING: Removing locks on a resource with an active stack job may cause " +
			"IaC state corruption. Verify no stack jobs are in progress before using this tool. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together.",
	}
}

// RemoveLocksHandler returns the typed tool handler for remove_cloud_resource_locks.
func RemoveLocksHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RemoveCloudResourceLocksInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RemoveCloudResourceLocksInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}

		text, err := RemoveLocks(ctx, serverAddress, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// rename_cloud_resource
// ---------------------------------------------------------------------------

// RenameCloudResourceInput defines the parameters for the
// rename_cloud_resource tool. Exactly one identification path must be
// provided: either id alone, or all of kind + org + env + slug.
type RenameCloudResourceInput struct {
	ID      string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
	Kind    string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
	Org     string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
	Env     string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
	Slug    string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
	NewName string `json:"new_name"       jsonschema:"required,The new display name for the cloud resource. The slug remains unchanged."`
}

// RenameTool returns the MCP tool definition for rename_cloud_resource.
func RenameTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "rename_cloud_resource",
		Description: "Rename a cloud resource on the Planton platform. " +
			"Changes the human-readable display name; the immutable slug is unaffected. " +
			"Returns the updated resource. " +
			"Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together.",
	}
}

// RenameHandler returns the typed tool handler for rename_cloud_resource.
func RenameHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RenameCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RenameCloudResourceInput) (*mcp.CallToolResult, any, error) {
		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}
		if err := validateIdentifier(id); err != nil {
			return nil, nil, err
		}
		if input.NewName == "" {
			return nil, nil, fmt.Errorf("'new_name' is required")
		}

		text, err := Rename(ctx, serverAddress, id, input.NewName)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_env_var_map
// ---------------------------------------------------------------------------

// GetEnvVarMapInput defines the parameters for the get_env_var_map tool.
type GetEnvVarMapInput struct {
	YAMLContent string `json:"yaml_content" jsonschema:"required,Raw YAML content of the cloud resource manifest (OpenMCF format with api_version, kind, metadata, and spec). The server parses this to identify the resource and resolve valueFrom references."`
}

// GetEnvVarMapTool returns the MCP tool definition for get_env_var_map.
func GetEnvVarMapTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_env_var_map",
		Description: "Extract the environment variable map from a cloud resource manifest. " +
			"Provide the raw YAML content of a cloud resource manifest (OpenMCF format). " +
			"The server parses the YAML, identifies the resource kind, extracts environment " +
			"variables and secrets, and resolves valueFrom references to plain string values. " +
			"Returns separate maps for variables and secrets.",
	}
}

// GetEnvVarMapHandler returns the typed tool handler for get_env_var_map.
func GetEnvVarMapHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetEnvVarMapInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetEnvVarMapInput) (*mcp.CallToolResult, any, error) {
		if input.YAMLContent == "" {
			return nil, nil, fmt.Errorf("'yaml_content' is required")
		}

		text, err := GetEnvVarMap(ctx, serverAddress, input.YAMLContent)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// resolve_value_references
// ---------------------------------------------------------------------------

// ResolveValueReferencesInput defines the parameters for the
// resolve_value_references tool. The kind field is always required because the
// RPC needs it for both authorization and resource transformation. The resource
// is identified by id alone, or by all of org + env + slug (kind is already
// provided separately).
type ResolveValueReferencesInput struct {
	Kind string `json:"kind" jsonschema:"required,PascalCase cloud resource kind (e.g. AwsEksCluster). Always required for reference resolution. Read cloud-resource-kinds://catalog for valid kinds."`
	ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of org, env, and slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with env, slug when id is not provided."`
	Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with org, slug when id is not provided."`
	Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with org, env when id is not provided."`
}

// ResolveValueReferencesTool returns the MCP tool definition for resolve_value_references.
func ResolveValueReferencesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_value_references",
		Description: "Resolve all valueFrom references in a cloud resource's spec. " +
			"The server loads the resource, finds all valueFrom references in its " +
			"specification, resolves them to concrete values, and returns the fully " +
			"transformed cloud resource as YAML. The response includes resolution " +
			"status, any errors, and diagnostics. " +
			"Identify the resource by 'id' or by all of 'org', 'env', and 'slug'. " +
			"The 'kind' field is always required.",
	}
}

// ResolveValueReferencesHandler returns the typed tool handler for resolve_value_references.
func ResolveValueReferencesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveValueReferencesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveValueReferencesInput) (*mcp.CallToolResult, any, error) {
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		kind, err := domains.ResolveKind(input.Kind)
		if err != nil {
			return nil, nil, err
		}

		// Custom validation: kind is always required (already checked above),
		// and the resource is identified by id alone or by org+env+slug.
		// This differs from validateIdentifier which treats kind as part of the
		// slug path — here kind is a separate, always-required field.
		hasID := input.ID != ""
		slugFields := [3]string{input.Org, input.Env, input.Slug}
		slugCount := 0
		for _, f := range slugFields {
			if f != "" {
				slugCount++
			}
		}

		switch {
		case hasID && slugCount > 0:
			return nil, nil, fmt.Errorf("provide either 'id' alone or all of 'org', 'env', and 'slug' — not both")
		case hasID:
			// ID path — valid
		case slugCount == 3:
			// Slug path — valid
		case slugCount > 0:
			return nil, nil, fmt.Errorf("when not using 'id', all of 'org', 'env', and 'slug' are required")
		default:
			return nil, nil, fmt.Errorf("provide 'id' or all of 'org', 'env', and 'slug' to identify the cloud resource")
		}

		id := ResourceIdentifier{
			ID:   input.ID,
			Kind: input.Kind,
			Org:  input.Org,
			Env:  input.Env,
			Slug: input.Slug,
		}

		text, err := ResolveValueReferences(ctx, serverAddress, kind, id)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
