package provider

import (
	"testing"
)

func TestParseClaudeResponse_StructuredJSON(t *testing.T) {
	input := `{"structured_output":{"response":"SELECT * FROM users"},"type":"result"}`

	result, err := parseClaudeResponse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestParseClaudeResponse_RawText(t *testing.T) {
	input := `SELECT * FROM users WHERE id = 1`

	result, err := parseClaudeResponse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != input {
		t.Errorf("expected %q, got %q", input, result)
	}
}

func TestParseClaudeResponse_EmptyResponse(t *testing.T) {
	input := ``

	_, err := parseClaudeResponse([]byte(input))
	if err != ErrParsing {
		t.Errorf("expected ErrParsing, got %v", err)
	}
}

func TestParseClaudeResponse_WhitespaceOnly(t *testing.T) {
	input := `   `

	_, err := parseClaudeResponse([]byte(input))
	if err != ErrParsing {
		t.Errorf("expected ErrParsing, got %v", err)
	}
}

func TestParseClaudeResponse_NestedJSON(t *testing.T) {
	input := `{"structured_output":{"response":"SELECT json_extract(data, '$.name') FROM users"}}`

	result, err := parseClaudeResponse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT json_extract(data, '$.name') FROM users"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestParseClaudeResponse_EmptyStructuredResponse(t *testing.T) {
	input := `{"structured_output":{"response":""}}`

	result, err := parseClaudeResponse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Falls back to trimmed raw content
	if result != `{"structured_output":{"response":""}}` {
		t.Errorf("unexpected fallback: %q", result)
	}
}

func TestParseClaudeResponse_FullResponse(t *testing.T) {
	// Test with a more complete response like the actual CLI returns
	input := `{"type":"result","subtype":"success","structured_output":{"response":"SELECT * FROM users WHERE name LIKE 'A%';"},"session_id":"abc123"}`

	result, err := parseClaudeResponse([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users WHERE name LIKE 'A%';"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
