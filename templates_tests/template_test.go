package templates_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

type testCase struct {
	name         string
	configFile   string
	goldenFile   string
	template     string
	outputFile   string
	ignoreTypes  []string
	ignoreGroups []string
}

const (
	// DefaultFilePerms are the default permissions for created files.
	DefaultFilePerms = 0o600
)

func TestTemplates(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// URL template test
		{
			name:       "go-env/url",
			configFile: "go-env/url.yaml",
			goldenFile: "go-env/url/url.go",
			template:   "https://raw.githubusercontent.com/safeblock-dev/envgen/main/templates/go-env",
			outputFile: "go-env/url/url.generated",
		},
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
			goldenFile: "markdown/options.md",
			template:   "../templates/markdown",
			outputFile: "markdown/options.generated",
		},
		// Ignore tests
		{
			name:         "example/ignore",
			configFile:   "example/ignore.yaml",
			goldenFile:   "example/ignore",
			template:     "../templates/example",
			outputFile:   "example/ignore.generated",
			ignoreGroups: []string{"Webserver"},
		},
		{
			name:         "go-env/ignore-types",
			configFile:   "go-env/ignore.yaml",
			goldenFile:   "go-env/ignore/types/types.go",
			template:     "../templates/go-env",
			outputFile:   "go-env/ignore/types/types.generated",
			ignoreTypes:  []string{"URL"},
			ignoreGroups: nil,
		},
		{
			name:         "go-env/ignore-groups",
			configFile:   "go-env/ignore.yaml",
			goldenFile:   "go-env/ignore/groups/groups.go",
			template:     "../templates/go-env",
			outputFile:   "go-env/ignore/groups/groups.generated",
			ignoreTypes:  nil,
			ignoreGroups: []string{"Database"},
		},
		{
			name:         "go-env/ignore-both",
			configFile:   "go-env/ignore.yaml",
			goldenFile:   "go-env/ignore/both/both.go",
			template:     "../templates/go-env",
			outputFile:   "go-env/ignore/both/both.generated",
			ignoreTypes:  []string{"Duration"},
			ignoreGroups: []string{"App"},
		},
		{
			name:         "markdown/ignore",
			configFile:   "markdown/ignore.yaml",
			goldenFile:   "markdown/ignore.md",
			template:     "../templates/example",
			outputFile:   "markdown/ignore.generated",
			ignoreGroups: []string{"Webserver"},
		},
		// Tags tests
		{
			name:       "go-env/tags",
			configFile: "go-env/tags.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/tags/tags.go",
			outputFile: "go-env/tags/tags.generated",
		},
		// Custom generate commands test
		{
			name:       "go-env/custom_generate",
			configFile: "go-env/custom_generate.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/custom_generate/config.go",
			outputFile: "go-env/custom_generate/config.generated",
		},
		// Meta only test
		{
			name:       "go-env/meta_only",
			configFile: "go-env/meta_only.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/meta_only/config.go",
			outputFile: "go-env/meta_only/config.generated",
		},
		// Generate only test
		{
			name:       "go-env/generate_only",
			configFile: "go-env/generate_only.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/generate_only/config.go",
			outputFile: "go-env/generate_only/config.generated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Generate file
			err := envgen.Generate(envgen.GenerateOptions{
				ConfigPath:   tt.configFile,
				OutputPath:   tt.outputFile,
				TemplatePath: tt.template,
				IgnoreTypes:  tt.ignoreTypes,
				IgnoreGroups: tt.ignoreGroups,
			})
			require.NoError(t, err)

			// Read generated file
			actual, err := os.ReadFile(tt.outputFile)
			require.NoError(t, err)

			// Update golden file if UPDATE_GOLDEN=1
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				err = os.WriteFile(tt.goldenFile, actual, DefaultFilePerms)
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
