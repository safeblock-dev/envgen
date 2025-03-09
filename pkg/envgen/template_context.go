package envgen

import (
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/safeblock-dev/envgen/pkg/envgen/templatefuncs"
)

// TemplateContext contains all the data needed for code generation.
// It provides access to configuration, paths, and template functions.
type TemplateContext struct {
	Config     *Config // Configuration from YAML file
	ConfigPath string  // Path to configuration file
	OutPath    string  // Path to output file
	TmplPath   string  // Path to template file
}

// NewTemplateContext creates a new template context with the provided configuration and paths.
// It initializes the context with all necessary data for template rendering.
func NewTemplateContext(cfg *Config, configPath, outPath, tmplPath string) *TemplateContext {
	return &TemplateContext{
		Config:     cfg,
		ConfigPath: configPath,
		OutPath:    outPath,
		TmplPath:   tmplPath,
	}
}

// ToRelativePath converts an absolute path to a relative path from the output file directory.
func (tc *TemplateContext) ToRelativePath(path string) string {
	// Check if path is a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	if !filepath.IsAbs(path) {
		// If path is already relative, make it relative to output file
		outDir := filepath.Dir(tc.OutPath)

		rel, err := filepath.Rel(outDir, path)
		if err == nil {
			return rel
		}
	}

	return path
}

// GetTemplateFuncs returns a map of functions available in templates.
// These functions can be used for string manipulation, type conversion,
// date formatting, and path operations.
func (tc *TemplateContext) GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// String transformations
		"title":  templatefuncs.Title,
		"upper":  strings.ToUpper,
		"lower":  strings.ToLower,
		"camel":  templatefuncs.ToCamelCase,
		"snake":  templatefuncs.ToSnakeCase,
		"kebab":  templatefuncs.ToKebabCase,
		"pascal": templatefuncs.ToPascalCase,
		"append": templatefuncs.Append,
		"uniq":   templatefuncs.Uniq,
		"slice":  templatefuncs.Slice,

		// Type conversions
		"toString": templatefuncs.ToString,
		"toInt":    templatefuncs.ToInt,
		"toBool":   templatefuncs.ToBool,

		// Date and time functions
		"now":        time.Now,
		"formatTime": templatefuncs.FormatTime,
		"date": func() string {
			return time.Now().Format("2006-01-02")
		},
		"datetime": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},

		// Conditional operations
		"default":  templatefuncs.DefaultValue,
		"coalesce": templatefuncs.Coalesce,
		"ternary":  templatefuncs.Ternary,

		// String operations
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"replace":   strings.ReplaceAll,
		"trim":      strings.TrimSpace,
		"join":      strings.Join,
		"split":     strings.Split,

		// Path operations
		"getDirName":  templatefuncs.GetDirName,  // Get directory name from path
		"getFileName": templatefuncs.GetFileName, // Get file name from path
		"getFileExt":  templatefuncs.GetFileExt,  // Get file extension
		"joinPaths":   templatefuncs.JoinPaths,   // Join path components

		// File paths
		"getConfigPath":   func() string { return tc.ToRelativePath(tc.ConfigPath) }, // Path to configuration file
		"getOutputPath":   func() string { return tc.ToRelativePath(tc.OutPath) },    // Path to output file
		"getTemplatePath": func() string { return tc.ToRelativePath(tc.TmplPath) },   // Path to template file

		// Configuration helpers
		"hasOption": tc.HasOption,

		// Type helpers
		"findType":   tc.Config.FindType,
		"getImports": tc.GetImports,
		"typeImport": func(typeName string) string {
			if t := tc.Config.FindType(typeName); t != nil {
				return t.Import
			}

			return ""
		},
	}
}

// HasOption checks if any field in the groups has a non-empty value for the specified option.
// This is used to conditionally include sections in templates based on field options.
func (tc *TemplateContext) HasOption(groups []Group, option string) bool {
	for _, group := range groups {
		for _, field := range group.Fields {
			if field.Options[option] != "" {
				return true
			}
		}
	}

	return false
}

// GetImports returns a list of unique imports from type definitions that are used in fields.
func (tc *TemplateContext) GetImports() []string {
	// Early return if no types defined
	if len(tc.Config.Types) == 0 {
		return nil
	}

	// Create a map of type names to their imports for O(1) lookup
	typeImports := make(map[string]string, len(tc.Config.Types))

	for _, t := range tc.Config.Types {
		if t.Import != "" {
			typeImports[t.Name] = t.Import
		}
	}

	// Use map to store unique imports
	uniqueImports := make(map[string]struct{})

	// Collect imports from used types
	for _, group := range tc.Config.Groups {
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
