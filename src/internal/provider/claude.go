package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

const claudeJSONSchema = `{"type":"object","properties":{"response":{"type":"string"}},"required":["response"]}`

// ClaudeClient implements Generator using the Claude CLI.
type ClaudeClient struct{}

// NewClaudeClient creates a new Claude CLI client.
func NewClaudeClient() *ClaudeClient {
	return &ClaudeClient{}
}

// Generate calls the Claude CLI with a system prompt and user prompt.
func (c *ClaudeClient) Generate(systemPrompt, userPrompt string) (string, error) {
	cmd := exec.Command("claude",
		"-p", userPrompt,
		"--append-system-prompt", systemPrompt,
		"--output-format", "json",
		"--json-schema", claudeJSONSchema,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Join(ErrCLIExecution, errors.New(stderr.String()))
	}

	result, err := parseClaudeResponse(stdout.Bytes())
	if err != nil {
		return "", err
	}

	return result, nil
}

// parseClaudeResponse extracts the response from Claude's JSON output.
func parseClaudeResponse(data []byte) (string, error) {
	// Claude returns: {"structured_output": {"response": "..."}, ...}
	var response struct {
		StructuredOutput struct {
			Response string `json:"response"`
		} `json:"structured_output"`
	}

	if err := json.Unmarshal(data, &response); err == nil && response.StructuredOutput.Response != "" {
		return response.StructuredOutput.Response, nil
	}

	// Fallback: try to extract raw content if structured parsing fails
	trimmed := strings.TrimSpace(string(data))
	if trimmed != "" {
		return trimmed, nil
	}

	return "", ErrParsing
}
