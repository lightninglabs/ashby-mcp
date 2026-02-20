package ashby

import "context"

// ListInterviewStages returns the interview stages for a job's
// interview plan.
func (c *Client) ListInterviewStages(
	ctx context.Context, interviewPlanID string,
) ([]InterviewStage, error) {

	return Paginate[InterviewStage](
		ctx, c, "interviewStage.list", map[string]any{
			"interviewPlanId": interviewPlanID,
		}, 0,
	)
}

// GetInterviewStage returns details for a single interview
// stage by ID.
func (c *Client) GetInterviewStage(
	ctx context.Context, stageID string,
) (*InterviewStage, error) {

	var resp struct {
		Success bool           `json:"success"`
		Results InterviewStage `json:"results"`
	}

	if err := c.Call(ctx, "interviewStage.info", map[string]any{
		"interviewStageId": stageID,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// ListInterviewPlans returns all interview plans.
func (c *Client) ListInterviewPlans(
	ctx context.Context,
) ([]InterviewPlan, error) {

	return Paginate[InterviewPlan](
		ctx, c, "interviewPlan.list", nil, 0,
	)
}

// ListInterviews returns interviews filtered by application
// ID.
func (c *Client) ListInterviews(
	ctx context.Context, applicationID string,
) ([]Interview, error) {

	params := make(map[string]any)
	if applicationID != "" {
		params["applicationId"] = applicationID
	}

	return Paginate[Interview](
		ctx, c, "interview.list", params, 0,
	)
}
