// Package infrahub groups MCP tool domains that belong to the Infrastructure
// Hub bounded context in the Planton API (apis/ai/planton/infrahub).
//
// Subpackages:
//   - cloudresource: lifecycle operations on cloud resources (apply, get, delete, destroy, list, rename, locks, etc.)
//   - infrachart:    infra chart template discovery and rendering (list, get, build)
//   - infrapipeline: pipeline observability and control (list, get, latest, run, cancel, gate resolution)
//   - infraproject:  infra project lifecycle management (search, get, apply, delete, slug check, undeploy)
//   - stackjob:      IaC stack job observability and lifecycle control (get, list, latest, rerun, cancel, resume, essentials)
//   - preset:        cloud object preset discovery (search, get)
package infrahub
