package ashby

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// defaultBaseURL is the Ashby API root.
	defaultBaseURL = "https://api.ashbyhq.com"

	// defaultTimeout is the HTTP request timeout.
	defaultTimeout = 30 * time.Second
)

// Caller abstracts HTTP POST calls to the Ashby API. This allows
// tests to substitute a fake without mocking the full HTTP stack.
type Caller interface {
	// Call makes a POST request to the given Ashby endpoint
	// with the JSON-encoded body and decodes the response into
	// result.
	Call(
		ctx context.Context, endpoint string,
		body, result any,
	) error
}

// Client communicates with the Ashby REST API using Basic Auth
// and JSON-encoded POST requests.
type Client struct {
	baseURL    string
	authHeader string
	httpClient *http.Client
}

// NewClient creates a new Ashby API client with the given key.
func NewClient(apiKey string) *Client {
	token := base64.StdEncoding.EncodeToString(
		[]byte(apiKey + ":"),
	)

	return &Client{
		baseURL:    defaultBaseURL,
		authHeader: "Basic " + token,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

// NewClientFromEnv creates an Ashby client using the
// ASHBY_API_KEY or ASHBY_KEY environment variable.
func NewClientFromEnv() (*Client, error) {
	key := os.Getenv("ASHBY_API_KEY")
	if key == "" {
		key = os.Getenv("ASHBY_KEY")
	}

	if key == "" {
		return nil, fmt.Errorf(
			"ASHBY_API_KEY environment variable not set",
		)
	}

	return NewClient(key), nil
}

// Call makes a POST request to the given Ashby endpoint, sending
// body as JSON and decoding the response into result. The result
// parameter should be a pointer to an APIResponse or
// PaginatedResponse.
func (c *Client) Call(
	ctx context.Context, endpoint string,
	body, result any,
) error {

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%s: marshal request: %w",
			endpoint, err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		c.baseURL+"/"+endpoint,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("%s: build request: %w",
			endpoint, err,
		)
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: read response: %w",
			endpoint, err,
		)
	}

	// Check for HTTP-level errors before JSON parsing.
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf(
			"%s: invalid or missing API key (401)",
			endpoint,
		)
	}

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf(
			"%s: API key lacks required permissions (403)",
			endpoint,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"%s: HTTP %d: %s",
			endpoint, resp.StatusCode, string(respBody),
		)
	}

	// Peek at the success field before full decoding.
	var envelope struct {
		Success   bool       `json:"success"`
		ErrorInfo *ErrorInfo `json:"errorInfo,omitempty"`
		Errors    []string   `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("%s: decode response: %w",
			endpoint, err,
		)
	}

	if !envelope.Success {
		msg := "unknown error"
		if envelope.ErrorInfo != nil {
			msg = envelope.ErrorInfo.Message
		} else if len(envelope.Errors) > 0 {
			msg = envelope.Errors[0]
		}

		return fmt.Errorf("%s: API error: %s",
			endpoint, msg,
		)
	}

	// Decode the full response into the caller's target.
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("%s: decode result: %w",
			endpoint, err,
		)
	}

	return nil
}
