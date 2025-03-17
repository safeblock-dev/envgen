package envgen

import "errors"

// Options contains options for the Generate function.
type Options struct {
	// ConfigPath is the path to the YAML configuration file
	ConfigPath string
	// OutputPath is the path where the generated file will be written
	OutputPath string
	// TemplatePath is the path to the template file, URL, or standard template name
	TemplatePath string
	// IgnoreTypes is a list of type names to ignore during generation
	IgnoreTypes []string
	// IgnoreGroups is a list of group names to ignore during generation
	IgnoreGroups []string
}

// Validate checks if all required options are set.
func (opts *Options) Validate() error {
	if opts.ConfigPath == "" {
		return errors.New("config path is required")
	}

	if opts.OutputPath == "" {
		return errors.New("output path is required")
	}

	if opts.TemplatePath == "" {
		return errors.New("template path is required")
	}

	return nil
}
