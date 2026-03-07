// Package runner provides MCP tools for managing RunnerRegistration resources.
// Runners are compute agents that execute infrastructure operations on behalf
// of the platform.
//
// Six tools are exposed:
//   - apply_runner_registration:       create or update
//   - get_runner_registration:         retrieve by ID
//   - delete_runner_registration:      delete by ID
//   - search_runner_registrations:     search within an organization
//   - generate_runner_credentials:     generate initial auth credentials (sensitive output)
//   - regenerate_runner_credentials:   rotate/regenerate auth credentials (sensitive output)
package runner
