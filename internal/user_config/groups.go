package user_config

import (
	"errors"
	"fmt"
)

// Group represents a group of environment variables.
// Example:
//
//	groups:
//	  - name: App                # Required: Group name
//	    description: App         # Optional: Group description
//	    prefix: APP_             # Optional: Environment variable prefix
//	    options:                 # Optional: Additional options
//	      go_name: Appuser       # Optional: Override struct name (Go-specific)
//	    fields:                  # Required: At least one field must be defined
//	      - name: port
//	        type: int
type Group struct {
	Name        string            `yaml:"name"`        // Required: Group name
	Description string            `yaml:"description"` // Optional: Group description
	Prefix      string            `yaml:"prefix"`      // Optional: Environment variable prefix
	Options     map[string]string `yaml:"options"`     // Optional: Group-specific options (go_name, etc)
	Fields      []Field           `yaml:"fields"`      // Required: At least one field must be defined
}

// Validate validates the group user_configuration.
// Returns an error if required fields are missing or if any field is invalid.
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.New("group name is required")
	}

	if len(g.Fields) == 0 {
		return fmt.Errorf("at least one field is required in group %q", g.Name)
	}

	// Check field names uniqueness
	fieldNames := make(map[string]struct{}, len(g.Fields))
	for _, field := range g.Fields {
		if _, exists := fieldNames[field.Name]; exists {
			return fmt.Errorf("duplicate field name %q in group %q", field.Name, g.Name)
		}

		fieldNames[field.Name] = struct{}{}
	}

	for i, field := range g.Fields {
		if err := field.Validate(); err != nil {
			return fmt.Errorf("invalid field %d in group %q: %w", i, g.Name, err)
		}
	}

	return nil
}
