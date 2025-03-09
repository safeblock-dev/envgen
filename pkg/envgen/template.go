package envgen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// LoadTemplateTTL is the timeout for loading a template from a URL.
	LoadTemplateTTL = 30 * time.Second
)

// LoadTemplate loads a template from either a file or URL.
// Automatically detects if the path is a URL (starts with http:// or https://)
// and uses the appropriate loading method.
func LoadTemplate(path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return LoadTemplateFromURL(path)
	}

	return LoadTemplateFromFile(path)
}

// LoadTemplateFromURL loads a template from a HTTP(S) URL.
// Makes a GET request to the specified URL and returns the response body as a string.
// Returns an error if the request fails or returns a non-200 status code.
func LoadTemplateFromURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), LoadTemplateTTL)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// LoadTemplateFromFile loads a template from a local file.
// Returns detailed error messages for common issues like file not found
// or permission denied. The file content is returned as a string.
func LoadTemplateFromFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template file not found: %s", path)
		}

		if os.IsPermission(err) {
			return "", fmt.Errorf("permission denied reading template file: %s", path)
		}

		return "", fmt.Errorf("failed to read template file %s: %w", path, err)
	}

	return string(content), nil
}
