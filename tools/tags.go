package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListTagsInput defines the input parameters for the list_tags
// tool. No parameters are required.
type ListTagsInput struct{}

// ListTagsOutput contains the list_tags results.
type ListTagsOutput struct {
	// Tags is the list of all candidate tags.
	Tags []ashby.Tag `json:"tags"`

	// Total is the number of tags returned.
	Total int `json:"total"`
}

// ListTags handles the list_tags MCP tool call.
func (h *Handler) ListTags(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListTagsInput,
) (*mcp.CallToolResult, ListTagsOutput, error) {

	tags, err := h.client.ListTags(ctx)
	if err != nil {
		return nil, ListTagsOutput{}, err
	}

	return nil, ListTagsOutput{
		Tags:  tags,
		Total: len(tags),
	}, nil
}

// AddCandidateTagInput defines the input parameters for the
// add_candidate_tag tool.
type AddCandidateTagInput struct {
	// CandidateID is the candidate to tag.
	CandidateID string `json:"candidateId" jsonschema:"description=The Ashby candidate ID"`

	// TagID is the tag to add.
	TagID string `json:"tagId" jsonschema:"description=The Ashby tag ID"`
}

// AddCandidateTagOutput confirms the tag was added.
type AddCandidateTagOutput struct {
	// Success indicates the tag was added.
	Success bool `json:"success"`
}

// AddCandidateTag handles the add_candidate_tag MCP tool call.
func (h *Handler) AddCandidateTag(
	ctx context.Context, req *mcp.CallToolRequest,
	input AddCandidateTagInput,
) (*mcp.CallToolResult, AddCandidateTagOutput, error) {

	err := h.client.AddCandidateTag(
		ctx, input.CandidateID, input.TagID,
	)
	if err != nil {
		return nil, AddCandidateTagOutput{}, err
	}

	return nil, AddCandidateTagOutput{Success: true}, nil
}
