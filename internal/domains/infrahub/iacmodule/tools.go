// Package iacmodule provides the MCP tools for the IaC Module domain,
// backed by the InfraHubSearchQueryController (ai.planton.search.v1.infrahub)
// and IacModuleQueryController (ai.planton.infrahub.iacmodule.v1) RPCs on the
// Planton backend.
//
// Two tools are exposed:
//   - search_iac_modules: search for IaC modules (official + org-scoped)
//   - get_iac_module:     retrieve full module details by ID
package iacmodule

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_iac_modules
// ---------------------------------------------------------------------------

// SearchIacModulesInput defines the parameters for the search_iac_modules tool.
type SearchIacModulesInput struct {
	Org         string `json:"org,omitempty"         jsonschema:"Organization identifier. When provided, results include both official and organization-specific IaC modules. When omitted, only official modules are returned."`
	SearchText  string `json:"search_text,omitempty" jsonschema:"Free-text search query to filter modules by name or description."`
	Kind        string `json:"kind,omitempty"        jsonschema:"PascalCase cloud resource kind to filter by (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds. When provided, only modules that can provision this resource type are returned."`
	Provisioner string `json:"provisioner,omitempty" jsonschema:"IaC provisioner to filter by (terraform, pulumi, or tofu). When omitted, modules for all provisioners are returned."`
	Provider    string `json:"provider,omitempty"    jsonschema:"Cloud provider to filter by (e.g. aws, gcp, azure). When omitted, modules for all providers are returned."`
	PageNum     int32  `json:"page_num,omitempty"    jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize    int32  `json:"page_size,omitempty"   jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_iac_modules.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_iac_modules",
		Description: "Search for IaC (Infrastructure as Code) modules. " +
			"IaC modules are the provisioning implementations that deploy cloud resources — " +
			"each module targets a specific cloud resource kind and IaC provisioner (Terraform, Pulumi, or OpenTofu). " +
			"When 'org' is provided, results include both official platform modules and organization-specific modules. " +
			"When 'org' is omitted, only official modules are returned. " +
			"Use the 'kind' filter to find modules that can provision a specific deployment component " +
			"(e.g. kind=AwsEksCluster returns all modules capable of deploying EKS clusters). " +
			"Use get_iac_module with the module ID from the results to retrieve full details.",
	}
}

// SearchHandler returns the typed tool handler for search_iac_modules.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchIacModulesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchIacModulesInput) (*mcp.CallToolResult, any, error) {
		text, err := Search(ctx, serverAddress, SearchInput{
			Org:         input.Org,
			SearchText:  input.SearchText,
			Kind:        input.Kind,
			Provisioner: input.Provisioner,
			Provider:    input.Provider,
			PageNum:     input.PageNum,
			PageSize:    input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_iac_module
// ---------------------------------------------------------------------------

// GetIacModuleInput defines the parameters for the get_iac_module tool.
type GetIacModuleInput struct {
	ID string `json:"id" jsonschema:"required,The IaC module ID obtained from search_iac_modules results."`
}

// GetTool returns the MCP tool definition for get_iac_module.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_iac_module",
		Description: "Retrieve the full details of an IaC module by ID. " +
			"Returns the complete module including metadata, provisioner type, cloud resource kind, " +
			"Git repository URL, version, and parameter schema. " +
			"Use search_iac_modules to discover module IDs.",
	}
}

// GetHandler returns the typed tool handler for get_iac_module.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetIacModuleInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetIacModuleInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Get(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
