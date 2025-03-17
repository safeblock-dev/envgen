package user_config

import (
	"errors"
	"fmt"
)

// TypeDefinition describes a type and its possible values.
// Example:
//
//	types:
//	  - name: LogLevel                  # Required: Type name for referencing in fields
//	    type: zerolog.Level             # Required: Type definition (built-in or custom)
//	    description: Log level          # Optional: Type description
//	    import: "github.com/rs/zerolog" # Optional: Import path for custom types
//	    values: [debug, info, no]       # Optional: Possible values for documentation
type TypeDefinition struct {
	Name        string   `yaml:"name"`        // Required: Type name for referencing in fields
	Type        string   `yaml:"type"`        // Required: Type definition (built-in or custom)
	Import      string   `yaml:"import"`      // Optional: Import path for custom types
	Description string   `yaml:"description"` // Optional: Type description
	Values      []string `yaml:"values"`      // Optional: Possible values for documentation
}

// HasValues checks if the type has predefined values.
// Returns true if the type has at least one value defined.
func (t *TypeDefinition) HasValues() bool {
	return len(t.Values) > 0
}

// Validate validates the type definition.
// Returns an error if required fields are missing.
func (t *TypeDefinition) Validate() error {
	if t.Name == "" {
		return errors.New("type name is required")
	}

	if t.Type == "" {
		return fmt.Errorf("type definition is required for type %q", t.Name)
	}

	return nil
}
