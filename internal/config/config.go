// Package config provides environment-variable-based configuration for
// mcp-server-planton.
//
// Every configurable value is read from an environment variable with a
// PLANTON_ prefix. Reasonable defaults are provided for development use.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Transport enumerates the supported communication modes between MCP clients
// and the MCP server.
type Transport string

const (
	// TransportStdio communicates over stdin/stdout. This is the primary mode
	// for local development: the MCP client (Cursor, Claude Desktop, etc.)
	// spawns the server as a child process.
	TransportStdio Transport = "stdio"

	// TransportHTTP serves MCP over Streamable HTTP. This is the mode for
	// remote / shared deployments where multiple users connect over the
	// network. Each request carries its own Bearer token.
	TransportHTTP Transport = "http"

	// TransportBoth runs STDIO and HTTP simultaneously. Useful for
	// development environments where you want both local and remote access.
	TransportBoth Transport = "both"
)

// LogFormat selects the structured log output encoding.
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// Well-known Planton API endpoints.
const (
	EndpointLive  = "api.live.planton.ai:443"
	EndpointTest  = "api.test.planton.cloud:443"
	EndpointLocal = "localhost:8080"
)

// Config holds all runtime configuration.
type Config struct {
	// ServerAddress is the gRPC dial target for the Planton backend
	// (e.g. "localhost:8080" or "api.live.planton.ai:443").
	ServerAddress string

	// APIKey optionally authenticates the MCP server's calls to the backend.
	// In STDIO mode this is loaded once from the environment at startup.
	// In HTTP mode every inbound request carries its own key via the
	// Authorization header, so this field is only used for STDIO.
	// When targeting an unauthenticated backend (e.g. the local dev server),
	// this may be empty.
	APIKey string

	// Transport selects the communication mode: stdio, http, or both.
	Transport Transport

	// HTTPPort is the TCP port the HTTP transport listens on.
	HTTPPort string

	// HTTPAuthEnabled controls whether HTTP requests require a valid
	// Authorization: Bearer token. Defaults to true.
	HTTPAuthEnabled bool

	// LogFormat controls the structured log encoding: "text" or "json".
	LogFormat LogFormat

	// LogLevel controls the minimum severity for emitted log records.
	LogLevel slog.Level
}

// LoadFromEnv reads configuration from the process environment.
//
// Environment variables:
//
//	PLANTON_API_KEY               – API key (required for stdio/both; per-request for http)
//	PLANTON_APIS_GRPC_ENDPOINT    – gRPC address override (takes precedence over environment preset)
//	PLANTON_CLOUD_ENVIRONMENT     – "live" | "test" | "local" (default "live")
//	PLANTON_MCP_TRANSPORT         – "stdio" | "http" | "both" (default "stdio")
//	PLANTON_MCP_HTTP_PORT         – HTTP listen port (default "8080")
//	PLANTON_MCP_HTTP_AUTH_ENABLED – "true" | "false" (default "true")
//	PLANTON_MCP_LOG_FORMAT        – "text" | "json" (default "text")
//	PLANTON_MCP_LOG_LEVEL         – "debug" | "info" | "warn" | "error" (default "info")
func LoadFromEnv() (*Config, error) {
	logLevel, err := ParseLogLevel(envOr("PLANTON_MCP_LOG_LEVEL", "info"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerAddress:   resolveEndpoint(),
		APIKey:          os.Getenv("PLANTON_API_KEY"),
		Transport:       Transport(strings.ToLower(envOr("PLANTON_MCP_TRANSPORT", "stdio"))),
		HTTPPort:        envOr("PLANTON_MCP_HTTP_PORT", "8080"),
		HTTPAuthEnabled: envOr("PLANTON_MCP_HTTP_AUTH_ENABLED", "true") == "true",
		LogFormat:       LogFormat(strings.ToLower(envOr("PLANTON_MCP_LOG_FORMAT", "text"))),
		LogLevel:        logLevel,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks invariants that should hold before the server starts.
func (c *Config) Validate() error {
	switch c.Transport {
	case TransportStdio, TransportHTTP, TransportBoth:
	default:
		return fmt.Errorf("invalid PLANTON_MCP_TRANSPORT %q: must be stdio, http, or both", c.Transport)
	}

	if c.ServerAddress == "" {
		return fmt.Errorf("server address must not be empty — set PLANTON_APIS_GRPC_ENDPOINT or PLANTON_CLOUD_ENVIRONMENT")
	}

	switch c.LogFormat {
	case LogFormatText, LogFormatJSON:
	default:
		return fmt.Errorf("invalid PLANTON_MCP_LOG_FORMAT %q: must be text or json", c.LogFormat)
	}

	needsKey := c.Transport == TransportStdio || c.Transport == TransportBoth
	if needsKey && c.APIKey == "" {
		return fmt.Errorf("PLANTON_API_KEY is required when transport is %q", c.Transport)
	}

	return nil
}

// ParseLogLevel converts a human-friendly level name to slog.Level.
func ParseLogLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid PLANTON_MCP_LOG_LEVEL %q: must be debug, info, warn, or error", s)
	}
}

// resolveEndpoint determines the gRPC endpoint using the following priority:
//  1. PLANTON_APIS_GRPC_ENDPOINT (explicit override)
//  2. PLANTON_CLOUD_ENVIRONMENT  (preset: live, test, local)
//  3. Default to the live endpoint
func resolveEndpoint() string {
	if ep := os.Getenv("PLANTON_APIS_GRPC_ENDPOINT"); ep != "" {
		return ep
	}
	switch strings.ToLower(os.Getenv("PLANTON_CLOUD_ENVIRONMENT")) {
	case "test":
		return EndpointTest
	case "local":
		return EndpointLocal
	default:
		return EndpointLive
	}
}

// envOr returns the value of the given environment variable, or fallback if
// the variable is unset or empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
