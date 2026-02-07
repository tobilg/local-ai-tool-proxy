package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// GeminiClient implements Generator using the Gemini CLI.
type GeminiClient struct{}

// NewGeminiClient creates a new Gemini CLI client.
func NewGeminiClient() *GeminiClient {
	return &GeminiClient{}
}

// Generate calls the Gemini CLI with a system prompt and user prompt.
func (g *GeminiClient) Generate(systemPrompt, userPrompt string) (string, error) {
	prompt := systemPrompt + "\n\n" + userPrompt

	cmd := exec.Command("gemini",
		"-p", prompt,
		"--output-format", "json",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Join(ErrCLIExecution, errors.New(stderr.String()))
	}

	result, err := parseGeminiResponse(stdout.Bytes())
	if err != nil {
		return "", err
	}

	return result, nil
}

// parseGeminiResponse extracts the response from Gemini's JSON output.
func parseGeminiResponse(data []byte) (string, error) {
	// Gemini returns: {"response": "...", ...}
	var response struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(data, &response); err == nil && response.Response != "" {
		return CleanResponse(response.Response), nil
	}

	// Fallback: try to extract raw content if structured parsing fails
	trimmed := strings.TrimSpace(string(data))
	if trimmed != "" {
		return CleanResponse(trimmed), nil
	}

	return "", ErrParsing
}
