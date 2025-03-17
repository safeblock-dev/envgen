package user_template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_template"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.tmpl")
	content := "test content"

	// Write content to file
	err := os.WriteFile(tmpFile, []byte(content), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid template",
			path:        tmpFile,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "path is empty",
		},
		{
			name:        "non-existent file",
			path:        "non/existent/path",
			expectError: true,
			errorMsg:    "failed to read local template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpl, err := user_template.New(t.Context(), tt.path)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				require.Nil(t, tmpl)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tmpl)
				require.Equal(t, filepath.Base(tt.path), tmpl.GetName())
				require.Equal(t, content, tmpl.Content)
			}
		})
	}
}

func TestTemplate_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		template    user_template.Template
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid template",
			template: user_template.Template{
				Name:         "test.tmpl",
				Content:      "test content",
				ResolvedPath: "/path/to/test.tmpl",
			},
			expectError: false,
		},
		{
			name: "empty name",
			template: user_template.Template{
				Content:      "test content",
				ResolvedPath: "/path/to/test.tmpl",
			},
			expectError: true,
			errorMsg:    "template name is empty",
		},
		{
			name: "empty content",
			template: user_template.Template{
				Name:         "test.tmpl",
				ResolvedPath: "/path/to/test.tmpl",
			},
			expectError: true,
			errorMsg:    "template content is empty",
		},
		{
			name: "empty resolved path",
			template: user_template.Template{
				Name:    "test.tmpl",
				Content: "test content",
			},
			expectError: true,
			errorMsg:    "resolved path is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.template.Validate()
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTemplate_GetName(t *testing.T) {
	t.Parallel()

	template := user_template.Template{
		Name: "test.tmpl",
	}

	require.Equal(t, "test.tmpl", template.GetName())
}

func TestTemplate_GetPath(t *testing.T) {
	t.Parallel()

	template := user_template.Template{
		ResolvedPath: "/path/to/test.tmpl",
	}

	require.Equal(t, "/path/to/test.tmpl", template.GetPath())
}

func TestTemplate_IsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "http url",
			path:     "http://example.com/template.tmpl",
			expected: true,
		},
		{
			name:     "https url",
			path:     "https://example.com/template.tmpl",
			expected: true,
		},
		{
			name:     "local path",
			path:     "/path/to/template.tmpl",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			template := user_template.Template{
				ResolvedPath: tt.path,
			}

			require.Equal(t, tt.expected, template.IsURL())
		})
	}
}
