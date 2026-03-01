// Package prompts provides cross-domain MCP prompt templates that encode
// multi-step platform workflows. Unlike tools (which execute single operations)
// and resources (which serve static data), prompts are conversation starters —
// pre-built message templates that guide an LLM through a recommended sequence
// of tool calls to accomplish a complex goal.
//
// Each prompt is purely static: handlers interpolate string arguments into
// pre-written text without making gRPC calls or accessing external state.
// This keeps prompts fast, testable, and free of runtime failure modes.
//
// Five prompts are exposed:
//   - debug_failed_deployment:   diagnose a failed infrastructure deployment
//   - assess_change_impact:      analyze blast radius before destructive changes
//   - explore_infrastructure:    top-down discovery of an organization's topology
//   - provision_cloud_resource:  guided creation and deployment of cloud infrastructure
//   - manage_access:             IAM discovery, audit, and policy management
package prompts
