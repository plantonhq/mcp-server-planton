// Package schemas embeds JSON schemas for cloud resource providers and
// credential types. Provider schemas are produced by the proto2schema codegen
// tool; credential schemas are hand-crafted from proto spec definitions.
// Both are consumed at runtime by MCP resource template handlers for
// per-kind/per-type schema discovery.
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
