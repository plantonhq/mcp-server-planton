package defaultrunner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	defaultrunnerbindingv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/defaultrunnerbinding/v1"
)

// ---------------------------------------------------------------------------
// apply_default_runner_binding
// ---------------------------------------------------------------------------

type ApplyInput struct {
	BindingObject map[string]any `json:"binding_object" jsonschema:"required,Full DefaultRunnerBinding object in OpenMCF envelope format."`
}

func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_default_runner_binding",
		Description: "Create or update a default runner binding that designates a runner as the default for an organization. " +
			"Pass the full DefaultRunnerBinding object as an OpenMCF envelope.",
	}
}

func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInput) (*mcp.CallToolResult, any, error) {
		if input.BindingObject == nil {
			return nil, nil, fmt.Errorf("'binding_object' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				var drb defaultrunnerbindingv1.DefaultRunnerBinding
				jsonBytes, err := json.Marshal(input.BindingObject)
				if err != nil {
					return "", fmt.Errorf("encoding binding object: %w", err)
				}
				if err := protojson.Unmarshal(jsonBytes, &drb); err != nil {
					return "", fmt.Errorf("invalid binding object: %w", err)
				}
				client := defaultrunnerbindingv1.NewDefaultRunnerBindingCommandControllerClient(conn)
				resp, err := client.Apply(ctx, &drb)
				if err != nil {
					return "", domains.RPCError(err, "default runner binding")
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
// get_default_runner_binding
// ---------------------------------------------------------------------------

type GetInput struct {
	ID string `json:"id" jsonschema:"required,DefaultRunnerBinding ID."`
}

func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_default_runner_binding",
		Description: "Retrieve a default runner binding by ID.",
	}
}

func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultrunnerbindingv1.NewDefaultRunnerBindingQueryControllerClient(conn)
				resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default runner binding %q", input.ID))
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
// resolve_default_runner_binding
// ---------------------------------------------------------------------------

type ResolveInput struct {
	Org string `json:"org" jsonschema:"required,Organization ID to resolve the default runner for."`
}

func ResolveTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_default_runner_binding",
		Description: "Resolve the effective default runner binding for an organization. " +
			"Returns the runner registration that is designated as the default.",
	}
}

func ResolveHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultrunnerbindingv1.NewDefaultRunnerBindingQueryControllerClient(conn)
				resp, err := client.Resolve(ctx, &defaultrunnerbindingv1.ResolveDefaultRunnerBindingRequest{
					Org: input.Org,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default runner binding for org %q", input.Org))
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
// delete_default_runner_binding
// ---------------------------------------------------------------------------

type DeleteInput struct {
	ID string `json:"id" jsonschema:"required,DefaultRunnerBinding ID to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_default_runner_binding",
		Description: "Delete a default runner binding by ID.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := defaultrunnerbindingv1.NewDefaultRunnerBindingCommandControllerClient(conn)
				resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("default runner binding %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
