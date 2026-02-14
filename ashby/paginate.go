package ashby

import (
	"context"
	"fmt"
	"time"
)

const (
	// defaultPageSize is the maximum items per page supported
	// by the Ashby API.
	defaultPageSize = 100

	// pageDelay is the delay between paginated requests to
	// avoid rate limiting.
	pageDelay = 200 * time.Millisecond
)

// PaginatedResponse is the common envelope returned by Ashby
// list endpoints. The type parameter T is the element type of
// the results slice.
type PaginatedResponse[T any] struct {
	// Success indicates whether the API call succeeded.
	Success bool `json:"success"`

	// Results holds the page of results.
	Results []T `json:"results"`

	// MoreDataAvailable signals additional pages exist.
	MoreDataAvailable bool `json:"moreDataAvailable"`

	// NextCursor is the opaque token for fetching the next
	// page.
	NextCursor string `json:"nextCursor"`
}

// PageResult holds a single page of results along with cursor
// metadata for the caller to decide whether to continue paging.
type PageResult[T any] struct {
	// Items contains the results from this page.
	Items []T

	// NextCursor is the cursor for the next page, empty if no
	// more data.
	NextCursor string

	// MoreDataAvailable indicates additional pages exist.
	MoreDataAvailable bool
}

// FetchPage fetches a single page from an Ashby list endpoint.
// This is the building block for both full pagination and
// cursor-passthrough in MCP tools.
func FetchPage[T any](
	ctx context.Context, c Caller, endpoint string,
	params map[string]any,
) (*PageResult[T], error) {

	var resp PaginatedResponse[T]
	if err := c.Call(ctx, endpoint, params, &resp); err != nil {
		return nil, err
	}

	return &PageResult[T]{
		Items:             resp.Results,
		NextCursor:        resp.NextCursor,
		MoreDataAvailable: resp.MoreDataAvailable,
	}, nil
}

// Paginate fetches all pages from an Ashby list endpoint,
// accumulating results up to limit. If limit is zero or
// negative, all results are fetched. A 200ms delay is inserted
// between page requests to respect rate limits.
func Paginate[T any](
	ctx context.Context, c Caller, endpoint string,
	params map[string]any, limit int,
) ([]T, error) {

	var all []T
	cursor := ""

	for {
		// Build per-page params, preserving caller's base
		// params.
		pageParams := make(map[string]any, len(params)+2)
		for k, v := range params {
			pageParams[k] = v
		}

		// Set page size, capping at limit if provided.
		pageSize := defaultPageSize
		if limit > 0 {
			remaining := limit - len(all)
			if remaining < pageSize {
				pageSize = remaining
			}
		}
		pageParams["per_page"] = pageSize

		if cursor != "" {
			pageParams["cursor"] = cursor
		}

		page, err := FetchPage[T](
			ctx, c, endpoint, pageParams,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"paginate %s: %w", endpoint, err,
			)
		}

		all = append(all, page.Items...)

		// Check termination conditions.
		if !page.MoreDataAvailable || page.NextCursor == "" {
			break
		}

		if limit > 0 && len(all) >= limit {
			break
		}

		cursor = page.NextCursor

		// Rate-limit delay between pages.
		select {
		case <-ctx.Done():
			return all, ctx.Err()
		case <-time.After(pageDelay):
		}
	}

	// Trim to exact limit if we overshot.
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}

	return all, nil
}
