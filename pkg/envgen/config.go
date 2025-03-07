package envgen

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Field represents an environment variable field configuration.
// Example:
//
//	fields:
//	  - name: port             # Required: Environment variable name
//	    type: int              # Required: Go type (any valid Go type including custom types)
//	    description: Port      # Optional: Field description
//	    default: "8080"        # Optional: Default value
//	    required: true         # Optional: Whether the field is required
//	    example: "8080"        # Optional: Example value for documentation
//	    options:               # Optional: Additional options
//	      import: "custom/pkg" # Optional: Import path for custom types
//	      name_field: Port     # Optional: Override struct field name
type Field struct {
	Name        string            `yaml:"name"`        // Environment variable name
	Description string            `yaml:"description"` // Field description
	Type        string            `yaml:"type"`        // Go type of the field (can be any valid Go type)
	Default     string            `yaml:"default"`     // Default value
	Required    bool              `yaml:"required"`    // Whether the field is required
	Example     string            `yaml:"example"`     // Example value
	Options     map[string]string `yaml:"options"`     // Field-specific options
}

// Validate validates the field configuration
func (f *Field) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("field name is required")
	}
	if f.Type == "" {
		return fmt.Errorf("field type is required for field %q", f.Name)
	}
	return nil
}

// Group represents a group of environment variables.
// Example:
//
//	groups:
//	  - name: App                # Required: Group name
//	    description: App         # Optional: Group description
//	    prefix: APP_             # Optional: Environment variable prefix
//	    options:                 # Optional: Additional options
//	      struct_name: AppConfig # Optional: Override struct name
//	    fields:                  # Required: List of fields
//	      - name: port
//	        type: int
type Group struct {
	Name        string            `yaml:"name"`        // Group name
	Description string            `yaml:"description"` // Group description
	Prefix      string            `yaml:"prefix"`      // Environment variable prefix (optional)
	Options     map[string]string `yaml:"options"`     // Group-specific options
	Fields      []Field           `yaml:"fields"`      // List of fields
}

// Validate validates the group configuration
func (g *Group) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("group name is required")
	}
	if len(g.Fields) == 0 {
		return fmt.Errorf("at least one field is required in group %q", g.Name)
	}
	for i, field := range g.Fields {
		if err := field.Validate(); err != nil {
			return fmt.Errorf("invalid field %d in group %q: %w", i, g.Name, err)
		}
	}
	return nil
}

// TypeDefinition describes a type and its possible values
type TypeDefinition struct {
	Name        string   `yaml:"name"`        // Type name
	Import      string   `yaml:"import"`      // Import path (optional)
	Description string   `yaml:"description"` // Type description
	Values      []string `yaml:"values"`      // Possible values (optional)
}

// Config represents the complete generation configuration.
// Example:
//
//		options: # Optional: template-specific options
//	   		...
//		types:               # Optional: Type definitions
//		  - name: LogLevel
//		    description: Log level
//		    values: [debug, info, warn, error]
//		groups:               # Required: List of groups
//		  - name: App
//		    fields:
//		      - name: log_level # Required: Field name
//		        type: LogLevel # Required: Type name
//	         options: # Optional: template-specific options
//	           ...
type Config struct {
	Options map[string]string `yaml:"options"` // Template-specific options
	Types   []TypeDefinition  `yaml:"types"`   // Type definitions
	Groups  []Group           `yaml:"groups"`  // List of groups
}

// FindType finds a type definition by name
func (c *Config) FindType(typeName string) *TypeDefinition {
	for _, t := range c.Types {
		if t.Name == typeName {
			return &t
		}
	}
	return nil
}

// HasValues checks if the type has predefined values
func (t *TypeDefinition) HasValues() bool {
	return len(t.Values) > 0
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.Groups) == 0 {
		return fmt.Errorf("at least one group is required")
	}
	if c.Options == nil {
		c.Options = make(map[string]string)
	}
	for i, group := range c.Groups {
		if err := group.Validate(); err != nil {
			return fmt.Errorf("invalid group %d: %w", i, err)
		}
	}
	return nil
}

// LoadConfig loads and parses configuration from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
