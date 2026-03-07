package connection

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_connection
// ---------------------------------------------------------------------------

type ApplyConnectionInput struct {
	ConnectionObject map[string]any `json:"connection_object" jsonschema:"required,Full connection object in OpenMCF envelope format: { api_version, kind, metadata: { name, org }, spec: { ... } }. Read connection-types://catalog to discover available kinds, then connection-schema://{kind} to learn the spec fields."`
}

func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_connection",
		Description: "Create or update a connection of any supported type (AWS, GCP, Azure, GitHub, Kubernetes, etc.). " +
			"Pass the full connection object as an OpenMCF envelope with api_version, kind, metadata, and spec. " +
			"Workflow: read connection-types://catalog to discover types, then connection-schema://{kind} for the spec schema, then call this tool. " +
			"The kind field in the connection_object determines which provider-specific backend is called.",
	}
}

func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyConnectionInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyConnectionInput) (*mcp.CallToolResult, any, error) {
		if input.ConnectionObject == nil {
			return nil, nil, fmt.Errorf("'connection_object' is required")
		}
		text, err := Apply(ctx, serverAddress, input.ConnectionObject)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_connection
// ---------------------------------------------------------------------------

type GetConnectionInput struct {
	Kind string `json:"kind" jsonschema:"required,PascalCase connection type (e.g. AwsProviderConnection, GcpProviderConnection). Read connection-types://catalog to discover available kinds."`
	ID   string `json:"id,omitempty" jsonschema:"Connection ID. Provide either id or (org + slug) to identify the connection."`
	Org  string `json:"org,omitempty" jsonschema:"Organization ID. Used with slug for lookup by org+slug."`
	Slug string `json:"slug,omitempty" jsonschema:"Connection slug within the organization. Used with org for lookup by org+slug."`
}

func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_connection",
		Description: "Retrieve a connection by ID or by organization+slug. " +
			"Requires kind to dispatch to the correct provider backend. " +
			"Provide either 'id' alone, or both 'org' and 'slug'.",
	}
}

func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetConnectionInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetConnectionInput) (*mcp.CallToolResult, any, error) {
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}

		var text string
		var err error

		switch {
		case input.ID != "":
			text, err = Get(ctx, serverAddress, input.Kind, input.ID)
		case input.Org != "" && input.Slug != "":
			text, err = GetByOrgBySlug(ctx, serverAddress, input.Kind, input.Org, input.Slug)
		default:
			return nil, nil, fmt.Errorf("provide either 'id' or both 'org' and 'slug'")
		}

		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_connection
// ---------------------------------------------------------------------------

type DeleteConnectionInput struct {
	Kind string `json:"kind" jsonschema:"required,PascalCase connection type (e.g. AwsProviderConnection, GcpProviderConnection)."`
	ID   string `json:"id" jsonschema:"required,Connection ID to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_connection",
		Description: "Delete a connection by ID. " +
			"Requires kind to dispatch to the correct provider backend. " +
			"WARNING: This permanently removes the connection. Cloud resources using this connection will lose connectivity.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteConnectionInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteConnectionInput) (*mcp.CallToolResult, any, error) {
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.Kind, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// search_connections
// ---------------------------------------------------------------------------

type SearchConnectionsInput struct {
	Org        string   `json:"org" jsonschema:"required,Organization ID to search within."`
	Env        string   `json:"env,omitempty" jsonschema:"Optional environment slug to narrow results."`
	Kinds      []string `json:"kinds,omitempty" jsonschema:"Optional list of PascalCase connection kinds to filter (e.g. ['AwsProviderConnection', 'GcpProviderConnection'])."`
	SearchText string   `json:"search_text,omitempty" jsonschema:"Optional free-text filter on connection name."`
}

func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_connections",
		Description: "Search connections within an organization. " +
			"Returns lightweight search records (id, name, kind) without full spec details. " +
			"Optionally filter by environment, connection kinds, or free-text search. " +
			"Use get_connection to retrieve the full connection details.",
	}
}

func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchConnectionsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchConnectionsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := Search(ctx, serverAddress, input.Org, input.Env, input.Kinds, input.SearchText)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// check_connection_slug
// ---------------------------------------------------------------------------

type CheckConnectionSlugInput struct {
	Org  string `json:"org" jsonschema:"required,Organization ID to check within."`
	Kind string `json:"kind" jsonschema:"required,PascalCase connection type (e.g. AwsProviderConnection)."`
	Slug string `json:"slug" jsonschema:"required,Slug to check for availability."`
}

func CheckSlugTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_connection_slug",
		Description: "Check if a connection slug is available within an organization for a given connection kind. " +
			"Returns { is_available: true/false }. " +
			"Use this before apply_connection to verify the slug is unique.",
	}
}

func CheckSlugHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckConnectionSlugInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CheckConnectionSlugInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := CheckSlugAvailability(ctx, serverAddress, input.Org, input.Kind, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
