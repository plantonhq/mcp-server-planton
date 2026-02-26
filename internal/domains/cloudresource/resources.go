package cloudresource

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantoncloud/mcp-server-planton/internal/domains"
)

const kindCatalogURI = "cloud-resource-kinds://catalog"

// KindCatalogResource returns the static MCP resource definition for the cloud
// resource kinds catalog. Agents read this resource to discover all supported
// provider kinds before fetching a per-kind schema or calling tools.
func KindCatalogResource() *mcp.Resource {
	return &mcp.Resource{
		URI:  kindCatalogURI,
		Name: "cloud_resource_kinds_catalog",
		Title: "Cloud Resource Kinds Catalog",
		Description: "Catalog of all supported cloud resource kinds grouped by cloud provider. " +
			"Each provider entry includes the api_version and a sorted list of PascalCase kind " +
			"strings. Use these kind values with the cloud-resource-schema://{kind} resource " +
			"template to fetch the full spec schema for a specific kind.",
		MIMEType: "application/json",
	}
}

// KindCatalogHandler returns a resource handler that serves the kind catalog
// JSON. The catalog is built once from the embedded provider registry and
// cached for the lifetime of the process.
func KindCatalogHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		data, err := buildKindCatalog()
		if err != nil {
			return nil, fmt.Errorf("kind catalog: %w", err)
		}
		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}

// SchemaTemplate returns the MCP resource template for per-kind schema
// discovery. Agents read a specific kind's schema before calling
// apply_cloud_resource to learn the expected cloud_object format.
func SchemaTemplate() *mcp.ResourceTemplate {
	return &mcp.ResourceTemplate{
		URITemplate: "cloud-resource-schema://{kind}",
		Name:        "cloud_resource_schema",
		Title:       "Cloud Resource Schema",
		Description: "JSON schema for a specific cloud resource kind. Returns the full spec " +
			"definition with field types, descriptions, validation rules, and defaults. " +
			"Example: cloud-resource-schema://AwsEksCluster",
		MIMEType: "application/json",
	}
}

// SchemaHandler returns a resource handler that serves per-kind JSON schemas
// from the embedded filesystem.
func SchemaHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		kind, err := parseSchemaURI(req.Params.URI)
		if err != nil {
			return nil, fmt.Errorf("cloud resource schema: %w", err)
		}

		data, err := loadProviderSchema(kind)
		if err != nil {
			return nil, err
		}

		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}
