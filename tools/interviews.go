package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListInterviewStagesInput defines the input parameters for
// the list_interview_stages tool.
type ListInterviewStagesInput struct {
	// InterviewPlanID is the interview plan to list stages
	// for. This is typically the job's interviewPlanId field.
	InterviewPlanID string `json:"interviewPlanId" jsonschema:"The interview plan ID (from job.interviewPlanId)"`
}

// ListInterviewStagesOutput contains the
// list_interview_stages results.
type ListInterviewStagesOutput struct {
	// Stages is the list of interview stages.
	Stages []ashby.InterviewStage `json:"stages"`

	// Total is the number of stages returned.
	Total int `json:"total"`
}

// ListInterviewStages handles the list_interview_stages MCP
// tool call.
func (h *Handler) ListInterviewStages(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListInterviewStagesInput,
) (*mcp.CallToolResult, ListInterviewStagesOutput, error) {

	stages, err := h.client.ListInterviewStages(
		ctx, input.InterviewPlanID,
	)
	if err != nil {
		return nil, ListInterviewStagesOutput{}, err
	}

	return nil, ListInterviewStagesOutput{
		Stages: stages,
		Total:  len(stages),
	}, nil
}

// ListInterviewsInput defines the input parameters for the
// list_interviews tool.
type ListInterviewsInput struct {
	// ApplicationID filters interviews by application.
	ApplicationID string `json:"applicationId,omitempty" jsonschema:"Filter by Ashby application ID"`
}

// ListInterviewsOutput contains the list_interviews results.
type ListInterviewsOutput struct {
	// Interviews is the list of interviews.
	Interviews []ashby.Interview `json:"interviews"`

	// Total is the number of interviews returned.
	Total int `json:"total"`
}

// ListInterviews handles the list_interviews MCP tool call.
func (h *Handler) ListInterviews(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListInterviewsInput,
) (*mcp.CallToolResult, ListInterviewsOutput, error) {

	interviews, err := h.client.ListInterviews(
		ctx, input.ApplicationID,
	)
	if err != nil {
		return nil, ListInterviewsOutput{}, err
	}

	return nil, ListInterviewsOutput{
		Interviews: interviews,
		Total:      len(interviews),
	}, nil
}
