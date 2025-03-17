package user_config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Create a temporary test file
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yaml")

		// Write test configuration
		err := os.WriteFile(testFile, []byte(`
options:
  go_package: config
types:
  - name: LogLevel
    type: string
    import: github.com/rs/zerolog
groups:
  - name: app
    fields:
      - name: port
        type: int
`), 0o600)
		require.NoError(t, err)

		// Test loading configuration
		cfg, err := user_config.New(testFile)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// Verify configuration
		require.Equal(t, "config", cfg.GetOption("go_package"))
		require.Len(t, cfg.GetTypes(), 1)
		require.Len(t, cfg.GetGroups(), 1)
	})

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()

		cfg, err := user_config.New("non_existent.yaml")
		require.Error(t, err)
		require.Nil(t, cfg)
		require.Contains(t, err.Error(), "failed to read user_config file")
	})

	t.Run("invalid yaml", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "invalid.yaml")

		err := os.WriteFile(testFile, []byte("invalid: yaml: content:"), 0o600)
		require.NoError(t, err)

		cfg, err := user_config.New(testFile)
		require.Error(t, err)
		require.Nil(t, cfg)
		require.Contains(t, err.Error(), "failed to parse user_config file")
	})
}

func TestGetOptions(t *testing.T) {
	t.Parallel()

	cfg := &user_config.Config{
		Options: map[string]string{
			"key": "value",
		},
	}

	options := cfg.GetOptions()
	require.NotNil(t, options)
	require.Equal(t, "value", options["key"])
}

func TestGetTypes(t *testing.T) {
	t.Parallel()

	cfg := &user_config.Config{
		Types: []user_config.TypeDefinition{
			{Name: "Type1", Type: "string"},
		},
	}

	types := cfg.GetTypes()
	require.Len(t, types, 1)
	require.Equal(t, "Type1", types[0].Name)
}

func TestGetGroups(t *testing.T) {
	t.Parallel()

	cfg := &user_config.Config{
		Groups: []user_config.Group{
			{Name: "Group1"},
		},
	}

	groups := cfg.GetGroups()
	require.Len(t, groups, 1)
	require.Equal(t, "Group1", groups[0].Name)
}

func TestFilterTypes(t *testing.T) {
	t.Parallel()

	baseTypes := []user_config.TypeDefinition{
		{Name: "Type1", Type: "string"},
		{Name: "Type2", Type: "int"},
		{Name: "Type3", Type: "bool"},
	}

	tests := []struct {
		name         string
		ignoreTypes  []string
		expectedLen  int
		expectedType string
	}{
		{
			name:         "filter one type",
			ignoreTypes:  []string{"Type1"},
			expectedLen:  2,
			expectedType: "Type2",
		},
		{
			name:         "filter multiple types",
			ignoreTypes:  []string{"Type1", "Type2"},
			expectedLen:  1,
			expectedType: "Type3",
		},
		{
			name:         "filter all types",
			ignoreTypes:  []string{"Type1", "Type2", "Type3"},
			expectedLen:  0,
			expectedType: "",
		},
		{
			name:         "filter no types",
			ignoreTypes:  []string{},
			expectedLen:  3,
			expectedType: "Type1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			types := make([]user_config.TypeDefinition, len(baseTypes))
			copy(types, baseTypes)

			cfg := &user_config.Config{
				Types: types,
			}

			cfg.FilterTypes(tt.ignoreTypes)
			types = cfg.GetTypes()
			require.Len(t, types, tt.expectedLen)

			if tt.expectedType != "" {
				found := false

				for _, t := range types {
					if t.Name == tt.expectedType {
						found = true

						break
					}
				}

				require.True(t, found, "expected type %s not found", tt.expectedType)
			}
		})
	}
}

func TestFilterGroups(t *testing.T) {
	t.Parallel()

	baseGroups := []user_config.Group{
		{Name: "Group1"},
		{Name: "Group2"},
		{Name: "Group3"},
	}

	tests := []struct {
		name          string
		ignoreGroups  []string
		expectedLen   int
		expectedGroup string
	}{
		{
			name:          "filter one group",
			ignoreGroups:  []string{"Group1"},
			expectedLen:   2,
			expectedGroup: "Group2",
		},
		{
			name:          "filter multiple groups",
			ignoreGroups:  []string{"Group1", "Group2"},
			expectedLen:   1,
			expectedGroup: "Group3",
		},
		{
			name:          "filter all groups",
			ignoreGroups:  []string{"Group1", "Group2", "Group3"},
			expectedLen:   0,
			expectedGroup: "",
		},
		{
			name:          "filter no groups",
			ignoreGroups:  []string{},
			expectedLen:   3,
			expectedGroup: "Group1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			groups := make([]user_config.Group, len(baseGroups))
			copy(groups, baseGroups)

			cfg := &user_config.Config{
				Groups: groups,
			}

			cfg.FilterGroups(tt.ignoreGroups)
			groups = cfg.GetGroups()
			require.Len(t, groups, tt.expectedLen)

			if tt.expectedGroup != "" {
				found := false

				for _, g := range groups {
					if g.Name == tt.expectedGroup {
						found = true

						break
					}
				}

				require.True(t, found, "expected group %s not found", tt.expectedGroup)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *user_config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &user_config.Config{
				Groups: []user_config.Group{
					{
						Name: "app",
						Fields: []user_config.Field{
							{
								Name: "port",
								Type: "int",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "no groups",
			cfg:     &user_config.Config{},
			wantErr: true,
		},
		{
			name: "invalid group",
			cfg: &user_config.Config{
				Groups: []user_config.Group{
					{
						Name: "app",
						Fields: []user_config.Field{
							{
								Name: "port",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "nil options",
			cfg: &user_config.Config{
				Groups: []user_config.Group{
					{
						Name: "app",
						Fields: []user_config.Field{
							{
								Name: "port",
								Type: "int",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			if tt.cfg.Options == nil {
				require.NotNil(t, tt.cfg.Options)
			}
		})
	}
}
