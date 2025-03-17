package user_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_config"
)

func TestFindType(t *testing.T) {
	t.Parallel()

	cfg := &user_config.Config{
		Types: []user_config.TypeDefinition{
			{
				Name: "LogLevel",
				Type: "string",
			},
		},
	}

	tests := []struct {
		name     string
		typeName string
		wantNil  bool
	}{
		{name: "existing type", typeName: "LogLevel", wantNil: false},
		{name: "non-existent type", typeName: "NonExistent", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := cfg.FindType(tt.typeName)
			if tt.wantNil {
				require.Nil(t, result)

				return
			}

			require.NotNil(t, result)
			require.Equal(t, tt.typeName, result.Name)
		})
	}
}

func TestHasOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  map[string]string
		option   string
		expected bool
	}{
		{
			name: "existing option",
			options: map[string]string{
				"go_package": "config",
			},
			option:   "go_package",
			expected: true,
		},
		{
			name: "empty option",
			options: map[string]string{
				"empty": "",
			},
			option:   "empty",
			expected: true,
		},
		{
			name: "non-existent option",
			options: map[string]string{
				"other": "value",
			},
			option:   "non_existent",
			expected: false,
		},
		{
			name:     "nil options",
			options:  nil,
			option:   "any",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &user_config.Config{
				Options: tt.options,
			}

			result := cfg.HasOption(tt.option)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHasGroupOption(t *testing.T) {
	t.Parallel()

	baseCfg := &user_config.Config{
		Groups: []user_config.Group{
			{
				Fields: []user_config.Field{
					{
						Options: map[string]string{
							"import": "custom/pkg",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		option   string
		expected bool
	}{
		{name: "existing option", option: "import", expected: true},
		{name: "non-existent option", option: "non_existent", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &user_config.Config{
				Groups: make([]user_config.Group, len(baseCfg.Groups)),
			}
			copy(cfg.Groups, baseCfg.Groups)

			result := cfg.HasGroupOption(tt.option)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		option   string
		options  map[string]string
		expected string
	}{
		{
			name:   "existing option",
			option: "go_package",
			options: map[string]string{
				"go_package": "config",
			},
			expected: "config",
		},
		{
			name:     "non-existent option",
			option:   "non_existent",
			options:  map[string]string{},
			expected: "",
		},
		{
			name:     "nil options",
			option:   "any",
			options:  nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &user_config.Config{
				Options: tt.options,
			}

			result := cfg.GetOption(tt.option)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetGroupOption(t *testing.T) {
	t.Parallel()

	baseCfg := &user_config.Config{
		Groups: []user_config.Group{
			{
				Fields: []user_config.Field{
					{
						Options: map[string]string{
							"import": "custom/pkg",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		option   string
		expected string
	}{
		{name: "existing option", option: "import", expected: "custom/pkg"},
		{name: "non-existent option", option: "non_existent", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of configuration for each test
			cfg := &user_config.Config{
				Groups: make([]user_config.Group, len(baseCfg.Groups)),
			}
			copy(cfg.Groups, baseCfg.Groups)

			result := cfg.GetGroupOption(tt.option)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetImports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cfg      *user_config.Config
		expected []string
	}{
		{
			name: "with imports",
			cfg: &user_config.Config{
				Types: []user_config.TypeDefinition{
					{
						Name:   "CustomType",
						Type:   "string",
						Import: "custom/pkg",
					},
				},
				Groups: []user_config.Group{
					{
						Fields: []user_config.Field{
							{
								Type: "CustomType",
							},
						},
					},
				},
			},
			expected: []string{"custom/pkg"},
		},
		{
			name: "no types",
			cfg: &user_config.Config{
				Groups: []user_config.Group{
					{
						Fields: []user_config.Field{
							{
								Type: "string",
							},
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "no groups",
			cfg: &user_config.Config{
				Types: []user_config.TypeDefinition{
					{
						Name:   "CustomType",
						Type:   "string",
						Import: "custom/pkg",
					},
				},
			},
			expected: nil,
		},
		{
			name:     "empty config",
			cfg:      &user_config.Config{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.cfg.GetImports()
			require.Equal(t, tt.expected, result)
		})
	}
}
