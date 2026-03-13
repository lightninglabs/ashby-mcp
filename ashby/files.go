package ashby

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ledongthuc/pdf"
)

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

// ResumeText holds the extracted plain text from a candidate's
// resume along with metadata about the file.
type ResumeText struct {
	// CandidateName is the candidate's full name.
	CandidateName string `json:"candidateName"`

	// FileName is the original resume filename.
	FileName string `json:"fileName"`

	// Text is the extracted plain text content.
	Text string `json:"text"`
}

// GetCandidateResume fetches a candidate's resume PDF and
// extracts the text content. It searches for the candidate by
// name or email, gets the pre-signed URL for the resume file,
// downloads the PDF, and returns the extracted plain text.
func (c *Client) GetCandidateResume(
	ctx context.Context, email, name string,
) (*ResumeText, error) {

	// Search for the candidate to get resume file handle.
	candidates, err := c.SearchCandidates(ctx, email, name)
	if err != nil {
		return nil, fmt.Errorf("search candidates: %w", err)
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf(
			"no candidates found for email=%q name=%q",
			email, name,
		)
	}

	// Use the first matching candidate.
	cand := candidates[0]

	if cand.ResumeFileHandle == nil ||
		cand.ResumeFileHandle.Handle == "" {

		return nil, fmt.Errorf(
			"candidate %s has no resume file",
			cand.Name,
		)
	}

	// Get the pre-signed download URL.
	fileURL, err := c.GetFileURL(
		ctx, cand.ResumeFileHandle.Handle,
	)
	if err != nil {
		return nil, fmt.Errorf("get file URL: %w", err)
	}

	// Download the PDF.
	pdfBytes, err := c.downloadFile(ctx, fileURL)
	if err != nil {
		return nil, fmt.Errorf("download resume: %w", err)
	}

	// Extract text from the PDF.
	text, err := extractPDFText(pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("extract PDF text: %w", err)
	}

	return &ResumeText{
		CandidateName: cand.Name,
		FileName:      cand.ResumeFileHandle.Name,
		Text:          text,
	}, nil
}

// downloadFile fetches a file from the given URL and returns
// the raw bytes.
func (c *Client) downloadFile(
	ctx context.Context, url string,
) ([]byte, error) {

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"HTTP %d downloading file", resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}

// extractPDFText reads a PDF from raw bytes and extracts the
// plain text content from all pages.
func extractPDFText(data []byte) (string, error) {
	reader, err := pdf.NewReader(
		bytes.NewReader(data), int64(len(data)),
	)
	if err != nil {
		return "", fmt.Errorf("open PDF: %w", err)
	}

	var buf strings.Builder

	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf(
				"extract page %d: %w", i, err,
			)
		}

		buf.WriteString(text)
	}

	return buf.String(), nil
}
