package user_template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_template"
)

func TestResolver_Template(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.tmpl")
	content := "test template content"

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
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "template path is empty",
		},
		{
			name:        "non-existent local file",
			path:        "/non/existent/path.tmpl",
			expectError: true,
			errorMsg:    "failed to read local template",
		},
		{
			name:        "invalid URL",
			path:        "https://invalid.url/template.tmpl",
			expectError: true,
			errorMsg:    "failed to fetch template from URL",
		},
		{
			name:        "existing local file",
			path:        tmpFile,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver, err := user_template.NewResolver()
			require.NoError(t, err)
			tmpl, err := resolver.Template(t.Context(), tt.path)

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

func TestResolver_ListAvailableTemplateNames(t *testing.T) {
	t.Parallel()

	resolver, err := user_template.NewResolver()
	require.NoError(t, err)

	names, err := resolver.ListAvailableTemplateNames(t.Context())

	// Since we can't guarantee the actual templates in the repository,
	// we'll just check that the function works without error
	require.NoError(t, err)
	require.NotNil(t, names)
}
