package pkg_test

import (
	"fmt"
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
			// Given
			rootDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)
			require.NoError(t, err)

			defer os.RemoveAll(rootDir) // clean up

			fmt.Println(rootDir)

			cmd.SetArgs(tt.args)

			// When
			err = cmd.Execute()

			// Then
			if tt.expectedError {
				assert.Error(t, err)
			}

		})
	}
}
