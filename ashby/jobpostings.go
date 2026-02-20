package ashby

import "context"

// ListJobPostings returns all job postings.
func (c *Client) ListJobPostings(
	ctx context.Context,
) ([]JobPosting, error) {

	return Paginate[JobPosting](
		ctx, c, "jobPosting.list", nil, 0,
	)
}

// GetJobPosting returns details for a single job posting by
// ID.
func (c *Client) GetJobPosting(
	ctx context.Context, jobPostingID string,
) (*JobPosting, error) {

	var resp struct {
		Success bool       `json:"success"`
		Results JobPosting `json:"results"`
	}

	if err := c.Call(ctx, "jobPosting.info", map[string]any{
		"jobPostingId": jobPostingID,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}
