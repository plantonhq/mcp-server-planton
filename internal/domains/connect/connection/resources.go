package connection

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

const catalogURI = "connection-types://catalog"

// CatalogResource returns the MCP resource definition for the connection types catalog.
func CatalogResource() *mcp.Resource {
	return &mcp.Resource{
		URI:      catalogURI,
		Name:     "connection_types_catalog",
		Title:    "Connection Types Catalog",
		MIMEType: "application/json",
		Description: "JSON catalog of all supported connection types, grouped by cloud provider. " +
			"Each entry includes the API version and available kinds. " +
			"Use the schema_uri_template to fetch per-type schema details.",
	}
}

// CatalogHandler returns the MCP resource handler for the connection types catalog.
func CatalogHandler() mcp.ResourceHandler {
	return func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		data, err := buildConnectionCatalog()
		if err != nil {
			return nil, fmt.Errorf("building connection catalog: %w", err)
		}
		return domains.ResourceResult(catalogURI, "application/json", string(data)), nil
	}
}

// SchemaTemplate returns the MCP resource template for per-type connection schemas.
func SchemaTemplate() *mcp.ResourceTemplate {
	return &mcp.ResourceTemplate{
		URITemplate: "connection-schema://{kind}",
		Name:        "connection_schema",
		Title:       "Connection Schema",
		MIMEType:    "application/json",
		Description: "JSON schema describing the spec fields for a specific connection type. " +
			"Replace {kind} with the PascalCase connection kind (e.g. AwsProviderConnection). " +
			"Discover available kinds from connection-types://catalog.",
	}
}

// SchemaHandler returns the MCP resource handler for per-type connection schemas.
func SchemaHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		kind, err := parseSchemaURI(req.Params.URI)
		if err != nil {
			return nil, err
		}
		data, err := loadConnectionSchema(kind)
		if err != nil {
			return nil, err
		}
		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}
