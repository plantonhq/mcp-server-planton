package serviceaccount

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"

	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// create_service_account_key
// ---------------------------------------------------------------------------

// CreateKeyInput defines the parameters for create_service_account_key.
type CreateKeyInput struct {
	ServiceAccountID string `json:"service_account_id" jsonschema:"required,ID of the service account to create a key for."`
}

// CreateKeyTool returns the MCP tool definition for create_service_account_key.
func CreateKeyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_service_account_key",
		Description: "Generate a new API key for a service account. " +
			"WARNING: The raw key value (pck_ prefix) is returned ONLY in this response — " +
			"it is never stored or retrievable again. Store it securely.",
	}
}

// CreateKeyHandler returns the typed tool handler for create_service_account_key.
func CreateKeyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateKeyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateKeyInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceAccountID == "" {
			return nil, nil, fmt.Errorf("'service_account_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := serviceaccountv1.NewServiceAccountCommandControllerClient(conn)
				resp, err := client.CreateKey(ctx, &serviceaccountv1.ServiceAccountId{Value: input.ServiceAccountID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("key for service account %q", input.ServiceAccountID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// revoke_service_account_key
// ---------------------------------------------------------------------------

// RevokeKeyInput defines the parameters for revoke_service_account_key.
type RevokeKeyInput struct {
	ServiceAccountID string `json:"service_account_id" jsonschema:"required,ID of the service account that owns the key."`
	APIKeyID         string `json:"api_key_id"         jsonschema:"required,ID of the API key to revoke."`
}

// RevokeKeyTool returns the MCP tool definition for revoke_service_account_key.
func RevokeKeyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "revoke_service_account_key",
		Description: "Revoke (delete) a specific API key owned by a service account. " +
			"The key is verified to belong to the service account before revocation.",
	}
}

// RevokeKeyHandler returns the typed tool handler for revoke_service_account_key.
func RevokeKeyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RevokeKeyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RevokeKeyInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceAccountID == "" {
			return nil, nil, fmt.Errorf("'service_account_id' is required")
		}
		if input.APIKeyID == "" {
			return nil, nil, fmt.Errorf("'api_key_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := serviceaccountv1.NewServiceAccountCommandControllerClient(conn)
				resp, err := client.RevokeKey(ctx, &serviceaccountv1.RevokeKeyRequest{
					ServiceAccountId: input.ServiceAccountID,
					ApiKeyId:         input.APIKeyID,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("API key %q for service account %q", input.APIKeyID, input.ServiceAccountID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_service_account_keys
// ---------------------------------------------------------------------------

// ListKeysInput defines the parameters for list_service_account_keys.
type ListKeysInput struct {
	ServiceAccountID string `json:"service_account_id" jsonschema:"required,ID of the service account."`
}

// ListKeysTool returns the MCP tool definition for list_service_account_keys.
func ListKeysTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_service_account_keys",
		Description: "List all API keys belonging to a service account. " +
			"Returns key metadata (ID, creation time, status) but not the raw key values.",
	}
}

// ListKeysHandler returns the typed tool handler for list_service_account_keys.
func ListKeysHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListKeysInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListKeysInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceAccountID == "" {
			return nil, nil, fmt.Errorf("'service_account_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := serviceaccountv1.NewServiceAccountQueryControllerClient(conn)
				resp, err := client.ListKeys(ctx, &serviceaccountv1.ServiceAccountId{Value: input.ServiceAccountID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("keys for service account %q", input.ServiceAccountID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
