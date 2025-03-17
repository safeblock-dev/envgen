package envgen

import (
	"context"
	"errors"
	"fmt"
)

// Generate generates code based on the provided options.
func Generate(ctx context.Context, opts Options) error {
	eg, err := New(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create envgen: %w", err)
	}

	return eg.Generate(ctx)
}

// Generate executes the template and writes the result to the output file.
func (e *Envgen) Generate(ctx context.Context) error {
	if e == nil {
		return errors.New("envgen instance is nil")
	}

	template, err := e.Template()
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	outFile, err := e.userOutput.Create()
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Execute template
	if err := template.Execute(outFile, e.userConfig); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := e.userOutput.Format(ctx); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
