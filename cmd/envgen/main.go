package main

import (
	"fmt"
	"os"

	"github.com/safeblock-dev/envgen/pkg/envgen"
	"github.com/spf13/cobra"
)

var (
	// Command flags
	configPath   string
	outputPath   string
	templatePath string

	// Version and build time are set during compilation
	version   = "dev"
	buildTime = "unknown"

	// Root command represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "envgen",
		Short: "Environment configuration generator",
		Long: `envgen is a tool for generating environment configuration files.
It supports multiple output formats and templates, making it easy to maintain
consistent configuration across different projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required flags
			if configPath == "" {
				return fmt.Errorf("config path is required")
			}
			if outputPath == "" {
				return fmt.Errorf("output path is required")
			}
			if templatePath == "" {
				return fmt.Errorf("template path is required")
			}

			// Generate configuration
			if err := envgen.Generate(configPath, outputPath, templatePath); err != nil {
				return err
			}

			fmt.Println("Successfully generated configuration files")
			return nil
		},
	}

	// Version command prints the current version of the tool
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("envgen version %s (built at %s)\n", version, buildTime)
		},
	}
)

func init() {
	// Add persistent flags to the root command
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to input YAML configuration file (required)")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "Path to output file (required)")
	rootCmd.PersistentFlags().StringVarP(&templatePath, "template", "t", "", "Path to template file or URL (required)")

	// Add version subcommand
	rootCmd.AddCommand(versionCmd)

	// Add usage examples
	rootCmd.Example = `  # Generate using local template
  envgen -c config.yaml -o config.go -t ./templates/config.tmpl

  # Generate using template from URL
  envgen --config config.yaml --out config.go --template https://raw.githubusercontent.com/user/repo/template.tmpl

  # Show version
  envgen version`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
