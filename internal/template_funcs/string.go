package template_funcs

import (
	"strings"
	"unicode"
)

// isWordBoundary checks if the current position in the string represents a word boundary.
func isWordBoundary(s string, i int, r rune) bool {
	if i == 0 {
		return false
	}

	prevRune := rune(s[i-1])

	return unicode.IsUpper(r) && !unicode.IsUpper(prevRune)
}

// isValidWordChar checks if the rune is a valid character for a word.
func isValidWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// appendWord adds a word to the result slice if it's not empty.
func appendWord(words []string, word []rune) []string {
	if len(word) > 0 {
		words = append(words, string(word))
	}

	return words
}

// SplitWords splits a string into words, supporting various delimiters.
// It handles camelCase, PascalCase, snake_case, and kebab-case.
func SplitWords(s string) []string {
	if s == "" {
		return nil
	}

	var (
		words []string
		word  []rune
	)

	for i, r := range s {
		if isValidWordChar(r) {
			if isWordBoundary(s, i, r) {
				words = appendWord(words, word)
				word = nil
			}

			word = append(word, unicode.ToLower(r))
		} else if len(word) > 0 {
			words = appendWord(words, word)
			word = nil
		}
	}

	return appendWord(words, word)
}

// Title converts the first letter of a string to uppercase.
func Title(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

// Example: "hello_world" -> "helloWorld".
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	words := SplitWords(s)
	for i := 1; i < len(words); i++ {
		words[i] = Title(words[i])
	}

	return strings.Join(words, "")
}

// Example: "hello_world" -> "HelloWorld".
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	words := SplitWords(s)
	for i := range words {
		words[i] = Title(words[i])
	}

	return strings.Join(words, "")
}

// Example: "helloWorld" -> "hello_world".
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	return strings.Join(SplitWords(s), "_")
}

// Example: "helloWorld" -> "hello-world".
func ToKebabCase(s string) string {
	if s == "" {
		return ""
	}

	return strings.Join(SplitWords(s), "-")
}

// StringAppend adds a value to a slice and returns a new slice containing all elements.
func StringAppend(slice []string, value string) []string {
	return append(slice, value)
}

// StringUniq removes duplicate values from a slice while preserving the original order of elements.
func StringUniq(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(items))

	for _, item := range items {
		if !seen[item] {
			seen[item] = true

			result = append(result, item)
		}
	}

	return result
}

// StringSlice creates a new slice from the given arguments.
func StringSlice(args ...string) []string {
	return args
}

// IsURL checks if a string is a valid URL.
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// Oneline removes all newlines from a string.
func Oneline(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "")
}
