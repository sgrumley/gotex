package testdata

import (
	"testing"
)

func TestSubtractNumbers(t *testing.T) {
	// Table of test cases
	tcs := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        1,
			b:        2,
			expected: 3,
		},
		{
			name:     "negative numbers",
			a:        -1,
			b:        -1,
			expected: -2,
		},
		{
			name:     "mixed numbers",
			a:        -1,
			b:        2,
			expected: 1,
		},
		{
			name:     "zero values",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "alternate case",
			a:        0,
			b:        0,
			expected: 0,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := AddNumbers(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("expected %d, but got %d", tc.expected, result)
			}
		})
	}
}
