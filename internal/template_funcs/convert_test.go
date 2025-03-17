package template_funcs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/template_funcs"
)

func TestToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer value",
			input:    123,
			expected: "123",
		},
		{
			name:     "boolean value",
			input:    true,
			expected: "true",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		def      []int
		expected int
	}{
		{name: "empty string", input: "", def: nil, expected: 0},
		{name: "empty string with default", input: "", def: []int{42}, expected: 42},
		{name: "positive number", input: "123", def: nil, expected: 123},
		{name: "negative number", input: "-123", def: nil, expected: -123},
		{name: "zero", input: "0", def: nil, expected: 0},
		{name: "invalid input", input: "abc", def: nil, expected: 0},
		{name: "invalid input with default", input: "abc", def: []int{42}, expected: 42},
		{name: "max int", input: "2147483647", def: nil, expected: 2147483647},
		{name: "min int", input: "-2147483648", def: nil, expected: -2147483648},
		{name: "with whitespace", input: " 123 ", def: nil, expected: 0},
		{name: "with plus sign", input: "+123", def: nil, expected: 123},
		{name: "decimal number", input: "123.45", def: nil, expected: 0},
		{name: "decimal number with default", input: "123.45", def: []int{42}, expected: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var result int
			if len(tt.def) > 0 {
				result = template_funcs.ToInt(tt.input, tt.def[0])
			} else {
				result = template_funcs.ToInt(tt.input)
			}

			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		defaultVal []bool
		expected   bool
	}{
		{
			name:     "true value",
			input:    "true",
			expected: true,
		},
		{
			name:     "false value",
			input:    "false",
			expected: false,
		},
		{
			name:     "1 value",
			input:    "1",
			expected: true,
		},
		{
			name:     "0 value",
			input:    "0",
			expected: false,
		},
		{
			name:     "yes value",
			input:    "yes",
			expected: true,
		},
		{
			name:     "no value",
			input:    "no",
			expected: false,
		},
		{
			name:     "on value",
			input:    "on",
			expected: true,
		},
		{
			name:     "off value",
			input:    "off",
			expected: false,
		},
		{
			name:     "invalid value",
			input:    "invalid",
			expected: false,
		},
		{
			name:       "invalid value with default",
			input:      "invalid",
			defaultVal: []bool{true},
			expected:   true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.ToBool(tt.input, tt.defaultVal...)
			require.Equal(t, tt.expected, result)
		})
	}
}
