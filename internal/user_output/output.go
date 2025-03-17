package user_output

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// DefaultDirPerms are the default permissions for created directories.
	DefaultDirPerms = 0o755
)

// Output represents the output configuration for generated code.
type Output struct {
	path string // Path to the output directory
}

// New creates a new Output instance with the specified path.
func New(path string) (*Output, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return &Output{
		path: path,
	}, nil
}

// GetPath returns the path to the output directory.
func (o *Output) GetPath() string {
	return o.path
}

// Create creates a new file in the output directory.
func (o *Output) Create() (*os.File, error) {
	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(o.path), DefaultDirPerms); err != nil {
		return nil, err
	}

	// Create output file
	out, err := os.Create(o.path)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (o *Output) Format(ctx context.Context) error {
	// Format generated code
	if strings.HasSuffix(o.path, ".go") {
		// Sanitize the path to prevent command injection
		safePath := filepath.Clean(o.path)
		cmd := exec.CommandContext(ctx, "go", "fmt", safePath)

		// Capture stderr output
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to format generated code: %w\nFormatting errors:\n%s", err, stderr.String())
		}
	}

	return nil
}
