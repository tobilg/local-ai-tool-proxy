package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	tests := []struct {
		name string
		flag string
	}{
		{"long flag", "--version"},
		{"short flag", "-v"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", tc.flag)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "local-ai-tool-proxy") {
				t.Errorf("expected output to contain 'local-ai-tool-proxy', got %q", outputStr)
			}
			if !strings.Contains(outputStr, "dev") {
				t.Errorf("expected output to contain 'dev' (default version), got %q", outputStr)
			}
		})
	}
}

func TestUnknownProvider(t *testing.T) {
	// Create a temp system prompt file
	dir := t.TempDir()
	promptPath := filepath.Join(dir, "system-prompt.txt")
	if err := os.WriteFile(promptPath, []byte("You are a helpful assistant."), 0644); err != nil {
		t.Fatalf("failed to create temp system prompt: %v", err)
	}

	// Set required env vars
	os.Setenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT", promptPath)
	os.Setenv("LOCAL_AI_TOOL_PROXY_PROVIDER", "invalid-provider")
	defer func() {
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT")
		os.Unsetenv("LOCAL_AI_TOOL_PROXY_PROVIDER")
	}()

	cmd := exec.Command("go", "run", ".")
	cmd.Env = append(os.Environ(),
		"LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT="+promptPath,
		"LOCAL_AI_TOOL_PROXY_PROVIDER=invalid-provider",
	)
	output, err := cmd.CombinedOutput()

	// Should exit with error
	if err == nil {
		t.Fatal("expected command to fail with unknown provider")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Unknown provider") {
		t.Errorf("expected error about unknown provider, got %q", outputStr)
	}
	if !strings.Contains(outputStr, "invalid-provider") {
		t.Errorf("expected error to mention 'invalid-provider', got %q", outputStr)
	}
}
