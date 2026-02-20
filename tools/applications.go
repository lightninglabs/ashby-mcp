package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListApplicationsInput defines the input parameters for the
// list_applications tool.
type ListApplicationsInput struct {
	// JobID filters applications by job.
	JobID string `json:"jobId,omitempty" jsonschema:"Filter by Ashby job ID"`

	// Status filters by application status: Active, Hired,
	// Archived, or Rejected.
	Status string `json:"status,omitempty" jsonschema:"Application status filter: Active Hired Archived or Rejected"`

	// Limit caps the maximum number of results.
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: 100)"`

	// Cursor is an opaque pagination token from a previous
	// response.
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor from a previous response"`

	// CreatedAfter filters to applications created after
	// this Unix epoch timestamp in milliseconds.
	CreatedAfter int64 `json:"createdAfter,omitempty" jsonschema:"Filter to applications created after this Unix epoch ms timestamp"`

	// UpdatedAfter filters to applications updated after
	// this Unix epoch timestamp in milliseconds.
	UpdatedAfter int64 `json:"updatedAfter,omitempty" jsonschema:"Filter to applications updated after this Unix epoch ms timestamp"`
}

// ListApplicationsOutput contains the list_applications
// results.
type ListApplicationsOutput struct {
	// Applications is the list of matching applications.
	Applications []ashby.Application `json:"applications"`

	// Total is the number returned in this response.
	Total int `json:"total"`

	// NextCursor is the cursor for fetching the next page.
	NextCursor string `json:"nextCursor,omitempty"`

	// MoreDataAvailable indicates additional pages exist.
	MoreDataAvailable bool `json:"moreDataAvailable,omitempty"`
}

// ListApplications handles the list_applications MCP tool call.
func (h *Handler) ListApplications(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListApplicationsInput,
) (*mcp.CallToolResult, ListApplicationsOutput, error) {

	result, err := h.client.ListApplications(
		ctx, ashby.ListApplicationsOpts{
			JobID:        input.JobID,
			Status:       input.Status,
			Limit:        input.Limit,
			Cursor:       input.Cursor,
			CreatedAfter: input.CreatedAfter,
			UpdatedAfter: input.UpdatedAfter,
		},
	)
	if err != nil {
		return nil, ListApplicationsOutput{}, err
	}

	return nil, ListApplicationsOutput{
		Applications:      result.Applications,
		Total:             len(result.Applications),
		NextCursor:        result.NextCursor,
		MoreDataAvailable: result.MoreDataAvailable,
	}, nil
}

// GetApplicationInput defines the input parameters for the
// get_application tool.
type GetApplicationInput struct {
	// ApplicationID is the Ashby application ID to look up.
	ApplicationID string `json:"applicationId" jsonschema:"The Ashby application ID"`

	// Expand controls which related data to include. Valid
	// values: applicationFormSubmissions, openings, referrals.
	Expand []string `json:"expand,omitempty" jsonschema:"Related data to include: applicationFormSubmissions openings referrals"`
}

// GetApplicationOutput contains the get_application results.
type GetApplicationOutput struct {
	// Application is the application details.
	Application *ashby.Application `json:"application"`
}

// GetApplication handles the get_application MCP tool call.
func (h *Handler) GetApplication(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetApplicationInput,
) (*mcp.CallToolResult, GetApplicationOutput, error) {

	app, err := h.client.GetApplication(
		ctx, input.ApplicationID, input.Expand,
	)
	if err != nil {
		return nil, GetApplicationOutput{}, err
	}

	return nil, GetApplicationOutput{Application: app}, nil
}

// ChangeApplicationStageInput defines the input parameters for
// the change_application_stage tool.
type ChangeApplicationStageInput struct {
	// ApplicationID is the application to move.
	ApplicationID string `json:"applicationId" jsonschema:"The Ashby application ID to move"`

	// InterviewStageID is the target interview stage.
	InterviewStageID string `json:"interviewStageId" jsonschema:"The target interview stage ID"`
}

// ChangeApplicationStageOutput confirms the stage change.
type ChangeApplicationStageOutput struct {
	// Success indicates whether the stage change succeeded.
	Success bool `json:"success"`
}

// ChangeApplicationStage handles the change_application_stage
// MCP tool call.
func (h *Handler) ChangeApplicationStage(
	ctx context.Context, req *mcp.CallToolRequest,
	input ChangeApplicationStageInput,
) (*mcp.CallToolResult, ChangeApplicationStageOutput, error) {

	err := h.client.ChangeApplicationStage(
		ctx, input.ApplicationID, input.InterviewStageID,
	)
	if err != nil {
		return nil, ChangeApplicationStageOutput{}, err
	}

	return nil, ChangeApplicationStageOutput{
		Success: true,
	}, nil
}

// CreateApplicationInput defines the input parameters for the
// create_application tool.
type CreateApplicationInput struct {
	// CandidateID is the candidate to create the application
	// for.
	CandidateID string `json:"candidateId" jsonschema:"The Ashby candidate ID"`

	// JobID is the job to apply to.
	JobID string `json:"jobId" jsonschema:"The Ashby job ID"`

	// Source is an optional source identifier.
	Source string `json:"source,omitempty" jsonschema:"Optional source identifier (e.g. Referral)"`
}

// CreateApplicationOutput contains the newly created
// application.
type CreateApplicationOutput struct {
	// Application is the newly created application.
	Application *ashby.Application `json:"application"`
}

// CreateApplication handles the create_application MCP tool
// call.
func (h *Handler) CreateApplication(
	ctx context.Context, req *mcp.CallToolRequest,
	input CreateApplicationInput,
) (*mcp.CallToolResult, CreateApplicationOutput, error) {

	app, err := h.client.CreateApplication(
		ctx, input.CandidateID, input.JobID, input.Source,
	)
	if err != nil {
		return nil, CreateApplicationOutput{}, err
	}

	return nil, CreateApplicationOutput{
		Application: app,
	}, nil
}

// ChangeApplicationSourceInput defines the input parameters
// for the change_application_source tool.
type ChangeApplicationSourceInput struct {
	// ApplicationID is the application to update.
	ApplicationID string `json:"applicationId" jsonschema:"The Ashby application ID to update"`

	// SourceID is the new source ID. Pass an empty string to
	// clear the source.
	SourceID string `json:"sourceId" jsonschema:"Source ID to assign, or empty string to unset"`
}

// ChangeApplicationSourceOutput confirms the source change.
type ChangeApplicationSourceOutput struct {
	// Success indicates whether the source change succeeded.
	Success bool `json:"success"`
}

// ChangeApplicationSource handles the
// change_application_source MCP tool call.
func (h *Handler) ChangeApplicationSource(
	ctx context.Context, req *mcp.CallToolRequest,
	input ChangeApplicationSourceInput,
) (*mcp.CallToolResult, ChangeApplicationSourceOutput, error) {

	err := h.client.ChangeApplicationSource(
		ctx, input.ApplicationID, input.SourceID,
	)
	if err != nil {
		return nil, ChangeApplicationSourceOutput{}, err
	}

	return nil, ChangeApplicationSourceOutput{
		Success: true,
	}, nil
}
