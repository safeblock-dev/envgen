package user_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_config"
)

func TestFieldValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		field   user_config.Field
		wantErr bool
	}{
		{
			name: "valid field",
			field: user_config.Field{
				Name: "port",
				Type: "int",
			},
			wantErr: false,
		},
		{
			name: "valid field with description",
			field: user_config.Field{
				Name:        "port",
				Type:        "int",
				Description: "Server port",
			},
			wantErr: false,
		},
		{
			name: "valid field with options",
			field: user_config.Field{
				Name: "port",
				Type: "int",
				Options: map[string]string{
					"default": "8080",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			field: user_config.Field{
				Type: "int",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			field: user_config.Field{
				Name: "",
				Type: "int",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			field: user_config.Field{
				Name: "port",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			field: user_config.Field{
				Name: "port",
				Type: "",
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

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFieldOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		field  user_config.Field
		option string
		value  string
		exists bool
	}{
		{
			name: "existing option",
			field: user_config.Field{
				Name: "port",
				Type: "int",
				Options: map[string]string{
					"import": "custom/pkg",
				},
			},
			option: "import",
			value:  "custom/pkg",
			exists: true,
		},
		{
			name: "non-existent option",
			field: user_config.Field{
				Name: "port",
				Type: "int",
			},
			option: "non_existent",
			value:  "",
			exists: false,
		},
		{
			name: "nil options",
			field: user_config.Field{
				Name: "port",
				Type: "int",
			},
			option: "any",
			value:  "",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, exists := tt.field.Options[tt.option]
			require.Equal(t, tt.exists, exists)

			if tt.exists {
				require.Equal(t, tt.value, value)
			}
		})
	}
}
