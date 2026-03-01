// Package policy provides MCP tools for IAM policy v2 access control.
//
// Seven tools are exposed:
//   - create_iam_policy:      grant a principal access to a resource with a specific relation
//   - delete_iam_policy:      revoke a specific access grant
//   - upsert_iam_policies:    declaratively sync a principal's relations on a resource
//   - revoke_org_access:      remove all of a user's access to an organization
//   - list_resource_access:   list who has access to a resource
//   - check_authorization:    check if a principal is authorized for a specific relation
//   - list_principals:        list principals (users or teams) in an org or environment
package policy
