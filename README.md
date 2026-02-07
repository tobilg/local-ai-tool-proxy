# local-ai-tool-proxy

A local HTTP proxy that bridges web applications with AI CLI tools using a configurable system prompt.

## Overview

This proxy allows browser-based applications to leverage your local AI CLI subscriptions for generating AI responses. It loads a system prompt from a file at startup, accepts user prompts via HTTP, and returns AI-generated responses.

### How It Works

```
┌─────────────────┐     HTTP POST      ┌──────────────────────┐     exec      ┌───────────────┐
│   Web Browser   │ ─────────────────► │  local-ai-tool-proxy │ ────────────► │    AI CLI     │
│                 │ ◄───────────────── │  localhost:4000      │ ◄──────────── │ (claude/etc.) │
└─────────────────┘     Response       └──────────────────────┘    Response   └───────────────┘
```

### Supported Providers

| Provider    | CLI Command | Install |
|-------------|-------------|---------|
| Claude Code | `claude` | [Installation Guide](https://docs.anthropic.com/en/docs/claude-cli) |
| Google Gemini | `gemini` | [Installation Guide](https://geminicli.com/docs/installation) |
| OpenAI Codex | `codex` | [Installation Guide](https://developers.openai.com/codex/cli/installation) |
| Continue | `cn` | `npm i -g @continuedev/cli` |
| OpenCode | `opencode` | [Installation Guide](https://opencode.ai/docs/cli) |

## Installation

### Homebrew (macOS Apple Silicon)

```bash
brew tap tobilg/local-ai-tool-proxy
brew install local-ai-tool-proxy
```

After installation, you can run it as a service:
```bash
brew services start local-ai-tool-proxy
```

Or run it manually:
```bash
local-ai-tool-proxy
```

### Build from source

```bash
# Clone the repository
git clone https://github.com/tobilg/local-ai-tool-proxy.git
cd local-ai-tool-proxy

# Build for your platform
make build

# Or build for all platforms
make build-all
```

### Pre-built binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/tobilg/local-ai-tool-proxy/releases):

| Platform | Binary |
|----------|--------|
| Windows | `local-ai-tool-proxy-windows-amd64.exe` |
| Linux (x86_64) | `local-ai-tool-proxy-linux-amd64` |
| Linux (ARM64, e.g. Raspberry Pi 5) | `local-ai-tool-proxy-linux-arm64` |
| macOS (Apple Silicon) | `local-ai-tool-proxy-darwin-arm64` |

## Usage

### Prerequisites

At least one of these CLI tools must be installed and authenticated:

```bash
# Claude CLI (Anthropic)
# Follow: https://docs.anthropic.com/en/docs/claude-cli

# Gemini CLI (Google)
# Follow: https://geminicli.com/docs/installation

# Codex CLI (OpenAI)
# Follow: https://developers.openai.com/codex/cli/installation

# Continue CLI
npm i -g @continuedev/cli

# OpenCode CLI
# Follow: https://opencode.ai/docs/cli
```

You also need a **system prompt file**. Create one with the instructions for your AI assistant:

```bash
echo "You are a helpful assistant." > /path/to/system-prompt.txt
```

### Running the proxy

```bash
# Run with default settings (Claude provider)
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt ./dist/local-ai-tool-proxy

# Run with a specific default provider
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt LOCAL_AI_TOOL_PROXY_PROVIDER=gemini ./dist/local-ai-tool-proxy

# Run with custom port
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt LOCAL_AI_TOOL_PROXY_PORT=8080 ./dist/local-ai-tool-proxy

# Run with custom allowed origin
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN="http://localhost:3000" ./dist/local-ai-tool-proxy

# Run with HTTPS/TLS (requires certificate and key files)
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt LOCAL_AI_TOOL_PROXY_TLS_CERT=/path/to/cert.pem LOCAL_AI_TOOL_PROXY_TLS_KEY=/path/to/key.pem ./dist/local-ai-tool-proxy
```

The proxy will start and display (with default settings):

```
Local AI Tool Proxy active at http://localhost:4000
Default provider: claude
System prompt: /path/to/system-prompt.txt
Allowed origin: http://localhost:3000
Available providers: claude, gemini, codex, continue, opencode
API docs: http://localhost:4000/openapi.json
Press Ctrl+C to stop
```

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `LOCAL_AI_TOOL_PROXY_PORT` | `4000` | Port the proxy listens on |
| `LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN` | `http://localhost:3000` | CORS allowed origin |
| `LOCAL_AI_TOOL_PROXY_PROVIDER` | `claude` | Default AI provider |
| `LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT` | *(required)* | Path to system prompt file |
| `LOCAL_AI_TOOL_PROXY_TLS_CERT` | - | Path to TLS certificate file (enables HTTPS) |
| `LOCAL_AI_TOOL_PROXY_TLS_KEY` | - | Path to TLS private key file (enables HTTPS) |

Valid providers: `claude`, `gemini`, `codex`, `continue`, `opencode`

### HTTPS/TLS Support

To run the proxy over HTTPS (required for Safari and strict browser security), provide both TLS certificate and key files:

```bash
# Generate self-signed certificates with mkcert (recommended for local development)
# Install mkcert: https://github.com/FiloSottile/mkcert
mkcert -install
mkcert localhost 127.0.0.1 ::1

# Run with the generated certificates
LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/path/to/system-prompt.txt LOCAL_AI_TOOL_PROXY_TLS_CERT=localhost+2.pem LOCAL_AI_TOOL_PROXY_TLS_KEY=localhost+2-key.pem ./dist/local-ai-tool-proxy
```

When TLS is enabled, the proxy will display:
```
Local AI Tool Proxy active at https://localhost:4000
...
TLS enabled: cert=localhost+2.pem, key=localhost+2-key.pem
```

## API

### GET /health

Health check endpoint to verify the proxy is running.

**Example Request:**

```bash
curl http://localhost:4000/health
```

**Example Response (200):**

Empty response with HTTP status 200.

---

### GET /openapi.json

Returns the OpenAPI v3 specification for this API.

**Example Request:**

```bash
curl http://localhost:4000/openapi.json
```

**Example Response (200):**

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Local AI Tool Proxy API",
    "version": "1.0.0"
  },
  "paths": { ... }
}
```

---

### GET /providers

Returns the list of available AI providers with their descriptions.

**Example Request:**

```bash
curl http://localhost:4000/providers
```

**Example Response (200):**

```json
{
  "providers": [
    {"name": "claude", "description": "Claude Code"},
    {"name": "gemini", "description": "Google Gemini"},
    {"name": "codex", "description": "OpenAI Codex"},
    {"name": "continue", "description": "Continue"},
    {"name": "opencode", "description": "OpenCode"}
  ]
}
```

---

### POST /prompt

Generate a response from a user prompt using the configured system prompt and AI provider.

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user` | string | Yes | The user prompt |
| `provider` | string | No | AI provider to use (defaults to configured provider) |

**Example Request:**

```bash
curl -X POST http://localhost:4000/prompt \
  -H "Content-Type: application/json" \
  -d '{
    "user": "What is the capital of France?"
  }'
```

**Example Response (200):**

```json
{
  "response": "The capital of France is Paris."
}
```

**Example Request with Provider Override:**

```bash
curl -X POST http://localhost:4000/prompt \
  -H "Content-Type: application/json" \
  -d '{
    "user": "Explain quantum computing in simple terms",
    "provider": "gemini"
  }'
```

**Error Responses:**

| Status | Description | Example |
|--------|-------------|---------|
| 400 | Invalid JSON or missing required fields | `{"error": "The 'user' field is required"}` |
| 400 | Unknown provider | `{"error": "Unknown provider: invalid"}` |
| 405 | Method not allowed | `{"error": "Method not allowed"}` |
| 500 | AI CLI execution failed | `{"error": "Failed to generate response"}` |

## Development

### Running tests

```bash
make test
```

### Build commands

```bash
make build              # Build for current platform
make build-all          # Build for all platforms
make build-windows      # Build for Windows
make build-linux        # Build for Linux
make build-darwin-amd64 # Build for macOS Intel
make build-darwin-arm64 # Build for macOS Apple Silicon
make clean              # Remove build artifacts
```

### Project structure

```
local-ai-tool-proxy/
├── src/
│   ├── cmd/local-ai-tool-proxy/    # Application entry point
│   └── internal/
│       ├── config/          # Configuration loading
│       ├── handler/         # HTTP handlers
│       └── provider/        # AI CLI provider implementations
├── dist/                    # Built binaries
├── Makefile
└── README.md
```

## Running as a systemd Service (Linux)

On Linux systems (including Raspberry Pi), you can run the proxy as a background service using systemd.

### 1. Install the binary

```bash
sudo cp local-ai-tool-proxy /usr/local/bin/
sudo chmod +x /usr/local/bin/local-ai-tool-proxy
```

### 2. Create a system prompt file

```bash
sudo mkdir -p /etc/local-ai-tool-proxy
echo "You are a helpful assistant." | sudo tee /etc/local-ai-tool-proxy/system-prompt.txt
```

### 3. Create the systemd unit file

```bash
sudo tee /etc/systemd/system/local-ai-tool-proxy.service << 'EOF'
[Unit]
Description=Local AI Tool Proxy
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/local-ai-tool-proxy
Restart=on-failure
RestartSec=5

# Configuration
Environment=LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=/etc/local-ai-tool-proxy/system-prompt.txt
Environment=LOCAL_AI_TOOL_PROXY_PORT=4000
Environment=LOCAL_AI_TOOL_PROXY_PROVIDER=claude
Environment=LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN=http://localhost:3000
# Optional: uncomment for HTTPS
# Environment=LOCAL_AI_TOOL_PROXY_TLS_CERT=/path/to/cert.pem
# Environment=LOCAL_AI_TOOL_PROXY_TLS_KEY=/path/to/key.pem

# Run as your user (required so the AI CLI can access your credentials)
User=%i
Group=%i

[Install]
WantedBy=multi-user.target
EOF
```

> **Note:** Replace `%i` with your actual username (e.g. `User=pi`), or use a user-level service instead (see below). The AI CLI tools need access to your user's authentication tokens, so the service must run as the user that has authenticated with the CLI.

### 4. Enable and start the service

```bash
sudo systemctl daemon-reload
sudo systemctl enable local-ai-tool-proxy
sudo systemctl start local-ai-tool-proxy
```

### 5. Check status and logs

```bash
sudo systemctl status local-ai-tool-proxy
journalctl -u local-ai-tool-proxy -f
```

### Alternative: user-level service (no sudo required)

If you prefer running the service as your own user without root:

```bash
mkdir -p ~/.config/systemd/user

tee ~/.config/systemd/user/local-ai-tool-proxy.service << 'EOF'
[Unit]
Description=Local AI Tool Proxy
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/local-ai-tool-proxy
Restart=on-failure
RestartSec=5
Environment=LOCAL_AI_TOOL_PROXY_SYSTEM_PROMPT=%h/.config/local-ai-tool-proxy/system-prompt.txt
Environment=LOCAL_AI_TOOL_PROXY_PORT=4000
Environment=LOCAL_AI_TOOL_PROXY_PROVIDER=claude
Environment=LOCAL_AI_TOOL_PROXY_ALLOWED_ORIGIN=http://localhost:3000

[Install]
WantedBy=default.target
EOF

# Create the config directory and system prompt
mkdir -p ~/.config/local-ai-tool-proxy
echo "You are a helpful assistant." > ~/.config/local-ai-tool-proxy/system-prompt.txt

# Enable and start
systemctl --user daemon-reload
systemctl --user enable local-ai-tool-proxy
systemctl --user start local-ai-tool-proxy

# Enable lingering so the service starts at boot (even without login)
loginctl enable-linger $USER
```

## Browser Security Notes

Modern browsers enforce strict security policies for requests from HTTPS sites to local HTTP servers. This proxy includes:

- **CORS headers** for cross-origin requests
- **Private Network Access** header (`Access-Control-Allow-Private-Network: true`) for browser compatibility
- **Optional HTTPS/TLS support** for browsers with strict mixed content policies (like Safari)

**Recommended:** Use HTTPS with [mkcert](https://github.com/FiloSottile/mkcert) for the best browser compatibility (see [HTTPS/TLS Support](#httpstls-support)).

**Alternative for Chrome:** Enable `chrome://flags/#allow-insecure-localhost` to allow HTTP connections to localhost.

## License

MIT
