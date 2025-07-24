package pkg_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCommand(t *testing.T) {
	tests := []TestData{
		{
			name:            "Should add new package with the old package manifest structure",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/old-structure",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"packages.yml",
				"_config/databases.yml",
			},
			keepTempDir: true,
		},
		{
			name:            "Should add new package with custom name with the old package manifest structure",
			args:            []string{"app", "app-hello"},
			testdataRootDir: "testdata/add/old-structure",
			releases: map[string]string{
				"app": "v6.0.0",
			},
			expectedFiles: []string{
				"packages.yml",
				"_config/app-hello.yml",
			},
			keepTempDir: true,
		},
		{
			name:            "Should add new package",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/standard-case",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
			keepTempDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testDataDir, err := os.Getwd()
			require.NoError(t, err)

			ghReleases := &GitHubReleasesMock{
				LatestReleases:            tt.releases,
				TestWorkingDirectory:      testDataDir,
				BoilerplateRepositoryPath: tt.testdataRootDir,
			}

			command := pkg.NewAddCommand(ghReleases)

			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)

			// Remove temp dir after test run
			if !tt.keepTempDir {
				defer func(path string) {
					err := os.RemoveAll(path)
					require.NoError(t, err)
				}(tempDir)
			}

			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataRootDirToTempDir(t, tt, tempDir)
			command.SetArgs(tt.args)

			err = os.Chdir(tempDir) // Works, but disables the possibility for parallel tests.
			require.NoError(t, err)

			// When
			err = command.Execute()

			// Then
			if tt.expectError {
				assert.Error(t, err, "expected an error")
				return
			}
			require.NoError(t, err)

			err = os.Chdir(testDataDir)
			require.NoError(t, err)

			// Compare package manifest file
			for _, expectedFile := range tt.expectedFiles {
				actualBytes, err := os.ReadFile(filepath.Join(tempDir, expectedFile))
				require.NoError(t, err)
				actual := string(actualBytes)

				lol := filepath.Join(tt.testdataRootDir, "expected", expectedFile)
				expectedBytes, err := os.ReadFile(lol)
				require.NoError(t, err)
				expected := string(expectedBytes)

				assert.Equal(t, expected, actual)
			}
		})
	}
}
