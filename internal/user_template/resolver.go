package user_template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/safeblock-dev/envgen/internal/github"
)

const (
	githubOwner = "safeblock-dev"
	githubRepo  = "envgen"
)

// Resolver handles template resolution from different sources.
type Resolver struct {
	github *github.Client
}

// NewResolver creates a new template resolver.
func NewResolver() (*Resolver, error) {
	return NewCustomResolver(githubOwner, githubRepo)
}

// NewCustomResolver creates a new template resolver.
func NewCustomResolver(githubOwner, githubRepo string) (*Resolver, error) {
	githubClient, err := github.New(githubOwner, githubRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return &Resolver{
		github: githubClient,
	}, nil
}

// Template resolves and loads a template from the given path.
func (r *Resolver) Template(ctx context.Context, path string) (*Template, error) {
	var tmpl Template

	if path == "" {
		return nil, &InvalidTemplateError{
			Path:    path,
			Message: "template path is empty",
		}
	}

	source := r.detectTemplateSource(ctx, path)

	var err error

	switch source {
	case TemplateSourceURL:
		tmpl, err = r.resolveURLTemplate(ctx, path)
	case TemplateSourceStandard:
		tmpl, err = r.resolveStandardTemplate(ctx, path)
	case TemplateSourceLocal:
		tmpl, err = r.resolveLocalTemplate(path)
	}

	if err != nil {
		return nil, err
	}

	if err := tmpl.Validate(); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// detectTemplateSource determines the source of a template based on its path.
func (r *Resolver) detectTemplateSource(ctx context.Context, path string) TemplateSource {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return TemplateSourceURL
	}

	// Check if it's a standard template by trying to find it in the repository
	templates, err := r.ListTemplates(ctx)
	if err == nil {
		for _, t := range templates {
			if t.Name == path {
				return TemplateSourceStandard
			}
		}
	}

	return TemplateSourceLocal
}

// resolveURLTemplate loads a template from a URL.
func (r *Resolver) resolveURLTemplate(ctx context.Context, url string) (Template, error) {
	content, err := r.github.GetFileFromURL(ctx, url)
	if err != nil {
		return Template{}, fmt.Errorf("failed to fetch template from URL: %w", err)
	}

	return Template{
		Name:         filepath.Base(url),
		Source:       TemplateSourceURL,
		Content:      content,
		ResolvedPath: url,
	}, nil
}

// resolveStandardTemplate loads a standard template from the repository.
func (r *Resolver) resolveStandardTemplate(ctx context.Context, name string) (Template, error) {
	content, err := r.github.GetFile(ctx, "templates/"+name)
	if err != nil {
		return Template{}, fmt.Errorf("failed to fetch standard template: %w", err)
	}

	url, err := r.github.GetStandardTemplateURL(ctx, name)
	if err != nil {
		return Template{}, fmt.Errorf("failed to get template URL: %w", err)
	}

	return Template{
		Name:         name,
		Source:       TemplateSourceStandard,
		Content:      content,
		ResolvedPath: url,
	}, nil
}

// resolveLocalTemplate loads a template from the local filesystem.
func (r *Resolver) resolveLocalTemplate(path string) (Template, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Template{}, fmt.Errorf("failed to read local template: %w", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return Template{
		Name:         filepath.Base(path),
		Source:       TemplateSourceLocal,
		Content:      string(content),
		ResolvedPath: absPath,
	}, nil
}

// ListTemplates returns a list of available templates from the repository.
func (r *Resolver) ListTemplates(ctx context.Context) ([]github.Content, error) {
	return r.github.GetList(ctx, "templates")
}

// ListAvailableTemplateNames returns a list of available template names.
func (r *Resolver) ListAvailableTemplateNames(ctx context.Context) ([]string, error) {
	templates, err := r.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var names []string

	for _, t := range templates {
		if t.Type == "file" {
			names = append(names, t.Name)
		}
	}

	return names, nil
}
