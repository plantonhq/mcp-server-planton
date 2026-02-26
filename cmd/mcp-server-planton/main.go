// Command mcp-server-planton is the Model Context Protocol server for the
// Planton Cloud platform.
//
// It exposes a set of tools that let MCP-capable AI clients (Cursor,
// Claude Desktop, Windsurf, etc.) create, update, read, and delete cloud
// resources managed by the Planton backend.
//
// # Usage
//
//	mcp-server-planton stdio       Start in stdio mode (stdin/stdout JSON-RPC)
//	mcp-server-planton http        Start in HTTP mode (Streamable HTTP)
//	mcp-server-planton both        Start both transports simultaneously
//
// When no subcommand is given, the transport is read from the
// PLANTON_MCP_TRANSPORT environment variable (default: "stdio").
//
// # Configuration
//
// All settings are read from environment variables — see the config package
// for the full list.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/plantoncloud/mcp-server-planton/pkg/mcpserver"
)

var validTransports = map[string]bool{
	"stdio": true,
	"http":  true,
	"both":  true,
}

func main() {
	cfg, err := mcpserver.DefaultConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		sub := os.Args[1]
		if !validTransports[sub] {
			fmt.Fprintf(os.Stderr, "unknown subcommand %q (expected stdio, http, or both)\n", sub)
			os.Exit(1)
		}
		cfg.Transport = sub
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := mcpserver.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}
