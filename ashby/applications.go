package ashby

import (
	"context"
)

// ListApplicationsOpts configures a call to list applications.
type ListApplicationsOpts struct {
	// JobID filters by job.
	JobID string

	// Status filters by application status (Active, Hired,
	// Archived, Rejected).
	Status string

	// Limit caps the total number of results.
	Limit int

	// Cursor is the pagination cursor for resuming from a
	// prior page.
	Cursor string
}

// ListApplicationsResult holds a page of applications with
// pagination metadata.
type ListApplicationsResult struct {
	// Applications is the list of results.
	Applications []Application

	// NextCursor is the cursor for fetching the next page.
	NextCursor string

	// MoreDataAvailable indicates additional pages exist.
	MoreDataAvailable bool
}

// ListApplications returns applications matching the given
// filters. When Cursor is set, a single page is fetched and
// pagination metadata is returned. When Cursor is empty, all
// pages are fetched up to Limit.
func (c *Client) ListApplications(
	ctx context.Context, opts ListApplicationsOpts,
) (*ListApplicationsResult, error) {

	params := make(map[string]any)
	if opts.JobID != "" {
		params["jobId"] = opts.JobID
	}
	if opts.Status != "" {
		params["status"] = opts.Status
	}

	// If a cursor is provided, fetch a single page for the
	// MCP tool's passthrough pagination.
	if opts.Cursor != "" {
		params["cursor"] = opts.Cursor

		page, err := FetchPage[Application](
			ctx, c, "application.list", params,
		)
		if err != nil {
			return nil, err
		}

		return &ListApplicationsResult{
			Applications:      page.Items,
			NextCursor:        page.NextCursor,
			MoreDataAvailable: page.MoreDataAvailable,
		}, nil
	}

	// No cursor: fetch all pages.
	apps, err := Paginate[Application](
		ctx, c, "application.list", params, opts.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &ListApplicationsResult{
		Applications: apps,
	}, nil
}

// GetApplication returns details for a single application by
// ID. The expand slice controls which related data to include
// (e.g. "applicationFormSubmissions", "openings", "referrals").
func (c *Client) GetApplication(
	ctx context.Context, appID string, expand []string,
) (*Application, error) {

	params := map[string]any{
		"applicationId": appID,
	}
	if len(expand) > 0 {
		params["expand"] = expand
	}

	var resp struct {
		Success bool        `json:"success"`
		Results Application `json:"results"`
	}

	if err := c.Call(
		ctx, "application.info", params, &resp,
	); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// ChangeApplicationStage moves an application to a different
// interview stage.
func (c *Client) ChangeApplicationStage(
	ctx context.Context, appID, stageID string,
) error {

	var resp struct {
		Success bool `json:"success"`
	}

	return c.Call(ctx, "application.changeStage", map[string]any{
		"applicationId":    appID,
		"interviewStageId": stageID,
	}, &resp)
}

// CreateApplication creates a new application linking a
// candidate to a job. Source is optional.
func (c *Client) CreateApplication(
	ctx context.Context, candidateID, jobID, source string,
) (*Application, error) {

	params := map[string]any{
		"candidateId": candidateID,
		"jobId":       jobID,
	}
	if source != "" {
		params["source"] = source
	}

	var resp struct {
		Success bool        `json:"success"`
		Results Application `json:"results"`
	}

	if err := c.Call(
		ctx, "application.create", params, &resp,
	); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}
