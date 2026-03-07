// Package search provides MCP tools that expose cross-domain search
// capabilities backed by the various *SearchQueryController gRPC services
// (ai.planton.search.v1.*) on the Planton backend.
//
// All tools in this package are read-only (query-side). They never mutate
// state; use the corresponding domain command tools for write operations.
package search
