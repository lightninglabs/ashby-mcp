package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
)

// ListSourcesInput defines the input parameters for the
// list_sources tool (none required).
type ListSourcesInput struct{}

// ListSourcesOutput contains the list_sources results.
type ListSourcesOutput struct {
	// Sources is the list of application sources.
	Sources []ashby.Source `json:"sources"`

	// Total is the number of sources returned.
	Total int `json:"total"`
}

// ListSources handles the list_sources MCP tool call.
func (h *Handler) ListSources(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListSourcesInput,
) (*mcp.CallToolResult, ListSourcesOutput, error) {

	sources, err := h.client.ListSources(ctx)
	if err != nil {
		return nil, ListSourcesOutput{}, err
	}

	return nil, ListSourcesOutput{
		Sources: sources,
		Total:   len(sources),
	}, nil
}

// ListArchiveReasonsInput defines the input parameters for the
// list_archive_reasons tool (none required).
type ListArchiveReasonsInput struct{}

// ListArchiveReasonsOutput contains the list_archive_reasons
// results.
type ListArchiveReasonsOutput struct {
	// ArchiveReasons is the list of archive reasons.
	ArchiveReasons []ashby.ArchiveReason `json:"archiveReasons"`

	// Total is the number of archive reasons returned.
	Total int `json:"total"`
}

// ListArchiveReasons handles the list_archive_reasons MCP tool
// call.
func (h *Handler) ListArchiveReasons(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListArchiveReasonsInput,
) (*mcp.CallToolResult, ListArchiveReasonsOutput, error) {

	reasons, err := h.client.ListArchiveReasons(ctx)
	if err != nil {
		return nil, ListArchiveReasonsOutput{}, err
	}

	return nil, ListArchiveReasonsOutput{
		ArchiveReasons: reasons,
		Total:          len(reasons),
	}, nil
}

// ListDepartmentsInput defines the input parameters for the
// list_departments tool (none required).
type ListDepartmentsInput struct{}

// ListDepartmentsOutput contains the list_departments results.
type ListDepartmentsOutput struct {
	// Departments is the list of departments.
	Departments []ashby.Department `json:"departments"`

	// Total is the number of departments returned.
	Total int `json:"total"`
}

// ListDepartments handles the list_departments MCP tool call.
func (h *Handler) ListDepartments(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListDepartmentsInput,
) (*mcp.CallToolResult, ListDepartmentsOutput, error) {

	depts, err := h.client.ListDepartments(ctx)
	if err != nil {
		return nil, ListDepartmentsOutput{}, err
	}

	return nil, ListDepartmentsOutput{
		Departments: depts,
		Total:       len(depts),
	}, nil
}

// ListLocationsInput defines the input parameters for the
// list_locations tool (none required).
type ListLocationsInput struct{}

// ListLocationsOutput contains the list_locations results.
type ListLocationsOutput struct {
	// Locations is the list of locations.
	Locations []ashby.Location `json:"locations"`

	// Total is the number of locations returned.
	Total int `json:"total"`
}

// ListLocations handles the list_locations MCP tool call.
func (h *Handler) ListLocations(
	ctx context.Context, req *mcp.CallToolRequest,
	input ListLocationsInput,
) (*mcp.CallToolResult, ListLocationsOutput, error) {

	locs, err := h.client.ListLocations(ctx)
	if err != nil {
		return nil, ListLocationsOutput{}, err
	}

	return nil, ListLocationsOutput{
		Locations: locs,
		Total:     len(locs),
	}, nil
}
