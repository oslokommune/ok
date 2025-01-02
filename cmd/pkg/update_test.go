package pkg_test

import (
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUpdateCommand(t *testing.T) {
	// Initialize the command
	cmd := pkg.UpdateCommand

	// Define test cases
	tests := []struct {
		name          string
		args          []string
		expectedError bool
	}{
		{
			name:          "Should work with no arguments",
			args:          []string{},
			expectedError: false,
		},
		{
			name:          "Should work with output folder",
			args:          []string{"out/app-common"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Chdir("testdata")
			require.NoError(t, err)

			cmd.SetArgs(tt.args)

			err = cmd.Execute()
			if tt.expectedError {
				assert.Error(t, err)
			}
		})
	}
}
