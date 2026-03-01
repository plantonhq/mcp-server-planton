package team

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// MemberInput represents a team member as provided by the MCP tool caller.
type MemberInput struct {
	MemberType string `json:"member_type" jsonschema:"required,Type of member: 'identity_account' or 'team'."`
	MemberID   string `json:"member_id"   jsonschema:"required,ID of the identity account or team to add."`
}

// ---------------------------------------------------------------------------
// create_team
// ---------------------------------------------------------------------------

// CreateTeamInput defines the parameters for the create_team tool.
type CreateTeamInput struct {
	Org         string        `json:"org"                    jsonschema:"required,Organization ID to create the team in."`
	Name        string        `json:"name"                   jsonschema:"required,Display name for the team."`
	Description string        `json:"description,omitempty"  jsonschema:"Optional description of the team."`
	Members     []MemberInput `json:"members,omitempty"      jsonschema:"Optional initial members. Each entry needs member_type ('identity_account' or 'team') and member_id."`
}

// CreateTool returns the MCP tool definition for create_team.
func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_team",
		Description: "Create a new team in an organization. " +
			"Teams group identity accounts (and other teams) for collective access control. " +
			"Optionally provide initial members at creation time.",
	}
}

// CreateHandler returns the typed tool handler for create_team.
func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateTeamInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateTeamInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		text, err := Create(ctx, serverAddress, input.Org, input.Name, input.Description, input.Members)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_team
// ---------------------------------------------------------------------------

// GetTeamInput defines the parameters for the get_team tool.
type GetTeamInput struct {
	TeamID string `json:"team_id" jsonschema:"required,Team ID (e.g. tm-xxx)."`
}

// GetTool returns the MCP tool definition for get_team.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_team",
		Description: "Get a team by ID. " +
			"Returns the full team including metadata, description, and current member list.",
	}
}

// GetHandler returns the typed tool handler for get_team.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetTeamInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetTeamInput) (*mcp.CallToolResult, any, error) {
		if input.TeamID == "" {
			return nil, nil, fmt.Errorf("'team_id' is required")
		}
		text, err := Get(ctx, serverAddress, input.TeamID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// update_team
// ---------------------------------------------------------------------------

// UpdateTeamInput defines the parameters for the update_team tool.
type UpdateTeamInput struct {
	TeamID      string        `json:"team_id"                jsonschema:"required,ID of the team to update."`
	Name        string        `json:"name,omitempty"         jsonschema:"New display name. Omit to keep unchanged."`
	Description string        `json:"description,omitempty"  jsonschema:"New description. Omit to keep unchanged."`
	Members     *[]MemberInput `json:"members,omitempty"     jsonschema:"New member list (replaces existing). Omit to keep unchanged. Send empty array to remove all members."`
}

// UpdateTool returns the MCP tool definition for update_team.
func UpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_team",
		Description: "Update an existing team. " +
			"Uses read-modify-write: fetches the current team, applies provided changes, then saves. " +
			"Only non-empty fields are updated. If 'members' is provided it replaces the entire member list.",
	}
}

// UpdateHandler returns the typed tool handler for update_team.
func UpdateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdateTeamInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdateTeamInput) (*mcp.CallToolResult, any, error) {
		if input.TeamID == "" {
			return nil, nil, fmt.Errorf("'team_id' is required")
		}
		text, err := Update(ctx, serverAddress, input.TeamID, UpdateFields{
			Name:        input.Name,
			Description: input.Description,
			Members:     input.Members,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_team
// ---------------------------------------------------------------------------

// DeleteTeamInput defines the parameters for the delete_team tool.
type DeleteTeamInput struct {
	TeamID string `json:"team_id" jsonschema:"required,ID of the team to delete."`
}

// DeleteTool returns the MCP tool definition for delete_team.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_team",
		Description: "Permanently delete a team. " +
			"This removes the team and all its access policies. This action cannot be undone.",
	}
}

// DeleteHandler returns the typed tool handler for delete_team.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteTeamInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteTeamInput) (*mcp.CallToolResult, any, error) {
		if input.TeamID == "" {
			return nil, nil, fmt.Errorf("'team_id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.TeamID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
