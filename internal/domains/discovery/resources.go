package discovery

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/schemas"
)

const catalogURI = "api-resource-kinds://catalog"

// CatalogResource returns the static MCP resource definition for the
// platform-wide API resource kinds catalog. Agents read this resource to
// discover all resource types available on the platform, grouped by domain,
// before drilling into domain-specific catalogs or calling tools.
func CatalogResource() *mcp.Resource {
	return &mcp.Resource{
		URI:  catalogURI,
		Name: "api_resource_kinds_catalog",
		Title: "API Resource Kinds Catalog",
		Description: "Catalog of all Planton Cloud API resource types grouped by domain " +
			"(ResourceManager, InfraHub, ServiceHub, ConfigManager, Connect, IAM). " +
			"Each entry includes the snake_case kind name and display name. " +
			"For cloud infrastructure types, see the cloud-resource-kinds://catalog resource. " +
			"For credential types, see the credential-types://catalog resource.",
		MIMEType: "application/json",
	}
}

// CatalogHandler returns a resource handler that serves the API resource
// kinds catalog JSON from the embedded filesystem.
func CatalogHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		data, err := schemas.ApiResourceKindFS.ReadFile("apiresourcekinds/catalog.json")
		if err != nil {
			return nil, fmt.Errorf("api resource kinds catalog: %w", err)
		}
		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}
