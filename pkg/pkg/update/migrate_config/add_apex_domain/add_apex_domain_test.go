package add_apex_domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"github.com/stretchr/testify/assert"
)

func TestAddApexDomainSupport(t *testing.T) {
	// Get the current working directory
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		inputFile    string
		expectedFile string
		metadata     migrate_config.VarFileMetadata
	}{
		{
			name:         "Should not run on ",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello-expected.yml",
			metadata:     migrate_config.VarFileMetadata{HasVersion: true, Template: "app", Version: "8.0.0"},
		},
		{
			name:         "Basic transformation",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello-expected.yml",
			metadata:     migrate_config.VarFileMetadata{HasVersion: true, Template: "app", Version: "8.0.0"},
		},
		{
			name:         "Values in app-hello.yml are false",
			inputFile:    "app-hello-false.yml",
			expectedFile: "app-hello-false-expected.yml",
			metadata:     migrate_config.VarFileMetadata{HasVersion: true, Template: "app", Version: "8.0.0"},
		},
		// Add more test cases here if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define paths for input and expected output files
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)
			expectedFile := filepath.Join(cwd, "testdata", tc.expectedFile)

			// Create a temporary copy of the input file
			tempDir, err := os.MkdirTemp("", "test_add_apex_domain_support")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			tempInputFile := filepath.Join(tempDir, tc.inputFile)
			err = copyFile(inputFile, tempInputFile)
			assert.NoError(t, err)

			// When
			err = AddApexDomainSupport(tempInputFile, tc.metadata)
			assert.NoError(t, err)

			// Then
			modifiedContent, err := os.ReadFile(tempInputFile)
			assert.NoError(t, err)

			expectedContent, err := os.ReadFile(expectedFile)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedContent), string(modifiedContent))
		})
	}
}

// Helper function to copy a file
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
