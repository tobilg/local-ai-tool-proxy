package config

import (
	"os"
	"path/filepath"
	"testing"
)

// createTempSystemPrompt creates a temp file with the given content and returns its path.
func createTempSystemPrompt(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "system-prompt.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp system prompt: %v", err)
	}
	return path
}

// setRequiredEnv sets the system prompt env var to a valid temp file.
func setRequiredEnv(t *testing.T) {
	t.Helper()
	path := createTempSystemPrompt(t, "You are a helpful assistant.")
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", path)
}

func TestLoad_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_PORT")
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN")
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_PROVIDER")
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_TLS_CERT")
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_TLS_KEY")

	setRequiredEnv(t)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 4000 {
		t.Errorf("expected default port 4000, got %d", cfg.Port)
	}
	if cfg.AllowedOrigin != "http://localhost:3000" {
		t.Errorf("expected default origin http://localhost:3000, got %s", cfg.AllowedOrigin)
	}
	if cfg.Provider != "claude" {
		t.Errorf("expected default provider claude, got %s", cfg.Provider)
	}
	if cfg.SystemPrompt != "You are a helpful assistant." {
		t.Errorf("expected system prompt content, got %s", cfg.SystemPrompt)
	}
	if cfg.TLSCert != "" {
		t.Errorf("expected empty TLS cert, got %s", cfg.TLSCert)
	}
	if cfg.TLSKey != "" {
		t.Errorf("expected empty TLS key, got %s", cfg.TLSKey)
	}
	if cfg.TLSEnabled() {
		t.Error("expected TLS to be disabled by default")
	}
}

func TestLoad_CustomPort(t *testing.T) {
	setRequiredEnv(t)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	os.Setenv("LOCAL_AI_TOOL_PROXY_PORT", "8080")
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	tests := []struct {
		name     string
		portVal  string
		expected int
	}{
		{"non-numeric", "abc", 4000},
		{"zero", "0", 4000},
		{"negative", "-1", 4000},
		{"too high", "65536", 4000},
		{"empty", "", 4000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setRequiredEnv(t)
			defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

			if tc.portVal == "" {
				os.Unsetenv("LOCAL_AI_TOOL_PROXY_PORT")
			} else {
				os.Setenv("LOCAL_AI_TOOL_PROXY_PORT", tc.portVal)
				defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_PORT")
			}

			cfg, err := Load()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Port != tc.expected {
				t.Errorf("expected port %d for %q, got %d", tc.expected, tc.portVal, cfg.Port)
			}
		})
	}
}

func TestLoad_CustomAllowedOrigin(t *testing.T) {
	setRequiredEnv(t)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	os.Setenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN", "https://example.com")
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AllowedOrigin != "https://example.com" {
		t.Errorf("expected origin https://example.com, got %s", cfg.AllowedOrigin)
	}
}

func TestLoad_CustomProvider(t *testing.T) {
	setRequiredEnv(t)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	os.Setenv("LOCAL_AI_TOOL_PROXY_PROVIDER", "gemini")
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_PROVIDER")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Provider != "gemini" {
		t.Errorf("expected provider gemini, got %s", cfg.Provider)
	}
}

func TestLoad_AllCustomValues(t *testing.T) {
	path := createTempSystemPrompt(t, "Custom system prompt.")
	os.Setenv("LOCAL_AI_TOOL_PROXY_PORT", "3000")
	os.Setenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN", "https://myapp.com")
	os.Setenv("LOCAL_AI_TOOL_PROXY_PROVIDER", "codex")
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", path)
	defer func() {
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_PORT")
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN")
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_PROVIDER")
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}
	if cfg.AllowedOrigin != "https://myapp.com" {
		t.Errorf("expected origin https://myapp.com, got %s", cfg.AllowedOrigin)
	}
	if cfg.Provider != "codex" {
		t.Errorf("expected provider codex, got %s", cfg.Provider)
	}
	if cfg.SystemPrompt != "Custom system prompt." {
		t.Errorf("expected system prompt 'Custom system prompt.', got %s", cfg.SystemPrompt)
	}
}

func TestLoad_TLSConfig(t *testing.T) {
	setRequiredEnv(t)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
	os.Setenv("LOCAL_AI_TOOL_PROXY_TLS_CERT", "/path/to/cert.pem")
	os.Setenv("LOCAL_AI_TOOL_PROXY_TLS_KEY", "/path/to/key.pem")
	defer func() {
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_TLS_CERT")
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_TLS_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TLSCert != "/path/to/cert.pem" {
		t.Errorf("expected TLS cert /path/to/cert.pem, got %s", cfg.TLSCert)
	}
	if cfg.TLSKey != "/path/to/key.pem" {
		t.Errorf("expected TLS key /path/to/key.pem, got %s", cfg.TLSKey)
	}
	if !cfg.TLSEnabled() {
		t.Error("expected TLS to be enabled when both cert and key are set")
	}
}

func TestTLSEnabled_OnlyCert(t *testing.T) {
	cfg := Config{TLSCert: "/path/to/cert.pem", TLSKey: ""}
	if cfg.TLSEnabled() {
		t.Error("TLS should not be enabled with only cert")
	}
}

func TestTLSEnabled_OnlyKey(t *testing.T) {
	cfg := Config{TLSCert: "", TLSKey: "/path/to/key.pem"}
	if cfg.TLSEnabled() {
		t.Error("TLS should not be enabled with only key")
	}
}

func TestTLSEnabled_BothSet(t *testing.T) {
	cfg := Config{TLSCert: "/path/to/cert.pem", TLSKey: "/path/to/key.pem"}
	if !cfg.TLSEnabled() {
		t.Error("TLS should be enabled when both cert and key are set")
	}
}

func TestLoad_MissingSystemPrompt(t *testing.T) {
	os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when system prompt env var is missing")
	}
}

func TestLoad_SystemPromptFileNotFound(t *testing.T) {
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", "/nonexistent/path/prompt.txt")
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when system prompt file does not exist")
	}
}

func TestLoad_EmptySystemPromptFile(t *testing.T) {
	path := createTempSystemPrompt(t, "")
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", path)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when system prompt file is empty")
	}
}

func TestLoad_SystemPromptFileSuccess(t *testing.T) {
	path := createTempSystemPrompt(t, "You are a SQL expert. Generate DuckDB queries.")
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", path)
	defer os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SystemPrompt != "You are a SQL expert. Generate DuckDB queries." {
		t.Errorf("unexpected system prompt: %s", cfg.SystemPrompt)
	}
	if cfg.SystemPromptPath != path {
		t.Errorf("expected system prompt path %s, got %s", path, cfg.SystemPromptPath)
	}
}
