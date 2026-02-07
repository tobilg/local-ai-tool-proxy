package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// ContinueClient implements Generator using the Continue CLI (cn).
type ContinueClient struct{}

// NewContinueClient creates a new Continue CLI client.
func NewContinueClient() *ContinueClient {
	return &ContinueClient{}
}

// Generate calls the Continue CLI with a system prompt and user prompt.
func (c *ContinueClient) Generate(systemPrompt, userPrompt string) (string, error) {
	prompt := systemPrompt + "\n\n" + userPrompt

	cmd := exec.Command("cn",
		"-p", prompt,
		"--format", "json",
		"--silent",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Join(ErrCLIExecution, errors.New(stderr.String()))
	}

	result, err := parseContinueResponse(stdout.Bytes())
	if err != nil {
		return "", err
	}

	return result, nil
}

// continueResponse represents the JSON response from Continue CLI.
type continueResponse struct {
	Response string `json:"response"`
	Status   string `json:"status"`
}

// parseContinueResponse extracts the response from Continue CLI's JSON output.
func parseContinueResponse(data []byte) (string, error) {
	// Continue CLI wraps plain text responses in:
	// {"response": "...", "status": "success", "note": "..."}

	var response continueResponse
	if err := json.Unmarshal(data, &response); err == nil && response.Response != "" {
		return CleanResponse(response.Response), nil
	}

	// Fallback: try to extract raw content if JSON parsing fails
	trimmed := strings.TrimSpace(string(data))
	if trimmed != "" {
		return CleanResponse(trimmed), nil
	}

	return "", ErrParsing
}
