package iacprovisionermapping

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	iacprovisionermappingv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/iacprovisionermapping/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_iac_provisioner_mapping
// ---------------------------------------------------------------------------

// ApplyInput defines the parameters for the apply_iac_provisioner_mapping tool.
type ApplyInput struct {
	MappingObject map[string]any `json:"mapping_object" jsonschema:"required,Full IacProvisionerMapping object in OpenMCF envelope format."`
}

// ApplyTool returns the MCP tool definition for apply_iac_provisioner_mapping.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_iac_provisioner_mapping",
		Description: "Create or update an IaC provisioner mapping (idempotent). " +
			"A mapping binds an IaC provisioner (e.g. Pulumi, Terraform) to an API resource selector " +
			"(kind + ID), overriding the platform default for that resource. " +
			"Pass the full IacProvisionerMapping object as an OpenMCF envelope.",
	}
}

// ApplyHandler returns the typed tool handler for apply_iac_provisioner_mapping.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInput) (*mcp.CallToolResult, any, error) {
		if input.MappingObject == nil {
			return nil, nil, fmt.Errorf("'mapping_object' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				var mapping iacprovisionermappingv1.IacProvisionerMapping
				jsonBytes, err := json.Marshal(input.MappingObject)
				if err != nil {
					return "", fmt.Errorf("encoding mapping object: %w", err)
				}
				if err := protojson.Unmarshal(jsonBytes, &mapping); err != nil {
					return "", fmt.Errorf("invalid mapping object: %w", err)
				}
				client := iacprovisionermappingv1.NewIacProvisionerMappingCommandControllerClient(conn)
				resp, err := client.Apply(ctx, &mapping)
				if err != nil {
					return "", domains.RPCError(err, "IaC provisioner mapping")
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
// get_iac_provisioner_mapping
// ---------------------------------------------------------------------------

// GetInput defines the parameters for the get_iac_provisioner_mapping tool.
type GetInput struct {
	ID string `json:"id" jsonschema:"required,IacProvisionerMapping ID."`
}

// GetTool returns the MCP tool definition for get_iac_provisioner_mapping.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_iac_provisioner_mapping",
		Description: "Retrieve an IaC provisioner mapping by ID. " +
			"Returns the full mapping including the resource selector and provisioner assignment.",
	}
}

// GetHandler returns the typed tool handler for get_iac_provisioner_mapping.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := iacprovisionermappingv1.NewIacProvisionerMappingQueryControllerClient(conn)
				resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("IaC provisioner mapping %q", input.ID))
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
// delete_iac_provisioner_mapping
// ---------------------------------------------------------------------------

// DeleteInput defines the parameters for the delete_iac_provisioner_mapping tool.
type DeleteInput struct {
	ID string `json:"id" jsonschema:"required,IacProvisionerMapping ID to delete."`
}

// DeleteTool returns the MCP tool definition for delete_iac_provisioner_mapping.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_iac_provisioner_mapping",
		Description: "Delete an IaC provisioner mapping by ID. " +
			"Once removed, the resource will fall back to the platform's default provisioner.",
	}
}

// DeleteHandler returns the typed tool handler for delete_iac_provisioner_mapping.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := iacprovisionermappingv1.NewIacProvisionerMappingCommandControllerClient(conn)
				resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("IaC provisioner mapping %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
