package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetFileURLInput defines the input parameters for the
// get_file_url tool.
type GetFileURLInput struct {
	// FileHandle is the opaque handle string from a
	// resumeFileHandle or fileHandles entry.
	FileHandle string `json:"fileHandle" jsonschema:"Opaque file handle string from a resumeFileHandle or fileHandles entry"`
}

// GetFileURLOutput contains the get_file_url results.
type GetFileURLOutput struct {
	// URL is the pre-signed download URL for the file.
	URL string `json:"url"`
}

// GetFileURL handles the get_file_url MCP tool call.
func (h *Handler) GetFileURL(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetFileURLInput,
) (*mcp.CallToolResult, GetFileURLOutput, error) {

	url, err := h.client.GetFileURL(ctx, input.FileHandle)
	if err != nil {
		return nil, GetFileURLOutput{}, err
	}

	return nil, GetFileURLOutput{URL: url}, nil
}
