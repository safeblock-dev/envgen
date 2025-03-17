package template_funcs

import (
	"time"
)

// FormatTime formats a time.Time value according to the specified layout.
// Uses the standard Go time formatting layout.
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// DefaultValue returns the first non-empty value from the provided list.
// Empty values are considered to be nil or empty strings.
// If all values are empty, returns the zero value.
func DefaultValue(values ...any) any {
	for _, v := range values {
		if !IsEmpty(v) {
			return v
		}
	}

	return nil
}

// Coalesce returns the first non-nil pointer value from a list of values.
// If all values are nil, returns nil.
func Coalesce(values ...*any) *any {
	for _, v := range values {
		if v != nil {
			return v
		}
	}

	return nil
}

// Ternary implements a ternary operator: condition ? trueVal : falseVal
// Returns trueVal if condition is true, falseVal otherwise.
func Ternary(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}

	return falseVal
}

// IsEmpty checks if a value is considered empty.
func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return val == 0
	case uint, uint8, uint16, uint32, uint64:
		return val == 0
	case float32, float64:
		return val == 0
	case bool:
		return !val
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return false
	}
}
