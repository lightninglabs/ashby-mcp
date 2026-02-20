package ashby

import "encoding/json"

// Job represents an Ashby job posting.
type Job struct {
	// ID is the unique identifier for the job.
	ID string `json:"id"`

	// Title is the display name of the job posting.
	Title string `json:"title"`

	// Status is the current state: Open, Closed, Archived, or
	// Draft.
	Status string `json:"status"`

	// Confidentiality indicates the visibility level of the job.
	Confidentiality string `json:"confidentiality,omitempty"`

	// EmploymentType is the type of employment (e.g.
	// FullTime).
	EmploymentType string `json:"employmentType,omitempty"`

	// Department contains the department this job belongs to.
	Department *Department `json:"department,omitempty"`

	// Location contains the location for this job.
	Location *Location `json:"location,omitempty"`

	// Team contains the team this job belongs to.
	Team *Team `json:"team,omitempty"`

	// CustomFields holds job-level custom field values.
	CustomFields []CustomField `json:"customFields,omitempty"`

	// JobPostingIds lists associated job posting identifiers.
	JobPostingIds []string `json:"jobPostingIds,omitempty"`

	// InterviewPlanID is the ID of the interview plan for this
	// job.
	InterviewPlanID string `json:"interviewPlanId,omitempty"`

	// CreatedAt is the ISO 8601 creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the ISO 8601 last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// Department represents an organizational department.
type Department struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the display name of the department.
	Name string `json:"name"`
}

// Location represents a physical or remote work location.
type Location struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the display name of the location.
	Name string `json:"name"`
}

// Team represents an organizational team.
type Team struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the display name of the team.
	Name string `json:"name"`
}

// Application represents a candidate's application to a job.
type Application struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Status is the current state: Active, Hired, Archived, or
	// Rejected.
	Status string `json:"status"`

	// CandidateID references the associated candidate.
	CandidateID string `json:"candidateId"`

	// JobID references the associated job.
	JobID string `json:"jobId"`

	// Candidate contains expanded candidate details when
	// requested.
	Candidate *Candidate `json:"candidate,omitempty"`

	// Job contains expanded job details when requested.
	Job *Job `json:"job,omitempty"`

	// CurrentInterviewStage is the current pipeline stage.
	CurrentInterviewStage *InterviewStage `json:"currentInterviewStage,omitempty"`

	// Source describes how the candidate was sourced.
	Source *Source `json:"source,omitempty"`

	// CustomFields holds application-level custom field values.
	CustomFields []CustomField `json:"customFields,omitempty"`

	// ApplicationFormSubmissions contains form response data
	// when expanded. Uses any because the schema varies by
	// form and the MCP SDK must not constrain item types.
	ApplicationFormSubmissions []any `json:"applicationFormSubmissions,omitempty"`

	// ResumeFileHandle holds resume metadata when expanded.
	ResumeFileHandle *FileHandle `json:"resumeFileHandle,omitempty"`

	// Openings contains expanded opening details.
	Openings []any `json:"openings,omitempty"`

	// Referrals contains expanded referral details.
	Referrals []any `json:"referrals,omitempty"`

	// HiringTeam lists the hiring team members.
	HiringTeam []HiringTeamMember `json:"hiringTeam,omitempty"`

	// CreatedAt is the ISO 8601 creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the ISO 8601 last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`

	// ArchivedAt is set when the application was archived.
	ArchivedAt string `json:"archivedAt,omitempty"`
}

// Candidate represents a person in the ATS.
type Candidate struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the candidate's full name.
	Name string `json:"name"`

	// PrimaryEmailAddress is the main email contact.
	PrimaryEmailAddress *EmailAddress `json:"primaryEmailAddress,omitempty"`

	// PhoneNumbers lists the candidate's phone numbers.
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`

	// Tags lists tags applied to this candidate.
	Tags []Tag `json:"tags,omitempty"`

	// CustomFields holds candidate-level custom field values.
	CustomFields []CustomField `json:"customFields,omitempty"`

	// SocialLinks lists the candidate's social/web profiles.
	SocialLinks []SocialLink `json:"socialLinks,omitempty"`

	// Source describes how the candidate was sourced.
	Source *Source `json:"source,omitempty"`

	// ResumeFileHandle is the primary resume file reference.
	ResumeFileHandle *FileHandle `json:"resumeFileHandle,omitempty"`

	// FileHandles lists all files attached to this candidate
	// (resumes, cover letters, portfolios, etc.).
	FileHandles []FileHandle `json:"fileHandles,omitempty"`

	// Position is the candidate's current job title.
	Position string `json:"position,omitempty"`

	// Company is the candidate's current employer.
	Company string `json:"company,omitempty"`

	// School is the candidate's current or most recent school.
	School string `json:"school,omitempty"`

	// ApplicationIds lists IDs of applications linked to this
	// candidate.
	ApplicationIds []string `json:"applicationIds,omitempty"`

	// ProfileURL is the Ashby profile URL for this candidate.
	ProfileURL string `json:"profileUrl,omitempty"`

	// CreatedAt is the ISO 8601 creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`

	// UpdatedAt is the ISO 8601 last update timestamp.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// InterviewStage represents a stage in the hiring pipeline.
