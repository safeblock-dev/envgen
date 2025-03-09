package templatefuncs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/pkg/envgen/templatefuncs"
)

func TestSplitWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{name: "empty string", input: "", expected: nil},
		{name: "single word", input: "hello", expected: []string{"hello"}},
		{name: "multiple words with spaces", input: "hello world test", expected: []string{"hello", "world", "test"}},
		{name: "snake case", input: "hello_world_test", expected: []string{"hello", "world", "test"}},
		{name: "kebab case", input: "hello-world-test", expected: []string{"hello", "world", "test"}},
		{name: "camel case", input: "helloWorldTest", expected: []string{"hello", "world", "test"}},
		{name: "pascal case", input: "HelloWorldTest", expected: []string{"hello", "world", "test"}},
		{name: "mixed case with numbers", input: "hello123World_test-case", expected: []string{"hello123", "world", "test", "case"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.SplitWords(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single word lowercase", input: "hello", expected: "Hello"},
		{name: "single word uppercase", input: "HELLO", expected: "HELLO"},
		{name: "mixed case", input: "hELLo", expected: "HELLo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.Title(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single word", input: "hello", expected: "hello"},
		{name: "snake case", input: "hello_world", expected: "helloWorld"},
		{name: "kebab case", input: "hello-world", expected: "helloWorld"},
		{name: "pascal case", input: "HelloWorld", expected: "helloWorld"},
		{name: "with numbers", input: "hello_world123_test", expected: "helloWorld123Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.ToCamelCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single word", input: "hello", expected: "Hello"},
		{name: "snake case", input: "hello_world", expected: "HelloWorld"},
		{name: "kebab case", input: "hello-world", expected: "HelloWorld"},
		{name: "camel case", input: "helloWorld", expected: "HelloWorld"},
		{name: "with numbers", input: "hello_world123_test", expected: "HelloWorld123Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.ToPascalCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single word", input: "hello", expected: "hello"},
		{name: "pascal case", input: "HelloWorld", expected: "hello_world"},
		{name: "camel case", input: "helloWorld", expected: "hello_world"},
		{name: "kebab case", input: "hello-world", expected: "hello_world"},
		{name: "with numbers", input: "HelloWorld123Test", expected: "hello_world123_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.ToSnakeCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToKebabCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single word", input: "hello", expected: "hello"},
		{name: "pascal case", input: "HelloWorld", expected: "hello-world"},
		{name: "camel case", input: "helloWorld", expected: "hello-world"},
		{name: "snake case", input: "hello_world", expected: "hello-world"},
		{name: "with numbers", input: "HelloWorld123Test", expected: "hello-world123-test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := templatefuncs.ToKebabCase(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()

	t.Run("string values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			slice    []any
			value    any
			expected []any
		}{
			{name: "empty slice", slice: nil, value: "test", expected: []any{"test"}},
			{name: "non-empty slice", slice: []any{"hello"}, value: "world", expected: []any{"hello", "world"}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := templatefuncs.Append(tt.slice, tt.value)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("int values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			slice    []any
			value    any
			expected []any
		}{
			{name: "empty slice", slice: nil, value: 42, expected: []any{42}},
			{name: "non-empty slice", slice: []any{1}, value: 2, expected: []any{1, 2}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := templatefuncs.Append(tt.slice, tt.value)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestUniq(t *testing.T) {
	t.Parallel()

	t.Run("string values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			input    []any
			expected []any
		}{
			{name: "empty slice", input: nil, expected: nil},
			{name: "no duplicates", input: []any{"a", "b", "c"}, expected: []any{"a", "b", "c"}},
			{name: "with duplicates", input: []any{"a", "b", "a", "c", "b"}, expected: []any{"a", "b", "c"}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := templatefuncs.Uniq(tt.input)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("int values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			input    []any
			expected []any
		}{
			{name: "empty slice", input: nil, expected: nil},
			{name: "no duplicates", input: []any{1, 2, 3}, expected: []any{1, 2, 3}},
			{name: "with duplicates", input: []any{1, 2, 1, 3, 2}, expected: []any{1, 2, 3}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := templatefuncs.Uniq(tt.input)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}
