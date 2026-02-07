package provider

import (
	"testing"
)

func TestCleanResponse_RemovesMarkdownCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sql code block",
			input:    "```sql\nSELECT * FROM users\n```",
			expected: "SELECT * FROM users",
		},
		{
			name:     "plain code block",
			input:    "```\nSELECT * FROM users\n```",
			expected: "SELECT * FROM users",
		},
		{
			name:     "no code block",
			input:    "SELECT * FROM users",
			expected: "SELECT * FROM users",
		},
		{
			name:     "with extra whitespace",
			input:    "  SELECT * FROM users  ",
			expected: "SELECT * FROM users",
		},
		{
			name:     "with backticks",
			input:    "SELECT COUNT(*) FROM `aws_iam`.actions",
			expected: "SELECT COUNT(*) FROM `aws_iam`.actions",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CleanResponse(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
