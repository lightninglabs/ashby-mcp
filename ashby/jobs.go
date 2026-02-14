package ashby

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListJobs returns all jobs, optionally filtered by status.
// Valid statuses are Open, Closed, Archived, and Draft.
// Filtering is done client-side because the Ashby API does not
// support server-side status filtering on job.list.
func (c *Client) ListJobs(
	ctx context.Context, status string, limit int,
) ([]Job, error) {

	params := map[string]any{}

	jobs, err := Paginate[Job](
		ctx, c, "job.list", params, limit,
	)
	if err != nil {
		return nil, err
	}

	// Client-side status filter.
	if status != "" {
		filtered := make([]Job, 0, len(jobs))
		for _, j := range jobs {
			if j.Status == status {
				filtered = append(filtered, j)
			}
		}

		return filtered, nil
	}

	return jobs, nil
}

// GetJob returns details for a single job by ID.
func (c *Client) GetJob(
	ctx context.Context, jobID string,
) (*Job, error) {

	var resp struct {
		Success bool `json:"success"`
		Results Job  `json:"results"`
	}

	err := c.Call(ctx, "job.info", map[string]any{
		"jobId": jobID,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// SearchJobs searches for jobs matching the given term.
func (c *Client) SearchJobs(
	ctx context.Context, term string, limit int,
) ([]Job, error) {

	var resp struct {
		Success bool              `json:"success"`
		Results []json.RawMessage `json:"results"`
	}

	err := c.Call(ctx, "job.search", map[string]any{
		"term": term,
	}, &resp)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var j Job
		if err := json.Unmarshal(raw, &j); err != nil {
			return nil, fmt.Errorf(
				"job.search: decode result: %w", err,
			)
		}
		jobs = append(jobs, j)
	}

	if limit > 0 && len(jobs) > limit {
		jobs = jobs[:limit]
	}

	return jobs, nil
}
