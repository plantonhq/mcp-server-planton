package identity

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// whoami
// ---------------------------------------------------------------------------

// WhoAmIInput takes no parameters — the identity is derived from the auth token.
type WhoAmIInput struct{}

// WhoAmITool returns the MCP tool definition for whoami.
func WhoAmITool() *mcp.Tool {
	return &mcp.Tool{
		Name: "whoami",
		Description: "Get the identity account of the currently authenticated user. " +
			"Returns the full identity account including ID, email, name, and account type. " +
			"Use this to confirm which user context is active before performing operations.",
	}
}

// WhoAmIHandler returns the typed tool handler for whoami.
func WhoAmIHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *WhoAmIInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ *WhoAmIInput) (*mcp.CallToolResult, any, error) {
		text, err := WhoAmI(ctx, serverAddress)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_identity_account
// ---------------------------------------------------------------------------

// GetIdentityAccountInput supports dual-resolution: by ID or by email.
// Exactly one must be provided.
type GetIdentityAccountInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"Identity account ID (e.g. ida-xxx). Provide this OR email, not both."`
	Email string `json:"email,omitempty" jsonschema:"Email address of the identity account. Provide this OR id, not both."`
}

// GetTool returns the MCP tool definition for get_identity_account.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_identity_account",
		Description: "Look up an identity account by ID or email address. " +
			"Provide exactly one of 'id' or 'email'. " +
			"Returns the full identity account including metadata, spec (email, name, account type), and status.",
	}
}

// GetHandler returns the typed tool handler for get_identity_account.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetIdentityAccountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetIdentityAccountInput) (*mcp.CallToolResult, any, error) {
		hasID := input.ID != ""
		hasEmail := input.Email != ""
		if hasID == hasEmail {
			return nil, nil, fmt.Errorf("provide exactly one of 'id' or 'email'")
		}
		var text string
		var err error
		if hasID {
			text, err = Get(ctx, serverAddress, input.ID)
		} else {
			text, err = GetByEmail(ctx, serverAddress, input.Email)
		}
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// invite_member
// ---------------------------------------------------------------------------

// InviteMemberInput defines the parameters for inviting a user to an org.
type InviteMemberInput struct {
	Org        string   `json:"org"          jsonschema:"required,Organization ID to invite the user to."`
	Email      string   `json:"email"        jsonschema:"required,Email address of the user to invite."`
	IAMRoleIDs []string `json:"iam_role_ids" jsonschema:"required,IAM role IDs to grant upon acceptance. Use list_iam_roles_for_resource_kind to discover available roles."`
}

// InviteTool returns the MCP tool definition for invite_member.
func InviteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "invite_member",
		Description: "Invite a user to an organization with specified IAM roles. " +
			"The user receives an invitation that, when accepted, grants the specified roles on the organization. " +
			"Requires iam_policy_update permission on the organization.",
	}
}

// InviteHandler returns the typed tool handler for invite_member.
func InviteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *InviteMemberInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *InviteMemberInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Email == "" {
			return nil, nil, fmt.Errorf("'email' is required")
		}
		if len(input.IAMRoleIDs) == 0 {
			return nil, nil, fmt.Errorf("'iam_role_ids' is required and must contain at least one role ID")
		}
		text, err := Invite(ctx, serverAddress, input.Org, input.Email, input.IAMRoleIDs)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_invitations
// ---------------------------------------------------------------------------

// ListInvitationsInput defines the parameters for listing user invitations.
type ListInvitationsInput struct {
	Org    string `json:"org"              jsonschema:"required,Organization ID to list invitations for."`
	Status string `json:"status,omitempty" jsonschema:"Invitation status filter. Valid values: pending, accepted, removed. Defaults to 'pending' if omitted."`
}

// ListInvitationsTool returns the MCP tool definition for list_invitations.
func ListInvitationsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_invitations",
		Description: "List user invitations for an organization, filtered by status. " +
			"Defaults to showing pending invitations. " +
			"Returns invitation details including email, assigned roles, and status.",
	}
}

// ListInvitationsHandler returns the typed tool handler for list_invitations.
func ListInvitationsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListInvitationsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListInvitationsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		status := input.Status
		if status == "" {
			status = "pending"
		}
		text, err := ListInvitations(ctx, serverAddress, input.Org, status)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
