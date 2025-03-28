package templates_test

import (
	"net/http"
	"net/http/httptest"
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
	fromURL      bool
}

const (
	// DefaultFilePerms are the default permissions for created files.
	DefaultFilePerms = 0o600
)

func TestTemplates(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// --------------------------------
		// Example template tests (.env)
		// --------------------------------
		{
			name:       "example/basic",
			configFile: "example/basic.yaml",
			goldenFile: "example/basic.env",
			template:   "../templates/example",
			outputFile: "example/basic.generated",
		},
		{
			name:       "example/minimal",
			configFile: "example/minimal.yaml",
			goldenFile: "example/minimal.env",
			template:   "../templates/example",
			outputFile: "example/minimal.generated",
		},
		{
			name:       "example/types",
			configFile: "example/types.yaml",
			goldenFile: "example/types.env",
			template:   "../templates/example",
			outputFile: "example/types.generated",
		},
		{
			name:       "example/prefix",
			configFile: "example/prefix.yaml",
			goldenFile: "example/prefix.env",
			template:   "../templates/example",
			outputFile: "example/prefix.generated",
		},
		{
			name:         "example/ignore-types",
			configFile:   "example/ignore.yaml",
			goldenFile:   "example/ignore-types.env",
			template:     "../templates/example",
			outputFile:   "example/ignore-types.generated",
			ignoreTypes:  []string{"URL"},
			ignoreGroups: nil,
		},
		{
			name:         "example/ignore-groups",
			configFile:   "example/ignore.yaml",
			goldenFile:   "example/ignore-groups.env",
			template:     "../templates/example",
			outputFile:   "example/ignore-groups.generated",
			ignoreTypes:  nil,
			ignoreGroups: []string{"DatabaseConfig"},
		},
		{
			name:         "example/ignore-both",
			configFile:   "example/ignore.yaml",
			goldenFile:   "example/ignore-both.env",
			template:     "../templates/example",
			outputFile:   "example/ignore-both.generated",
			ignoreTypes:  []string{"Duration"},
			ignoreGroups: []string{"App"},
		},

		// --------------------------------
		// Go-env template tests (Go structs)
		// --------------------------------
		{
			name:       "go-env/url",
			configFile: "go-env/url.yaml",
			goldenFile: "go-env/url/url.go",
			template:   "../templates/go-env",
			outputFile: "go-env/url/url.generated",
			fromURL:    true,
		},
		{
			name:       "go-env/basic",
			configFile: "go-env/basic.yaml",
			goldenFile: "go-env/basic/basic.go",
			template:   "../templates/go-env",
			outputFile: "go-env/basic/basic.generated",
		},
		{
			name:       "go-env/minimal",
			configFile: "go-env/minimal.yaml",
			goldenFile: "go-env/minimal/minimal.go",
			template:   "../templates/go-env",
			outputFile: "go-env/minimal/minimal.generated",
		},
		{
			name:       "go-env/types",
			configFile: "go-env/types.yaml",
			goldenFile: "go-env/types/types.go",
			template:   "../templates/go-env",
			outputFile: "go-env/types/types.generated",
		},
		{
			name:       "go-env/prefix",
			configFile: "go-env/prefix.yaml",
			goldenFile: "go-env/prefix/prefix.go",
			template:   "../templates/go-env",
			outputFile: "go-env/prefix/prefix.generated",
		},
		{
			name:       "go-env/options",
			configFile: "go-env/options.yaml",
			goldenFile: "go-env/options/options.go",
			template:   "../templates/go-env",
			outputFile: "go-env/options/options.generated",
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
			ignoreGroups: []string{"DatabaseConfig"},
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
			name:       "go-env/tags",
			configFile: "go-env/tags.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/tags/tags.go",
			outputFile: "go-env/tags/tags.generated",
		},
		{
			name:       "go-env/meta",
			configFile: "go-env/meta.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/meta/config.go",
			outputFile: "go-env/meta/config.generated",
		},
		{
			name:       "go-env/empty_meta",
			configFile: "go-env/empty_meta.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/empty_meta/config.go",
			outputFile: "go-env/empty_meta/config.generated",
		},
		{
			name:       "go-env/skip_env_tag",
			configFile: "go-env/skip_env_tag.yaml",
			template:   "../templates/go-env",
			goldenFile: "go-env/skip_env_tag/skip_env_tag.go",
			outputFile: "go-env/skip_env_tag/skip_env_tag.generated",
		},

		// --------------------------------
		// Markdown template tests (Documentation)
		// --------------------------------
		{
			name:       "markdown/basic",
			configFile: "markdown/basic.yaml",
			goldenFile: "markdown/basic.md",
			template:   "../templates/markdown",
			outputFile: "markdown/basic.generated",
		},
		{
			name:       "markdown/minimal",
			configFile: "markdown/minimal.yaml",
			goldenFile: "markdown/minimal.md",
			template:   "../templates/markdown",
			outputFile: "markdown/minimal.generated",
		},
		{
			name:       "markdown/types",
			configFile: "markdown/types.yaml",
			goldenFile: "markdown/types.md",
			template:   "../templates/markdown",
			outputFile: "markdown/types.generated",
		},
		{
			name:       "markdown/prefix",
			configFile: "markdown/prefix.yaml",
			goldenFile: "markdown/prefix.md",
			template:   "../templates/markdown",
			outputFile: "markdown/prefix.generated",
		},
		{
			name:       "markdown/options",
			configFile: "markdown/options.yaml",
			goldenFile: "markdown/options.md",
			template:   "../templates/markdown",
			outputFile: "markdown/options.generated",
		},
		{
			name:       "markdown/column_visibility",
			configFile: "markdown/column_visibility.yaml",
			template:   "../templates/markdown",
			goldenFile: "markdown/column_visibility.md",
			outputFile: "markdown/column_visibility.generated",
		},
		{
			name:         "markdown/ignore-types",
			configFile:   "markdown/ignore.yaml",
			goldenFile:   "markdown/ignore-types.md",
			template:     "../templates/example",
			outputFile:   "markdown/ignore-types.generated",
			ignoreTypes:  []string{"URL"},
			ignoreGroups: nil,
		},
		{
			name:         "markdown/ignore-groups",
			configFile:   "markdown/ignore.yaml",
			goldenFile:   "markdown/ignore-groups.md",
			template:     "../templates/example",
			outputFile:   "markdown/ignore-groups.generated",
			ignoreTypes:  nil,
			ignoreGroups: []string{"DatabaseConfig"},
		},
		{
			name:         "markdown/ignore-both",
			configFile:   "markdown/ignore.yaml",
			goldenFile:   "markdown/ignore-both.md",
			template:     "../templates/example",
			outputFile:   "markdown/ignore-both.generated",
			ignoreTypes:  []string{"Duration"},
			ignoreGroups: []string{"App"},
		},

		// --------------------------------
		// Go-env-example template tests (.env.example)
		// --------------------------------
		{
			name:       "go-env-example/skip_env_tag",
			configFile: "go-env-example/skip_env_tag.yaml",
			template:   "../templates/go-env-example",
			goldenFile: "go-env-example/skip_env_tag.env",
			outputFile: "go-env-example/skip_env_tag.generated",
		},
		{
			name:       "go-env-example/basic",
			configFile: "go-env-example/basic.yaml",
			template:   "../templates/go-env-example",
			goldenFile: "go-env-example/basic.env",
			outputFile: "go-env-example/basic.generated",
		},
		{
			name:       "go-env-example/minimal",
			configFile: "go-env-example/minimal.yaml",
			template:   "../templates/go-env-example",
			goldenFile: "go-env-example/minimal.env",
			outputFile: "go-env-example/minimal.generated",
		},
		{
			name:       "go-env-example/prefix",
			configFile: "go-env-example/prefix.yaml",
			template:   "../templates/go-env-example",
			goldenFile: "go-env-example/prefix.env",
			outputFile: "go-env-example/prefix.generated",
		},
		{
			name:       "go-env-example/types",
			configFile: "go-env-example/types.yaml",
			template:   "../templates/go-env-example",
			goldenFile: "go-env-example/types.env",
			outputFile: "go-env-example/types.generated",
		},
		{
			name:         "go-env-example/ignore-types",
			configFile:   "go-env-example/ignore.yaml",
			goldenFile:   "go-env-example/ignore-types.env",
			template:     "../templates/example",
			outputFile:   "go-env-example/ignore-types.generated",
			ignoreTypes:  []string{"URL"},
			ignoreGroups: nil,
		},
		{
			name:         "go-env-example/ignore-groups",
			configFile:   "go-env-example/ignore.yaml",
			goldenFile:   "go-env-example/ignore-groups.env",
			template:     "../templates/example",
			outputFile:   "go-env-example/ignore-groups.generated",
			ignoreTypes:  nil,
			ignoreGroups: []string{"DatabaseConfig"},
		},
		{
			name:         "go-env-example/ignore-both",
			configFile:   "go-env-example/ignore.yaml",
			goldenFile:   "go-env-example/ignore-both.env",
			template:     "../templates/example",
			outputFile:   "go-env-example/ignore-both.generated",
			ignoreTypes:  []string{"Duration"},
			ignoreGroups: []string{"App"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.fromURL {
				server := CreateTestServer(t, tt.template)
				defer server.Close()

				tt.template = server.URL
			}

			// Generate file
			err := envgen.Generate(t.Context(), envgen.Options{
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

// CreateTestServer creates a test server that serves content from a file.
func CreateTestServer(t *testing.T, filePath string) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the file
		data, _ := os.ReadFile(filePath)

		// Write file content to the response
		_, _ = w.Write(data)
	})

	// Create and start test server
	return httptest.NewServer(handler)
}
