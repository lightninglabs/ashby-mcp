package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListCandidatesInput defines the input parameters for the
// list_candidates tool.
type ListCandidatesInput struct {
	// Limit caps the maximum number of results returned.
	Limit int `json:"limit,omitempty" jsonschema:"description=Maximum number of results to return (default: 100)"`
}

// ListCandidatesOutput contains the list_candidates results.
type ListCandidatesOutput struct {
	// Candidates is the list of candidates.
	Candidates []ashby.Candidate `json:"candidates"`

	// Total is the number of candidates returned.
	Total int `json:"total"`
}

// ListCandidates handles the list_candidates MCP tool call.
func (h *Handler) ListCandidates(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListCandidatesInput,
) (*mcp.CallToolResult, ListCandidatesOutput, error) {

	cands, err := h.client.ListCandidates(
		ctx, input.Limit,
	)
	if err != nil {
		return nil, ListCandidatesOutput{}, err
	}

	return nil, ListCandidatesOutput{
		Candidates: cands,
		Total:      len(cands),
	}, nil
}

// SearchCandidatesInput defines the input parameters for the
// search_candidates tool.
type SearchCandidatesInput struct {
	// Email searches by email address.
	Email string `json:"email,omitempty" jsonschema:"description=Search by email address"`

	// Name searches by candidate name.
	Name string `json:"name,omitempty" jsonschema:"description=Search by candidate name"`
}

// SearchCandidatesOutput contains the search_candidates
// results.
type SearchCandidatesOutput struct {
	// Candidates is the list of matching candidates.
	Candidates []ashby.Candidate `json:"candidates"`

	// Total is the number of matching candidates.
	Total int `json:"total"`
}

// SearchCandidates handles the search_candidates MCP tool call.
func (h *Handler) SearchCandidates(
	ctx context.Context, req *mcp.CallToolRequest,
	input SearchCandidatesInput,
) (*mcp.CallToolResult, SearchCandidatesOutput, error) {

	cands, err := h.client.SearchCandidates(
		ctx, input.Email, input.Name,
	)
	if err != nil {
		return nil, SearchCandidatesOutput{}, err
	}

	return nil, SearchCandidatesOutput{
		Candidates: cands,
		Total:      len(cands),
	}, nil
}

// GetCandidateInput defines the input parameters for the
// get_candidate tool.
type GetCandidateInput struct {
	// CandidateID is the Ashby candidate ID to look up.
	CandidateID string `json:"candidateId" jsonschema:"description=The Ashby candidate ID"`
}

// GetCandidateOutput contains the get_candidate results.
type GetCandidateOutput struct {
	// Candidate is the candidate details.
	Candidate *ashby.Candidate `json:"candidate"`
}

// GetCandidate handles the get_candidate MCP tool call.
func (h *Handler) GetCandidate(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetCandidateInput,
) (*mcp.CallToolResult, GetCandidateOutput, error) {

	cand, err := h.client.GetCandidate(
		ctx, input.CandidateID,
	)
	if err != nil {
		return nil, GetCandidateOutput{}, err
	}

	return nil, GetCandidateOutput{Candidate: cand}, nil
}

// CreateCandidateInput defines the input parameters for the
// create_candidate tool.
type CreateCandidateInput struct {
	// Name is the candidate's full name.
	Name string `json:"name" jsonschema:"description=Candidate full name"`

	// Email is the candidate's email address.
	Email string `json:"email" jsonschema:"description=Candidate email address"`

	// Phone is an optional phone number.
	Phone string `json:"phone,omitempty" jsonschema:"description=Optional phone number"`
}

// CreateCandidateOutput contains the newly created candidate.
type CreateCandidateOutput struct {
	// Candidate is the newly created candidate.
	Candidate *ashby.Candidate `json:"candidate"`
}

// CreateCandidate handles the create_candidate MCP tool call.
func (h *Handler) CreateCandidate(
	ctx context.Context, req *mcp.CallToolRequest,
	input CreateCandidateInput,
) (*mcp.CallToolResult, CreateCandidateOutput, error) {

	cand, err := h.client.CreateCandidate(
		ctx, input.Name, input.Email, input.Phone,
	)
	if err != nil {
		return nil, CreateCandidateOutput{}, err
	}

	return nil, CreateCandidateOutput{Candidate: cand}, nil
}
