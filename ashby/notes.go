package ashby

import "context"

// CreateCandidateNote adds an HTML-formatted note to a
// candidate.
func (c *Client) CreateCandidateNote(
	ctx context.Context, candidateID, body string,
) error {

	var resp struct {
		Success bool `json:"success"`
	}

	return c.Call(ctx, "candidate.createNote", map[string]any{
		"candidateId": candidateID,
		"note":        body,
	}, &resp)
}

// ListCandidateNotes returns all notes for a candidate.
func (c *Client) ListCandidateNotes(
	ctx context.Context, candidateID string,
) ([]Note, error) {

	return Paginate[Note](
		ctx, c, "candidate.listNotes", map[string]any{
			"candidateId": candidateID,
		}, 0,
	)
}
