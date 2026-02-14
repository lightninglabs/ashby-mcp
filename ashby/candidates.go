package ashby

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListCandidates returns all candidates with pagination.
func (c *Client) ListCandidates(
	ctx context.Context, limit int,
) ([]Candidate, error) {

	return Paginate[Candidate](
		ctx, c, "candidate.list", nil, limit,
	)
}

// SearchCandidates searches for candidates by email and/or
// name. At least one of email or name must be provided.
func (c *Client) SearchCandidates(
	ctx context.Context, email, name string,
) ([]Candidate, error) {

	params := make(map[string]any)
	if email != "" {
		params["email"] = email
	}
	if name != "" {
		params["name"] = name
	}

	var resp struct {
		Success bool              `json:"success"`
		Results []json.RawMessage `json:"results"`
	}

	if err := c.Call(
		ctx, "candidate.search", params, &resp,
	); err != nil {
		return nil, err
	}

	candidates := make([]Candidate, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var cand Candidate
		if err := json.Unmarshal(raw, &cand); err != nil {
			return nil, fmt.Errorf(
				"candidate.search: decode: %w", err,
			)
		}
		candidates = append(candidates, cand)
	}

	return candidates, nil
}

// GetCandidate returns details for a single candidate by ID.
func (c *Client) GetCandidate(
	ctx context.Context, candidateID string,
) (*Candidate, error) {

	var resp struct {
		Success bool      `json:"success"`
		Results Candidate `json:"results"`
	}

	if err := c.Call(ctx, "candidate.info", map[string]any{
		"candidateId": candidateID,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// CreateCandidate creates a new candidate record. Name and
// email are required; phone is optional.
func (c *Client) CreateCandidate(
	ctx context.Context, name, email, phone string,
) (*Candidate, error) {

	params := map[string]any{
		"name":  name,
		"email": email,
	}
	if phone != "" {
		params["phoneNumber"] = phone
	}

	var resp struct {
		Success bool      `json:"success"`
		Results Candidate `json:"results"`
	}

	if err := c.Call(
		ctx, "candidate.create", params, &resp,
	); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}
