// Package defaultrunner provides MCP tools for managing DefaultRunnerBinding
// resources, which designate a default runner for an organization.
//
// Six tools are exposed:
//   - apply_default_runner_binding:                create or update
//   - get_default_runner_binding:                  retrieve by ID
//   - get_default_runner_binding_by_selector:      retrieve by resource selector (kind + ID)
//   - resolve_default_runner_binding:              resolve the effective default for an org
//   - delete_default_runner_binding:               delete by ID
//   - delete_default_runner_binding_by_selector:   delete by resource selector (kind + ID)
package defaultrunner
