// Package variablegroup provides MCP tools for managing VariableGroup
// resources, which bundle related configuration variables into a single
// named, scoped group with optional external source references.
//
// Eight tools are exposed:
//   - apply_variable_group:                create or update a group (envelope)
//   - get_variable_group:                  retrieve by ID or by org+scope+slug
//   - delete_variable_group:               delete by ID or by org+scope+slug
//   - upsert_variable_group_entry:         add or update a single entry
//   - delete_variable_group_entry:         remove a single entry
//   - refresh_variable_group_entry:        refresh one entry from its source
//   - refresh_all_variable_group_entries:  refresh all sourced entries
//   - resolve_variable_group_entry:        quick value lookup (returns plain string)
package variablegroup
