package config

import (
	"log/slog"
	"strings"
	"testing"
)

func validConfig() *Config {
	return &Config{
		ServerAddress:   "localhost:8080",
		APIKey:          "test-key",
		Transport:       TransportStdio,
		HTTPPort:        "8080",
		HTTPAuthEnabled: true,
		LogFormat:       LogFormatText,
		LogLevel:        slog.LevelInfo,
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_AllTransports(t *testing.T) {
	for _, transport := range []Transport{TransportStdio, TransportHTTP, TransportBoth} {
		cfg := validConfig()
		cfg.Transport = transport
		if err := cfg.Validate(); err != nil {
			t.Errorf("transport %q: unexpected error: %v", transport, err)
		}
	}
}

func TestValidate_InvalidTransport(t *testing.T) {
	cfg := validConfig()
	cfg.Transport = "websocket"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid transport")
	}
	if !strings.Contains(err.Error(), "PLANTON_MCP_TRANSPORT") {
		t.Fatalf("expected transport error, got: %v", err)
	}
}

func TestValidate_EmptyServerAddress(t *testing.T) {
	cfg := validConfig()
	cfg.ServerAddress = ""
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty server address")
	}
	if !strings.Contains(err.Error(), "server address") {
		t.Fatalf("expected server address error, got: %v", err)
	}
}

func TestValidate_InvalidLogFormat(t *testing.T) {
	cfg := validConfig()
	cfg.LogFormat = "yaml"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid log format")
	}
	if !strings.Contains(err.Error(), "PLANTON_MCP_LOG_FORMAT") {
		t.Fatalf("expected log format error, got: %v", err)
	}
}

func TestValidate_StdioRequiresAPIKey(t *testing.T) {
	cfg := validConfig()
	cfg.Transport = TransportStdio
	cfg.APIKey = ""
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing API key with stdio transport")
	}
	if !strings.Contains(err.Error(), "PLANTON_API_KEY") {
		t.Fatalf("expected API key error, got: %v", err)
	}
}

func TestValidate_BothRequiresAPIKey(t *testing.T) {
	cfg := validConfig()
	cfg.Transport = TransportBoth
	cfg.APIKey = ""
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing API key with both transport")
	}
}

func TestValidate_HTTPDoesNotRequireAPIKey(t *testing.T) {
	cfg := validConfig()
	cfg.Transport = TransportHTTP
	cfg.APIKey = ""
	if err := cfg.Validate(); err != nil {
		t.Fatalf("HTTP transport should not require API key, got: %v", err)
	}
}

func TestParseLogLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"Error", slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseLogLevel(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseLogLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseLogLevel_Invalid(t *testing.T) {
	_, err := ParseLogLevel("trace")
	if err == nil {
		t.Fatal("expected error for invalid log level")
	}
	if !strings.Contains(err.Error(), "PLANTON_MCP_LOG_LEVEL") {
		t.Fatalf("expected log level error, got: %v", err)
	}
}
