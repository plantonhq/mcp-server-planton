// Package service provides the MCP tools for the ServiceHub Service domain,
// backed by the ServiceQueryController, ServiceCommandController
// (ai.planton.servicehub.service.v1), and ApiResourceSearchQueryController
// (ai.planton.search.v1.apiresource) RPCs on the Planton backend.
//
// Seven tools are exposed:
//   - search_services:              free-text search with org/pagination filters
//   - get_service:                  retrieve a service by ID or org+slug
//   - apply_service:                create or update a service (idempotent)
//   - delete_service:               remove the service record
//   - disconnect_service_git_repo:  remove the webhook from the Git provider
//   - configure_service_webhook:    create or refresh the webhook on the Git provider
//   - list_service_branches:        list Git branches from the connected repository
package service

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_services
// ---------------------------------------------------------------------------

// SearchServicesInput defines the parameters for the search_services tool.
type SearchServicesInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within. Use list_organizations to discover available organizations."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text search query to filter services by name or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_services.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_services",
		Description: "Search services within an organization. " +
			"A service connects a Git repository to Planton Cloud's CI/CD pipeline system. " +
			"Returns lightweight search records with service IDs, names, and metadata. " +
			"Use get_service with a service ID from the results to retrieve full details.",
	}
}

// SearchHandler returns the typed tool handler for search_services.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchServicesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchServicesInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := Search(ctx, serverAddress, SearchInput{
			Org:        input.Org,
			SearchText: input.SearchText,
			PageNum:    input.PageNum,
			PageSize:   input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_service
// ---------------------------------------------------------------------------

// GetServiceInput defines the parameters for the get_service tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetServiceInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The service ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Service slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_service.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_service",
		Description: "Retrieve the full details of a service by its ID or by org+slug. " +
			"A service connects a Git repository to Planton Cloud's CI/CD pipeline system, " +
			"defining where code lives, how to build it, and where to deploy it. " +
			"Returns the complete service including metadata, spec (Git repo, pipeline config, ingress, deployment targets), " +
			"and status (per-environment deployment tracking). " +
			"The output JSON can be modified and passed to apply_service for updates.",
	}
}

// GetHandler returns the typed tool handler for get_service.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetServiceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetServiceInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := Get(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// apply_service
// ---------------------------------------------------------------------------

// ApplyServiceInput defines the parameters for the apply_service tool.
type ApplyServiceInput struct {
	Service map[string]any `json:"service" jsonschema:"required,The full Service resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('Service'), 'metadata' (with 'name' and 'org'), and 'spec' (with 'git_repo' and 'pipeline_configuration'). The output of get_service can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_service.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_service",
		Description: "Create or update a service (idempotent). " +
			"Accepts the full Service resource as a JSON object. " +
			"A service connects a Git repository to Planton Cloud's CI/CD pipeline — defining the code location, " +
			"build method, and deployment targets. " +
			"For new services, provide api_version, kind, metadata (name, org), and spec (git_repo, pipeline_configuration). " +
			"For updates, retrieve the service with get_service, modify the desired fields, and pass the result here. " +
			"Returns the applied service with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_service.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyServiceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyServiceInput) (*mcp.CallToolResult, any, error) {
		if len(input.Service) == 0 {
			return nil, nil, fmt.Errorf("'service' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.Service)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_service
// ---------------------------------------------------------------------------

// DeleteServiceInput defines the parameters for the delete_service tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteServiceInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The service ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Service slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_service.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_service",
		Description: "Delete a service record from the platform. " +
			"WARNING: This removes the service definition and disconnects webhooks from the Git provider. " +
			"It does NOT delete deployed cloud resources (microservices, databases, etc.) — " +
			"those must be removed separately. " +
			"Identify the service by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_service.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteServiceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteServiceInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := Delete(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// disconnect_service_git_repo
// ---------------------------------------------------------------------------

// DisconnectServiceGitRepoInput defines the parameters for the
// disconnect_service_git_repo tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DisconnectServiceGitRepoInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The service ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Service slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DisconnectGitRepoTool returns the MCP tool definition for disconnect_service_git_repo.
func DisconnectGitRepoTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "disconnect_service_git_repo",
		Description: "Disconnect the Git repository from a service by removing its webhook on GitHub/GitLab. " +
			"After disconnection, new commits no longer trigger pipelines. " +
			"The service definition remains in Planton Cloud and can be reconnected later via configure_service_webhook. " +
			"Useful for temporarily pausing automation, preparing to move a repository, or cleaning up before repo deletion. " +
			"Identify the service by ID or by org+slug.",
	}
}

// DisconnectGitRepoHandler returns the typed tool handler for disconnect_service_git_repo.
func DisconnectGitRepoHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DisconnectServiceGitRepoInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DisconnectServiceGitRepoInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := DisconnectGitRepo(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// configure_service_webhook
// ---------------------------------------------------------------------------

// ConfigureServiceWebhookInput defines the parameters for the
// configure_service_webhook tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type ConfigureServiceWebhookInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The service ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Service slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// ConfigureWebhookTool returns the MCP tool definition for configure_service_webhook.
func ConfigureWebhookTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "configure_service_webhook",
		Description: "Create or refresh the webhook on GitHub/GitLab for a service. " +
			"This registers (or updates) the webhook so that push and pull_request events trigger pipelines. " +
			"Useful for recovering from accidentally deleted webhooks, refreshing webhook configuration, " +
			"or troubleshooting webhook delivery issues. " +
			"Identify the service by ID or by org+slug.",
	}
}

// ConfigureWebhookHandler returns the typed tool handler for configure_service_webhook.
func ConfigureWebhookHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ConfigureServiceWebhookInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ConfigureServiceWebhookInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := ConfigureWebhook(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_service_branches
// ---------------------------------------------------------------------------

// ListServiceBranchesInput defines the parameters for the
// list_service_branches tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type ListServiceBranchesInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The service ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Service slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// ListBranchesTool returns the MCP tool definition for list_service_branches.
func ListBranchesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_service_branches",
		Description: "List all Git branches in the repository connected to a service. " +
			"Uses the GitHub/GitLab API to fetch branch information. " +
			"Useful for selecting a branch for pipeline configuration, " +
			"validating that a branch exists before triggering a pipeline, " +
			"or displaying available branches to the user. " +
			"Identify the service by ID or by org+slug.",
	}
}

// ListBranchesHandler returns the typed tool handler for list_service_branches.
func ListBranchesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListServiceBranchesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListServiceBranchesInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := ListBranches(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// shared validation
// ---------------------------------------------------------------------------

// validateIdentification checks that exactly one identification path is
// provided: either 'id' alone, or both 'org' and 'slug'.
func validateIdentification(id, org, slug string) error {
	hasID := id != ""
	hasOrg := org != ""
	hasSlug := slug != ""

	switch {
	case hasID && (hasOrg || hasSlug):
		return fmt.Errorf("provide either 'id' alone or both 'org' and 'slug' — not both paths")
	case hasID:
		return nil
	case hasOrg && hasSlug:
		return nil
	case hasOrg:
		return fmt.Errorf("'slug' is required when using 'org' for identification")
	case hasSlug:
		return fmt.Errorf("'org' is required when using 'slug' for identification")
	default:
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the service")
	}
}
