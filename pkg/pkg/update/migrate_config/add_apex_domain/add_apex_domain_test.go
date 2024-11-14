package add_apex_domain

import (
	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
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
			name:         "Basic transformation",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("9.0.0")},
		},
		{
			name:         "Values in app-hello.yml are false",
			inputFile:    "app-hello-false.yml",
			expectedFile: "app-hello-false-expected.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("9.0.0")},
		},
		{
			name:         "Should not transform when the template is not 'app'",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello.yml",
			jsonSchema:   metadata.JsonSchema{Template: "scaffold", Version: semver.MustParse("9.0.0")},
		},
		{
			name:         "Should not transform when the package version is less than 9.0.0",
			inputFile:    "app-hello.yml",
			expectedFile: "app-hello.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("8.9.0")},
		},
		{
			name:         "Should not transform when the var file already contains Subdomain",
			inputFile:    "app-hello-subdomain.yml",
			expectedFile: "app-hello-subdomain.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("9.0.0")},
		},
		{
			name:         "Should not transform if the var file already contains ApexDomain",
			inputFile:    "app-hello-apexdomain.yml",
			expectedFile: "app-hello-apexdomain.yml",
			jsonSchema:   metadata.JsonSchema{Template: "app", Version: semver.MustParse("9.0.0")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define paths for input and expected output files
			inputFile := filepath.Join(cwd, "testdata", tc.inputFile)
			expectedFile := filepath.Join(cwd, "testdata", tc.expectedFile)

			// Create a temporary copy of the input file
			tempDir, err := os.MkdirTemp("", "test_add_apex_domain_support")
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
			err = AddApexDomainSupport(tempInputFile, tc.jsonSchema)
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
			inputFile:             "app-hello-with-subdomain-and-apex.yml",
			expectedIsTransformed: true,
		},
		{
			name:                  "File with only Subdomain set",
			inputFile:             "app-hello-with-subdomain.yml",
			expectedIsTransformed: true,
		},
		{
			name:                  "File with neither Subdomain or Apex set",
			inputFile:             "app-hello.yml",
			expectedIsTransformed: false,
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

// Helper function to copy a file
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
