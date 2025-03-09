package templates_test

import (
	"os"
	"testing"

	"github.com/safeblock-dev/envgen/pkg/envgen"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name       string
	configFile string
	goldenFile string
	template   string
	outputFile string
}

func TestTemplates(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Basic tests
		{
			name:       "example/basic",
			configFile: "example/basic.yaml",
			goldenFile: "example/basic.env",
			template:   "../templates/example",
			outputFile: "example/basic.generated",
		},
		{
			name:       "go-env/basic",
			configFile: "go-env/basic.yaml",
			goldenFile: "go-env/basic/basic.go",
			template:   "../templates/go-env",
			outputFile: "go-env/basic/basic.generated",
		},
		{
			name:       "markdown/basic",
			configFile: "markdown/basic.yaml",
			goldenFile: "markdown/basic.md",
			template:   "../templates/markdown",
			outputFile: "markdown/basic.generated",
		},
		// Minimal configuration tests
		{
			name:       "example/minimal",
			configFile: "example/minimal.yaml",
			goldenFile: "example/minimal.env",
			template:   "../templates/example",
			outputFile: "example/minimal.generated",
		},
		{
			name:       "go-env/minimal",
			configFile: "go-env/minimal.yaml",
			goldenFile: "go-env/minimal/minimal.go",
			template:   "../templates/go-env",
			outputFile: "go-env/minimal/minimal.generated",
		},
		{
			name:       "markdown/minimal",
			configFile: "markdown/minimal.yaml",
			goldenFile: "markdown/minimal.md",
			template:   "../templates/markdown",
			outputFile: "markdown/minimal.generated",
		},
		// Custom types tests
		{
			name:       "example/types",
			configFile: "example/types.yaml",
			goldenFile: "example/types.env",
			template:   "../templates/example",
			outputFile: "example/types.generated",
		},
		{
			name:       "go-env/types",
			configFile: "go-env/types.yaml",
			goldenFile: "go-env/types/types.go",
			template:   "../templates/go-env",
			outputFile: "go-env/types/types.generated",
		},
		{
			name:       "markdown/types",
			configFile: "markdown/types.yaml",
			goldenFile: "markdown/types.md",
			template:   "../templates/markdown",
			outputFile: "markdown/types.generated",
		},
		// Prefix tests
		{
			name:       "example/prefix",
			configFile: "example/prefix.yaml",
			goldenFile: "example/prefix.env",
			template:   "../templates/example",
			outputFile: "example/prefix.generated",
		},
		{
			name:       "go-env/prefix",
			configFile: "go-env/prefix.yaml",
			goldenFile: "go-env/prefix/prefix.go",
			template:   "../templates/go-env",
			outputFile: "go-env/prefix/prefix.generated",
		},
		{
			name:       "markdown/prefix",
			configFile: "markdown/prefix.yaml",
			goldenFile: "markdown/prefix.md",
			template:   "../templates/markdown",
			outputFile: "markdown/prefix.generated",
		},
		// Options tests
		{
			name:       "go-env/options",
			configFile: "go-env/options.yaml",
			goldenFile: "go-env/options/options.go",
			template:   "../templates/go-env",
			outputFile: "go-env/options/options.generated",
		},
		{
			name:       "markdown/options",
			configFile: "markdown/options.yaml",
			goldenFile: "markdown/options.golden.md",
			template:   "../templates/markdown",
			outputFile: "markdown/options.generated",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Generate file
			err := envgen.Generate(tt.configFile, tt.outputFile, tt.template)
			require.NoError(t, err)

			// Read generated file
			actual, err := os.ReadFile(tt.outputFile)
			require.NoError(t, err)

			// Update golden file if UPDATE_GOLDEN=1
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				err = os.WriteFile(tt.goldenFile, actual, 0644)
				require.NoError(t, err)
				return
			}

			// Read golden file
			expected, err := os.ReadFile(tt.goldenFile)
			require.NoError(t, err)

			// Compare results
			require.Equal(t, string(expected), string(actual))
		})
	}
}
