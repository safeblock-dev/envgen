package templatefuncs

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToString(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		defaultVal []int
		expected   int
	}{
		{
			name:     "valid integer",
			input:    "123",
			expected: 123,
		},
		{
			name:     "invalid integer",
			input:    "abc",
			expected: 0,
		},
		{
			name:       "invalid integer with default",
			input:      "abc",
			defaultVal: []int{42},
			expected:   42,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToInt(tt.input, tt.defaultVal...)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToBool(tt.input, tt.defaultVal...)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"test"},
			expected: []string{"test"},
		},
		{
			name:     "multiple elements",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Slice(tt.input...)
			require.Equal(t, tt.expected, result)
		})
	}
}
