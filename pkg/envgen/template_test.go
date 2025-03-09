package envgen_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

const (
	testTemplateContent = "test template content"
)

type testCase struct {
	name        string
	path        string
	wantContent string
	wantErr     bool
}

func setupTestFile(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.template")
	err := os.WriteFile(testFile, []byte(testTemplateContent), testFilePerm)
	require.NoError(t, err, "Failed to create test file")

	return testFile
}

func TestLoadTemplateFromFile(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			name:        "valid file",
			path:        setupTestFile(t),
			wantContent: testTemplateContent,
			wantErr:     false,
		},
		{
			name:    "non-existent file",
			path:    "non_existent.template",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			content, err := envgen.LoadTemplateFromFile(tt.path)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantContent, content)
		})
	}
}

func TestLoadTemplateFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(testTemplateContent))
			},
			wantContent: testTemplateContent,
			wantErr:     false,
		},
		{
			name: "error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.handler)
			t.Cleanup(server.Close)

			content, err := envgen.LoadTemplateFromURL(server.URL)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantContent, content)
		})
	}
}

func TestLoadTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "valid file",
			setup:       setupTestFile,
			wantContent: testTemplateContent,
			wantErr:     false,
		},
		{
			name: "valid url",
			setup: func(t *testing.T) string {
				t.Helper()

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("test url content"))
				}))
				t.Cleanup(server.Close)

				return server.URL
			},
			wantContent: "test url content",
			wantErr:     false,
		},
		{
			name: "invalid url",
			setup: func(t *testing.T) string {
				t.Helper()

				return "http://invalid-url-that-does-not-exist.com"
			},
			wantErr: true,
		},
		{
			name: "empty path",
			setup: func(t *testing.T) string {
				t.Helper()

				return ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := tt.setup(t)

			content, err := envgen.LoadTemplate(path)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantContent, content)
		})
	}
}
