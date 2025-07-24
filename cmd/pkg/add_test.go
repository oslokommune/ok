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
			name:            "Should add package",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/standard-case",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
		},
		{
			name:            "Should fail if output directory already exists, using default dir",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/dir-already-exists",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
			expectError:        true,
			expectErrorMessage: "folder already exists: databases",
		},
		{
			name:            "Should add package with the old package manifest structure",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/old-structure",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"packages.yml",
				"_config/databases.yml",
			},
		},
		{
			name:            "Should add package with the old package manifest structure with custom name",
			args:            []string{"app", "app-hello"},
			testdataRootDir: "testdata/add/old-structure-custom-stack-name",
			releases: map[string]string{
				"app": "v6.0.0",
			},
			expectedFiles: []string{
				"packages.yml",
				"_config/app-hello.yml",
			},
		},
		{
			name:            "Should fail if output directory already exists, using dir from argument",
			args:            []string{"app", "app-hello"},
			testdataRootDir: "testdata/add/dir-already-exists",
			expectedFiles: []string{
				"app/packages.yml",
				"app/package-config.yml",
			},
			expectError:        true,
			expectErrorMessage: "folder already exists: app-hello",
		},
		{
			name:            "Should add package with specified var file",
			args:            []string{"databases", "--var-file", "non-serverless"},
			testdataRootDir: "testdata/add/specified-var-file",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectedFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
		},
		{
			name:            "Should show error if var file does not exist",
			args:            []string{"databases", "--var-file", "some-missing-var-file"},
			testdataRootDir: "testdata/add/var-file-missing",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectError:        true,
			expectErrorMessage: "package-config-some-missing-var-file.yml: no such file or directory",
			keepTempDir:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testWorkingDirectory, err := os.Getwd()
			require.NoError(t, err)

			ghReleases := &GitHubReleasesMock{
				LatestReleases:            tt.releases,
				TestWorkingDirectory:      testWorkingDirectory,
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
			copyTestdataRootDirToTempDir(t, tt, testWorkingDirectory, tempDir)
			command.SetArgs(tt.args)

			err = os.Chdir(tempDir) // Works, but disables the possibility for parallel tests.
			require.NoError(t, err)
			defer func() {
				err = os.Chdir(testWorkingDirectory)
			}()

			// When
			err = command.Execute()

			// Then
			if tt.expectError {
				assert.Error(t, err, "expected an error")
				assert.Contains(t, err.Error(), tt.expectErrorMessage)

				return
			}
			require.NoError(t, err)

			err = os.Chdir(testWorkingDirectory)
			require.NoError(t, err)

			// Compare package manifest file
			for _, expectedFile := range tt.expectedFiles {
				actualBytes, err := os.ReadFile(filepath.Join(tempDir, expectedFile))
				require.NoError(t, err)
				actual := string(actualBytes)

				expectedFileFullPath := filepath.Join(tt.testdataRootDir, "expected", expectedFile)
				expectedBytes, err := os.ReadFile(expectedFileFullPath)
				require.NoError(t, err)
				expected := string(expectedBytes)

				assert.Equal(t, expected, actual)
			}
		})
	}
}
