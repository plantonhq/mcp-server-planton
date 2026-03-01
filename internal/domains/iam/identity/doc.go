// Package identity provides MCP tools for identity account lookup
// and user invitation management.
//
// Four tools are exposed:
//   - whoami:                 retrieve the currently authenticated user
//   - get_identity_account:   look up a user by ID or email
//   - invite_member:          invite a user to an organization with specified roles
//   - list_invitations:       list pending (or other status) invitations for an org
package identity
