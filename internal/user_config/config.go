package user_config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the complete generation user_configuration.
// Example:
//
//	options:                    # Optional: Template-specific options
//	  go_package: user_config       # Optional: Go package name
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

	path string `yaml:"-"` // Path to user_configuration file (not serialized)
}

// New loads and parses user_configuration from file.
// Returns an error if the file cannot be read or parsed.
func New(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read user_config file: %w", err)
	}

	var cfg Config
	cfg.path = path

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse user_config file: %w", err)
	}

	return &cfg, nil
}

// FilterTypes removes ignored types from the user_configuration.
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

// FilterGroups removes ignored groups from the user_configuration.
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

// Validate validates the user_configuration.
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

// GetOptions returns the template-specific options.
func (c *Config) GetOptions() map[string]string {
	if c.Options == nil {
		return nil
	}

	options := make(map[string]string, len(c.Options))
	for k, v := range c.Options {
		options[k] = v
	}

	return options
}

// GetTypes returns the type definitions.
func (c *Config) GetTypes() []TypeDefinition {
	if c.Types == nil {
		return nil
	}

	types := make([]TypeDefinition, len(c.Types))
	copy(types, c.Types)

	return types
}

// GetGroups returns the groups defined in the configuration.
func (c *Config) GetGroups() []Group {
	if c.Groups == nil {
		return nil
	}

	groups := make([]Group, len(c.Groups))
	copy(groups, c.Groups)

	return groups
}
