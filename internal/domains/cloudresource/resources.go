package cloudresource

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantoncloud/mcp-server-planton/internal/domains"
)

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
