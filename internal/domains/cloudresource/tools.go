// Package cloudresource provides the MCP tools and resource templates for the
// CloudResource domain, backed by the CloudResourceCommandController and
// CloudResourceQueryController RPCs on the Planton backend.
//
// Three tools are exposed:
//   - apply_cloud_resource: create or update (accepts opaque cloud_object map;
//     typed validation via generated parsers in gen/cloudresource/)
//   - get_cloud_resource: retrieve by ID or by (kind, org, env, slug)
//   - delete_cloud_resource: remove by ID or by (kind, org, env, slug)
package cloudresource

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantoncloud/mcp-server-planton/gen/cloudresource"
	"github.com/plantoncloud/mcp-server-planton/internal/domains"
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
