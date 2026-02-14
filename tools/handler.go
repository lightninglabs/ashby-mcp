package tools

import "github.com/lightninglabs/ashby-mcp/ashby"

// Handler provides MCP tool handlers backed by an Ashby API
// client. Each exported method implements a single MCP tool's
// logic.
type Handler struct {
	client *ashby.Client
}

// NewHandler creates a new Handler wrapping the given Ashby
// client.
func NewHandler(client *ashby.Client) *Handler {
	return &Handler{client: client}
}
