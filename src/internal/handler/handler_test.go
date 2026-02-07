package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tobilg/local-ai-tool-proxy/src/internal/provider"
)

// mockGenerator implements provider.Generator for testing.
type mockGenerator struct {
	response string
	err      error
}

func (m *mockGenerator) Generate(systemPrompt, userPrompt string) (string, error) {
	return m.response, m.err
}

func newTestHandler(mock *mockGenerator) *Handler {
	providers := map[string]provider.Generator{
		"claude": mock,
	}
	return New(providers, "claude", "http://localhost:3000", "You are a test assistant.")
}

func newTestHandlerWithProviders(providers map[string]provider.Generator, defaultProvider string) *Handler {
	return New(providers, defaultProvider, "http://localhost:3000", "You are a test assistant.")
}

func TestHandlePrompt_CORSHeaders(t *testing.T) {
	handler := newTestHandler(&mockGenerator{response: "Hello"})

	req := httptest.NewRequest(http.MethodOptions, "/prompt", nil)
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":          "http://localhost:3000",
		"Access-Control-Allow-Methods":         "POST, GET, OPTIONS",
		"Access-Control-Allow-Headers":         "Content-Type",
		"Access-Control-Allow-Private-Network": "true",
	}

	for header, expected := range expectedHeaders {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("header %s: expected %q, got %q", header, expected, got)
		}
	}
}

func TestHandlePrompt_OptionsRequest(t *testing.T) {
	handler := newTestHandler(&mockGenerator{})

	req := httptest.NewRequest(http.MethodOptions, "/prompt", nil)
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", w.Code)
	}
}

func TestHandlePrompt_InvalidJSON(t *testing.T) {
	handler := newTestHandler(&mockGenerator{})

	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error != "Invalid JSON" {
		t.Errorf("expected error 'Invalid JSON', got %q", resp.Error)
	}
}

func TestHandlePrompt_MissingFields(t *testing.T) {
	handler := newTestHandler(&mockGenerator{})

	body, _ := json.Marshal(Request{})
	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error != "The 'user' field is required" {
		t.Errorf("unexpected error: %q", resp.Error)
	}
}

func TestHandlePrompt_Success(t *testing.T) {
	expectedResponse := "Hello, world!"
	handler := newTestHandler(&mockGenerator{response: expectedResponse})

	body, _ := json.Marshal(Request{
		User: "Say hello",
	})
	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ResponseText != expectedResponse {
		t.Errorf("expected response %q, got %q", expectedResponse, resp.ResponseText)
	}
	if resp.Error != "" {
		t.Errorf("unexpected error: %q", resp.Error)
	}
}

func TestHandlePrompt_ProviderError(t *testing.T) {
	handler := newTestHandler(&mockGenerator{err: errors.New("CLI failed")})

	body, _ := json.Marshal(Request{
		User: "Say hello",
	})
	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error != "Failed to generate response" {
		t.Errorf("expected error 'Failed to generate response', got %q", resp.Error)
	}
}

func TestHandlePrompt_MethodNotAllowed(t *testing.T) {
	handler := newTestHandler(&mockGenerator{})

	req := httptest.NewRequest(http.MethodGet, "/prompt", nil)
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestHandlePrompt_ProviderSelection(t *testing.T) {
	claudeMock := &mockGenerator{response: "claude response"}
	geminiMock := &mockGenerator{response: "gemini response"}

	providers := map[string]provider.Generator{
		"claude": claudeMock,
		"gemini": geminiMock,
	}
	handler := newTestHandlerWithProviders(providers, "claude")

	// Test default provider (claude)
	body, _ := json.Marshal(Request{
		User: "Hello",
	})
	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ResponseText != "claude response" {
		t.Errorf("expected default provider (claude), got response %q", resp.ResponseText)
	}

	// Test explicit provider selection (gemini)
	body, _ = json.Marshal(Request{
		User:     "Hello",
		Provider: "gemini",
	})
	req = httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ResponseText != "gemini response" {
		t.Errorf("expected gemini provider, got response %q", resp.ResponseText)
	}
}

func TestHandlePrompt_UnknownProvider(t *testing.T) {
	handler := newTestHandler(&mockGenerator{response: "Hello"})

	body, _ := json.Marshal(Request{
		User:     "Hello",
		Provider: "unknown",
	})
	req := httptest.NewRequest(http.MethodPost, "/prompt", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandlePrompt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error != "Unknown provider: unknown" {
		t.Errorf("expected 'Unknown provider: unknown', got %q", resp.Error)
	}
}

func TestHandleProviders_Success(t *testing.T) {
	providers := map[string]provider.Generator{
		"claude": &mockGenerator{},
		"gemini": &mockGenerator{},
	}
	handler := newTestHandlerWithProviders(providers, "claude")

	req := httptest.NewRequest(http.MethodGet, "/providers", nil)
	w := httptest.NewRecorder()

	handler.HandleProviders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp struct {
		Providers []ProviderInfo `json:"providers"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(resp.Providers))
	}

	// Verify each provider has name and description
	for _, p := range resp.Providers {
		if p.Name == "" {
			t.Error("provider name should not be empty")
		}
		if p.Description == "" {
			t.Error("provider description should not be empty")
		}
	}
}

func TestHandleProviders_MethodNotAllowed(t *testing.T) {
	handler := newTestHandler(&mockGenerator{})

	req := httptest.NewRequest(http.MethodPost, "/providers", nil)
	w := httptest.NewRecorder()

	handler.HandleProviders(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
