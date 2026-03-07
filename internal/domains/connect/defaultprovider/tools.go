package defaultprovider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	defaultproviderconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/defaultproviderconnection/v1"
)

// ---------------------------------------------------------------------------
// apply_default_provider_connection
// ---------------------------------------------------------------------------

type ApplyInput struct {
	ConnectionObject map[string]any `json:"connection_object" jsonschema:"required,Full DefaultProviderConnection object in OpenMCF envelope format."`
}

func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_default_provider_connection",
		Description: "Create or update a default provider connection that binds a credential as the default for an organization or environment. " +
			"Pass the full DefaultProviderConnection object as an OpenMCF envelope.",
	}
}

func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInput) (*mcp.CallToolResult, any, error) {
		if input.ConnectionObject == nil {
			return nil, nil, fmt.Errorf("'connection_object' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				var dpc defaultproviderconnectionv1.DefaultProviderConnection
				jsonBytes, err := json.Marshal(input.ConnectionObject)
				if err != nil {
					return "", fmt.Errorf("encoding connection object: %w", err)
				}
				if err := protojson.Unmarshal(jsonBytes, &dpc); err != nil {
					return "", fmt.Errorf("invalid connection object: %w", err)
				}
				client := defaultproviderconnectionv1.NewDefaultProviderConnectionCommandControllerClient(conn)
				resp, err := client.Apply(ctx, &dpc)
				if err != nil {
					return "", domains.RPCError(err, "default provider connection")
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
// get_default_provider_connection
// ---------------------------------------------------------------------------

type GetInput struct {
	ID string `json:"id" jsonschema:"required,DefaultProviderConnection ID."`
}

func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_default_provider_connection",
		Description: "Retrieve a default provider connection by ID.",
	}
}

func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultproviderconnectionv1.NewDefaultProviderConnectionQueryControllerClient(conn)
				resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default provider connection %q", input.ID))
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
// resolve_default_provider_connection
// ---------------------------------------------------------------------------

type ResolveInput struct {
	Org         string `json:"org" jsonschema:"required,Organization ID."`
	Provider    string `json:"provider" jsonschema:"required,Cloud resource provider (e.g. 'aws', 'gcp', 'azure')."`
	Environment string `json:"environment,omitempty" jsonschema:"Optional environment slug. If provided, resolves the env-level default; otherwise resolves the org-level default."`
}

func ResolveTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_default_provider_connection",
		Description: "Resolve the effective default provider connection for an organization and cloud provider. " +
			"Optionally specify an environment to resolve the env-level default. " +
			"Falls back to the org-level default if no env-level default is set.",
	}
}

func ResolveHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Provider == "" {
			return nil, nil, fmt.Errorf("'provider' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultproviderconnectionv1.NewDefaultProviderConnectionQueryControllerClient(conn)
				resp, err := client.Resolve(ctx, &defaultproviderconnectionv1.ResolveDefaultProviderConnectionRequest{
					Org:         input.Org,
					Environment: input.Environment,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default provider connection for %s in org %q", input.Provider, input.Org))
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
// delete_default_provider_connection
// ---------------------------------------------------------------------------

type DeleteInput struct {
	ID string `json:"id" jsonschema:"required,DefaultProviderConnection ID to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_default_provider_connection",
		Description: "Delete a default provider connection by ID. " +
			"WARNING: Cloud resources relying on this default will need a new default or explicit connection.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultproviderconnectionv1.NewDefaultProviderConnectionCommandControllerClient(conn)
				resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default provider connection %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
