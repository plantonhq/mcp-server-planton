// Package mcpserver provides the public API for embedding the Planton MCP
// server into other Go programs.
//
// The package exposes a minimal surface: [Config] holds runtime settings,
// [DefaultConfig] reads them from environment variables, and [Run] starts
// the server.
package mcpserver

import (
	"log/slog"
	"strings"

	"github.com/plantoncloud/mcp-server-planton/internal/config"
)

// Config holds runtime settings for the MCP server. All fields use plain
// Go types so that callers outside the mcp-server module can construct a
// Config without importing internal packages.
type Config struct {
	// ServerAddress is the gRPC dial target for the Planton backend
	// (e.g. "localhost:8080" or "api.live.planton.ai:443").
	ServerAddress string

	// APIKey authenticates calls to the backend. Required when Transport
	// is "stdio" or "both"; ignored for "http" (per-request Bearer tokens).
	APIKey string

	// Transport selects the communication mode: "stdio", "http", or "both".
	Transport string

	// HTTPPort is the TCP port for the HTTP transport (e.g. "8080").
	HTTPPort string

	// HTTPAuthEnabled controls whether HTTP requests require a valid
	// Authorization: Bearer token.
	HTTPAuthEnabled bool

	// LogFormat selects the structured log encoding: "text" or "json".
	LogFormat string

	// LogLevel sets the minimum log severity: "debug", "info", "warn",
	// or "error".
	LogLevel string
}

// DefaultConfig returns a Config populated from environment variables.
// It applies the same defaults and validation as the standalone binary.
//
// Environment variables (see internal/config for the full list):
//
//	PLANTON_API_KEY               – API key (required for stdio/both)
//	PLANTON_APIS_GRPC_ENDPOINT    – gRPC address override
//	PLANTON_CLOUD_ENVIRONMENT     – "live" | "test" | "local" (default "live")
//	PLANTON_MCP_TRANSPORT         – "stdio" | "http" | "both" (default "stdio")
//	PLANTON_MCP_HTTP_PORT         – HTTP listen port (default "8080")
//	PLANTON_MCP_HTTP_AUTH_ENABLED – "true" | "false" (default "true")
//	PLANTON_MCP_LOG_FORMAT        – "text" | "json" (default "text")
//	PLANTON_MCP_LOG_LEVEL         – "debug" | "info" | "warn" | "error" (default "info")
func DefaultConfig() (*Config, error) {
	ic, err := config.LoadFromEnv()
	if err != nil {
		return nil, err
	}
	return fromInternal(ic), nil
}

// fromInternal maps an internal config to the public representation.
func fromInternal(ic *config.Config) *Config {
	return &Config{
		ServerAddress:   ic.ServerAddress,
		APIKey:          ic.APIKey,
		Transport:       string(ic.Transport),
		HTTPPort:        ic.HTTPPort,
		HTTPAuthEnabled: ic.HTTPAuthEnabled,
		LogFormat:       string(ic.LogFormat),
		LogLevel:        logLevelString(ic.LogLevel),
	}
}

// toInternal converts the public Config to the internal representation,
// parsing string fields into typed values and validating invariants.
func (c *Config) toInternal() (*config.Config, error) {
	logLevel, err := config.ParseLogLevel(c.LogLevel)
	if err != nil {
		return nil, err
	}

	ic := &config.Config{
		ServerAddress:   c.ServerAddress,
		APIKey:          c.APIKey,
		Transport:       config.Transport(strings.ToLower(c.Transport)),
		HTTPPort:        c.HTTPPort,
		HTTPAuthEnabled: c.HTTPAuthEnabled,
		LogFormat:       config.LogFormat(strings.ToLower(c.LogFormat)),
		LogLevel:        logLevel,
	}

	if err := ic.Validate(); err != nil {
		return nil, err
	}
	return ic, nil
}

// logLevelString converts an slog.Level back to its human-friendly name.
func logLevelString(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return "info"
	}
}
