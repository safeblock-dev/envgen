package envgen

import (
	"context"
	"fmt"
	"text/template"

	"github.com/safeblock-dev/envgen/internal/user_config"
	"github.com/safeblock-dev/envgen/internal/user_output"
	"github.com/safeblock-dev/envgen/internal/user_template"
)

// Envgen represents the main structure for code generation.
type Envgen struct {
	userTemplate *user_template.Template // Template for code generation
	userConfig   *user_config.Config     // Configuration for code generation
	userOutput   *user_output.Output     // Output configuration for generated code
}

// New creates a new Envgen instance with the specified options.
func New(ctx context.Context, opts Options) (*Envgen, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	envgen := new(Envgen)

	err := envgen.SetConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to add user config: %w", err)
	}

	err = envgen.SetOutput(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to add output config: %w", err)
	}

	err = envgen.SetTemplate(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to add template config: %w", err)
	}

	return envgen, nil
}

// SetConfig sets the configuration for code generation.
func (e *Envgen) SetConfig(opts Options) error {
	// Read and parse configuration
	cfg, err := user_config.New(opts.ConfigPath)
	if err != nil {
		return err
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Filter out ignored types and groups
	cfg.FilterTypes(opts.IgnoreTypes)
	cfg.FilterGroups(opts.IgnoreGroups)

	e.userConfig = cfg

	return nil
}

// SetOutput sets the output configuration for generated code.
func (e *Envgen) SetOutput(opts Options) error {
	userOutput, err := user_output.New(opts.OutputPath)
	if err != nil {
		return err
	}

	e.userOutput = userOutput

	return nil
}

// SetTemplate sets the template for code generation.
func (e *Envgen) SetTemplate(ctx context.Context, opts Options) error {
	userTemplate, err := user_template.New(ctx, opts.TemplatePath)
	if err != nil {
		return err
	}

	e.userTemplate = userTemplate

	return nil
}

// Template returns the compiled template for code generation.
func (e *Envgen) Template() (*template.Template, error) {
	const templateName = "envgen"

	// Create template
	tmpl, err := template.New(templateName).
		Funcs(e.Funcs()).
		Parse(e.userTemplate.GetContent())

	return tmpl, err
}
