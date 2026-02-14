package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// PipelineDashboardInput defines the input parameters for the
// pipeline_dashboard tool. No parameters are required.
type PipelineDashboardInput struct{}

// JobPipelineStats holds aggregated pipeline statistics for a
// single job.
type JobPipelineStats struct {
	// JobID is the job identifier.
	JobID string `json:"jobId"`

	// JobTitle is the human-readable job name.
	JobTitle string `json:"jobTitle"`

	// TotalApplications is the count of all applications.
	TotalApplications int `json:"totalApplications"`

	// ActiveApplications is the count of active applications.
	ActiveApplications int `json:"activeApplications"`

	// ByStatus maps application status to count.
	ByStatus map[string]int `json:"byStatus"`

	// ByStage maps interview stage title to count.
	ByStage map[string]int `json:"byStage"`
}

// PipelineDashboardOutput contains the aggregated pipeline
// dashboard data.
type PipelineDashboardOutput struct {
	// Jobs is the per-job pipeline breakdown.
	Jobs []JobPipelineStats `json:"jobs"`

	// TotalJobs is the number of open jobs.
	TotalJobs int `json:"totalJobs"`

	// TotalApplications is the total count across all jobs.
	TotalApplications int `json:"totalApplications"`
}

// PipelineDashboard handles the pipeline_dashboard MCP tool
// call. It fetches all open jobs and their applications, then
// aggregates counts by status and interview stage.
func (h *Handler) PipelineDashboard(
	ctx context.Context, req *mcp.CallToolRequest,
	input PipelineDashboardInput,
) (*mcp.CallToolResult, PipelineDashboardOutput, error) {

	// Fetch all open jobs.
	jobs, err := h.client.ListJobs(ctx, "Open", 0)
	if err != nil {
		return nil, PipelineDashboardOutput{},
			fmt.Errorf("list jobs: %w", err)
	}

	var (
		stats    []JobPipelineStats
		totalAll int
	)

	// For each job, fetch applications and aggregate.
	for _, job := range jobs {
		result, err := h.client.ListApplications(
			ctx, ashby.ListApplicationsOpts{
				JobID: job.ID,
			},
		)
		if err != nil {
			return nil, PipelineDashboardOutput{},
				fmt.Errorf("list apps for %s: %w",
					job.ID, err,
				)
		}

		apps := result.Applications
		byStatus := make(map[string]int)
		byStage := make(map[string]int)
		active := 0

		for _, app := range apps {
			byStatus[app.Status]++

			if app.Status == "Active" {
				active++
			}

			if app.CurrentInterviewStage != nil {
				stage := app.CurrentInterviewStage.Title
				byStage[stage]++
			}
		}

		stats = append(stats, JobPipelineStats{
			JobID:              job.ID,
			JobTitle:           job.Title,
			TotalApplications:  len(apps),
			ActiveApplications: active,
			ByStatus:           byStatus,
			ByStage:            byStage,
		})

		totalAll += len(apps)
	}

	return nil, PipelineDashboardOutput{
		Jobs:              stats,
		TotalJobs:         len(jobs),
		TotalApplications: totalAll,
	}, nil
}
