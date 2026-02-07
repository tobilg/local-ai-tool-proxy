package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultPort          = 4000
	defaultAllowedOrigin = "http://localhost:3000"
	defaultProvider      = "claude"
)

// Config holds the application configuration.
type Config struct {
	Port             int
	AllowedOrigin    string
	Provider         string
	SystemPromptPath string
	SystemPrompt     string
	TLSCert          string
	TLSKey           string
}

// TLSEnabled returns true if both TLS cert and key are configured.
func (c Config) TLSEnabled() bool {
	return c.TLSCert != "" && c.TLSKey != ""
}

// Load loads configuration from environment variables with sensible defaults.
func Load() (Config, error) {
	cfg := Config{
		Port:          defaultPort,
		AllowedOrigin: defaultAllowedOrigin,
		Provider:      defaultProvider,
	}

	if portStr := os.Getenv("LOCAL_AI_TOOL_PROXY_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
			cfg.Port = port
		}
	}

	if origin := os.Getenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN"); origin != "" {
		cfg.AllowedOrigin = origin
	}

	if provider := os.Getenv("LOCAL_AI_TOOL_PROXY_PROVIDER"); provider != "" {
		cfg.Provider = provider
	}

	cfg.TLSCert = os.Getenv("LOCAL_AI_TOOL_PROXY_TLS_CERT")
	cfg.TLSKey = os.Getenv("LOCAL_AI_TOOL_PROXY_TLS_KEY")

	// System prompt file is required
	cfg.SystemPromptPath = os.Getenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	if cfg.SystemPromptPath == "" {
		return Config{}, fmt.Errorf("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT environment variable is required")
	}

	data, err := os.ReadFile(cfg.SystemPromptPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read system prompt file: %w", err)
	}

	cfg.SystemPrompt = strings.TrimSpace(string(data))
	if cfg.SystemPrompt == "" {
		return Config{}, fmt.Errorf("system prompt file is empty: %s", cfg.SystemPromptPath)
	}

	return cfg, nil
}
