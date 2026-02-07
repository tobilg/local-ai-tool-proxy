package handler

import (
	"net/http"
)

// OpenAPISpec is the OpenAPI v3 specification for the Local AI Tool Proxy API.
const OpenAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Local AI Tool Proxy API",
    "description": "A local HTTP proxy that bridges web applications with AI CLI tools to generate responses using a configurable system prompt.",
    "version": "1.0.0",
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "servers": [
    {
      "url": "http://localhost:4000",
      "description": "Local development server"
    }
  ],
  "paths": {
    "/prompt": {
      "post": {
        "summary": "Generate Response",
        "description": "Generate a response from a user prompt using an AI CLI tool with the configured system prompt.",
        "operationId": "generateResponse",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Request"
              },
              "examples": {
                "basic": {
                  "summary": "Basic prompt",
                  "value": {
                    "user": "What is the capital of France?"
                  }
                },
                "with_provider": {
                  "summary": "Prompt with specific provider",
                  "value": {
                    "user": "Explain quantum computing in simple terms",
                    "provider": "gemini"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successfully generated response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Response"
                },
                "example": {
                  "response": "The capital of France is Paris."
                }
              }
            }
          },
          "400": {
            "description": "Bad request - invalid JSON, missing fields, or unknown provider",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "examples": {
                  "invalid_json": {
                    "summary": "Invalid JSON",
                    "value": {
                      "error": "Invalid JSON"
                    }
                  },
                  "missing_fields": {
                    "summary": "Missing required fields",
                    "value": {
                      "error": "The 'user' field is required"
                    }
                  },
                  "unknown_provider": {
                    "summary": "Unknown provider",
                    "value": {
                      "error": "Unknown provider: invalid"
                    }
                  }
                }
              }
            }
          },
          "405": {
            "description": "Method not allowed - only POST and OPTIONS are supported",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "Method not allowed"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error - AI CLI execution failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                },
                "example": {
                  "error": "Failed to generate response"
                }
              }
            }
          }
        }
      },
      "options": {
        "summary": "CORS Preflight",
        "description": "Handle CORS preflight requests for cross-origin access.",
        "operationId": "promptOptions",
        "responses": {
          "200": {
            "description": "CORS preflight response with appropriate headers"
          }
        }
      }
    },
    "/openapi.json": {
      "get": {
        "summary": "OpenAPI Specification",
        "description": "Returns the OpenAPI v3 specification for this API.",
        "operationId": "getOpenAPISpec",
        "responses": {
          "200": {
            "description": "OpenAPI v3 specification",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object"
                }
              }
            }
          }
        }
      }
    },
    "/providers": {
      "get": {
        "summary": "List Providers",
        "description": "Returns the list of available AI providers with their descriptions.",
        "operationId": "listProviders",
        "responses": {
          "200": {
            "description": "List of available providers",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProvidersResponse"
                },
                "example": {
                  "providers": [
                    {"name": "claude", "description": "Anthropic Claude CLI"},
                    {"name": "gemini", "description": "Google Gemini CLI"},
                    {"name": "codex", "description": "OpenAI Codex CLI"},
                    {"name": "continue", "description": "Continue CLI"},
                    {"name": "opencode", "description": "OpenCode CLI"}
                  ]
                }
              }
            }
          }
        }
      }
    },
    "/health": {
      "get": {
        "summary": "Health Check",
        "description": "Returns HTTP 200 if the proxy is running. Used by tools to check if the proxy is available.",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "Proxy is running"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Request": {
        "type": "object",
        "required": ["user"],
        "properties": {
          "user": {
            "type": "string",
            "description": "The user prompt to send to the AI provider",
            "example": "What is the capital of France?"
          },
          "provider": {
            "type": "string",
            "description": "AI provider to use for response generation. If omitted, uses the default configured provider.",
            "enum": ["claude", "gemini", "codex", "continue", "opencode"],
            "example": "claude"
          }
        }
      },
      "Response": {
        "type": "object",
        "properties": {
          "response": {
            "type": "string",
            "description": "Generated response from the AI provider",
            "example": "The capital of France is Paris."
          },
          "error": {
            "type": "string",
            "description": "Error message if the request failed"
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "string",
            "description": "Error message describing what went wrong",
            "example": "Failed to generate response"
          }
        }
      },
      "ProvidersResponse": {
        "type": "object",
        "properties": {
          "providers": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/ProviderInfo"
            },
            "description": "List of available AI providers"
          }
        }
      },
      "ProviderInfo": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "description": "Technical name of the provider",
            "example": "claude"
          },
          "description": {
            "type": "string",
            "description": "Human-readable description of the provider",
            "example": "Anthropic Claude CLI"
          }
        }
      }
    }
  }
}`

// HandleOpenAPI serves the OpenAPI v3 specification.
func (h *Handler) HandleOpenAPI(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(OpenAPISpec))
}
