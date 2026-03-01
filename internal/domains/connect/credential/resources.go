package credential

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

const catalogURI = "credential-types://catalog"

// CatalogResource returns the static MCP resource definition for the
// credential types catalog. Agents read this resource to discover all
// supported credential types before fetching a per-type schema.
func CatalogResource() *mcp.Resource {
	return &mcp.Resource{
		URI:   catalogURI,
		Name:  "credential_types_catalog",
		Title: "Credential Types Catalog",
		Description: "Catalog of all supported credential types grouped by cloud provider. " +
			"Each provider entry includes the api_version and a sorted list of PascalCase kind " +
			"strings. Use these kind values with the credential-schema://{kind} resource " +
			"template to fetch the full spec schema for a specific type.",
		MIMEType: "application/json",
	}
}

// CatalogHandler returns a resource handler that serves the credential
// catalog JSON. Built once from the embedded registry and cached.
func CatalogHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		data, err := buildCredentialCatalog()
		if err != nil {
			return nil, fmt.Errorf("credential catalog: %w", err)
		}
		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}

// SchemaTemplate returns the MCP resource template for per-type schema
// discovery. Agents read a specific type's schema before calling
// apply_credential to learn the expected credential_object format.
func SchemaTemplate() *mcp.ResourceTemplate {
	return &mcp.ResourceTemplate{
		URITemplate: "credential-schema://{kind}",
		Name:        "credential_schema",
		Title:       "Credential Schema",
		Description: "JSON schema for a specific credential type. Returns the full spec " +
			"definition with field types, descriptions, and validation rules. " +
			"Example: credential-schema://AwsCredential",
		MIMEType: "application/json",
	}
}

// SchemaHandler returns a resource handler that serves per-type JSON
// schemas from the embedded filesystem.
func SchemaHandler() mcp.ResourceHandler {
	return func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		kind, err := parseSchemaURI(req.Params.URI)
		if err != nil {
			return nil, fmt.Errorf("credential schema: %w", err)
		}

		data, err := loadCredentialSchema(kind)
		if err != nil {
			return nil, err
		}

		return domains.ResourceResult(req.Params.URI, "application/json", string(data)), nil
	}
}
