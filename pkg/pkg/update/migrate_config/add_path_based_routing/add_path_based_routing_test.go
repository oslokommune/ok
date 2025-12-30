package add_path_based_routing

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestMigrateToNewConfigStructure(t *testing.T) {
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
			name:         "Basic transformation",
			inputFile:    "app-too-tikki.yml",
			expectedFile: "app-too-tikki-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse(RequiredVersion)},
		},
		{
			name:         "Don't add ApexDomain if it's not set",
			inputFile:    "app-too-tikki-no-apex-domain.yml",
			expectedFile: "app-too-tikki-no-apex-domain-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse(RequiredVersion)},
		},
		{
			name:         "Enabled must be same as in original",
			inputFile:    "app-too-tikki-alb-disabled.yml",
			expectedFile: "app-too-tikki-alb-disabled-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse(RequiredVersion)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define paths for input and expected output files
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)
			expectedFile := filepath.Join(cwd, "testdata", tc.expectedFile)

			tempDir, err := os.MkdirTemp("", "test_"+cleanTestName(tc.name))
			assert.NoError(t, err)
			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					require.NoError(t, err)
				}
			}(tempDir)

			tempInputFile := filepath.Join(tempDir, tc.inputFile)
			err = copyFile(inputFile, tempInputFile)
			require.NoError(t, err)

			// When
			err = AddPathBasedRouting(tempInputFile, tc.jsonSchema)
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
		name                  string
		inputFile             string
		expectedIsTransformed bool
	}{
		{
			name:                  "File with Subdomain and Apex set",
			inputFile:             "app-too-tikki-expected.yml",
			expectedIsTransformed: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)

			// When
			result, err := isMigrated(inputFile)
			assert.NoError(t, err)

			// Then
			assert.Equal(t, tc.expectedIsTransformed, result)
		})
	}
}

// cleanTestName keeps only alphanumeric characters in a string
func cleanTestName(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	return reg.ReplaceAllString(name, "")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
