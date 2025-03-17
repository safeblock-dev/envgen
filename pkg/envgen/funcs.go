package envgen

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/safeblock-dev/envgen/internal/template_funcs"
)

// Funcs returns a map of functions available in templates.
// These functions can be used for string manipulation, type conversion,
// date formatting, and path operations.
//
//nolint:funlen // This function is a map of all available template functions
func (e *Envgen) Funcs() template.FuncMap {
	if e == nil {
		return nil
	}

	return template.FuncMap{
		// String transformations
		"title":  template_funcs.Title,
		"upper":  strings.ToUpper,
		"lower":  strings.ToLower,
		"camel":  template_funcs.ToCamelCase,
		"snake":  template_funcs.ToSnakeCase,
		"kebab":  template_funcs.ToKebabCase,
		"pascal": template_funcs.ToPascalCase,
		"append": template_funcs.StringAppend,
		"uniq":   template_funcs.StringUniq,
		"slice":  template_funcs.StringSlice,

		// Type conversions
		"toString": template_funcs.ToString,
		"toInt":    template_funcs.ToInt,
		"toBool":   template_funcs.ToBool,

		// Date and time functions
		"now":        time.Now,
		"formatTime": template_funcs.FormatTime,
		"date": func() string {
			return time.Now().Format("2006-01-02")
		},
		"datetime": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},

		// Conditional operations
		"default":  template_funcs.DefaultValue,
		"coalesce": template_funcs.Coalesce,
		"ternary":  template_funcs.Ternary,

		// String operations
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"isURL":     template_funcs.IsURL,
		"join":      strings.Join,
		"oneline":   template_funcs.Oneline,
		"replace":   strings.ReplaceAll,
		"split":     strings.Split,
		"trim":      strings.TrimSpace,

		// Path operations
		"pathDir":  filepath.Dir,  // Get directory name from path
		"pathBase": filepath.Base, // Get file name from path
		"pathExt":  filepath.Ext,  // Get file extension
		"pathRel": func(outputPath, path string) string {
			relPath, err := filepath.Rel(outputPath, path)
			if err != nil {
				return ""
			}

			return relPath
		},

		// File paths
		"getConfigPath":   e.userConfig.GetPath,   // Path to configuration file
		"getOutputPath":   e.userOutput.GetPath,   // Path to output file
		"getTemplatePath": e.userTemplate.GetPath, // Path or URL to template file

		// Golang
		"goCommentGenerate": e.goCommentGenerate,

		// Configuration helpers
		"hasOption":       e.userConfig.HasOption,
		"hasGroupOption":  e.userConfig.HasGroupOption,
		"getOption":       e.userConfig.GetOption,
		"getGroupOption":  e.userConfig.GetGroupOption,
		"processTemplate": e.ProcessTemplate,

		// Type helpers
		"findType":   e.userConfig.FindType,
		"getImports": e.userConfig.GetImports,
	}
}

func (e *Envgen) ProcessTemplate(content string) string {
	if e == nil || content == "" {
		return content
	}

	// Create template with functions
	tmpl, err := template.New("process").Funcs(e.Funcs()).Parse(content)
	if err != nil {
		log.Println("error parsing template:", err)

		return content
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		log.Println("error executing template:", err)

		return content
	}

	return buf.String()
}

func (e *Envgen) goCommentGenerate(configPath, outputFile, templatePath string) string {
	var err error

	if configPath == "" {
		configPath, err = filepath.Rel(
			filepath.Dir(e.userOutput.GetPath()), e.userConfig.GetPath(),
		)
		if err != nil {
			log.Println("error getting relative config path:", err)
		}
	}

	if outputFile == "" {
		outputFile = filepath.Base(e.userOutput.GetPath())
	}

	if templatePath == "" {
		templatePath, err = filepath.Rel(
			filepath.Dir(e.userOutput.GetPath()), e.userTemplate.GetPath(),
		)
		if err != nil {
			log.Println("error getting relative template path:", err)
		}
	}

	path := fmt.Sprintf("//go:generate envgen -c %s -o %s -t %s", configPath, outputFile, templatePath)

	return path
}
