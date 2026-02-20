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

	mcp.AddTool(s, &mcp.Tool{
		Name: "set_job_status",
		Description: "Set the status of an Ashby job " +
			"(Open, Closed, or Archived).",
		Annotations: writeIdempotent,
	}, h.SetJobStatus)

	mcp.AddTool(s, &mcp.Tool{
		Name: "update_job",
		Description: "Update mutable fields on an Ashby " +
			"job: title, departmentId, locationIds, " +
			"employmentType.",
		Annotations: writeIdempotent,
	}, h.UpdateJob)

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
		Name: "change_application_source",
		Description: "Set or clear the source on an Ashby " +
			"application. Pass an empty sourceId to " +
			"unset the source.",
		Annotations: writeNonIdempotent,
	}, h.ChangeApplicationSource)

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

	mcp.AddTool(s, &mcp.Tool{
		Name: "update_candidate",
		Description: "Update mutable fields on an " +
			"existing Ashby candidate: name, email, " +
			"phoneNumber, linkedInUrl, websiteUrl, " +
			"githubUrl, twitterHandle, " +
			"alternativeEmailAddresses, sourceId, " +
			"creditedToUserId.",
		Annotations: writeIdempotent,
	}, h.UpdateCandidate)

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
	// User tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_users",
		Description: "List Ashby team members, optionally " +
			"filtered by name.",
		Annotations: readOnly,
	}, h.ListUsers)

	mcp.AddTool(s, &mcp.Tool{
		Name: "search_users",
		Description: "Search Ashby users by name or " +
			"email.",
		Annotations: readOnly,
	}, h.SearchUsers)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user",
		Description: "Get details for an Ashby user by ID.",
		Annotations: readOnly,
	}, h.GetUser)

	// =============================================================
	// Lookup list tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_sources",
		Description: "List all Ashby application sources.",
		Annotations: readOnly,
	}, h.ListSources)

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_archive_reasons",
		Description: "List all Ashby application archive " +
			"reasons.",
		Annotations: readOnly,
	}, h.ListArchiveReasons)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_departments",
		Description: "List all Ashby departments.",
		Annotations: readOnly,
	}, h.ListDepartments)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_locations",
		Description: "List all Ashby locations.",
		Annotations: readOnly,
	}, h.ListLocations)

	// =============================================================
	// Job posting tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_job_postings",
		Description: "List all Ashby public job postings.",
		Annotations: readOnly,
	}, h.ListJobPostings)

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_job_posting",
		Description: "Get details for a specific Ashby " +
			"job posting by ID.",
		Annotations: readOnly,
	}, h.GetJobPosting)

	// =============================================================
	// Opening tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_openings",
		Description: "List all Ashby headcount openings.",
		Annotations: readOnly,
	}, h.ListOpenings)

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_opening",
		Description: "Get details for a specific Ashby " +
			"opening by ID.",
		Annotations: readOnly,
	}, h.GetOpening)

	mcp.AddTool(s, &mcp.Tool{
		Name: "search_openings",
		Description: "Search Ashby openings by keyword.",
		Annotations: readOnly,
	}, h.SearchOpenings)

	// =============================================================
	// Interview plan tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "list_interview_plans",
		Description: "List all Ashby interview plans.",
		Annotations: readOnly,
	}, h.ListInterviewPlans)

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_interview_stage",
		Description: "Get details for a specific Ashby " +
			"interview stage by ID.",
		Annotations: readOnly,
	}, h.GetInterviewStage)

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
	// File tools.
	// =============================================================

	mcp.AddTool(s, &mcp.Tool{
		Name: "get_file_url",
		Description: "Retrieve a pre-signed download URL " +
			"for a candidate file (resume, cover " +
			"letter) using its opaque handle string " +
			"from a resumeFileHandle or fileHandles " +
			"entry.",
		Annotations: readOnly,
	}, h.GetFileURL)

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
