package user_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_config"
)

func TestGroupValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		group   user_config.Group
		wantErr bool
	}{
		{
			name: "valid group",
			group: user_config.Group{
				Name: "app",
				Fields: []user_config.Field{
					{
						Name: "port",
						Type: "int",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			group: user_config.Group{
				Fields: []user_config.Field{
					{
						Name: "port",
						Type: "int",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no fields",
			group: user_config.Group{
				Name: "app",
			},
			wantErr: true,
		},
		{
			name: "invalid field",
			group: user_config.Group{
				Name: "app",
				Fields: []user_config.Field{
					{
						Name: "port",
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

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGroupOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		group  user_config.Group
		option string
		value  string
		exists bool
	}{
		{
			name: "existing option",
			group: user_config.Group{
				Name: "app",
				Fields: []user_config.Field{
					{
						Name: "port",
						Type: "int",
					},
				},
				Options: map[string]string{
					"go_name": "AppConfig",
				},
			},
			option: "go_name",
			value:  "AppConfig",
			exists: true,
		},
		{
			name: "non-existent option",
			group: user_config.Group{
				Name: "app",
				Fields: []user_config.Field{
					{
						Name: "port",
						Type: "int",
					},
				},
				Options: map[string]string{
					"go_name": "AppConfig",
				},
			},
			option: "non_existent",
			value:  "",
			exists: false,
		},
		{
			name: "nil options",
			group: user_config.Group{
				Name: "app",
				Fields: []user_config.Field{
					{
						Name: "port",
						Type: "int",
					},
				},
			},
			option: "any",
			value:  "",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, exists := tt.group.Options[tt.option]
			require.Equal(t, tt.exists, exists)

			if tt.exists {
				require.Equal(t, tt.value, value)
			}
		})
	}
}
