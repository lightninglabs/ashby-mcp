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
