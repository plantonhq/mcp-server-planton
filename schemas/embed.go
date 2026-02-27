// Package schemas embeds the generated JSON schemas for all OpenMCF cloud
// resource providers. These schemas are produced by the proto2schema codegen
// tool (Stage 1) and consumed at runtime by MCP resource template handlers
// for per-kind schema discovery.
package schemas

import "embed"

// FS contains the generated provider schemas and the provider registry.
//
// Directory layout:
//
//	providers/registry.json       — kind-to-schema-path index
//	providers/{cloud}/{kind}.json — per-provider JSON schema
//	shared/metadata.json          — shared metadata field definitions
//
//go:embed providers shared
var FS embed.FS
