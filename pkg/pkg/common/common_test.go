package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		configName string
		expected   string
	}{
		{
			name:       "no prefix",
			prefix:     "",
			configName: "config",
			expected:   "config.yml",
		},
		{
			name:       "with prefix",
			prefix:     "prefix",
			configName: "config",
			expected:   "prefix/config.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VarFile(tt.prefix, tt.configName)
			require.Equal(t, tt.expected, result)
		})
	}
}
