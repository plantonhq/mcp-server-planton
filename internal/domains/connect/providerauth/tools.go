package providerauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	pcav1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/providerconnectionauthorization/v1"
)

// ---------------------------------------------------------------------------
// apply_provider_connection_authorization
// ---------------------------------------------------------------------------

type ApplyInput struct {
	AuthorizationObject map[string]any `json:"authorization_object" jsonschema:"required,Full ProviderConnectionAuthorization object in OpenMCF envelope format."`
}

func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_provider_connection_authorization",
		Description: "Create or update a provider connection authorization. " +
			"Pass the full ProviderConnectionAuthorization object as an OpenMCF envelope. " +
			"Controls which credentials can be used in which environments.",
	}
}

func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInput) (*mcp.CallToolResult, any, error) {
		if input.AuthorizationObject == nil {
			return nil, nil, fmt.Errorf("'authorization_object' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				var pca pcav1.ProviderConnectionAuthorization
				jsonBytes, err := json.Marshal(input.AuthorizationObject)
				if err != nil {
					return "", fmt.Errorf("encoding authorization object: %w", err)
				}
				if err := protojson.Unmarshal(jsonBytes, &pca); err != nil {
					return "", fmt.Errorf("invalid authorization object: %w", err)
				}
				client := pcav1.NewProviderConnectionAuthorizationCommandControllerClient(conn)
				resp, err := client.Apply(ctx, &pca)
				if err != nil {
					return "", domains.RPCError(err, "provider connection authorization")
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
// get_provider_connection_authorization
// ---------------------------------------------------------------------------

type GetInput struct {
	ID         string `json:"id,omitempty"         jsonschema:"Authorization ID. Provide this OR the semantic key fields (org + provider + connection), not both."`
	Org        string `json:"org,omitempty"        jsonschema:"Organization ID (for semantic key lookup)."`
	Provider   string `json:"provider,omitempty"   jsonschema:"Cloud provider (e.g. 'aws', 'gcp') (for semantic key lookup)."`
	Connection string `json:"connection,omitempty" jsonschema:"Connection/credential slug (for semantic key lookup)."`
}

func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_provider_connection_authorization",
		Description: "Retrieve a provider connection authorization by ID or by semantic key (org + provider + connection). " +
			"Provide either 'id' alone, or all three of 'org', 'provider', and 'connection'.",
	}
}

func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInput) (*mcp.CallToolResult, any, error) {
		hasID := input.ID != ""
		hasSemantic := input.Org != "" || input.Provider != "" || input.Connection != ""
		if hasID == hasSemantic {
			return nil, nil, fmt.Errorf("provide either 'id' alone, or all three of 'org', 'provider', and 'connection'")
		}

		var text string
		var err error

		if hasID {
			text, err = domains.WithConnection(ctx, serverAddress,
				func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
					client := pcav1.NewProviderConnectionAuthorizationQueryControllerClient(conn)
					resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ID})
					if err != nil {
						return "", domains.RPCError(err, fmt.Sprintf("provider connection authorization %q", input.ID))
					}
					return domains.MarshalJSON(resp)
				})
		} else {
			if input.Org == "" || input.Provider == "" || input.Connection == "" {
				return nil, nil, fmt.Errorf("all three of 'org', 'provider', and 'connection' are required for semantic key lookup")
			}
			providerEnum, resolveErr := domains.ResolveProvider(input.Provider)
			if resolveErr != nil {
				return nil, nil, resolveErr
			}
			text, err = domains.WithConnection(ctx, serverAddress,
				func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
					client := pcav1.NewProviderConnectionAuthorizationQueryControllerClient(conn)
					resp, err := client.GetBySemanticKey(ctx, &pcav1.GetBySemanticKeyRequest{
						Org:        input.Org,
						Provider:   providerEnum,
						Connection: input.Connection,
					})
					if err != nil {
						return "", domains.RPCError(err,
							fmt.Sprintf("provider connection authorization for %s/%s in org %q", input.Provider, input.Connection, input.Org))
					}
					return domains.MarshalJSON(resp)
				})
		}
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_provider_connection_authorization
// ---------------------------------------------------------------------------

type DeleteInput struct {
	ID string `json:"id" jsonschema:"required,Authorization ID to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_provider_connection_authorization",
		Description: "Delete a provider connection authorization by ID. " +
			"After deletion, the credential will no longer be usable in any environment " +
			"unless a new authorization is created.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := pcav1.NewProviderConnectionAuthorizationCommandControllerClient(conn)
				resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("provider connection authorization %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
