package user_config

import "sort"

// Functions that the template uses

// GetPath returns the path to the user_configuration file.
func (c *Config) GetPath() string {
	return c.path
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

// HasOption checks if the specified option exists in the user_configuration.
// This is used to conditionally include sections in templates based on user_configuration options.
func (c *Config) HasOption(option string) bool {
	if c.Options == nil {
		return false
	}

	_, ok := c.Options[option]

	return ok
}

// HasGroupOption checks if any field in the groups has a non-empty value for the specified option.
// This is used to conditionally include sections in templates based on field options.
func (c *Config) HasGroupOption(option string) bool {
	for _, group := range c.Groups {
		for _, field := range group.Fields {
			if field.Options != nil {
				_, ok := field.Options[option]

				return ok
			}
		}
	}

	return false
}

// GetOption returns the value of the specified option from the user_configuration.
// If the option is not found, returns an empty string.
func (c *Config) GetOption(option string) string {
	if c.Options == nil {
		return ""
	}

	return c.Options[option]
}

// GetGroupOption returns the value of the specified option from the first field that has it.
// If no field has the option, returns an empty string.
func (c *Config) GetGroupOption(option string) string {
	for _, group := range c.Groups {
		for _, field := range group.Fields {
			if field.Options != nil {
				if value, ok := field.Options[option]; ok {
					return value
				}
			}
		}
	}

	return ""
}

// GetImports returns a list of unique imports from type definitions that are used in fields.
func (c *Config) GetImports() []string {
	// Early return if no types defined
	if len(c.Types) == 0 {
		return nil
	}

	// Create a map of type names to their imports for O(1) lookup
	typeImports := make(map[string]string, len(c.Types))

	for _, t := range c.Types {
		if t.Import != "" {
			typeImports[t.Name] = t.Import
		}
	}

	// Use map to store unique imports
	uniqueImports := make(map[string]struct{})

	// Collect imports from used types
	for _, group := range c.Groups {
		for _, field := range group.Fields {
			if imp, exists := typeImports[field.Type]; exists {
				uniqueImports[imp] = struct{}{}
			}
		}
	}

	// Convert unique imports to slice
	if len(uniqueImports) == 0 {
		return nil
	}

	imports := make([]string, 0, len(uniqueImports))
	for imp := range uniqueImports {
		imports = append(imports, imp)
	}

	// Sort imports for consistent output
	sort.Strings(imports)

	return imports
}
