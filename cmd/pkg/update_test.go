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

type UpdateTestData struct {
	TestData

	releases map[string]string
}

func TestUpdateCommand(t *testing.T) {
	tests := []UpdateTestData{
		{
			TestData: TestData{
				name:            "Should bump the Ref field for the specified packages",
				args:            []string{"app-hello", "load-balancing-alb-main"},
				testdataRootDir: "testdata/update/bump-ref-field",
				expectError:     false,
				expectedFiles: []string{
					"packages.yml",
					"config/app-hello.yml",
					"config/common-config.yml",
				},
			},
			releases: map[string]string{
				"app":                "v9.0.0",
				"load-balancing-alb": "v4.0.0",
				"app-common":         "v7.0.0",
			},
		},
		{
			TestData: TestData{
				name:            "Should bump the Ref field only for semver-version package Refs",
				args:            []string{},
				testdataRootDir: "testdata/update/bump-ref-field-semver-only",
				expectError:     false,
				expectedFiles: []string{
					"packages.yml",
				}},
			releases: map[string]string{
				"app": "v9.0.0",
			},
		},
		{
			TestData: TestData{
				name:            "Should bump schema version in var files",
				args:            []string{"app-hello"},
				testdataRootDir: "testdata/update/bump-schema-version",
				expectError:     false,
				expectedFiles: []string{
					"packages.yml",
					"config/app-hello.yml",
					"common-config.yml",
				},
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
		},
		{
			TestData: TestData{
				name:            "Should migrate schema declaration from dir based to HTTPS based",
				args:            []string{"app-hello"},
				testdataRootDir: "testdata/update/migrate-schema-declaration-format",
				expectError:     false,
				expectedFiles: []string{
					"packages.yml",
					"config/app-hello.yml",
					"common-config.yml",
				},
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
		},
		{
			TestData: TestData{
				name:            "Should update ok packages recursively",
				args:            []string{"--recursive"},
				testdataRootDir: "testdata/update/recursive",
				expectError:     false,
				expectedFiles: []string{
					"app-common/packages.yml",
					"app-common/config.yml",
					"app-hello/packages.yml",
					"app-hello/config.yml",
				},
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
			testDir, err := os.Getwd()
			require.NoError(t, err)

			ghReleases := &GitHubReleasesMock{
				LatestReleases: tt.releases,
			}

			command := pkg.NewUpdateCommand(ghReleases)

			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)

			// Remove temp dir after test run
			defer func(path string) {
				err := os.RemoveAll(path)
				require.NoError(t, err)
			}(tempDir)

			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataRootDirToTempDir(t, tt.TestData, tempDir)
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

			err = os.Chdir(testDir)
			require.NoError(t, err)

			// Compare package manifest file
			for _, file := range tt.expectedFiles {
				actualBytes, err := os.ReadFile(filepath.Join(tempDir, file))
				require.NoError(t, err)
				actual := string(actualBytes)

				expectedBytes, err := os.ReadFile(filepath.Join(tt.testdataRootDir, "expected", file))
				require.NoError(t, err)
				expected := string(expectedBytes)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

type GitHubReleasesMock struct {
	LatestReleases map[string]string
}

func (g *GitHubReleasesMock) GetLatestReleases() (map[string]string, error) {
	return g.LatestReleases, nil
}
