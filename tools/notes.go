package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// CreateCandidateNoteInput defines the input parameters for the
// create_candidate_note tool.
type CreateCandidateNoteInput struct {
	// CandidateID is the candidate to add a note to.
	CandidateID string `json:"candidateId" jsonschema:"description=The Ashby candidate ID"`

	// Body is the HTML-formatted note content.
	Body string `json:"body" jsonschema:"description=HTML-formatted note content"`
}

// CreateCandidateNoteOutput confirms the note was created.
type CreateCandidateNoteOutput struct {
	// Success indicates the note was created.
	Success bool `json:"success"`
}

// CreateCandidateNote handles the create_candidate_note MCP
// tool call.
func (h *Handler) CreateCandidateNote(
	ctx context.Context, req *mcp.CallToolRequest,
	input CreateCandidateNoteInput,
) (*mcp.CallToolResult, CreateCandidateNoteOutput, error) {

	err := h.client.CreateCandidateNote(
		ctx, input.CandidateID, input.Body,
	)
	if err != nil {
		return nil, CreateCandidateNoteOutput{}, err
	}

	return nil, CreateCandidateNoteOutput{
		Success: true,
	}, nil
}

// ListCandidateNotesInput defines the input parameters for the
// list_candidate_notes tool.
type ListCandidateNotesInput struct {
	// CandidateID is the candidate to list notes for.
	CandidateID string `json:"candidateId" jsonschema:"description=The Ashby candidate ID"`
}

// ListCandidateNotesOutput contains the list_candidate_notes
// results.
type ListCandidateNotesOutput struct {
	// Notes is the list of notes for the candidate.
	Notes []ashby.Note `json:"notes"`

	// Total is the number of notes returned.
	Total int `json:"total"`
}

// ListCandidateNotes handles the list_candidate_notes MCP tool
// call.
func (h *Handler) ListCandidateNotes(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListCandidateNotesInput,
) (*mcp.CallToolResult, ListCandidateNotesOutput, error) {

	notes, err := h.client.ListCandidateNotes(
		ctx, input.CandidateID,
	)
	if err != nil {
		return nil, ListCandidateNotesOutput{}, err
	}

	return nil, ListCandidateNotesOutput{
		Notes: notes,
		Total: len(notes),
	}, nil
}
