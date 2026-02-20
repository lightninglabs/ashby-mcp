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

// SetJobStatus changes the status of a job. Valid values for
// status are Open, Closed, and Archived.
func (c *Client) SetJobStatus(
	ctx context.Context, jobID, status string,
) (*Job, error) {

	var resp struct {
		Success bool `json:"success"`
		Results Job  `json:"results"`
	}

	if err := c.Call(ctx, "job.setStatus", map[string]any{
		"jobId":  jobID,
		"status": status,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// UpdateJobOpts holds optional fields that may be updated on a
// job record.
type UpdateJobOpts struct {
	// Title is the job's display title.
	Title string

	// DepartmentID references the department for this job.
	DepartmentID string

	// LocationIds lists the associated location IDs.
	LocationIds []string

	// EmploymentType is the type of employment
	// (e.g. "FullTime").
	EmploymentType string
}

// UpdateJob updates mutable fields on an existing job. Only
// fields with non-zero values are sent.
func (c *Client) UpdateJob(
	ctx context.Context, jobID string, opts UpdateJobOpts,
) (*Job, error) {

	params := map[string]any{
		"jobId": jobID,
	}

	if opts.Title != "" {
		params["title"] = opts.Title
	}
	if opts.DepartmentID != "" {
		params["departmentId"] = opts.DepartmentID
	}
	if len(opts.LocationIds) > 0 {
		params["locationIds"] = opts.LocationIds
	}
	if opts.EmploymentType != "" {
		params["employmentType"] = opts.EmploymentType
	}

	var resp struct {
		Success bool `json:"success"`
		Results Job  `json:"results"`
	}

	if err := c.Call(
		ctx, "job.update", params, &resp,
	); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}
