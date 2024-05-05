package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsNonASCII(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "String contains non-ASCII characters",
			input:    "Hello, 世界",
			expected: true,
		},
		{
			name:     "String does not contain non-ASCII characters",
			input:    "Hello, world",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := containsNonASCII(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
