package github_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/github"
)

// setupMockServer creates a test server and configures a client to use it.
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *github.Client) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := github.New("safeblock-dev", "envgen")
	require.NoError(t, err)
	// Override client's base URLs to use mock server
	client.SetBaseURL(server.URL)
	client.SetRawBaseURL(server.URL)

	return server, client
}

func TestClient_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		owner       string
		repo        string
		expectError error
	}{
		{
			name:        "valid_values",
			owner:       "testowner",
			repo:        "testrepo",
			expectError: nil,
		},
		{
			name:        "empty_owner",
			owner:       "",
			repo:        "testrepo",
			expectError: github.ErrEmptyOwner,
		},
		{
			name:        "empty_repo",
			owner:       "testowner",
			repo:        "",
			expectError: github.ErrEmptyRepo,
		},
		{
			name:        "both_empty",
			owner:       "",
			repo:        "",
			expectError: github.ErrEmptyOwner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := github.New(tt.owner, tt.repo)

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
				require.Nil(t, client)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, client)
			require.NotNil(t, client.Client())
			require.NotNil(t, client.Transport())
		})
	}
}

func TestClient_GetLatestCommit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedSHA    string
		expectError    bool
	}{
		{
			name:           "success_valid_response",
			responseStatus: http.StatusOK,
			responseBody:   `{"sha": "abc123"}`,
			expectedSHA:    "abc123",
			expectError:    false,
		},
		{
			name:           "error_server_error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "",
			expectedSHA:    "",
			expectError:    true,
		},
		{
			name:           "error_invalid_json",
			responseStatus: http.StatusOK,
			responseBody:   `{"invalid": json}`,
			expectedSHA:    "",
			expectError:    true,
		},
		{
			name:           "error_empty_response",
			responseStatus: http.StatusOK,
			responseBody:   `{}`,
			expectedSHA:    "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
				}

				if !strings.Contains(r.URL.Path, "/commits/main") {
					t.Errorf("expected path to contain /commits/main, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.responseStatus)

				if _, err := w.Write([]byte(tt.responseBody)); err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			})

			sha, err := client.GetLatestCommit(t.Context())

			if tt.expectError {
				require.Error(t, err)
				require.Empty(t, sha)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedSHA, sha)
			}
		})
	}
}

func TestClient_GetList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		responseStatus int
		responseBody   string
		expectedFiles  []github.Content
		expectError    bool
	}{
		{
			name:           "success_single_file",
			path:           "templates",
			responseStatus: http.StatusOK,
			responseBody:   `[{"name":"file1.txt","path":"templates/file1.txt","type":"file","download_url":"http://example.com/file1.txt"}]`,
			expectedFiles: []github.Content{
				{
					Name:        "file1.txt",
					Path:        "templates/file1.txt",
					Type:        "file",
					DownloadURL: "http://example.com/file1.txt",
				},
			},
			expectError: false,
		},
		{
			name:           "success_empty_directory",
			path:           "empty",
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			expectedFiles:  []github.Content{},
			expectError:    false,
		},
		{
			name:           "error_not_found",
			path:           "nonexistent",
			responseStatus: http.StatusNotFound,
			responseBody:   "",
			expectedFiles:  nil,
			expectError:    true,
		},
		{
			name:           "error_invalid_json",
			path:           "invalid",
			responseStatus: http.StatusOK,
			responseBody:   `{"invalid": json}`,
			expectedFiles:  nil,
			expectError:    true,
		},
		{
			name:           "error_malformed_response",
			path:           "malformed",
			responseStatus: http.StatusOK,
			responseBody:   `[{"name": 123}]`,
			expectedFiles:  nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
				}

				if !strings.Contains(r.URL.Path, tt.path) {
					t.Errorf("expected path to contain %s, got %s", tt.path, r.URL.Path)
				}

				w.WriteHeader(tt.responseStatus)

				if _, err := w.Write([]byte(tt.responseBody)); err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			})

			files, err := client.GetList(t.Context(), tt.path)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, files)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedFiles, files)
			}
		})
	}
}

func TestClient_GetFileFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		url            string
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "success_valid_file",
			url:            "/file.txt",
			responseStatus: http.StatusOK,
			responseBody:   "file content",
			expectError:    false,
		},
		{
			name:           "error_not_found",
			url:            "/nonexistent.txt",
			responseStatus: http.StatusNotFound,
			responseBody:   "",
			expectError:    true,
		},
		{
			name:           "error_server_error",
			url:            "/error.txt",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "",
			expectError:    true,
		},
		{
			name:           "success_empty_file",
			url:            "/empty.txt",
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
				}

				if r.URL.Path != tt.url {
					t.Errorf("expected URL path %s, got %s", tt.url, r.URL.Path)
				}

				w.WriteHeader(tt.responseStatus)

				if tt.responseStatus == http.StatusOK {
					if _, err := w.Write([]byte(tt.responseBody)); err != nil {
						t.Errorf("failed to write response: %v", err)
					}
				}
			})

			content, err := client.GetFileFromURL(t.Context(), server.URL+tt.url)

			if tt.expectError {
				require.Error(t, err)
				require.Empty(t, content)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.responseBody, content)
			}
		})
	}
}

func TestClient_GetStandardTemplateURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		templateName string
		expectError  bool
	}{
		{
			name:         "success_valid_template",
			templateName: "example.tmpl",
			expectError:  false,
		},
		{
			name:         "success_empty_template_name",
			templateName: "",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
				}

				w.WriteHeader(http.StatusOK)

				if _, err := w.Write([]byte(`{"sha": "test"}`)); err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			})

			url, err := client.GetStandardTemplateURL(t.Context(), tt.templateName)

			if tt.expectError {
				require.Error(t, err)
				require.Empty(t, url)
			} else {
				require.NoError(t, err)

				expectedURL := fmt.Sprintf("%s/%s/%s/%s/templates/%s",
					server.URL,
					"safeblock-dev",
					"envgen",
					"test",
					tt.templateName,
				)
				require.Equal(t, expectedURL, url)
			}
		})
	}
}
