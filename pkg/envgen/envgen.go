package envgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Generate generates code based on configuration and template.
// Parameters:
//   - configPath: Path to the YAML configuration file
//   - outPath: Output path for generated files (package name defaults to directory name)
//   - templatePath: Path to the template file (embedded or filesystem)
func Generate(configPath, outPath, templatePath string) error {
	// Validate input parameters
	if err := validateInputs(configPath, outPath, templatePath); err != nil {
		return fmt.Errorf("invalid input parameters: %w", err)
	}

	// Read and parse configuration
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create template context
	ctx := NewTemplateContext(cfg, configPath, outPath, templatePath)

	// Generate code
	if err := generateCode(ctx); err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

// validateInputs validates input parameters
func validateInputs(configPath, outPath, templatePath string) error {
	if configPath == "" {
		return fmt.Errorf("config path is required")
	}
	if outPath == "" {
		return fmt.Errorf("output path is required")
	}
	if templatePath == "" {
		return fmt.Errorf("template path is required")
	}

	return nil
}

// generateCode generates the final code using the template and configuration
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
		if err := exec.Command("go", "fmt", ctx.OutPath).Run(); err != nil {
			return fmt.Errorf("failed to format generated code: %w", err)
		}
	}

	return nil
}
