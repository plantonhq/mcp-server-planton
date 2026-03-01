package policy

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// create_iam_policy
// ---------------------------------------------------------------------------

type CreateIamPolicyInput struct {
	PrincipalKind string `json:"principal_kind" jsonschema:"required,Kind of principal (e.g. 'identity_account', 'team')."`
	PrincipalID   string `json:"principal_id"   jsonschema:"required,ID of the principal."`
	ResourceKind  string `json:"resource_kind"  jsonschema:"required,Kind of resource (e.g. 'organization', 'aws_credential')."`
	ResourceID    string `json:"resource_id"    jsonschema:"required,ID of the resource."`
	Relation      string `json:"relation"       jsonschema:"required,Relation/permission to grant (e.g. 'admin', 'viewer', 'editor'). Use list_iam_roles_for_resource_kind to discover valid relations."`
}

func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_iam_policy",
		Description: "Grant a principal access to a resource with a specific relation. " +
			"Creates a single IAM policy (principal -> relation -> resource). " +
			"Idempotent: if the policy already exists, the call succeeds without duplicating.",
	}
}

func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateIamPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *CreateIamPolicyInput) (*mcp.CallToolResult, any, error) {
		if in.PrincipalKind == "" || in.PrincipalID == "" || in.ResourceKind == "" || in.ResourceID == "" || in.Relation == "" {
			return nil, nil, fmt.Errorf("all fields (principal_kind, principal_id, resource_kind, resource_id, relation) are required")
		}
		text, err := Create(ctx, serverAddress, in.PrincipalKind, in.PrincipalID, in.ResourceKind, in.ResourceID, in.Relation)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_iam_policy
// ---------------------------------------------------------------------------

type DeleteIamPolicyInput struct {
	PrincipalKind string `json:"principal_kind" jsonschema:"required,Kind of principal."`
	PrincipalID   string `json:"principal_id"   jsonschema:"required,ID of the principal."`
	ResourceKind  string `json:"resource_kind"  jsonschema:"required,Kind of resource."`
	ResourceID    string `json:"resource_id"    jsonschema:"required,ID of the resource."`
	Relation      string `json:"relation"       jsonschema:"required,Relation to revoke."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_iam_policy",
		Description: "Revoke a specific access grant (principal -> relation -> resource). " +
			"Idempotent: succeeds even if the policy doesn't exist.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteIamPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *DeleteIamPolicyInput) (*mcp.CallToolResult, any, error) {
		if in.PrincipalKind == "" || in.PrincipalID == "" || in.ResourceKind == "" || in.ResourceID == "" || in.Relation == "" {
			return nil, nil, fmt.Errorf("all fields (principal_kind, principal_id, resource_kind, resource_id, relation) are required")
		}
		text, err := Delete(ctx, serverAddress, in.PrincipalKind, in.PrincipalID, in.ResourceKind, in.ResourceID, in.Relation)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// upsert_iam_policies
// ---------------------------------------------------------------------------

type UpsertIamPoliciesInput struct {
	PrincipalKind string   `json:"principal_kind" jsonschema:"required,Kind of principal."`
	PrincipalID   string   `json:"principal_id"   jsonschema:"required,ID of the principal."`
	ResourceKind  string   `json:"resource_kind"  jsonschema:"required,Kind of resource."`
	ResourceID    string   `json:"resource_id"    jsonschema:"required,ID of the resource."`
	Relations     []string `json:"relations"      jsonschema:"required,Desired set of relations. After this call the principal will have EXACTLY these relations on the resource—extras removed, missing ones added."`
}

func UpsertTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "upsert_iam_policies",
		Description: "Declaratively sync a principal's relations on a resource. " +
			"After this call the principal has EXACTLY the specified relations — " +
			"extras are removed and missing ones are added. Idempotent.",
	}
}

func UpsertHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpsertIamPoliciesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *UpsertIamPoliciesInput) (*mcp.CallToolResult, any, error) {
		if in.PrincipalKind == "" || in.PrincipalID == "" || in.ResourceKind == "" || in.ResourceID == "" {
			return nil, nil, fmt.Errorf("principal_kind, principal_id, resource_kind, and resource_id are required")
		}
		if len(in.Relations) == 0 {
			return nil, nil, fmt.Errorf("'relations' must contain at least one relation")
		}
		text, err := Upsert(ctx, serverAddress, in.PrincipalKind, in.PrincipalID, in.ResourceKind, in.ResourceID, in.Relations)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// revoke_org_access
// ---------------------------------------------------------------------------

type RevokeOrgAccessInput struct {
	IdentityAccountID string `json:"identity_account_id" jsonschema:"required,ID of the user whose access to revoke."`
	OrganizationID    string `json:"organization_id"     jsonschema:"required,ID of the organization to revoke access from."`
}

func RevokeOrgAccessTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "revoke_org_access",
		Description: "Remove ALL of a user's access to an organization and its child resources. " +
			"This is a nuclear revocation — it removes every IAM policy the user has " +
			"within the organization scope. Use with caution.",
	}
}

func RevokeOrgAccessHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RevokeOrgAccessInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *RevokeOrgAccessInput) (*mcp.CallToolResult, any, error) {
		if in.IdentityAccountID == "" {
			return nil, nil, fmt.Errorf("'identity_account_id' is required")
		}
		if in.OrganizationID == "" {
			return nil, nil, fmt.Errorf("'organization_id' is required")
		}
		text, err := RevokeOrgAccess(ctx, serverAddress, in.IdentityAccountID, in.OrganizationID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_resource_access
// ---------------------------------------------------------------------------

type ListResourceAccessInput struct {
	ResourceKind     string `json:"resource_kind"               jsonschema:"required,Kind of resource to query access for (e.g. 'organization', 'vpc')."`
	ResourceID       string `json:"resource_id"                 jsonschema:"required,ID of the resource."`
	IncludeInherited bool   `json:"include_inherited,omitempty" jsonschema:"Include inherited access from parent resources. Defaults to false."`
}

func ListResourceAccessTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_resource_access",
		Description: "List who has access to a resource. " +
			"Returns principals grouped with their assigned roles. " +
			"Set include_inherited=true to also show access inherited from parent resources.",
	}
}

func ListResourceAccessHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListResourceAccessInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *ListResourceAccessInput) (*mcp.CallToolResult, any, error) {
		if in.ResourceKind == "" || in.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_kind' and 'resource_id' are required")
		}
		text, err := ListResourceAccess(ctx, serverAddress, in.ResourceKind, in.ResourceID, in.IncludeInherited)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// check_authorization
// ---------------------------------------------------------------------------

type CheckAuthorizationInput struct {
	PrincipalKind string `json:"principal_kind" jsonschema:"required,Kind of principal."`
	PrincipalID   string `json:"principal_id"   jsonschema:"required,ID of the principal."`
	ResourceKind  string `json:"resource_kind"  jsonschema:"required,Kind of resource."`
	ResourceID    string `json:"resource_id"    jsonschema:"required,ID of the resource."`
	Relation      string `json:"relation"       jsonschema:"required,Relation/permission to check."`
}

func CheckAuthorizationTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_authorization",
		Description: "Check if a principal is authorized for a specific relation on a resource. " +
			"Returns {is_authorized: true/false}. " +
			"Any authenticated user can perform this check.",
	}
}

func CheckAuthorizationHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckAuthorizationInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *CheckAuthorizationInput) (*mcp.CallToolResult, any, error) {
		if in.PrincipalKind == "" || in.PrincipalID == "" || in.ResourceKind == "" || in.ResourceID == "" || in.Relation == "" {
			return nil, nil, fmt.Errorf("all fields (principal_kind, principal_id, resource_kind, resource_id, relation) are required")
		}
		text, err := CheckAuthorization(ctx, serverAddress, in.PrincipalKind, in.PrincipalID, in.ResourceKind, in.ResourceID, in.Relation)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_principals
// ---------------------------------------------------------------------------

type ListPrincipalsInput struct {
	OrgID         string `json:"org_id"                  jsonschema:"required,Organization ID."`
	Env           string `json:"env,omitempty"            jsonschema:"Optional environment slug. If omitted, returns org-level principals."`
	PrincipalKind string `json:"principal_kind"           jsonschema:"required,Type of principal to list: 'identity_account' or 'team'."`
	PageNumber    int32  `json:"page_number,omitempty"    jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize      int32  `json:"page_size,omitempty"      jsonschema:"Page size. Defaults to server default."`
}

func ListPrincipalsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_principals",
		Description: "List principals (users or teams) with access to an organization or environment. " +
			"Set principal_kind to 'identity_account' for users or 'team' for teams. " +
			"Optionally provide an environment slug to scope to that environment.",
	}
}

func ListPrincipalsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListPrincipalsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in *ListPrincipalsInput) (*mcp.CallToolResult, any, error) {
		if in.OrgID == "" {
			return nil, nil, fmt.Errorf("'org_id' is required")
		}
		if in.PrincipalKind == "" {
			return nil, nil, fmt.Errorf("'principal_kind' is required")
		}
		text, err := ListPrincipals(ctx, serverAddress, in.OrgID, in.Env, in.PrincipalKind, in.PageNumber, in.PageSize)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
