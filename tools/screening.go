package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ScreenCandidatesInput defines the input parameters for the
// screen_candidates tool.
type ScreenCandidatesInput struct {
	// JobID is the Ashby job ID to screen candidates for.
	JobID string `json:"jobId" jsonschema:"Ashby job ID to screen candidates for"`

	// Status filters by application status (default: Active).
	Status string `json:"status,omitempty" jsonschema:"Application status filter (default: Active)"`

	// Limit caps the number of applications to screen.
	Limit int `json:"limit,omitempty" jsonschema:"Maximum applications to screen"`

	// MinTier filters results to this tier or above: strong,
	// moderate, weak, no_signal.
	MinTier string `json:"minTier,omitempty" jsonschema:"Minimum tier to include: strong moderate weak no_signal"`

	// Enrich fetches expanded application details for better
	// scoring accuracy. This is slower due to per-application
	// API calls.
	Enrich bool `json:"enrich,omitempty" jsonschema:"Fetch expanded details for better scoring (slower)"`
}

// ScreenedCandidate holds the screening result for a single
// candidate.
type ScreenedCandidate struct {
	// CandidateID is the Ashby candidate ID.
	CandidateID string `json:"candidateId"`

	// CandidateName is the candidate's full name.
	CandidateName string `json:"candidateName"`

	// ApplicationID is the Ashby application ID.
	ApplicationID string `json:"applicationId"`

	// JobTitle is the name of the job applied to.
	JobTitle string `json:"jobTitle"`

	// Stage is the current interview stage.
	Stage string `json:"stage"`

	// Tier is the classification: strong, moderate, weak, or
	// no_signal.
	Tier string `json:"tier"`

	// Score is the detailed scoring breakdown.
	Score ScoreResult `json:"score"`
}

// TierSummary counts candidates in each tier.
type TierSummary struct {
	// Strong is the count of strong candidates.
	Strong int `json:"strong"`

	// Moderate is the count of moderate candidates.
	Moderate int `json:"moderate"`

	// Weak is the count of weak candidates.
	Weak int `json:"weak"`

	// NoSignal is the count of candidates with no signal.
	NoSignal int `json:"noSignal"`
}

// ScreenCandidatesOutput contains the screening results.
type ScreenCandidatesOutput struct {
	// ScreenedAt is the ISO 8601 timestamp of the screening.
	ScreenedAt string `json:"screenedAt"`

	// TotalScreened is the number of candidates screened.
	TotalScreened int `json:"totalScreened"`

	// Summary counts candidates in each tier.
	Summary TierSummary `json:"summary"`

	// Candidates is the list of screened candidates, sorted
	// by score descending.
	Candidates []ScreenedCandidate `json:"candidates"`
}

// tierRank returns a numeric rank for tier ordering, where
// lower values indicate better tiers.
func tierRank(tier string) int {
	switch tier {
	case "strong":
		return 0
	case "moderate":
		return 1
	case "weak":
		return 2
	default:
		return 3
	}
}

// tierMeetsMinimum checks whether a tier meets or exceeds the
// minimum tier threshold.
func tierMeetsMinimum(tier, minTier string) bool {
	if minTier == "" {
		return true
	}

	return tierRank(tier) <= tierRank(minTier)
}

// ScreenCandidates handles the screen_candidates MCP tool call.
// It fetches applications for a job, extracts searchable text,
// scores against Lightning Labs criteria, and returns ranked
// results.
func (h *Handler) ScreenCandidates(
	ctx context.Context, req *mcp.CallToolRequest,
	input ScreenCandidatesInput,
) (*mcp.CallToolResult, ScreenCandidatesOutput, error) {

	if input.JobID == "" {
		return nil, ScreenCandidatesOutput{},
			fmt.Errorf("jobId is required")
	}

	status := input.Status
	if status == "" {
		status = "Active"
	}

	// Fetch applications.
	result, err := h.client.ListApplications(
		ctx, ashby.ListApplicationsOpts{
			JobID:  input.JobID,
			Status: status,
			Limit:  input.Limit,
		},
	)
	if err != nil {
		return nil, ScreenCandidatesOutput{},
			fmt.Errorf("list applications: %w", err)
	}

	var screened []ScreenedCandidate
	var summary TierSummary

	for _, app := range result.Applications {
		// Optionally enrich with expanded details.
		var appData map[string]any
		if input.Enrich {
			enriched, err := h.client.GetApplication(
				ctx, app.ID, []string{
					"applicationFormSubmissions",
					"openings",
					"referrals",
				},
			)
			if err != nil {
				return nil, ScreenCandidatesOutput{},
					fmt.Errorf(
						"enrich %s: %w",
						app.ID, err,
					)
			}

			// Marshal/unmarshal to get a generic map
			// for text extraction.
			raw, err := json.Marshal(enriched)
			if err != nil {
				return nil, ScreenCandidatesOutput{},
					fmt.Errorf(
						"marshal app %s: %w",
						app.ID, err,
					)
			}
			if err := json.Unmarshal(
				raw, &appData,
			); err != nil {
				return nil, ScreenCandidatesOutput{},
					fmt.Errorf(
						"unmarshal app %s: %w",
						app.ID, err,
					)
			}
		} else {
			raw, err := json.Marshal(app)
			if err != nil {
				return nil, ScreenCandidatesOutput{},
					fmt.Errorf(
						"marshal app %s: %w",
						app.ID, err,
					)
			}
			if err := json.Unmarshal(
				raw, &appData,
			); err != nil {
				return nil, ScreenCandidatesOutput{},
					fmt.Errorf(
						"unmarshal app %s: %w",
						app.ID, err,
					)
			}
		}

		text := ExtractText(appData)
		score := ScoreCandidate(text)
		tier := ClassifyTier(score.Pct)

		// Apply tier filter.
		if !tierMeetsMinimum(tier, input.MinTier) {
			continue
		}

		// Extract display fields.
		candidateName := "Unknown"
		candidateID := app.CandidateID
		if app.Candidate != nil {
			candidateName = app.Candidate.Name
			candidateID = app.Candidate.ID
		}

		jobTitle := ""
		if app.Job != nil {
			jobTitle = app.Job.Title
		}

		stage := ""
		if app.CurrentInterviewStage != nil {
			stage = app.CurrentInterviewStage.Title
		}

		screened = append(screened, ScreenedCandidate{
			CandidateID:   candidateID,
			CandidateName: candidateName,
			ApplicationID: app.ID,
			JobTitle:      jobTitle,
			Stage:         stage,
			Tier:          tier,
			Score:         score,
		})

		// Update tier counts.
		switch tier {
		case "strong":
			summary.Strong++
		case "moderate":
			summary.Moderate++
		case "weak":
			summary.Weak++
		default:
			summary.NoSignal++
		}
	}

	// Sort by score descending.
	sort.Slice(screened, func(i, j int) bool {
		return screened[i].Score.Pct > screened[j].Score.Pct
	})

	return nil, ScreenCandidatesOutput{
		ScreenedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalScreened: len(screened),
		Summary:       summary,
		Candidates:    screened,
	}, nil
}
