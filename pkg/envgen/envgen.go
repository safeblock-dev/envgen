package envgen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	// DefaultDirPerms are the default permissions for created directories.
	DefaultDirPerms = 0o755
)

// GenerateOptions contains options for the Generate function.
type GenerateOptions struct {
	// ConfigPath is the path to the YAML configuration file
	ConfigPath string
	// OutputPath is the path where the generated file will be written
	OutputPath string
	// TemplatePath is the path to the template file
	TemplatePath string
	// IgnoreTypes is a list of type names to ignore during generation
	IgnoreTypes []string
	// IgnoreGroups is a list of group names to ignore during generation
	IgnoreGroups []string
}

// Validate checks if all required options are set.
func (opts *GenerateOptions) Validate() error {
	if opts.ConfigPath == "" {
		return errors.New("config path is required")
	}

	if opts.OutputPath == "" {
		return errors.New("output path is required")
	}

	if opts.TemplatePath == "" {
		return errors.New("template path is required")
	}

	return nil
}

// Generate generates code based on configuration and template.
func Generate(opts GenerateOptions) error {
	// Validate options
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	// Read and parse configuration
	cfg, err := LoadConfig(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Filter out ignored types and groups
	cfg.FilterTypes(opts.IgnoreTypes)
	cfg.FilterGroups(opts.IgnoreGroups)

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), DefaultDirPerms); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create template context
	ctx := NewTemplateContext(cfg, opts.ConfigPath, opts.OutputPath, opts.TemplatePath)

	// Generate code
	if err := generateCode(ctx); err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

// generateCode generates the final code using the template and configuration.
func generateCode(ctx *TemplateContext) error {
	// Load template content
	templateContent, err := LoadTemplate(ctx.TmplPath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Create template
	tmpl, err := template.New("envgen").
		Funcs(ctx.GetTemplateFuncs()).
		Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	out, err := os.Create(ctx.OutPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Execute template
	if err := tmpl.Execute(out, ctx.Config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format generated code
	if strings.HasSuffix(ctx.OutPath, ".go") {
		// Sanitize the path to prevent command injection
		safePath := filepath.Clean(ctx.OutPath)
		cmd := exec.CommandContext(context.Background(), "go", "fmt", safePath)

		// Capture stderr output
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to format generated code: %w\nFormatting errors:\n%s", err, stderr.String())
		}
	}

	return nil
}
