package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

const (
	// defaultDirPerm is the default permission for created directories.
	defaultDirPerm = 0o755
)

var (
	configPath   string
	outputPath   string
	templatePath string
	ignoreTypes  []string
	ignoreGroups []string
)

// NewGenerateCmd creates a new generate command.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate"},
		Short:   "Generate configuration files",
		Long: `Generate configuration files based on the provided template and configuration.
It supports multiple output formats and templates, making it easy to maintain
consistent configuration across different projects.

The template can be:
- A standard template name (e.g. 'go-env')
- A local file path (e.g. './templates/config.tmpl')
- A URL (e.g. 'https://example.com/templates/config.tmpl')`,
		RunE: runGenerate,
	}

	// Add flags
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to input YAML configuration file")
	cmd.Flags().StringVarP(&outputPath, "out", "o", "", "Path to output file")
	cmd.Flags().StringVarP(&templatePath, "template", "t", "", "Template name, path, or URL")
	cmd.Flags().StringSliceVar(&ignoreTypes, "ignore-types", nil, "Types to ignore (comma-separated)")
	cmd.Flags().StringSliceVar(&ignoreGroups, "ignore-groups", nil, "Groups to ignore (comma-separated)")

	// Mark required flags
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("out")
	_ = cmd.MarkFlagRequired("template")

	return cmd
}

func runGenerate(cmd *cobra.Command, _ []string) error {
	// Create output directory if it doesn't exist
	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate configuration
	if err := envgen.Generate(cmd.Context(), envgen.Options{
		ConfigPath:   configPath,
		OutputPath:   outputPath,
		TemplatePath: templatePath,
		IgnoreTypes:  ignoreTypes,
		IgnoreGroups: ignoreGroups,
	}); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	fmt.Printf("Generated %s\n", outputPath)

	return nil
}
