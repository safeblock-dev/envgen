package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/safeblock-dev/envgen/commands"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "envgen",
	Short:   "Environment configuration generator",
	Version: version,
	Long: `envgen is a tool for generating environment configuration files.
It supports multiple output formats and templates, making it easy to maintain
consistent configuration across different projects.`,
	Example: `  # Generate using standard template
  envgen gen -c config.yaml -o config.go -t go-env

  # Generate using local template
  envgen gen -c config.yaml -o config.go -t ./templates/config.tmpl

  # Generate using template from URL
  envgen gen -c config.yaml -o config.go -t https://raw.githubusercontent.com/safeblock-dev/envgen/main/templates/go-env

  # Generate ignoring specific types and groups
  envgen gen -c config.yaml -o config.go -t go-env --ignore-types Duration,URL --ignore-groups Database

  # List available standard templates
  envgen ls`,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(commands.NewGenerateCmd())
	rootCmd.AddCommand(commands.NewTemplatesCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
