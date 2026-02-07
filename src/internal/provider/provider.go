package provider

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrCLIExecution = errors.New("CLI execution failed")
	ErrParsing      = errors.New("failed to parse response")
)

// Generator defines the interface for AI prompt generation providers.
type Generator interface {
	Generate(systemPrompt, userPrompt string) (string, error)
}

// CleanResponse removes any markdown code blocks or extra formatting from the response.
func CleanResponse(s string) string {
	// Remove markdown code blocks like ```lang ... ``` or ``` ... ```
	codeBlockRegex := regexp.MustCompile("(?s)```(?:\\w+)?\\s*(.+?)\\s*```")
	if matches := codeBlockRegex.FindStringSubmatch(s); len(matches) > 1 {
		s = matches[1]
	}

	return strings.TrimSpace(s)
}
