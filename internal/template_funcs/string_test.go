package template_funcs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/template_funcs"
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
		{name: "unicode letters", input: "привет_мир", expected: []string{"привет", "мир"}},
		{name: "unicode mixed case", input: "ПриветМир", expected: []string{"привет", "мир"}},
		{name: "unicode with latin", input: "helloПривет_worldМир", expected: []string{"hello", "привет", "world", "мир"}},
		{name: "unicode numbers", input: "тест123World", expected: []string{"тест123", "world"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.SplitWords(test.input)
			require.Equal(t, test.expected, result)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.Title(test.input)
			require.Equal(t, test.expected, result)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToCamelCase(test.input)
			require.Equal(t, test.expected, result)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToPascalCase(test.input)
			require.Equal(t, test.expected, result)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToSnakeCase(test.input)
			require.Equal(t, test.expected, result)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToKebabCase(test.input)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestStringUniq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{name: "empty slice", input: nil, expected: nil},
		{name: "no duplicates", input: []string{"a", "b", "c"}, expected: []string{"a", "b", "c"}},
		{name: "with duplicates", input: []string{"a", "b", "a", "c", "b"}, expected: []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.StringUniq(test.input)
			require.Equal(t, test.expected, result)
		})
	}
}
