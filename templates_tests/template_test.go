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
			goldenFile: "example/basic.golden",
			template:   "../templates/example",
			outputFile: "example/basic.generated",
		},
		{
			name:       "go-env/basic",
			configFile: "go-env/basic.yaml",
			goldenFile: "go-env/basic.golden",
			template:   "../templates/go-env",
			outputFile: "go-env/basic.generated",
		},
		// Minimal configuration tests
		{
			name:       "example/minimal",
			configFile: "example/minimal.yaml",
			goldenFile: "example/minimal.golden",
			template:   "../templates/example",
			outputFile: "example/minimal.generated",
		},
		{
			name:       "go-env/minimal",
			configFile: "go-env/minimal.yaml",
			goldenFile: "go-env/minimal.golden",
			template:   "../templates/go-env",
			outputFile: "go-env/minimal.generated",
		},
		// Custom types tests
		{
			name:       "example/types",
			configFile: "example/types.yaml",
			goldenFile: "example/types.golden",
			template:   "../templates/example",
			outputFile: "example/types.generated",
		},
		{
			name:       "go-env/types",
			configFile: "go-env/types.yaml",
			goldenFile: "go-env/types.golden",
			template:   "../templates/go-env",
			outputFile: "go-env/types.generated",
		},
		// Prefix tests
		{
			name:       "example/prefix",
			configFile: "example/prefix.yaml",
			goldenFile: "example/prefix.golden",
			template:   "../templates/example",
			outputFile: "example/prefix.generated",
		},
		{
			name:       "go-env/prefix",
			configFile: "go-env/prefix.yaml",
			goldenFile: "go-env/prefix.golden",
			template:   "../templates/go-env",
			outputFile: "go-env/prefix.generated",
		},
		// Options tests (only for go-env)
		{
			name:       "go-env/options",
			configFile: "go-env/options.yaml",
			goldenFile: "go-env/options.golden",
			template:   "../templates/go-env",
			outputFile: "go-env/options.generated",
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
