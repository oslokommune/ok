package format_test

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                        string
		packageManifestFilename         string
		expectedPackageManifestFilename string
	}{
		{
			testName:                        "Should format packages.yml",
			packageManifestFilename:         "packages.yml",
			expectedPackageManifestFilename: "packagesFormatted.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			// Given
			packagesFilename, err := common.GetTestdataFilepath(tc.packageManifestFilename)
			require.Nil(t, err)

			// When
			err = format.Run(packagesFilename)
			require.NoError(t, err)

			// Then
			assertEqualFileContents(t, tc.expectedPackageManifestFilename, tc.packageManifestFilename)
		})
	}
}

func assertEqualFileContents(t *testing.T, expectedFile string, actualFile string) {
	expectedTestdataFilepath, err := common.GetTestdataFilepath(expectedFile)
	require.Nil(t, err)

	expected, err := os.ReadFile(expectedTestdataFilepath)
	assert.NoError(t, err)

	actualTestdataFilepath, err := common.GetTestdataFilepath(actualFile)
	require.Nil(t, err)

	actual, err := os.ReadFile(actualTestdataFilepath)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}
