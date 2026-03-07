// Package iam is the parent package for the IAM bounded context,
// which manages identity, access control, team membership, API keys,
// service accounts, and role definitions.
//
// Sub-packages:
//   - identity:       Identity account lookup, whoami, and user invitations (4 tools)
//   - team:           Team CRUD (4 tools)
//   - policy:         IAM policy v2 access control (7 tools)
//   - role:           IAM role read-only lookups (2 tools)
//   - apikey:         API key lifecycle management (3 tools)
//   - serviceaccount: Machine identity CRUD and key management (8 tools)
package iam
