// Package azure provides MCP tools for Azure cloud operations via the Planton control plane.
package azure

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	azurecompute "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/azure/v1/compute"
	azurestorage "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/azure/v1/storage"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudopsctx "github.com/plantonhq/mcp-server-planton/internal/domains/cloudops"
	"google.golang.org/grpc"
)

type listVirtualMachinesInput struct {
	Org               string `json:"org"                        jsonschema:"required,Organization slug."`
	Env               string `json:"env,omitempty"              jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection        string `json:"connection,omitempty"       jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
	ResourceGroup     string `json:"resource_group,omitempty"    jsonschema:"Azure resource group to filter VMs. Optional."`
}

func ListVirtualMachinesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_azure_virtual_machines",
		Description: "List Azure virtual machines in a subscription, optionally filtered by resource group. Returns VM details including name, status, size, and location. Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListVirtualMachinesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *listVirtualMachinesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *listVirtualMachinesInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := azurecompute.NewVirtualMachineQueryControllerClient(conn).List(ctx, &azurecompute.ListVirtualMachinesRequest{
					Context:       opsCtx,
					ResourceGroup: input.ResourceGroup,
				})
				if err != nil {
					return "", domains.RPCError(err, "list Azure virtual machines")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type listBlobContainersInput struct {
	Org                string `json:"org"                        jsonschema:"required,Organization slug."`
	Env                string `json:"env,omitempty"              jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind  string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug  string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection         string `json:"connection,omitempty"       jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
	StorageAccountName string `json:"storage_account_name"      jsonschema:"required,Azure storage account name to list containers from."`
	ResourceGroup      string `json:"resource_group,omitempty"  jsonschema:"Azure resource group containing the storage account. Optional."`
}

func ListBlobContainersTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_azure_blob_containers",
		Description: "List blob containers in an Azure Storage account. Returns container names and metadata. Requires the storage account name; optionally filter by resource group. Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListBlobContainersHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *listBlobContainersInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *listBlobContainersInput) (*mcp.CallToolResult, any, error) {
		if input.StorageAccountName == "" {
			return nil, nil, fmt.Errorf("'storage_account_name' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := azurestorage.NewBlobContainerQueryControllerClient(conn).List(ctx, &azurestorage.ListBlobContainersRequest{
					Context:            opsCtx,
					StorageAccountName: input.StorageAccountName,
					ResourceGroup:      input.ResourceGroup,
				})
				if err != nil {
					return "", domains.RPCError(err, "list Azure blob containers")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
