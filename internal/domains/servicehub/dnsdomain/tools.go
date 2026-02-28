// Package dnsdomain provides the MCP tools for the ServiceHub DnsDomain
// domain, backed by the DnsDomainQueryController and
// DnsDomainCommandController (ai.planton.servicehub.dnsdomain.v1) RPCs on the
// Planton backend.
//
// Three tools are exposed:
//   - get_dns_domain:    retrieve a DNS domain by ID or org+slug
//   - apply_dns_domain:  create or update a DNS domain (idempotent)
//   - delete_dns_domain: remove a DNS domain record
package dnsdomain

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// get_dns_domain
// ---------------------------------------------------------------------------

// GetDnsDomainInput defines the parameters for the get_dns_domain tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetDnsDomainInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The DNS domain ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"DNS domain slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_dns_domain.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_dns_domain",
		Description: "Retrieve the full details of a DNS domain by its ID or by org+slug. " +
			"A DNS domain registers a domain name within an organization for use by services " +
			"that need custom ingress hostnames. " +
			"Returns the complete domain including metadata, spec (domain_name, description), " +
			"and audit status. " +
			"The output JSON can be modified and passed to apply_dns_domain for updates.",
	}
}

// GetHandler returns the typed tool handler for get_dns_domain.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetDnsDomainInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetDnsDomainInput) (*mcp.CallToolResult, any, error) {
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
// apply_dns_domain
// ---------------------------------------------------------------------------

// ApplyDnsDomainInput defines the parameters for the apply_dns_domain tool.
type ApplyDnsDomainInput struct {
	DnsDomain map[string]any `json:"dns_domain" jsonschema:"required,The full DnsDomain resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('DnsDomain'), 'metadata' (with 'name' and 'org'), and 'spec' (with 'domain_name'). The output of get_dns_domain can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_dns_domain.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_dns_domain",
		Description: "Create or update a DNS domain (idempotent). " +
			"Accepts the full DnsDomain resource as a JSON object. " +
			"A DNS domain registers a domain name within an organization, " +
			"making it available for services to use in their ingress configuration. " +
			"For new domains, provide api_version, kind, metadata (name, org), and spec (domain_name). " +
			"For updates, retrieve the domain with get_dns_domain, modify the desired fields, and pass the result here. " +
			"Returns the applied domain with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_dns_domain.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyDnsDomainInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyDnsDomainInput) (*mcp.CallToolResult, any, error) {
		if len(input.DnsDomain) == 0 {
			return nil, nil, fmt.Errorf("'dns_domain' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.DnsDomain)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_dns_domain
// ---------------------------------------------------------------------------

// DeleteDnsDomainInput defines the parameters for the delete_dns_domain tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteDnsDomainInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The DNS domain ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"DNS domain slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_dns_domain.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_dns_domain",
		Description: "Delete a DNS domain record from the platform. " +
			"WARNING: Services using this domain for ingress hostnames will lose their custom domain routing. " +
			"Ensure no services reference this domain before deleting. " +
			"Identify the domain by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_dns_domain.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteDnsDomainInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteDnsDomainInput) (*mcp.CallToolResult, any, error) {
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
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the DNS domain")
	}
}
