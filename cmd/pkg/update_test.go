package pkg_test

import (
	"fmt"
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateCommand(t *testing.T) {
	tests := []TestData{
		{
			name:            "Should bump the Ref field for the specified packages",
			args:            []string{"app-hello", "load-balancing-alb-main"},
			testdataRootDir: "testdata/update/bump-ref-field",
			releases: map[string]string{
				"app":                "v9.0.0",
				"load-balancing-alb": "v4.0.0",
				"app-common":         "v7.0.0",
			},
			expectError: false,
			expectFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"config/common-config.yml",
			},
		},
		{
			name:            "Should bump the Ref field only for semver-version package Refs",
			args:            []string{},
			testdataRootDir: "testdata/update/bump-ref-field-semver-only",
			expectError:     false,
			expectFiles: []string{
				"packages.yml",
			},
			releases: map[string]string{
				"app": "v9.0.0",
			},
		},
		{
			name:            "Should bump schema version in var files",
			args:            []string{"app-hello"},
			testdataRootDir: "testdata/update/bump-schema-version",
			expectError:     false,
			expectFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"common-config.yml",
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
		},
		{
			name:            "Should migrate schema declaration from dir based to HTTPS based",
			args:            []string{"app-hello"},
			testdataRootDir: "testdata/update/migrate-schema-declaration-format",
			expectError:     false,
			expectFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"common-config.yml",
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
		},
		{
			name:            "Should update ok packages recursively",
			args:            []string{"--recursive"},
			testdataRootDir: "testdata/update/recursive",
			expectError:     false,
			expectFiles: []string{
				"app-common/packages.yml",
				"app-common/config.yml",
				"app-hello/packages.yml",
				"app-hello/config.yml",
			},
			releases: map[string]string{
				"app":        "v9.0.0",
				"app-common": "v4.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testWorkingDirectory, err := os.Getwd()
			require.NoError(t, err)

			ghReleases := &GitHubReleasesMock{
				LatestReleases: tt.releases,
			}

			command := pkg.NewUpdateCommand(ghReleases)

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
				return
			}
			require.NoError(t, err)

			err = os.Chdir(testWorkingDirectory)
			require.NoError(t, err)

			// Compare package manifest file
			for _, expectedFile := range tt.expectFiles {
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
