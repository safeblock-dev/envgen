package envgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

func TestNewTemplateContext(t *testing.T) {
	t.Parallel()

	cfg := &envgen.Config{}
	configPath := "/path/to/config.yaml"
	outPath := "/path/to/output.go"
	tmplPath := "/path/to/template.tmpl"

	tc := envgen.NewTemplateContext(cfg, configPath, outPath, tmplPath)

	require.Equal(t, cfg, tc.Config, "Config not set correctly")
	require.Equal(t, configPath, tc.ConfigPath, "ConfigPath not set correctly")
	require.Equal(t, outPath, tc.OutPath, "OutPath not set correctly")
	require.Equal(t, tmplPath, tc.TmplPath, "TmplPath not set correctly")
}

func TestToRelativePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{
			name:     "absolute path",
			basePath: "/path/to/output.go",
			path:     "/absolute/path/to/file",
			want:     "/absolute/path/to/file",
		},
		{
			name:     "relative path",
			basePath: "/path/to/output.go",
			path:     "relative/path/to/file",
			want:     "relative/path/to/file",
		},
		{
			name:     "http url",
			basePath: "/path/to/output.go",
			path:     "http://example.com/template",
			want:     "http://example.com/template",
		},
		{
			name:     "https url",
			basePath: "/path/to/output.go",
			path:     "https://example.com/template",
			want:     "https://example.com/template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tc := envgen.NewTemplateContext(&envgen.Config{}, "/config.yaml", tt.basePath, "/template.tmpl")
			got := tc.ToRelativePath(tt.path)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHasOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		groups []envgen.Group
		option string
		want   bool
	}{
		{
			name: "option exists",
			groups: []envgen.Group{
				{
					Fields: []envgen.Field{
						{
							Options: map[string]string{
								"test": "value",
							},
						},
					},
				},
			},
			option: "test",
			want:   true,
		},
		{
			name: "option doesn't exist",
			groups: []envgen.Group{
				{
					Fields: []envgen.Field{
						{
							Options: map[string]string{
								"other": "value",
							},
						},
					},
				},
			},
			option: "test",
			want:   false,
		},
		{
			name:   "empty groups",
			groups: []envgen.Group{},
			option: "test",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tc := envgen.NewTemplateContext(&envgen.Config{}, "", "", "")
			got := tc.HasOption(tt.groups, tt.option)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetImports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *envgen.Config
		want   []string
	}{
		{
			name: "with imports",
			config: &envgen.Config{
				Types: []envgen.TypeDefinition{
					{
						Name:   "TestType",
						Import: "test/package",
					},
				},
				Groups: []envgen.Group{
					{
						Fields: []envgen.Field{
							{
								Type: "TestType",
							},
						},
					},
				},
			},
			want: []string{"test/package"},
		},
		{
			name: "no imports",
			config: &envgen.Config{
				Types: []envgen.TypeDefinition{
					{
						Name: "TestType",
					},
				},
				Groups: []envgen.Group{
					{
						Fields: []envgen.Field{
							{
								Type: "TestType",
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name:   "empty config",
			config: &envgen.Config{},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tc := envgen.NewTemplateContext(tt.config, "", "", "")
			got := tc.GetImports()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetTemplateFuncs(t *testing.T) {
	t.Parallel()

	tc := envgen.NewTemplateContext(&envgen.Config{}, "/config.yaml", "/output.go", "/template.tmpl")
	funcs := tc.GetTemplateFuncs()

	// Check for required functions
	requiredFuncs := []string{
		"title", "upper", "lower", "camel", "snake", "kebab", "pascal",
		"append", "uniq", "slice", "toString", "toInt", "toBool",
		"now", "formatTime", "date", "datetime", "default", "coalesce",
		"ternary", "contains", "hasPrefix", "hasSuffix", "replace",
		"trim", "join", "split", "getDirName", "getFileName", "getFileExt",
		"joinPaths", "getConfigPath", "getOutputPath", "getTemplatePath",
		"hasOption", "findType", "getImports", "typeImport",
	}

	for _, name := range requiredFuncs {
		require.Contains(t, funcs, name, "Required function %s not found in template functions", name)
	}
}
