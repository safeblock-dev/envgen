package envgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

const (
	testFilePerm = 0o644
)

func TestFieldValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		field   envgen.Field
		wantErr bool
	}{
		{
			name: "valid field",
			field: envgen.Field{
				Name: "test",
				Type: "string",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			field: envgen.Field{
				Type: "string",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			field: envgen.Field{
				Name: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.field.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGroupValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		group   envgen.Group
		wantErr bool
	}{
		{
			name: "valid group",
			group: envgen.Group{
				Name: "test",
				Fields: []envgen.Field{
					{
						Name: "field1",
						Type: "string",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			group: envgen.Group{
				Fields: []envgen.Field{
					{
						Name: "field1",
						Type: "string",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no fields",
			group: envgen.Group{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "invalid field",
			group: envgen.Group{
				Name: "test",
				Fields: []envgen.Field{
					{
						Name: "field1",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.group.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  envgen.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: envgen.Config{
				Groups: []envgen.Group{
					{
						Name: "test",
						Fields: []envgen.Field{
							{
								Name: "field1",
								Type: "string",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "no groups",
			config:  envgen.Config{},
			wantErr: true,
		},
		{
			name: "invalid group",
			config: envgen.Config{
				Groups: []envgen.Group{
					{
						Name: "test",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFindType(t *testing.T) {
	t.Parallel()

	baseTypes := []envgen.TypeDefinition{
		{
			Name: "TestType",
			Type: "string",
		},
	}

	tests := []struct {
		name     string
		typeName string
		want     *envgen.TypeDefinition
	}{
		{
			name:     "existing type",
			typeName: "TestType",
			want: &envgen.TypeDefinition{
				Name: "TestType",
				Type: "string",
			},
		},
		{
			name:     "non-existent type",
			typeName: "NonExistentType",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new config for each test
			config := envgen.Config{
				Types: make([]envgen.TypeDefinition, len(baseTypes)),
			}
			copy(config.Types, baseTypes)

			got := config.FindType(tt.typeName)
			if tt.want == nil {
				require.Nil(t, got)

				return
			}

			require.NotNil(t, got)
			require.Equal(t, tt.want.Name, got.Name)
			require.Equal(t, tt.want.Type, got.Type)
		})
	}
}

func TestFilterTypes(t *testing.T) {
	t.Parallel()

	baseTypes := []envgen.TypeDefinition{
		{
			Name: "Type1",
			Type: "string",
		},
		{
			Name: "Type2",
			Type: "int",
		},
	}

	tests := []struct {
		name        string
		ignoreTypes []string
		want        []string
	}{
		{
			name:        "filter one type",
			ignoreTypes: []string{"Type1"},
			want:        []string{"Type2"},
		},
		{
			name:        "filter all types",
			ignoreTypes: []string{"Type1", "Type2"},
			want:        []string{},
		},
		{
			name:        "filter no types",
			ignoreTypes: []string{},
			want:        []string{"Type1", "Type2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new config for each test
			config := envgen.Config{
				Types: make([]envgen.TypeDefinition, len(baseTypes)),
			}
			copy(config.Types, baseTypes)

			config.FilterTypes(tt.ignoreTypes)

			got := make([]string, len(config.Types))
			for i, typ := range config.Types {
				got[i] = typ.Name
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestFilterGroups(t *testing.T) {
	t.Parallel()

	baseGroups := []envgen.Group{
		{
			Name: "Group1",
			Fields: []envgen.Field{
				{
					Name: "field1",
					Type: "string",
				},
			},
		},
		{
			Name: "Group2",
			Fields: []envgen.Field{
				{
					Name: "field2",
					Type: "int",
				},
			},
		},
	}

	tests := []struct {
		name         string
		ignoreGroups []string
		want         []string
	}{
		{
			name:         "filter one group",
			ignoreGroups: []string{"Group1"},
			want:         []string{"Group2"},
		},
		{
			name:         "filter all groups",
			ignoreGroups: []string{"Group1", "Group2"},
			want:         []string{},
		},
		{
			name:         "filter no groups",
			ignoreGroups: []string{},
			want:         []string{"Group1", "Group2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new config for each test
			config := envgen.Config{
				Groups: make([]envgen.Group, len(baseGroups)),
			}
			copy(config.Groups, baseGroups)

			config.FilterGroups(tt.ignoreGroups)

			got := make([]string, len(config.Groups))
			for i, g := range config.Groups {
				got[i] = g.Name
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestHasValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  envgen.TypeDefinition
		want bool
	}{
		{
			name: "has values",
			typ: envgen.TypeDefinition{
				Values: []string{"value1", "value2"},
			},
			want: true,
		},
		{
			name: "no values",
			typ: envgen.TypeDefinition{
				Values: []string{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.typ.HasValues()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	// Create temporary file with test configuration
	tmpDir := t.TempDir()
	testConfig := `
types:
  - name: TestType
    type: string
    import: test/package
groups:
  - name: TestGroup
    fields:
      - name: test_field
        type: TestType
`
	testFile := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(testConfig), testFilePerm)
	require.NoError(t, err, "Failed to create test config file")

	// Create file with invalid YAML
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")
	err = os.WriteFile(invalidFile, []byte("invalid: yaml: content:"), testFilePerm)
	require.NoError(t, err, "Failed to create invalid config file")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid config",
			path:    testFile,
			wantErr: false,
		},
		{
			name:    "non-existent file",
			path:    "non_existent.yaml",
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			path:    invalidFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := envgen.LoadConfig(tt.path)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)
		})
	}
}