type InterviewStage struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Title is the display name (e.g. "Phone Screen").
	Title string `json:"title"`

	// Type categorizes the stage (e.g. "PhoneScreen",
	// "TakeHome").
	Type string `json:"type,omitempty"`

	// OrderInInterviewPlan is the position in the pipeline.
	OrderInInterviewPlan int `json:"orderInInterviewPlan,omitempty"`

	// InterviewPlanID references the parent plan.
	InterviewPlanID string `json:"interviewPlanId,omitempty"`
}

// Interview represents a scheduled interview event.
type Interview struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// InterviewStageID references the pipeline stage.
	InterviewStageID string `json:"interviewStageId,omitempty"`

	// ApplicationID references the application.
	ApplicationID string `json:"applicationId,omitempty"`

	// Status is the interview state (e.g. "Scheduled").
	Status string `json:"status,omitempty"`

	// StartTime is the ISO 8601 scheduled start.
	StartTime string `json:"startTime,omitempty"`

	// EndTime is the ISO 8601 scheduled end.
	EndTime string `json:"endTime,omitempty"`

	// Interviewers lists the interview participants.
	Interviewers []Interviewer `json:"interviewers,omitempty"`
}

// Interviewer represents a person conducting an interview.
type Interviewer struct {
	// Name is the interviewer's full name.
	Name string `json:"name"`

	// Email is the interviewer's email address.
	Email string `json:"email,omitempty"`

	// UserID references the Ashby user.
	UserID string `json:"userId,omitempty"`
}

// Tag represents a candidate tag.
type Tag struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Title is the display name of the tag.
	Title string `json:"title"`
}

// Note represents a note attached to a candidate.
type Note struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Body contains the HTML-formatted note content.
	Body string `json:"body"`

	// Author is the user who created the note.
	Author *NoteAuthor `json:"author,omitempty"`

	// CreatedAt is the ISO 8601 creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`
}

// NoteAuthor identifies who created a note.
type NoteAuthor struct {
	// ID is the unique user identifier.
	ID string `json:"id"`

	// Name is the author's full name.
	Name string `json:"name"`

	// Email is the author's email address.
	Email string `json:"email,omitempty"`
}

// Source describes how a candidate or application was sourced.
type Source struct {
	// ID is the unique identifier.
	ID string `json:"id,omitempty"`

	// Title is the display name (e.g. "LinkedIn", "Referral").
	Title string `json:"title"`
}

// CustomField holds a custom field value on an entity.
type CustomField struct {
	// ID is the unique field identifier.
	ID string `json:"id,omitempty"`

	// Title is the display name of the field.
	Title string `json:"title"`

	// Value is the raw field value.
	Value any `json:"value"`

	// ValueLabel is the human-readable value label.
	ValueLabel string `json:"valueLabel,omitempty"`
}

// EmailAddress represents an email contact.
type EmailAddress struct {
	// Value is the email address string.
	Value string `json:"value"`

	// Type classifies the email (e.g. "Personal", "Work").
	Type string `json:"type,omitempty"`

	// IsPrimary indicates whether this is the primary email.
	IsPrimary bool `json:"isPrimary,omitempty"`
}

// PhoneNumber represents a phone contact.
type PhoneNumber struct {
	// Value is the phone number string.
	Value string `json:"value"`

	// Type classifies the number (e.g. "Mobile", "Home").
	Type string `json:"type,omitempty"`
}

// SocialLink represents a social or web profile link.
type SocialLink struct {
	// URL is the profile URL.
	URL string `json:"url"`

	// Type classifies the link (e.g. "LinkedIn", "GitHub").
	Type string `json:"type,omitempty"`
}

// FileHandle represents a file attachment reference.
type FileHandle struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the filename.
	Name string `json:"name"`

	// Handle is the opaque token used with the file.info
	// endpoint to retrieve a pre-signed download URL.
	Handle string `json:"handle,omitempty"`
}

// HiringTeamMember represents a member of the hiring team.
type HiringTeamMember struct {
	// UserID references the Ashby user.
	UserID string `json:"userId"`

	// Name is the team member's full name.
	Name string `json:"name"`

	// Email is the team member's email address.
	Email string `json:"email,omitempty"`

	// Role is the team member's role (e.g. "Hiring Manager").
	Role string `json:"role,omitempty"`
}

// ApplicationHistory represents a history event for an
// application.
type ApplicationHistory struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Action describes what happened.
	Action string `json:"action"`

	// OccurredAt is the ISO 8601 timestamp.
	OccurredAt string `json:"occurredAt"`

	// Actor is who performed the action.
	Actor *NoteAuthor `json:"actor,omitempty"`

	// Details holds action-specific data.
	Details json.RawMessage `json:"details,omitempty"`
}

// APIResponse is the common response envelope from Ashby. The
// Results field holds the decoded payload for single-object
// endpoints.
type APIResponse struct {
	// Success indicates whether the API call succeeded.
	Success bool `json:"success"`

	// Results holds the response payload.
	Results json.RawMessage `json:"results,omitempty"`

	// ErrorInfo contains error details on failure.
	ErrorInfo *ErrorInfo `json:"errorInfo,omitempty"`

	// Errors holds error messages when ErrorInfo is absent.
	Errors []string `json:"errors,omitempty"`
}

// ErrorInfo contains details about an API error.
type ErrorInfo struct {
	// Message is the human-readable error description.
	Message string `json:"message"`

	// Code is the machine-readable error code.
	Code string `json:"code,omitempty"`
}
