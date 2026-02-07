package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tobilg/local-ai-tool-proxy/src/internal/provider"
)

// Request represents the incoming request payload.
type Request struct {
	User     string `json:"user"`
	Provider string `json:"provider,omitempty"`
}

// Response represents the response payload.
type Response struct {
	ResponseText string `json:"response,omitempty"`
	Error        string `json:"error,omitempty"`
}

// ProviderInfo represents a provider with its metadata.
type ProviderInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// providerDescriptions maps provider names to human-readable descriptions.
var providerDescriptions = map[string]string{
	"claude":   "Claude Code",
	"gemini":   "Google Gemini",
	"codex":    "OpenAI Codex",
	"continue": "Continue",
	"opencode": "OpenCode",
}

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	providers       map[string]provider.Generator
	defaultProvider string
	allowedOrigin   string
	systemPrompt    string
}

// New creates a new Handler with the given dependencies.
func New(providers map[string]provider.Generator, defaultProvider, allowedOrigin, systemPrompt string) *Handler {
	return &Handler{
		providers:       providers,
		defaultProvider: defaultProvider,
		allowedOrigin:   allowedOrigin,
		systemPrompt:    systemPrompt,
	}
}

// HandlePrompt handles POST /prompt requests.
func (h *Handler) HandlePrompt(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Invalid JSON: %v", err)
		h.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.User == "" {
		log.Printf("[ERROR] Missing required field: user=%q", req.User)
		h.sendError(w, "The 'user' field is required", http.StatusBadRequest)
		return
	}

	// Determine which provider to use
	providerName := req.Provider
	if providerName == "" {
		providerName = h.defaultProvider
	}

	p, ok := h.providers[providerName]
	if !ok {
		log.Printf("[ERROR] Unknown provider: %s", providerName)
		h.sendError(w, fmt.Sprintf("Unknown provider: %s", providerName), http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Generating response using %s for prompt: %q", providerName, req.User)

	result, err := p.Generate(h.systemPrompt, req.User)
	if err != nil {
		log.Printf("[ERROR] %s CLI failed: %v", providerName, err)
		h.sendError(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Successfully generated response")
	h.sendJSON(w, Response{ResponseText: result})
}

// setCORSHeaders sets the required CORS and Private Network Access headers.
func (h *Handler) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", h.allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Private-Network", "true")
}

// sendError sends an error response as JSON.
func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Error: message})
}

// sendJSON sends a successful JSON response.
func (h *Handler) sendJSON(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleHealth handles GET /health requests for health checks.
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleProviders handles GET /providers requests.
func (h *Handler) HandleProviders(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providers := make([]ProviderInfo, 0, len(h.providers))
	for name := range h.providers {
		description := providerDescriptions[name]
		if description == "" {
			description = name
		}
		providers = append(providers, ProviderInfo{
			Name:        name,
			Description: description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]ProviderInfo{"providers": providers})
}
