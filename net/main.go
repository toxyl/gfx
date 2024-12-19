package net

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return data, nil
}

// IsURL checks if the provided path is a remote URL or a local path.
func IsURL(path string) bool {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return false
	}

	// A remote URL should have a scheme and a hostname (e.g., http://example.com)
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
