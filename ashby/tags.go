package ashby

import "context"

// ListTags returns all candidate tags.
func (c *Client) ListTags(
	ctx context.Context,
) ([]Tag, error) {

	return Paginate[Tag](
		ctx, c, "candidateTag.list", nil, 0,
	)
}

// AddCandidateTag adds a tag to a candidate. Both the candidate
// ID and the tag ID must reference existing records.
func (c *Client) AddCandidateTag(
	ctx context.Context, candidateID, tagID string,
) error {

	var resp struct {
		Success bool `json:"success"`
	}

	return c.Call(ctx, "candidate.addTag", map[string]any{
		"candidateId": candidateID,
		"tagId":       tagID,
	}, &resp)
}
