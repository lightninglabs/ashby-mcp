package ashby

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListUsers returns all team members, optionally filtered by
// name.
func (c *Client) ListUsers(
	ctx context.Context, name string,
) ([]User, error) {

	params := make(map[string]any)
	if name != "" {
		params["name"] = name
	}

	return Paginate[User](
		ctx, c, "user.list", params, 0,
	)
}

// SearchUsers searches for users by name or email.
func (c *Client) SearchUsers(
	ctx context.Context, name, email string,
) ([]User, error) {

	params := make(map[string]any)
	if name != "" {
		params["name"] = name
	}
	if email != "" {
		params["email"] = email
	}

	var resp struct {
		Success bool              `json:"success"`
		Results []json.RawMessage `json:"results"`
	}

	if err := c.Call(
		ctx, "user.search", params, &resp,
	); err != nil {
		return nil, err
	}

	users := make([]User, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var u User
		if err := json.Unmarshal(raw, &u); err != nil {
			return nil, fmt.Errorf(
				"user.search: decode: %w", err,
			)
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUser returns details for a single user by ID.
func (c *Client) GetUser(
	ctx context.Context, userID string,
) (*User, error) {

	var resp struct {
		Success bool `json:"success"`
		Results User `json:"results"`
	}

	if err := c.Call(ctx, "user.info", map[string]any{
		"userId": userID,
	}, &resp); err != nil {
		return nil, err
	}

	return &resp.Results, nil
}
