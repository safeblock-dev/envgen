package user_output_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/user_output"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "valid path",
			path:     "/path/to/output",
			wantPath: "/path/to/output",
			wantErr:  false,
		},
		{
			name:     "empty path",
			path:     "",
			wantPath: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := user_output.New(tt.path)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, output)

			wantPath, err := filepath.Abs(tt.wantPath)
			require.NoError(t, err)

			require.Equal(t, wantPath, output.GetPath())
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	// Create temporary directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create file",
			path:    filepath.Join(tempDir, "test.txt"),
			wantErr: false,
		},
		{
			name:    "create file in subdirectory",
			path:    filepath.Join(tempDir, "subdir", "test.txt"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := user_output.New(tt.path)
			require.NoError(t, err)

			file, err := output.Create()
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, file)
			file.Close()
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	// Create temporary directory for tests
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
	}{
		{
			name:    "format go file",
			path:    filepath.Join(tempDir, "test.go"),
			content: "package test\n",
			wantErr: false,
		},
		{
			name:    "format non-go file",
			path:    filepath.Join(tempDir, "test.txt"),
			content: "",
			wantErr: false,
		},
		{
			name:    "format invalid go file",
			path:    filepath.Join(tempDir, "invalid.go"),
			content: "package test\nfunc invalid { syntax error",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := user_output.New(tt.path)
			require.NoError(t, err)

			// Create and write file before formatting
			file, err := output.Create()
			require.NoError(t, err)
			require.NotNil(t, file)

			if tt.content != "" {
				_, err = file.WriteString(tt.content)
				require.NoError(t, err)
				err = file.Close()
				require.NoError(t, err)
			}

			err = output.Format(t.Context())
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}
