package user_template

import (
	"context"
	"strings"
)

// Template represents a template with its metadata.
type Template struct {
	// Name is the template filename without path
	Name string
	// Source indicates where the template comes from
	Source TemplateSource
	// Content is the template content
	Content string
	// ResolvedPath is the resolved path or URL to the template
	ResolvedPath string
}

// Source represents the source of a template.
type TemplateSource int

const (
	// TemplateSourceLocal indicates a template from the local filesystem.
	TemplateSourceLocal TemplateSource = iota
	// TemplateSourceURL indicates a template from a URL.
	TemplateSourceURL
	// TemplateSourceStandard indicates a standard template.
	TemplateSourceStandard
)

func New(ctx context.Context, path string) (*Template, error) {
	resolver, err := NewResolver()
	if err != nil {
		return nil, err
	}

	template, err := resolver.Template(ctx, path)

	return template, err
}

// Validate checks if the template has all required fields.
func (t Template) Validate() error {
	if t.Name == "" {
		return &InvalidTemplateError{
			Path:    t.ResolvedPath,
			Message: "template name is empty",
		}
	}

	if t.Content == "" {
		return &InvalidTemplateError{
			Path:    t.ResolvedPath,
			Message: "template content is empty",
		}
	}

	if t.ResolvedPath == "" {
		return &InvalidTemplateError{
			Path:    t.Name,
			Message: "resolved path is empty",
		}
	}

	return nil
}

// GetName returns the template name (filename without path).
func (t Template) GetName() string {
	return t.Name
}

// GetPath returns the resolved path to the template (URL or local path).
func (t Template) GetPath() string {
	return t.ResolvedPath
}

// IsURL returns true if the template is loaded from URL.
func (t Template) IsURL() bool {
	return strings.HasPrefix(t.ResolvedPath, "http://") || strings.HasPrefix(t.ResolvedPath, "https://")
}

// GetContent returns the template content.
func (t Template) GetContent() string {
	return t.Content
}
