// Package infrahub groups MCP tool domains that belong to the Infrastructure
// Hub bounded context in the Planton API (apis/ai/planton/infrahub).
//
// Subpackages:
//   - cloudresource: lifecycle operations on cloud resources (apply, get, delete, destroy, list, rename, locks, etc.)
//   - infrachart:    infra chart template discovery and rendering (list, get, build)
//   - stackjob:      observability for IaC stack jobs (get, list, latest)
//   - preset:        cloud object preset discovery (search, get)
package infrahub
