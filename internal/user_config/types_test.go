package user_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_config"
)

func TestHasValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{name: "with values", values: []string{"debug", "info"}, expected: true},
		{name: "empty values", values: []string{}, expected: false},
		{name: "nil values", values: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			td := user_config.TypeDefinition{
				Values: tt.values,
			}

			result := td.HasValues()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeDefinitionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		td      user_config.TypeDefinition
		wantErr bool
	}{
		{
			name: "valid type",
			td: user_config.TypeDefinition{
				Name: "LogLevel",
				Type: "string",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			td: user_config.TypeDefinition{
				Type: "string",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			td: user_config.TypeDefinition{
				Name: "LogLevel",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.td.Validate()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}
