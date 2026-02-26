// Package cloudresource provides the MCP tool and resource template for the
// CloudResource domain, backed by the CloudResourceCommandController and
// CloudResourceQueryController RPCs on the Planton backend.
//
// The apply_cloud_resource tool accepts an opaque cloud_object map. Typed
// validation is performed inside the handler using the generated parsers in
// gen/cloudresource/, so the tool schema stays small regardless of how many
// provider kinds exist.
package cloudresource

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantoncloud/mcp-server-planton/gen/cloudresource"
	"github.com/plantoncloud/mcp-server-planton/internal/domains"
)

// ApplyCloudResourceInput defines the parameters for the apply_cloud_resource tool.
type ApplyCloudResourceInput struct {
	CloudObject map[string]any `json:"cloud_object" jsonschema:"required,The full OpenMCF cloud resource object. Must contain api_version, kind, metadata (with name, org, env), and spec. Read cloud-resource-kinds://catalog for available kinds, then cloud-resource-schema://{kind} for the spec format."`
}

// ApplyTool returns the MCP tool definition for apply_cloud_resource.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_cloud_resource",
		Description: "Create or update a cloud resource on the Planton platform (idempotent). " +
			"The cloud_object must follow the OpenMCF format: " +
			"{ api_version, kind, metadata: { name, org, env }, spec: { ... } }. " +
			"Step 1: Read cloud-resource-kinds://catalog to discover supported kinds and api_versions. " +
			"Step 2: Read cloud-resource-schema://{kind} to get the full spec definition. " +
			"Step 3: Call this tool with the assembled cloud_object.",
	}
}

// ApplyHandler returns the typed tool handler for apply_cloud_resource.
// serverAddress is captured at registration time; the API key is read from
// context at call time.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyCloudResourceInput) (*mcp.CallToolResult, any, error) {
		co := input.CloudObject
		if co == nil {
			return nil, nil, fmt.Errorf("cloud_object is required")
		}

		kindStr, err := extractKindFromCloudObject(co)
		if err != nil {
			return nil, nil, err
		}

		parseFn, ok := cloudresource.GetParser(kindStr)
		if !ok {
			return nil, nil, fmt.Errorf(
				"unsupported cloud resource kind %q — read cloud-resource-kinds://catalog for all valid kinds",
				kindStr,
			)
		}

		normalizedObject, err := parseFn(co)
		if err != nil {
			return nil, nil, fmt.Errorf("cloud_object validation failed: %w", err)
		}

		cr, err := buildCloudResource(co, kindStr, normalizedObject)
		if err != nil {
			return nil, nil, err
		}

		text, err := Apply(ctx, serverAddress, cr)
		if err != nil {
			return nil, nil, err
		}

		return domains.TextResult(text)
	}
}
