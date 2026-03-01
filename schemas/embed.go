// Package schemas embeds JSON schemas for cloud resource providers,
// credential types, and the platform-wide API resource kind catalog.
// Provider schemas are produced by the proto2schema codegen tool; credential
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

// CredentialFS contains the credential type schemas.
//
// Directory layout:
//
//	credentials/registry.json       — kind-to-schema-path index
//	credentials/{kind}.json         — per-type JSON schema
//
//go:embed credentials
var CredentialFS embed.FS

// ApiResourceKindFS contains the platform-wide API resource kind catalog.
//
// Directory layout:
//
//	apiresourcekinds/catalog.json — all resource kinds grouped by domain
//
//go:embed apiresourcekinds
var ApiResourceKindFS embed.FS
