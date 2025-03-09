package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/safeblock-dev/envgen/pkg/envgen"
)

var (
	// Command flags.
	configPath   string
	outputPath   string
	templatePath string
	ignoreTypes  []string
	ignoreGroups []string

	// Version is set during compilation.
	version   = "unknown"
	buildTime = time.Now().Format(time.RFC3339)

	// Root command represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "envgen",
		Short: "Environment configuration generator",
		Long: `envgen is a tool for generating environment configuration files.
It supports multiple output formats and templates, making it easy to maintain
consistent configuration across different projects.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Generate configuration
			if err := envgen.Generate(envgen.GenerateOptions{
				ConfigPath:   configPath,
				OutputPath:   outputPath,
				TemplatePath: templatePath,
				IgnoreTypes:  ignoreTypes,
				IgnoreGroups: ignoreGroups,
			}); err != nil {
				return err
			}

			fmt.Println("Successfully generated configuration files")

			return nil
		},
	}

	// Version command prints the current version of the tool.
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("envgen version %s (built at %s)\n", version, buildTime)
		},
	}
)

func init() { //nolint: gochecknoinits
	// Add persistent flags to the root command
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to input YAML configuration file")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "Path to output file")
	rootCmd.PersistentFlags().StringVarP(&templatePath, "template", "t", "", "Path to template file or URL")
	rootCmd.PersistentFlags().StringSliceVar(&ignoreTypes, "ignore-types", nil, "Comma-separated list of types to ignore")
	rootCmd.PersistentFlags().StringSliceVar(&ignoreGroups, "ignore-groups", nil, "Comma-separated list of groups to ignore")

	// Mark required flags
	_ = rootCmd.MarkPersistentFlagRequired("config")
	_ = rootCmd.MarkPersistentFlagRequired("out")
	_ = rootCmd.MarkPersistentFlagRequired("template")

	// Add version subcommand
	rootCmd.AddCommand(versionCmd)

	// Add usage examples
	rootCmd.Example = `  # Generate using local template
  envgen -c config.yaml -o config.go -t ./templates/config.tmpl

  # Generate using template from URL
  envgen --c config.yaml --o config.go --t https://raw.githubusercontent.com/safeblock-dev/envgen/refs/heads/main/templates/go-env

  # Generate ignoring specific types and groups
  envgen -c config.yaml -o config.go -t ./templates/config.tmpl --ignore-types Duration,URL --ignore-groups Database

  # Show version
  envgen version`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
