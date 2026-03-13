package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetFileURLInput defines the input parameters for the
// get_file_url tool.
type GetFileURLInput struct {
	// FileHandle is the opaque handle string from a
	// resumeFileHandle or fileHandles entry.
	FileHandle string `json:"fileHandle" jsonschema:"Opaque file handle string from a resumeFileHandle or fileHandles entry"`
}

// GetFileURLOutput contains the get_file_url results.
type GetFileURLOutput struct {
	// URL is the pre-signed download URL for the file.
	URL string `json:"url"`
}

// GetFileURL handles the get_file_url MCP tool call.
func (h *Handler) GetFileURL(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetFileURLInput,
) (*mcp.CallToolResult, GetFileURLOutput, error) {

	url, err := h.client.GetFileURL(ctx, input.FileHandle)
	if err != nil {
		return nil, GetFileURLOutput{}, err
	}

	return nil, GetFileURLOutput{URL: url}, nil
}

// GetCandidateResumeInput defines the input parameters for the
// get_candidate_resume tool.
type GetCandidateResumeInput struct {
	// Email is the candidate's email address to search by.
	Email string `json:"email,omitempty" jsonschema:"Candidate email address to search by"`

	// Name is the candidate's name to search by.
	Name string `json:"name,omitempty" jsonschema:"Candidate name to search by"`
}

// GetCandidateResumeOutput contains the extracted resume text
// and metadata.
type GetCandidateResumeOutput struct {
	// CandidateName is the candidate's full name.
	CandidateName string `json:"candidateName"`

	// FileName is the original resume filename.
	FileName string `json:"fileName"`

	// Text is the extracted plain text content of the resume.
	Text string `json:"text"`
}

// GetCandidateResume handles the get_candidate_resume MCP tool
// call. It searches for a candidate by email or name, fetches
// the resume PDF, extracts the plain text, and returns it.
func (h *Handler) GetCandidateResume(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetCandidateResumeInput,
) (*mcp.CallToolResult, GetCandidateResumeOutput, error) {

	resume, err := h.client.GetCandidateResume(
		ctx, input.Email, input.Name,
	)
	if err != nil {
		return nil, GetCandidateResumeOutput{}, err
	}

	return nil, GetCandidateResumeOutput{
		CandidateName: resume.CandidateName,
		FileName:      resume.FileName,
		Text:          resume.Text,
	}, nil
}
