package use_schema_uri

import (
	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/testhelper"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

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
		jsonSchema   metadata.JsonSchema
	}{
		{
			name:         "Should migrate",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("9.0.0")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define paths for input and expected output files
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)
			expectedFile := filepath.Join(cwd, "testdata", tc.expectedFile)

			// Create a temporary copy of the input file
			tempDir, err := os.MkdirTemp("", testhelper.TestNameToDir(tc.name))
			assert.NoError(t, err)
			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					require.NoError(t, err)
				}
			}(tempDir)

			tempInputFile := filepath.Join(tempDir, tc.inputFile)
			err = copyFile(inputFile, tempInputFile)
			assert.NoError(t, err)

			// When
			err = ReplaceDirWithUri(tempInputFile, tc.jsonSchema)
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

func TestIsTransformed(t *testing.T) {
	// Get the current working directory
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	testCases := []struct {
		name                string
		inputFile           string
		expectShouldMigrate bool
	}{
		{
			name:                "Should migrate file with dir based schema declaration",
			inputFile:           "app-hello.yml",
			expectShouldMigrate: true,
		},
		{
			name:                "Should not migrate file with URL based schema declaration",
			inputFile:           "app-hello-expected.yml",
			expectShouldMigrate: false,
		},
		{
			name:                "Should not migrate file with missing schema declaration",
			inputFile:           "app-hello-expected.yml",
			expectShouldMigrate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)

			// When
			result, err := shouldMigrate(inputFile)
			assert.NoError(t, err)

			// Then
			assert.Equal(t, tc.expectShouldMigrate, result)
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
