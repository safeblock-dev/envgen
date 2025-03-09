package templatefuncs

import (
	"strings"
	"unicode"
)

// SplitWords splits a string into words, supporting various delimiters.
// It handles camelCase, PascalCase, snake_case, and kebab-case.
func SplitWords(s string) []string { //nolint:cyclop
	if s == "" {
		return nil
	}

	var words []string

	var word []rune

	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(rune(s[i-1])) {
				if len(word) > 0 {
					words = append(words, string(word))
					word = nil
				}
			}

			word = append(word, unicode.ToLower(r))
		} else if len(word) > 0 {
			words = append(words, string(word))
			word = nil
		}
	}

	if len(word) > 0 {
		words = append(words, string(word))
	}

	return words
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

// Append adds a value to a slice and returns a new slice containing all elements.
func Append(slice []any, value any) []any {
	return append(slice, value)
}

// Uniq removes duplicate values from a slice while preserving the original order of elements.
func Uniq(items []any) []any {
	if len(items) == 0 {
		return nil
	}

	seen := make(map[any]bool)
	result := make([]any, 0, len(items))

	for _, item := range items {
		if !seen[item] {
			seen[item] = true

			result = append(result, item)
		}
	}

	return result
}
