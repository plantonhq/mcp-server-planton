// Package secretversion provides the MCP tools for the SecretVersion domain,
// backed by the SecretVersionCommandController and SecretVersionQueryController
// RPCs (ai.planton.configmanager.secretversion.v1) on the Planton backend.
//
// Two tools are exposed:
//   - create_secret_version:  store a new set of encrypted key-value pairs
//   - list_secret_versions:   list version metadata for a secret (no data)
//
// Reading decrypted secret data is intentionally excluded from the MCP tool
// surface. This is a security boundary: agents can write secret values but
// cannot read them back. See design decision DD-2 (Phase 2B plan).
package secretversion

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// create_secret_version
// ---------------------------------------------------------------------------

// CreateSecretVersionInput defines the parameters for the
// create_secret_version tool.
type CreateSecretVersionInput struct {
	SecretID string            `json:"secret_id" jsonschema:"required,The parent secret's ID. Use get_secret or list_secrets to find the secret ID."`
	Data     map[string]string `json:"data"      jsonschema:"required,Key-value pairs to store as the secret version. Values are encrypted via envelope encryption and stored in the secret's backend. Example: {\"DB_PASSWORD\": \"s3cret\", \"API_KEY\": \"abc123\"}."`
}

// CreateTool returns the MCP tool definition for create_secret_version.
func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_secret_version",
		Description: "Store a new version of encrypted key-value data for a secret. " +
			"Each call creates an immutable version — previous versions are preserved. " +
			"The data is encrypted via envelope encryption and stored in the secret's backend. " +
			"Use apply_secret first to create the parent secret if it does not exist. " +
			"Use list_secret_versions to see existing versions before creating a new one.",
	}
}

// CreateHandler returns the typed tool handler for create_secret_version.
func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateSecretVersionInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateSecretVersionInput) (*mcp.CallToolResult, any, error) {
		if input.SecretID == "" {
			return nil, nil, fmt.Errorf("'secret_id' is required")
		}
		if len(input.Data) == 0 {
			return nil, nil, fmt.Errorf("'data' is required and must contain at least one key-value pair")
		}
		text, err := Create(ctx, serverAddress, input.SecretID, input.Data)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_secret_versions
// ---------------------------------------------------------------------------

// ListSecretVersionsInput defines the parameters for the
// list_secret_versions tool.
type ListSecretVersionsInput struct {
	SecretID string `json:"secret_id" jsonschema:"required,The parent secret's ID. Use get_secret or list_secrets to find the secret ID."`
}

// ListTool returns the MCP tool definition for list_secret_versions.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_secret_versions",
		Description: "List all versions of a secret. " +
			"Returns version metadata only (timestamps, backend version ID) — " +
			"encrypted data is not included for security. " +
			"Use this to understand a secret's version history before creating a new version " +
			"or to verify that a create_secret_version call succeeded.",
	}
}

// ListHandler returns the typed tool handler for list_secret_versions.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListSecretVersionsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListSecretVersionsInput) (*mcp.CallToolResult, any, error) {
		if input.SecretID == "" {
			return nil, nil, fmt.Errorf("'secret_id' is required")
		}
		text, err := List(ctx, serverAddress, input.SecretID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
