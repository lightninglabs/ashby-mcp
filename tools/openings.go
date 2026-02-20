package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListOpeningsInput defines the input parameters for the
// list_openings tool (none required).
type ListOpeningsInput struct{}

// ListOpeningsOutput contains the list_openings results.
type ListOpeningsOutput struct {
	// Openings is the list of openings.
	Openings []ashby.Opening `json:"openings"`

	// Total is the number of openings returned.
	Total int `json:"total"`
}

// ListOpenings handles the list_openings MCP tool call.
func (h *Handler) ListOpenings(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListOpeningsInput,
) (*mcp.CallToolResult, ListOpeningsOutput, error) {

	openings, err := h.client.ListOpenings(ctx)
	if err != nil {
		return nil, ListOpeningsOutput{}, err
	}

	return nil, ListOpeningsOutput{
		Openings: openings,
		Total:    len(openings),
	}, nil
}

// GetOpeningInput defines the input parameters for the
// get_opening tool.
type GetOpeningInput struct {
	// OpeningID is the Ashby opening ID to look up.
	OpeningID string `json:"openingId" jsonschema:"The Ashby opening ID"`
}

// GetOpeningOutput contains the get_opening results.
type GetOpeningOutput struct {
	// Opening is the opening details.
	Opening *ashby.Opening `json:"opening"`
}

// GetOpening handles the get_opening MCP tool call.
func (h *Handler) GetOpening(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetOpeningInput,
) (*mcp.CallToolResult, GetOpeningOutput, error) {

	opening, err := h.client.GetOpening(ctx, input.OpeningID)
	if err != nil {
		return nil, GetOpeningOutput{}, err
	}

	return nil, GetOpeningOutput{Opening: opening}, nil
}

// SearchOpeningsInput defines the input parameters for the
// search_openings tool.
type SearchOpeningsInput struct {
	// Term is the search query string.
	Term string `json:"term" jsonschema:"Search term to match against openings"`
}

// SearchOpeningsOutput contains the search_openings results.
type SearchOpeningsOutput struct {
	// Openings is the list of matching openings.
	Openings []ashby.Opening `json:"openings"`

	// Total is the number of matching openings.
	Total int `json:"total"`
}

// SearchOpenings handles the search_openings MCP tool call.
func (h *Handler) SearchOpenings(
	ctx context.Context, req *mcp.CallToolRequest,
	input SearchOpeningsInput,
) (*mcp.CallToolResult, SearchOpeningsOutput, error) {

	openings, err := h.client.SearchOpenings(ctx, input.Term)
	if err != nil {
		return nil, SearchOpeningsOutput{}, err
	}

	return nil, SearchOpeningsOutput{
		Openings: openings,
		Total:    len(openings),
	}, nil
}
