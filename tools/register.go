package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// readOnly is a convenience helper that returns a pointer to
// false, used for the DestructiveHint field on read-only tool
// annotations.
func ptrBool(b bool) *bool {
	return &b
}

// RegisterAll registers all Ashby MCP tools on the given
// server. Tools are organized into read-only query tools and
// write tools, each annotated with appropriate hints.
func RegisterAll(s *mcp.Server, h *Handler) {
	readOnly := &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: ptrBool(false),
	}
	writeIdempotent := &mcp.ToolAnnotations{
		IdempotentHint: true,
	}
	writeNonIdempotent := &mcp.ToolAnnotations{}

	// =============================================================
	// Job tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_jobs",
		Description: "List Ashby jobs with optional status " +
			"filter (Open, Closed, Archived, Draft).",
		Annotations: readOnly,
	}, h.ListJobs)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_job",
		Description: "Get detailed information about a specific Ashby job.",
		Annotations: readOnly,
	}, h.GetJob)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_jobs",
		Description: "Search Ashby jobs by title or keyword.",
		Annotations: readOnly,
	}, h.SearchJobs)

	// =============================================================
	// Application tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_applications",
		Description: "List Ashby applications with optional " +
			"filters (jobId, status). Supports cursor " +
			"pagination.",
		Annotations: readOnly,
	}, h.ListApplications)

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_application",
		Description: "Get detailed Ashby application info. " +
			"Use expand to include form submissions, " +
			"openings, or referrals.",
		Annotations: readOnly,
	}, h.GetApplication)

	mcp.AddTool(s, &mcp.Tool{
		Name: "change_application_stage",
		Description: "Move an Ashby application to a " +
			"different interview stage.",
		Annotations: writeIdempotent,
	}, h.ChangeApplicationStage)

	mcp.AddTool(s, &mcp.Tool{
		Name: "create_application",
		Description: "Create a new Ashby application " +
			"linking a candidate to a job.",
		Annotations: writeNonIdempotent,
	}, h.CreateApplication)

	// =============================================================
	// Candidate tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_candidates",
		Description: "List Ashby candidates with " +
			"pagination.",
		Annotations: readOnly,
	}, h.ListCandidates)

	mcp.AddTool(s, &mcp.Tool{
		Name: "search_candidates",
		Description: "Search Ashby candidates by email " +
			"and/or name.",
		Annotations: readOnly,
	}, h.SearchCandidates)

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_candidate",
		Description: "Get detailed Ashby candidate " +
			"profile by ID.",
		Annotations: readOnly,
	}, h.GetCandidate)

	mcp.AddTool(s, &mcp.Tool{
		Name: "create_candidate",
		Description: "Create a new candidate record in " +
			"Ashby.",
		Annotations: writeNonIdempotent,
	}, h.CreateCandidate)

	// =============================================================
	// Tag tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_tags",
		Description: "List all Ashby candidate tags.",
		Annotations: readOnly,
	}, h.ListTags)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_candidate_tag",
		Description: "Add a tag to an Ashby candidate.",
		Annotations: writeIdempotent,
	}, h.AddCandidateTag)

	// =============================================================
	// Note tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "create_candidate_note",
		Description: "Add an HTML-formatted note to an " +
			"Ashby candidate.",
		Annotations: writeNonIdempotent,
	}, h.CreateCandidateNote)

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_candidate_notes",
		Description: "List notes for an Ashby candidate.",
		Annotations: readOnly,
	}, h.ListCandidateNotes)

	// =============================================================
	// Interview tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_interview_stages",
		Description: "List interview stages for a job's " +
			"interview plan.",
		Annotations: readOnly,
	}, h.ListInterviewStages)

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_interviews",
		Description: "List Ashby interviews, optionally " +
			"filtered by application.",
		Annotations: readOnly,
	}, h.ListInterviews)

	// =============================================================
	// Analytics & screening tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "pipeline_dashboard",
		Description: "Get aggregated pipeline statistics " +
			"across all open Ashby jobs. Shows total " +
			"applications, active count, and breakdown " +
			"by status and interview stage per job.",
		Annotations: readOnly,
	}, h.PipelineDashboard)

	mcp.AddTool(s, &mcp.Tool{
		Name: "screen_candidates",
		Description: "Score candidates against Lightning " +
			"Labs hiring criteria using weighted " +
			"keyword matching. Returns ranked results " +
			"with per-category breakdowns and tier " +
			"classifications (strong/moderate/weak/" +
			"no_signal). Use enrich=true for better " +
			"accuracy (slower, fetches expanded " +
			"application details).",
		Annotations: readOnly,
	}, h.ScreenCandidates)
}
