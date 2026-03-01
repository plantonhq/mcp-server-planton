package apikey

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// create_api_key
// ---------------------------------------------------------------------------

type CreateApiKeyInput struct {
	Name         string `json:"name"                    jsonschema:"required,Display name for the API key."`
	NeverExpires bool   `json:"never_expires,omitempty"  jsonschema:"Set to true for a non-expiring key. Defaults to false."`
	ExpiresAt    string `json:"expires_at,omitempty"     jsonschema:"Expiration time in RFC 3339 format (e.g. '2026-12-31T23:59:59Z'). Ignored if never_expires is true."`
}

func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_api_key",
		Description: "Create a new API key for the authenticated user. " +
			"IMPORTANT: The raw key value is shown ONLY in the create response and cannot be retrieved afterward. " +
			"Store it securely. Set never_expires=true or provide expires_at for a time-limited key.",
	}
}

func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateApiKeyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *CreateApiKeyInput) (*mcp.CallToolResult, any, error) {
		if in.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		text, err := Create(ctx, serverAddress, in.Name, in.NeverExpires, in.ExpiresAt)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_api_keys
// ---------------------------------------------------------------------------

type ListApiKeysInput struct{}

func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_api_keys",
		Description: "List all API keys belonging to the currently authenticated user. " +
			"Returns key metadata (name, fingerprint, expiration) but NOT the raw key values.",
	}
}

func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListApiKeysInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ *ListApiKeysInput) (*mcp.CallToolResult, any, error) {
		text, err := List(ctx, serverAddress)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_api_key
// ---------------------------------------------------------------------------

type DeleteApiKeyInput struct {
	ApiKeyID string `json:"api_key_id" jsonschema:"required,ID of the API key to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_api_key",
		Description: "Permanently revoke and delete an API key. " +
			"This immediately invalidates the key. This action cannot be undone.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteApiKeyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *DeleteApiKeyInput) (*mcp.CallToolResult, any, error) {
		if in.ApiKeyID == "" {
			return nil, nil, fmt.Errorf("'api_key_id' is required")
		}
		text, err := Delete(ctx, serverAddress, in.ApiKeyID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
