package ashby

import "context"

// ListSources returns all application sources.
func (c *Client) ListSources(
	ctx context.Context,
) ([]Source, error) {

	return Paginate[Source](
		ctx, c, "source.list", nil, 0,
	)
}

// ListArchiveReasons returns all application archive reasons.
func (c *Client) ListArchiveReasons(
	ctx context.Context,
) ([]ArchiveReason, error) {

	return Paginate[ArchiveReason](
		ctx, c, "archiveReason.list", nil, 0,
	)
}

// ListDepartments returns all departments.
func (c *Client) ListDepartments(
	ctx context.Context,
) ([]Department, error) {

	return Paginate[Department](
		ctx, c, "department.list", nil, 0,
	)
}

// ListLocations returns all locations.
func (c *Client) ListLocations(
	ctx context.Context,
) ([]Location, error) {

	return Paginate[Location](
		ctx, c, "location.list", nil, 0,
	)
}
