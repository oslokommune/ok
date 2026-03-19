package install

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBoilerplateVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard version",
			input:    "boilerplate version v0.12.1",
			expected: "v0.12.1",
		},
		{
			name:     "latest version",
			input:    "boilerplate version latest",
			expected: "latest",
		},
		{
			name:     "with trailing newline",
			input:    "boilerplate version v0.12.0\n",
			expected: "v0.12.0",
		},
		{
			name:     "just version string",
			input:    "v1.0.0",
			expected: "v1.0.0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseBoilerplateVersion(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
