// Package schemas embeds JSON schemas for cloud resource providers,
// connection types, and the platform-wide API resource kind catalog.
// Provider schemas are produced by the proto2schema codegen tool; connection
// schemas and the resource kind catalog are hand-crafted. All are consumed at
// runtime by MCP resource handlers for agent discovery.
package schemas

import "embed"

// FS contains the cloud resource provider schemas.
//
// Directory layout:
//
//	providers/registry.json       — kind-to-schema-path index
//	providers/{cloud}/{kind}.json — per-provider JSON schema
//	shared/metadata.json          — shared metadata field definitions
//
//go:embed providers shared
var FS embed.FS

// ConnectionFS contains the connection type schemas.
//
// Directory layout:
//
//	connections/registry.json       — kind-to-schema-path index
//	connections/{kind}.json         — per-type JSON schema
//
//go:embed connections
var ConnectionFS embed.FS

// ApiResourceKindFS contains the platform-wide API resource kind catalog.
//
// Directory layout:
//
//	apiresourcekinds/catalog.json — all resource kinds grouped by domain
//
//go:embed apiresourcekinds
var ApiResourceKindFS embed.FS
