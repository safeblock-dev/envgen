package templatefuncs

import (
	"fmt"
	"strings"
)

// ToString converts any value to its string representation.
// Uses fmt.Sprintf with %v format to handle any type.
func ToString(v any) string {
	return fmt.Sprintf("%v", v)
}

// ToInt converts a string to an integer with an optional default value.
// If the conversion fails, returns the default value (0 if not specified).
func ToInt(s string, def ...int) int {
	var defaultVal int
	if len(def) > 0 {
		defaultVal = def[0]
	}

	var result int
	if _, err := fmt.Sscanf(s, "%d", &result); err != nil {
		return defaultVal
	}
	return result
}

// ToBool converts a string to a boolean with an optional default value.
// Recognizes common boolean string representations:
// - true: "true", "1", "yes", "on"
// - false: "false", "0", "no", "off"
// If the conversion fails, returns the default value (false if not specified).
func ToBool(s string, def ...bool) bool {
	var defaultVal bool
	if len(def) > 0 {
		defaultVal = def[0]
	}

	switch strings.ToLower(s) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultVal
	}
}

// Slice creates a new slice from the given arguments
func Slice(args ...string) []string {
	return args
}
