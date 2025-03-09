package envgen

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Field represents an environment variable field configuration.
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

// Validate validates the field configuration.
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

// Group represents a group of environment variables.
// Example:
//
//	groups:
//	  - name: App                # Required: Group name
//	    description: App         # Optional: Group description
//	    prefix: APP_             # Optional: Environment variable prefix
//	    options:                 # Optional: Additional options
//	      go_name: AppConfig     # Optional: Override struct name (Go-specific)
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

// Validate validates the group configuration.
// Returns an error if required fields are missing or if any field is invalid.
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.New("group name is required")
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

// Config represents the complete generation configuration.
// Example:
//
//	options:                    # Optional: Template-specific options
//	  go_package: config       # Optional: Go package name
//	types:                     # Optional: Type definitions
//	  - name: LogLevel        # Required: Type name for referencing in fields
//	    type: zerolog.Level   # Required: Type definition
//	groups:                    # Required: At least one group must be defined
//	  - name: App             # Required: Group name
//	    fields:               # Required: At least one field must be defined
//	      - name: log_level   # Required: Field name
//	        type: LogLevel    # Required: Field type
type Config struct {
	Options map[string]string `yaml:"options"` // Optional: Template-specific options
	Types   []TypeDefinition  `yaml:"types"`   // Optional: Type definitions
	Groups  []Group           `yaml:"groups"`  // Required: At least one group must be defined
}

// FindType finds a type definition by name.
// Returns nil if the type is not found.
func (c *Config) FindType(typeName string) *TypeDefinition {
	for _, t := range c.Types {
		if t.Name == typeName {
			return &t
		}
	}

	return nil
}

// FilterTypes removes ignored types from the configuration.
// If ignoreTypes is empty, no filtering is performed.
func (c *Config) FilterTypes(ignoreTypes []string) {
	if len(ignoreTypes) == 0 {
		return
	}

	ignoreSet := make(map[string]struct{}, len(ignoreTypes))
	for _, typeName := range ignoreTypes {
		ignoreSet[typeName] = struct{}{}
	}

	filtered := make([]TypeDefinition, 0, len(c.Types))

	for _, t := range c.Types {
		if _, ignored := ignoreSet[t.Name]; !ignored {
			filtered = append(filtered, t)
		}
	}

	c.Types = filtered
}

// FilterGroups removes ignored groups from the configuration.
// If ignoreGroups is empty, no filtering is performed.
func (c *Config) FilterGroups(ignoreGroups []string) {
	if len(ignoreGroups) == 0 {
		return
	}

	ignoreSet := make(map[string]struct{}, len(ignoreGroups))
	for _, groupName := range ignoreGroups {
		ignoreSet[groupName] = struct{}{}
	}

	filtered := make([]Group, 0, len(c.Groups))

	for _, g := range c.Groups {
		if _, ignored := ignoreSet[g.Name]; !ignored {
			filtered = append(filtered, g)
		}
	}

	c.Groups = filtered
}

// HasValues checks if the type has predefined values.
// Returns true if the type has at least one value defined.
func (t *TypeDefinition) HasValues() bool {
	return len(t.Values) > 0
}

// Validate validates the configuration.
// Returns an error if required fields are missing or if any group is invalid.
func (c *Config) Validate() error {
	if len(c.Groups) == 0 {
		return errors.New("at least one group is required")
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

// LoadConfig loads and parses configuration from file.
// Returns an error if the file cannot be read or parsed.
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
