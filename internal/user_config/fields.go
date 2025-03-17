package user_config

import (
	"errors"
	"fmt"
)

// Field represents an environment variable field user_configuration.
// Example:
//
//	fields:
//	  - name: port             # Required: Environment variable name
//	    type: int              # Required: Field type (built-in or custom type)
//	    description: Port      # Optional: Field description
//	    default: "8080"        # Optional: Default value
//	    required: true         # Optional: Whether the field is required
//	    example: "8080"        # Optional: Example value for documentation
//	    options:               # Optional: Additional options
//	      import: "custom/pkg" # Optional: Import path for custom types
//	      name_field: Port     # Optional: Override struct field name
type Field struct {
	Name        string            `yaml:"name"`        // Required: Environment variable name
	Type        string            `yaml:"type"`        // Required: Field type (built-in or custom type)
	Description string            `yaml:"description"` // Optional: Field description
	Default     string            `yaml:"default"`     // Optional: Default value
	Required    bool              `yaml:"required"`    // Optional: Whether the field is required
	Example     string            `yaml:"example"`     // Optional: Example value for documentation
	Options     map[string]string `yaml:"options"`     // Optional: Field-specific options (import, name_field, etc)
}

// Validate validates the field user_configuration.
// Returns an error if required fields are missing.
func (f *Field) Validate() error {
	if f.Name == "" {
		return errors.New("field name is required")
	}

	if f.Type == "" {
		return fmt.Errorf("field type is required for field %q", f.Name)
	}

	return nil
}
