package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/safeblock-dev/envgen/internal/github/retry"
)

// ErrEmptyOwner indicates that the repository owner was not provided.
var ErrEmptyOwner = errors.New("repository owner cannot be empty")

// ErrEmptyRepo indicates that the repository name was not provided.
var ErrEmptyRepo = errors.New("repository name cannot be empty")

const (
	// APITimeout is the timeout for GitHub API requests.
	APITimeout = 10 * time.Second

	// APIBaseURL is the base URL for GitHub API.
	APIBaseURL = "https://api.github.com"

	// RawBaseURL is the base URL for raw GitHub content.
	RawBaseURL = "https://raw.githubusercontent.com"
)

// Content represents a file or directory in a GitHub repository.
type Content struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"` //nolint: tagliatelle
}

// Commit represents a commit in a GitHub repository.
type Commit struct {
	SHA string `json:"sha"`
}

// Client handles GitHub API requests.
type Client struct {
	client     *http.Client
	transport  *retry.RoundTripper
	owner      string
	repo       string
	baseURL    string
	rawBaseURL string
}

// New creates a new GitHub client.
// Returns ErrEmptyOwner if owner is empty.
// Returns ErrEmptyRepo if repo is empty.
func New(owner, repo string) (*Client, error) {
	if owner == "" {
		return nil, ErrEmptyOwner
	}

	if repo == "" {
		return nil, ErrEmptyRepo
	}

	transport := retry.New(http.DefaultTransport).
		WithRetryCondition(func(resp *http.Response, err error) bool {
			return err != nil || (resp != nil && (resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests))
		})

	client := &http.Client{
		Timeout:       APITimeout,
		Transport:     transport,
		CheckRedirect: nil, // use default behavior
		Jar:           nil, // no cookie jar needed
	}

	return &Client{
		client:     client,
		transport:  transport,
		owner:      owner,
		repo:       repo,
		baseURL:    APIBaseURL,
		rawBaseURL: RawBaseURL,
	}, nil
}

// GetLatestCommit returns the latest commit SHA from the main branch.
func (c *Client) GetLatestCommit(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/main", c.baseURL, c.owner, c.repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit: %w", err)
	}
	defer resp.Body.Close()

	var commit Commit
	if err := decodeJSON(resp, &commit); err != nil {
		return "", err
	}

	return commit.SHA, nil
}

// GetList returns a list of files in a directory.
func (c *Client) GetList(ctx context.Context, path string) ([]Content, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.baseURL, c.owner, c.repo, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}
	defer resp.Body.Close()

	var contents []Content
	if err := decodeJSON(resp, &contents); err != nil {
		return nil, err
	}

	return contents, nil
}

// GetFile downloads a file from the repository.
func (c *Client) GetFile(ctx context.Context, path string) (string, error) {
	commit, err := c.GetLatestCommit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", c.rawBaseURL, c.owner, c.repo, commit, path)

	return c.GetFileFromURL(ctx, url)
}

// GetFileFromURL downloads a file from a URL.
func (c *Client) GetFileFromURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(content), nil
}

// GetStandardTemplateURL returns the URL for a standard template using the latest commit.
func (c *Client) GetStandardTemplateURL(ctx context.Context, name string) (string, error) {
	commit, err := c.GetLatestCommit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s/%s/templates/%s", c.rawBaseURL, c.owner, c.repo, commit, name), nil
}

// Client returns the http client for the client.
func (c *Client) Client() *http.Client {
	return c.client
}

// Transport returns the transport for the client.
func (c *Client) Transport() *retry.RoundTripper {
	return c.transport
}

// doRequest performs an HTTP request and handles common error cases.
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		return nil, fmt.Errorf("github: unexpected status code: %d (url: %s)", resp.StatusCode, req.URL.String())
	}

	return resp, nil
}

// decodeJSON decodes a JSON response into the given value.
func decodeJSON(resp *http.Response, v interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// SetBaseURL sets the base URL for the GitHub API.
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetRawBaseURL sets the base URL for raw GitHub content.
func (c *Client) SetRawBaseURL(url string) {
	c.rawBaseURL = url
}
