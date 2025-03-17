package user_template

import "fmt"

// Template errors.
type (
	// ErrTemplateNotFound represents a template not found error.
	TemplateNotFoundError struct {
		Name string
	}

	// ErrInvalidTemplate represents an invalid template error.
	InvalidTemplateError struct {
		Path    string
		Message string
	}
)

func (e *TemplateNotFoundError) Error() string {
	return "template not found: " + e.Name
}

func (e *InvalidTemplateError) Error() string {
	return fmt.Sprintf("invalid template: %s: %s", e.Path, e.Message)
}
