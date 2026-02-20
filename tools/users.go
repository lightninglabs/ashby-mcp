package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListUsersInput defines the input parameters for the
// list_users tool.
type ListUsersInput struct {
	// Name optionally filters users by name.
	Name string `json:"name,omitempty" jsonschema:"Optional name filter"`
}

// ListUsersOutput contains the list_users results.
type ListUsersOutput struct {
	// Users is the list of team members.
	Users []ashby.User `json:"users"`

	// Total is the number of users returned.
	Total int `json:"total"`
}

// ListUsers handles the list_users MCP tool call.
func (h *Handler) ListUsers(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListUsersInput,
) (*mcp.CallToolResult, ListUsersOutput, error) {

	users, err := h.client.ListUsers(ctx, input.Name)
	if err != nil {
		return nil, ListUsersOutput{}, err
	}

	return nil, ListUsersOutput{
		Users: users,
		Total: len(users),
	}, nil
}

// SearchUsersInput defines the input parameters for the
// search_users tool.
type SearchUsersInput struct {
	// Name searches by user name.
	Name string `json:"name,omitempty" jsonschema:"Search by name"`

	// Email searches by email address.
	Email string `json:"email,omitempty" jsonschema:"Search by email address"`
}

// SearchUsersOutput contains the search_users results.
type SearchUsersOutput struct {
	// Users is the list of matching users.
	Users []ashby.User `json:"users"`

	// Total is the number of matching users.
	Total int `json:"total"`
}

// SearchUsers handles the search_users MCP tool call.
func (h *Handler) SearchUsers(
	ctx context.Context, req *mcp.CallToolRequest,
	input SearchUsersInput,
) (*mcp.CallToolResult, SearchUsersOutput, error) {

	users, err := h.client.SearchUsers(
		ctx, input.Name, input.Email,
	)
	if err != nil {
		return nil, SearchUsersOutput{}, err
	}

	return nil, SearchUsersOutput{
		Users: users,
		Total: len(users),
	}, nil
}

// GetUserInput defines the input parameters for the get_user
// tool.
type GetUserInput struct {
	// UserID is the Ashby user ID to look up.
	UserID string `json:"userId" jsonschema:"The Ashby user ID"`
}

// GetUserOutput contains the get_user results.
type GetUserOutput struct {
	// User is the user details.
	User *ashby.User `json:"user"`
}

// GetUser handles the get_user MCP tool call.
func (h *Handler) GetUser(
	ctx context.Context, req *mcp.CallToolRequest,
	input GetUserInput,
) (*mcp.CallToolResult, GetUserOutput, error) {

	user, err := h.client.GetUser(ctx, input.UserID)
	if err != nil {
		return nil, GetUserOutput{}, err
	}

	return nil, GetUserOutput{User: user}, nil
}
