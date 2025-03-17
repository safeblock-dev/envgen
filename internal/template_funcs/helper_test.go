package template_funcs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/safeblock-dev/envgen/internal/template_funcs"
)

func TestFormatTime(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 1, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		name     string
		time     time.Time
		layout   string
		expected string
	}{
		{name: "RFC3339", time: now, layout: time.RFC3339, expected: "2024-01-01T12:34:56Z"},
		{name: "custom format", time: now, layout: "2006-01-02 15:04:05", expected: "2024-01-01 12:34:56"},
		{name: "date only", time: now, layout: "2006-01-02", expected: "2024-01-01"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := template_funcs.FormatTime(test.time, test.layout)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestDefaultValue(t *testing.T) {
	t.Parallel()

	t.Run("string values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			values   []any
			expected any
		}{
			{name: "all empty", values: []any{"", "", ""}, expected: nil},
			{name: "first non-empty", values: []any{"value", "", ""}, expected: "value"},
			{name: "second non-empty", values: []any{"", "value", ""}, expected: "value"},
			{name: "no values", values: []any{}, expected: nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.DefaultValue(tt.values...)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("int values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			values   []any
			expected any
		}{
			{name: "all zero", values: []any{0, 0, 0}, expected: nil},
			{name: "first non-zero", values: []any{1, 0, 0}, expected: 1},
			{name: "second non-zero", values: []any{0, 2, 0}, expected: 2},
			{name: "no values", values: []any{}, expected: nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.DefaultValue(tt.values...)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestCoalesce(t *testing.T) {
	t.Parallel()

	t.Run("string values", func(t *testing.T) {
		t.Parallel()

		var nilVal *any

		emptyStr := any("")
		valueStr := any("value")
		emptyStrPtr := &emptyStr
		valueStrPtr := &valueStr

		tests := []struct {
			name     string
			values   []*any
			expected *any
		}{
			{name: "all nil", values: []*any{nilVal, nilVal, nilVal}, expected: nilVal},
			{name: "first non-nil", values: []*any{valueStrPtr, nilVal, nilVal}, expected: valueStrPtr},
			{name: "empty string is valid", values: []*any{nilVal, emptyStrPtr, valueStrPtr}, expected: emptyStrPtr},
			{name: "no values", values: []*any{}, expected: nilVal},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.Coalesce(tt.values...)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("int values", func(t *testing.T) {
		t.Parallel()

		var nilVal *any

		zeroInt := any(0)
		valueInt := any(42)
		zeroIntPtr := &zeroInt
		valueIntPtr := &valueInt

		tests := []struct {
			name     string
			values   []*any
			expected *any
		}{
			{name: "all nil", values: []*any{nilVal, nilVal, nilVal}, expected: nilVal},
			{name: "first non-nil", values: []*any{valueIntPtr, nilVal, nilVal}, expected: valueIntPtr},
			{name: "zero is valid", values: []*any{nilVal, zeroIntPtr, valueIntPtr}, expected: zeroIntPtr},
			{name: "no values", values: []*any{}, expected: nilVal},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.Coalesce(tt.values...)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestTernary(t *testing.T) {
	t.Parallel()

	t.Run("string values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			condition bool
			trueVal   any
			falseVal  any
			expected  any
		}{
			{name: "true condition", condition: true, trueVal: "yes", falseVal: "no", expected: "yes"},
			{name: "false condition", condition: false, trueVal: "yes", falseVal: "no", expected: "no"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.Ternary(tt.condition, tt.trueVal, tt.falseVal)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("int values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			condition bool
			trueVal   any
			falseVal  any
			expected  any
		}{
			{name: "true condition", condition: true, trueVal: 42, falseVal: 0, expected: 42},
			{name: "false condition", condition: false, trueVal: 42, falseVal: 0, expected: 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := template_funcs.Ternary(tt.condition, tt.trueVal, tt.falseVal)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}
