package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListJobsInput defines the input parameters for the list_jobs
// tool.
type ListJobsInput struct {
	// Status filters jobs by their current status: Open,
	// Closed, Archived, or Draft.
	Status string `json:"status,omitempty" jsonschema:"Job status filter: Open Closed Archived or Draft"`

	// Limit caps the maximum number of results returned.
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: all)"`
}

// ListJobsOutput contains the list_jobs results.
type ListJobsOutput struct {
	// Jobs is the list of matching jobs.
	Jobs []ashby.Job `json:"jobs"`

	// Total is the number of jobs returned.
	Total int `json:"total"`
}

// ListJobs handles the list_jobs MCP tool call.
func (h *Handler) ListJobs(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListJobsInput,
) (*mcp.CallToolResult, ListJobsOutput, error) {

	jobs, err := h.client.ListJobs(
		ctx, input.Status, input.Limit,
	)
	if err != nil {
		return nil, ListJobsOutput{}, err
	}

	return nil, ListJobsOutput{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

// GetJobInput defines the input parameters for the get_job
// tool.
type GetJobInput struct {
	// JobID is the Ashby job ID to look up.
	JobID string `json:"jobId" jsonschema:"The Ashby job ID"`
}

// GetJobOutput contains the get_job results.
type GetJobOutput struct {
	// Job is the job details.
	Job *ashby.Job `json:"job"`
}

// GetJob handles the get_job MCP tool call.
func (h *Handler) GetJob(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetJobInput,
) (*mcp.CallToolResult, GetJobOutput, error) {

	job, err := h.client.GetJob(ctx, input.JobID)
	if err != nil {
		return nil, GetJobOutput{}, err
	}

	return nil, GetJobOutput{Job: job}, nil
}

// SearchJobsInput defines the input parameters for the
// search_jobs tool.
type SearchJobsInput struct {
	// Term is the search query string.
	Term string `json:"term" jsonschema:"Search term to match against job titles"`

	// Limit caps the maximum number of results returned.
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of results to return"`
}

// SearchJobsOutput contains the search_jobs results.
type SearchJobsOutput struct {
	// Jobs is the list of matching jobs.
	Jobs []ashby.Job `json:"jobs"`

	// Total is the number of jobs returned.
	Total int `json:"total"`
}

// SearchJobs handles the search_jobs MCP tool call.
func (h *Handler) SearchJobs(
	ctx context.Context, req *mcp.CallToolRequest,
	input SearchJobsInput,
) (*mcp.CallToolResult, SearchJobsOutput, error) {

	jobs, err := h.client.SearchJobs(
		ctx, input.Term, input.Limit,
	)
	if err != nil {
		return nil, SearchJobsOutput{}, err
	}

	return nil, SearchJobsOutput{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

// SetJobStatusInput defines the input parameters for the
// set_job_status tool.
type SetJobStatusInput struct {
	// JobID is the Ashby job ID to update.
	JobID string `json:"jobId" jsonschema:"The Ashby job ID"`

	// Status is the new status: Open, Closed, or Archived.
	Status string `json:"status" jsonschema:"New job status: Open Closed or Archived"`
}

// SetJobStatusOutput contains the updated job after the status
// change.
type SetJobStatusOutput struct {
	// Job is the updated job record.
	Job *ashby.Job `json:"job"`
}

// SetJobStatus handles the set_job_status MCP tool call.
func (h *Handler) SetJobStatus(
	ctx context.Context, req *mcp.CallToolRequest,
	input SetJobStatusInput,
) (*mcp.CallToolResult, SetJobStatusOutput, error) {

	job, err := h.client.SetJobStatus(
		ctx, input.JobID, input.Status,
	)
	if err != nil {
		return nil, SetJobStatusOutput{}, err
	}

	return nil, SetJobStatusOutput{Job: job}, nil
}

// UpdateJobInput defines the input parameters for the
// update_job tool.
type UpdateJobInput struct {
	// JobID is the Ashby job ID to update.
	JobID string `json:"jobId" jsonschema:"The Ashby job ID"`

	// Title is the updated job title.
	Title string `json:"title,omitempty" jsonschema:"Updated job title"`

	// DepartmentID is the updated department ID.
	DepartmentID string `json:"departmentId,omitempty" jsonschema:"Updated department ID"`

	// LocationIds is the updated list of location IDs.
	LocationIds []string `json:"locationIds,omitempty" jsonschema:"Updated location IDs"`

	// EmploymentType is the updated employment type
	// (e.g. FullTime).
	EmploymentType string `json:"employmentType,omitempty" jsonschema:"Updated employment type"`
}

// UpdateJobOutput contains the updated job record.
type UpdateJobOutput struct {
	// Job is the updated job record.
	Job *ashby.Job `json:"job"`
}

// UpdateJob handles the update_job MCP tool call.
func (h *Handler) UpdateJob(
	ctx context.Context, req *mcp.CallToolRequest,
	input UpdateJobInput,
) (*mcp.CallToolResult, UpdateJobOutput, error) {

	job, err := h.client.UpdateJob(
		ctx, input.JobID, ashby.UpdateJobOpts{
			Title:          input.Title,
			DepartmentID:   input.DepartmentID,
			LocationIds:    input.LocationIds,
			EmploymentType: input.EmploymentType,
		},
	)
	if err != nil {
		return nil, UpdateJobOutput{}, err
	}

	return nil, UpdateJobOutput{Job: job}, nil
}
