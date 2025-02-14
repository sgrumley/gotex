package testdata

import (
	"testing"
)

// TestAddNumbers is an example of test cases represented as anonymous structs
func TestAddNumbers(t *testing.T) {
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

type tc struct {
	name     string
	a        int
	b        int
	expected int
}

// TestAddTyped is an example of test cases represented as typed structs
func TestAddNumbersTyped(t *testing.T) {
	tcs := []tc{
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

// TestAddNumbersFail is an example of a failing test
func TestAddNumbersFail(t *testing.T) {
	tcs := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "failed",
			a:        1,
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

// TestAddNumbersImplicit is an example of test cases represented as anonymous structs implicitly initialised
// NOTE: This is currently not supported and only the function will show up in gotex
func TestAddNumbersImplicit(t *testing.T) {
	// Table of test cases
	tcs := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 1, 2, 3},
		{"negative numbers", -1, -1, -2},
		{"mixed numbers", -1, 2, 1},
		{"zero values", 0, 0, 0},
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

// TestAddNumbers is an example of test cases represented as a map
func TestAddNumbersMapCase(t *testing.T) {
	tcs := map[string]struct {
		a        int
		b        int
		expected int
	}{
		"positive numbers": {
			a:        1,
			b:        2,
			expected: 3,
		},
		"negative numbers": {
			a:        -1,
			b:        -1,
			expected: -2,
		},
		"mixed numbers": {
			a:        -1,
			b:        2,
			expected: 1,
		},
		"zero values": {
			a:        0,
			b:        0,
			expected: 0,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			result := AddNumbers(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("expected %d, but got %d", tc.expected, result)
			}
		})
	}
}

// BUG: each call to t.Run() is treated as function and shows up twice in functions
func TestSingle(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		a := 1
		b := 1
		expected := 2

		result := AddNumbers(a, b)
		if result != expected {
			t.Errorf("expected %d, but got %d", tc.expected, result)
		}
	})

	t.Run("single repeated", func(t *testing.T) {
		a := 1
		b := 1
		expected := 2

		result := AddNumbers(a, b)
		if result != expected {
			t.Errorf("expected %d, but got %d", tc.expected, result)
		}
	})
}
