# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

local-ai-tool-proxy is a local HTTP proxy server (Go) that bridges web applications with AI CLI tools. It accepts a user prompt via HTTP POST, combines it with a configurable system prompt loaded from a file, executes local AI CLI commands, and returns the AI-generated response.

## Build and Test Commands

```bash
make build        # Build for current platform (output: dist/local-ai-tool-proxy)
make test         # Run all tests with verbose output
make build-all    # Build for all platforms (windows, linux, darwin-amd64, darwin-arm64)
make clean        # Remove build artifacts
```

Run individual tests:
```bash
go test -v ./src/internal/provider/...   # Test provider package only
go test -v ./src/internal/handler/...    # Test handler package only
go test -run TestClaudeClient ./src/internal/provider/  # Run specific test
```

## Architecture

**Entry point**: [main.go](src/cmd/local-ai-tool-proxy/main.go) - Initializes all providers, sets up HTTP routes (`/prompt`, `/openapi.json`), and handles graceful shutdown.

**Provider pattern**: Each AI CLI (claude, gemini, codex, continue, opencode) implements the `Generator` interface in [provider.go](src/internal/provider/provider.go):
```go
type Generator interface {
    Generate(systemPrompt, userPrompt string) (string, error)
}
```
Provider implementations call external CLI binaries via `os/exec` and parse their responses. Use `CleanResponse()` helper to strip markdown formatting from responses.

**Handler**: [handler.go](src/internal/handler/handler.go) handles HTTP requests, CORS (including Private Network Access headers for browser compatibility), and provider selection (per-request override or default from config).

**Config**: [config.go](src/internal/config/config.go) loads from environment variables:
- `LOCAL_AI_TOOL_PROXY_PORT` (default: 4000)
- `LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN` (default: http://localhost:3000)
- `LOCAL_AI_TOOL_PROXY_PROVIDER` (default: claude)
- `LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT` (required: path to system prompt file)
- `LOCAL_AI_TOOL_PROXY_TLS_CERT` (optional: path to TLS certificate file)
- `LOCAL_AI_TOOL_PROXY_TLS_KEY` (optional: path to TLS private key file)

When both TLS cert and key are provided, the server runs over HTTPS instead of HTTP.

## Adding a New Provider

1. Create `src/internal/provider/<name>.go` implementing `Generator`
2. Create corresponding `<name>_test.go`
3. Register in `main.go` providers map
4. Update valid provider names in error message
