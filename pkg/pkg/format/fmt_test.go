package format_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                        string
		packageManifestFilename         string
		expectedPackageManifestFilename string
	}{
		{
			testName:                        "Should format packages.yml",
			packageManifestFilename:         "packagesBadFormatting.yml",
			expectedPackageManifestFilename: "packagesFormatted.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			// Given
			packagesFilename, err := common.GetTestdataFilepath(tc.packageManifestFilename)
			require.Nil(t, err)

			// Create temporary file copy of the packages file
			// which we can format
			tempF, err := createTemporaryPackagesFile(packagesFilename)
			require.Nil(t, err)
			defer os.Remove(tempF.Name())

			// When
			err = format.Run(tempF.Name())
			require.NoError(t, err)

			// Then
			expectedPackageFilePath, err := common.GetTestdataFilepath(tc.expectedPackageManifestFilename)
			require.Nil(t, err)
			assertEqualFileContents(t, expectedPackageFilePath, tempF.Name())
		})
	}
}

func assertEqualFileContents(t *testing.T, expectedPath string, actualPath string) {
	expected, err := os.ReadFile(expectedPath)
	assert.NoError(t, err)

	actual, err := os.ReadFile(actualPath)
	assert.NoError(t, err)

	assert.Equalf(t, string(expected), string(actual), "Expected file contents to be equal, temp file: %s (correct %s)", actualPath, expectedPath)
}

func createTemporaryPackagesFile(inputPath string) (file *os.File, err error) {
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	tempF, err := os.CreateTemp("", "packages.yml")
	if err != nil {
		return nil, err
	}
	defer func() {
		// if we fail, we should clean up the file before returning
		if err != nil {
			tempF.Close()
			os.Remove(tempF.Name())
		}
	}()

	_, err = io.Copy(tempF, bytes.NewReader(inputData))
	if err != nil {
		return nil, err
	}

	err = tempF.Close()
	if err != nil {
		return nil, err
	}

	return tempF, nil
}
