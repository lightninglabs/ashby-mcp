package ashby

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListOpenings returns all openings.
func (c *Client) ListOpenings(
	ctx context.Context,
) ([]Opening, error) {

	return Paginate[Opening](
		ctx, c, "opening.list", nil, 0,
	)
}

// GetOpening returns details for a single opening by ID.
func (c *Client) GetOpening(
	ctx context.Context, openingID string,
) (*Opening, error) {

	var resp struct {
		Success bool    `json:"success"`
		Results Opening `json:"results"`
	}

	if err := c.Call(ctx, "opening.info", map[string]any{
		"openingId": openingID,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

// SearchOpenings searches for openings using the given term.
func (c *Client) SearchOpenings(
	ctx context.Context, term string,
) ([]Opening, error) {

	var resp struct {
		Success bool              `json:"success"`
		Results []json.RawMessage `json:"results"`
	}

	if err := c.Call(ctx, "opening.search", map[string]any{
		"term": term,
	}, &resp); err != nil {
		return nil, err
	}

	openings := make([]Opening, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var o Opening
		if err := json.Unmarshal(raw, &o); err != nil {
			return nil, fmt.Errorf(
				"opening.search: decode: %w", err,
			)
		}
		openings = append(openings, o)
	}

	return openings, nil
}
