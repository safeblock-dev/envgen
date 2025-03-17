package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/safeblock-dev/envgen/internal/user_template"
)

// NewTemplatesCmd creates a new templates command.
func NewTemplatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"templates", "list"},
		Short:   "List available standard templates",
		RunE:    runTemplates,
	}
}

func runTemplates(cmd *cobra.Command, _ []string) error {
	resolver, err := user_template.NewResolver()
	if err != nil {
		return err
	}

	templates, err := resolver.ListTemplates(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No standard templates available")

		return nil
	}

	fmt.Println("Available standard templates:")

	for _, t := range templates {
		if t.Type != "file" {
			continue
		}

		fmt.Printf("  %s\n", t.Name)
	}

	return nil
}
