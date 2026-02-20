package ashby

import "context"

// GetFileURL retrieves a pre-signed download URL for a file
// using its opaque handle string. The handle is obtained from
// a FileHandle's Handle field (e.g. from resumeFileHandle or
// fileHandles on a candidate or application).
func (c *Client) GetFileURL(
	ctx context.Context, fileHandle string,
) (string, error) {

	var resp struct {
		Success bool `json:"success"`
		Results struct {
			URL string `json:"url"`
		} `json:"results"`
	}

	if err := c.Call(ctx, "file.info", map[string]any{
		"fileHandle": fileHandle,
	}, &resp); err != nil {
		return "", err
	}

	return resp.Results.URL, nil
}
