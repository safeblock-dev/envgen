package user_template_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_template"
)

func TestInvalidTemplateError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		message  string
		expected string
	}{
		{
			name:     "with path",
			path:     "/path/to/template.tmpl",
			message:  "invalid template",
			expected: "invalid template: /path/to/template.tmpl: invalid template",
		},
		{
			name:     "empty path",
			path:     "",
			message:  "invalid template",
			expected: "invalid template: : invalid template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := &user_template.InvalidTemplateError{
				Path:    tt.path,
				Message: tt.message,
			}

			require.Equal(t, tt.expected, err.Error())
		})
	}
}
