package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListJobPostingsInput defines the input parameters for the
// list_job_postings tool (none required).
type ListJobPostingsInput struct{}

// ListJobPostingsOutput contains the list_job_postings results.
type ListJobPostingsOutput struct {
	// JobPostings is the list of job postings.
	JobPostings []ashby.JobPosting `json:"jobPostings"`

	// Total is the number of job postings returned.
	Total int `json:"total"`
}

// ListJobPostings handles the list_job_postings MCP tool call.
func (h *Handler) ListJobPostings(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListJobPostingsInput,
) (*mcp.CallToolResult, ListJobPostingsOutput, error) {

	postings, err := h.client.ListJobPostings(ctx)
	if err != nil {
		return nil, ListJobPostingsOutput{}, err
	}

	return nil, ListJobPostingsOutput{
		JobPostings: postings,
		Total:       len(postings),
	}, nil
}

// GetJobPostingInput defines the input parameters for the
// get_job_posting tool.
type GetJobPostingInput struct {
	// JobPostingID is the Ashby job posting ID to look up.
	JobPostingID string `json:"jobPostingId" jsonschema:"The Ashby job posting ID"`
}

// GetJobPostingOutput contains the get_job_posting results.
type GetJobPostingOutput struct {
	// JobPosting is the job posting details.
	JobPosting *ashby.JobPosting `json:"jobPosting"`
}

// GetJobPosting handles the get_job_posting MCP tool call.
func (h *Handler) GetJobPosting(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetJobPostingInput,
) (*mcp.CallToolResult, GetJobPostingOutput, error) {

	posting, err := h.client.GetJobPosting(
		ctx, input.JobPostingID,
	)
	if err != nil {
		return nil, GetJobPostingOutput{}, err
	}

	return nil, GetJobPostingOutput{JobPosting: posting}, nil
}
