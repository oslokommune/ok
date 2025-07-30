package pkg_test

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
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
			name:            "Should add Terraform package",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/databases",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
		},
		{
			name:            "Should add package in output folder",
			args:            []string{"app", "app-hello"},
			testdataRootDir: "testdata/add/app-hello",
			releases: map[string]string{
				"app": "v6.0.0",
			},
			expectFiles: []string{
				"app-hello/packages.yml",
				"app-hello/package-config.yml",
			},
		},
		{
			name:                        "Should add GitHub actions package",
			args:                        []string{"docker-build-push"},
			testdataRootDir:             "testdata/add/docker-build-push",
			workingDirectoryFromRootDir: "workflows/_config/dev",
			releases: map[string]string{
				"docker-build-push": "v2.3.2",
			},
			expectFiles: []string{
				"workflows/_config/dev/packages.yml",
				"workflows/_config/dev/docker-build-push.yml",
			},
			keepTempDir: true,
		},
		{
			name:                        "Should add GitHub actions package with named var file",
			args:                        []string{"docker-build-push", "my-app-docker-build-push"},
			testdataRootDir:             "testdata/add/docker-build-push-named",
			workingDirectoryFromRootDir: "workflows/_config/dev",
			releases: map[string]string{
				"docker-build-push": "v2.3.2",
			},
			expectFiles: []string{
				"workflows/_config/dev/packages.yml",
				"workflows/_config/dev/my-app-docker-build-push.yml",
			},
		},
		{
			name:            "Should add package in output folder, without data folder",
			args:            []string{"databases", "my-db"},
			testdataRootDir: "testdata/add/my-db",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectFiles: []string{
				"my-db/packages.yml",
				"my-db/package-config.yml",
			},
		},
		{
			name:            "Should fail if output directory already exists, using default dir",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/dir-already-exists",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectFiles: []string{
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
			expectFiles: []string{
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
			expectFiles: []string{
				"packages.yml",
				"_config/app-hello.yml",
			},
		},
		{
			name:            "Should fail if output directory already exists, using dir from argument",
			args:            []string{"app", "app-hello"},
			testdataRootDir: "testdata/add/dir-already-exists",
			expectFiles: []string{
				"app/packages.yml",
				"app/package-config.yml",
			},
			expectError:        true,
			expectErrorMessage: "folder already exists: app-hello",
		},
		{
			name:            "Should add package with specified var file",
			args:            []string{"databases", "--var-file", "non-serverless"},
			testdataRootDir: "testdata/add/var-file-specified",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectFiles: []string{
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
		},
		{
			name:            "Should add package without var file",
			args:            []string{"databases", "--no-var-file"},
			testdataRootDir: "testdata/add/var-file-disabled",
			releases: map[string]string{
				"databases": "v4.0.0",
			},
			expectFiles: []string{
				"databases/packages.yml",
			},
			expectNoFiles: []string{
				"databases/package-config.yml",
			},
		},
		{
			name:            "Should use base URI",
			args:            []string{"databases"},
			testdataRootDir: "testdata/add/base-url",
			baseUrl:         "boilerplate-repo",
			expectFiles: []string{
				"databases/packages.yml",
				"databases/package-config.yml",
			},
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

			if len(tt.baseUrl) > 0 {
				err = os.Setenv(common.BaseUrlEnvName, tt.baseUrl)
				require.NoError(t, err)
			}

			var workingDirectory string
			if len(tt.workingDirectoryFromRootDir) == 0 {
				workingDirectory = tempDir
			} else {
				workingDirectory = filepath.Join(tempDir, tt.workingDirectoryFromRootDir)
			}
			err = os.Chdir(workingDirectory) // Works, but disables the possibility for parallel tests.
			require.NoError(t, err)
			defer func() {
				err = os.Chdir(testWorkingDirectory)
			}()
			fmt.Println("workingDirectory: ", workingDirectory)

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

			for _, expectFile := range tt.expectFiles {
				actualBytes, err := os.ReadFile(filepath.Join(tempDir, expectFile))
				require.NoError(t, err)
				actual := string(actualBytes)

				expectFileFullPath := filepath.Join(tt.testdataRootDir, "expected", expectFile)
				expectBytes, err := os.ReadFile(expectFileFullPath)
				require.NoError(t, err)
				expected := string(expectBytes)

				assert.Equal(t, expected, actual)
			}

			for _, expectNoFile := range tt.expectNoFiles {
				expectNoFilePath := filepath.Join(tempDir, expectNoFile)

				_, err := os.Stat(expectNoFilePath)
				assert.True(
					t,
					os.IsNotExist(err),
					"expected file %s to NOT exist, but it does", expectNoFilePath,
				)
			}
		})
	}
}
