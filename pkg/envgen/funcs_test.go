package envgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

func TestEnvgen_ProcessTemplate(t *testing.T) {
	t.Parallel()

	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Create test files
	configPath := filepath.Join(tmpDir, "config.yaml")
	outputPath := filepath.Join(tmpDir, "output.go")
	templatePath := filepath.Join(tmpDir, "template.tmpl")

	// Create minimal config file with required group
	configContent := `groups:
  - name: App
    description: Application settings
    fields:
      - name: debug
        type: bool
        description: Enable debug mode`
	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	// Create empty output file
	err = os.WriteFile(outputPath, []byte(""), 0o600)
	require.NoError(t, err)

	// Create template file
	err = os.WriteFile(templatePath, []byte("package main"), 0o600)
	require.NoError(t, err)

	// Create a properly initialized Envgen instance
	envgenClient, err := envgen.New(t.Context(), envgen.Options{
		ConfigPath:   configPath,
		OutputPath:   outputPath,
		TemplatePath: templatePath,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		envgen   *envgen.Envgen
		content  string
		expected string
	}{
		{
			name:     "empty content",
			envgen:   envgenClient,
			content:  "",
			expected: "",
		},
		{
			name:     "nil envgen",
			envgen:   nil,
			content:  "test",
			expected: "test",
		},
		{
			name:     "simple template with functions",
			envgen:   envgenClient,
			content:  "{{ upper \"hello\" }}",
			expected: "HELLO",
		},
		{
			name:     "invalid template",
			envgen:   envgenClient,
			content:  "{{ invalid }}",
			expected: "{{ invalid }}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.envgen.ProcessTemplate(tt.content)
			require.Equal(t, tt.expected, result)
		})
	}
}
