// Package infraproject provides the MCP tools for the InfraProject domain,
// backed by the InfraProjectQueryController, InfraProjectCommandController
// (ai.planton.infrahub.infraproject.v1), and InfraHubSearchQueryController
// (ai.planton.search.v1.infrahub) RPCs on the Planton backend.
//
// Six tools are exposed:
//   - search_infra_projects:       free-text search with org/env/pagination filters
//   - get_infra_project:           retrieve a project by ID or org+slug
//   - apply_infra_project:         create or update a project (idempotent)
//   - delete_infra_project:        remove the project record
//   - check_infra_project_slug:    check slug availability within an org
//   - undeploy_infra_project:      tear down deployed resources, keep the record
package infraproject

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_infra_projects
// ---------------------------------------------------------------------------

// SearchInfraProjectsInput defines the parameters for the
// search_infra_projects tool.
type SearchInfraProjectsInput struct {
	Org        string `json:"org"                jsonschema:"required,Organization identifier to search within. Use list_organizations to discover available organizations."`
	Env        string `json:"env,omitempty"       jsonschema:"Environment slug to filter by. When provided, only infra-chart sourced projects deployed to this environment are returned. When omitted, all projects (both infra-chart and git-repo sourced) are returned."`
	SearchText string `json:"search_text,omitempty" jsonschema:"Free-text search query to filter projects by name or description."`
	PageNum    int32  `json:"page_num,omitempty"  jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty" jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_infra_projects.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_infra_projects",
		Description: "Search infra projects within an organization. " +
			"Infra projects are deployable infrastructure compositions sourced from infra charts or Git repositories. " +
			"Returns lightweight search records with project IDs, names, and metadata. " +
			"Use get_infra_project with a project ID from the results to retrieve full details.",
	}
}

// SearchHandler returns the typed tool handler for search_infra_projects.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchInfraProjectsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchInfraProjectsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := Search(ctx, serverAddress, SearchInput{
			Org:        input.Org,
			Env:        input.Env,
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
// get_infra_project
// ---------------------------------------------------------------------------

// GetInfraProjectInput defines the parameters for the get_infra_project tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetInfraProjectInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The infra project ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Project slug for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_infra_project.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_infra_project",
		Description: "Retrieve the full details of an infra project by its ID or by org+slug. " +
			"Returns the complete project including metadata, spec (source type, chart/git config, parameters), " +
			"and status (rendered YAML, cloud resource DAG, pipeline ID). " +
			"The output JSON can be modified and passed to apply_infra_project for updates.",
	}
}

// GetHandler returns the typed tool handler for get_infra_project.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInfraProjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInfraProjectInput) (*mcp.CallToolResult, any, error) {
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
// apply_infra_project
// ---------------------------------------------------------------------------

// ApplyInfraProjectInput defines the parameters for the apply_infra_project tool.
type ApplyInfraProjectInput struct {
	InfraProject map[string]any `json:"infra_project" jsonschema:"required,The full InfraProject resource as a JSON object. Must include 'metadata' (with 'name' and 'org') and 'spec' (with 'source' and the corresponding source config). The output of get_infra_project can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_infra_project.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_infra_project",
		Description: "Create or update an infra project (idempotent). " +
			"Accepts the full InfraProject resource as a JSON object. " +
			"For new projects, provide metadata (name, org) and spec (source type with chart or git-repo config). " +
			"For updates, retrieve the project with get_infra_project, modify the desired fields, and pass the result here. " +
			"Returns the applied project with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_infra_project.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInfraProjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInfraProjectInput) (*mcp.CallToolResult, any, error) {
		if len(input.InfraProject) == 0 {
			return nil, nil, fmt.Errorf("'infra_project' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.InfraProject)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_infra_project
// ---------------------------------------------------------------------------

// DeleteInfraProjectInput defines the parameters for the
// delete_infra_project tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteInfraProjectInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The infra project ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Project slug for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_infra_project.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_infra_project",
		Description: "Delete an infra project record from the platform. " +
			"WARNING: This removes the database record only — it does NOT tear down deployed cloud resources. " +
			"To tear down infrastructure first, use undeploy_infra_project before deleting. " +
			"Identify the project by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_infra_project.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInfraProjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInfraProjectInput) (*mcp.CallToolResult, any, error) {
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
// check_infra_project_slug
// ---------------------------------------------------------------------------

// CheckInfraProjectSlugInput defines the parameters for the
// check_infra_project_slug tool.
type CheckInfraProjectSlugInput struct {
	Org  string `json:"org"  jsonschema:"required,Organization identifier to check the slug within."`
	Slug string `json:"slug" jsonschema:"required,The slug to check for availability."`
}

// CheckSlugTool returns the MCP tool definition for check_infra_project_slug.
func CheckSlugTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_infra_project_slug",
		Description: "Check whether an infra project slug is available within an organization. " +
			"Returns true if no project with the given slug exists in the organization. " +
			"Use this before apply_infra_project to avoid slug conflicts.",
	}
}

// CheckSlugHandler returns the typed tool handler for check_infra_project_slug.
func CheckSlugHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckInfraProjectSlugInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CheckInfraProjectSlugInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := CheckSlugAvailability(ctx, serverAddress, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// undeploy_infra_project
// ---------------------------------------------------------------------------

// UndeployInfraProjectInput defines the parameters for the
// undeploy_infra_project tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type UndeployInfraProjectInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The infra project ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Project slug for lookup within an organization. Must be paired with 'org'."`
}

// UndeployTool returns the MCP tool definition for undeploy_infra_project.
func UndeployTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "undeploy_infra_project",
		Description: "Tear down all cloud resources deployed by an infra project while keeping the project record. " +
			"This triggers an undeploy pipeline that destroys the infrastructure managed by the project. " +
			"The project record is preserved and can be redeployed later via apply_infra_project. " +
			"Use delete_infra_project after undeploying if you also want to remove the record.",
	}
}

// UndeployHandler returns the typed tool handler for undeploy_infra_project.
func UndeployHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UndeployInfraProjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UndeployInfraProjectInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := Undeploy(ctx, serverAddress, input.ID, input.Org, input.Slug)
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
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the infra project")
	}
}
