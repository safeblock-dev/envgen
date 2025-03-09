package envgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

func TestGenerateOptionsValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    envgen.GenerateOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: envgen.GenerateOptions{
				ConfigPath:   "config.yaml",
				OutputPath:   "output.go",
				TemplatePath: "template.tmpl",
			},
			wantErr: false,
		},
		{
			name: "missing config path",
			opts: envgen.GenerateOptions{
				OutputPath:   "output.go",
				TemplatePath: "template.tmpl",
			},
			wantErr: true,
		},
		{
			name: "missing output path",
			opts: envgen.GenerateOptions{
				ConfigPath:   "config.yaml",
				TemplatePath: "template.tmpl",
			},
			wantErr: true,
		},
		{
			name: "missing template path",
			opts: envgen.GenerateOptions{
				ConfigPath: "config.yaml",
				OutputPath: "output.go",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.opts.Validate()
			if tt.wantErr {
				require.Error(t, err, "Validate() should return an error")
			} else {
				require.NoError(t, err, "Validate() should not return an error")
			}
		})
	}
}
