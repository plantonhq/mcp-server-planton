// Package role provides read-only MCP tools for IAM role lookups.
//
// IAM roles are reference data managed by platform administrators.
// These tools allow querying available roles for granting access to resources.
//
// Two tools are exposed:
//   - get_iam_role:                       retrieve a role definition by ID
//   - list_iam_roles_for_resource_kind:   list available roles for a resource kind and principal type
package role
