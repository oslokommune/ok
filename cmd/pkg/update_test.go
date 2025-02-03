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
			packageManifest: "testdata/bump-ref-field/input/packages.yml",
			configDir:       "testdata/bump-ref-field/input/config",
			releases: map[string]string{
				"app":                "v9.0.0",
				"load-balancing-alb": "v4.0.0",
				"app-common":         "v7.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-ref-field/expected/packages.yml",
			expectedConfigDir:       "testdata/bump-ref-field/expected/config",
		},
		{
			name:            "Should bump the Ref field only for semver-version package Refs",
			args:            []string{},
			packageManifest: "testdata/bump-ref-field-semver-only/input/packages.yml",
			configDir:       "testdata/bump-ref-field-semver-only/input/config",
			releases: map[string]string{
				"app": "v9.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-ref-field-semver-only/expected/packages.yml",
		},
		{
			name:            "Should bump schema version in var files",
			args:            []string{"app-hello"},
			packageManifest: "testdata/bump-schema-version/input/packages.yml",
			configDir:       "testdata/bump-schema-version/input/config",
			releases: map[string]string{
				"app": "v9.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-schema-version/expected/packages.yml",
			expectedConfigDir:       "testdata/bump-schema-version/expected/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			ghReleases := &GitHubReleasesMock{
				LatestReleases: tt.releases,
			}
			cmd := pkg.NewUpdateCommand(ghReleases)

			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)
			defer func(path string) {
				err := os.RemoveAll(path)
				require.NoError(t, err)
			}(tempDir)

			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataToTempDir(t, tt, tempDir)
			cmd.SetArgs(tt.args)

			testDir, err := os.Getwd()
			require.NoError(t, err)

			err = os.Chdir(tempDir) // Works, but disables the possibility for parallel tests.
			require.NoError(t, err)

			// When
			err = cmd.Execute()

			// Then
			if tt.expectError {
				assert.Error(t, err, "expected an error")
				return
			}

			err = os.Chdir(testDir)
			require.NoError(t, err)

			// Compare package manifest file
			actualBytes, err := os.ReadFile(filepath.Join(tempDir, "packages.yml"))
			require.NoError(t, err)
			actual := string(actualBytes)

			expectedBytes, err := os.ReadFile(tt.expectedPackageManifest)
			require.NoError(t, err)
			expected := string(expectedBytes)

			assert.Equal(t, expected, actual)

			// Compare var files:
			// Given
			// testadat/some-test/expected/config/app-hello.yml, we want to compare it to
			//                     tempDir/config/app-hello.yml
			if tt.expectedConfigDir == "" {
				return
			}

			err = filepath.Walk(tt.expectedConfigDir, func(path string, fileInfo os.FileInfo, err error) error {
				require.NoError(t, err)

				if fileInfo.IsDir() {
					return nil
				}

				actualFilename := filepath.Join(tempDir, "config", fileInfo.Name())
				actualVarFile := actualFilename

				varFileBytes, err := os.ReadFile(actualVarFile)
				require.NoError(t, err)
				varFile := string(varFileBytes)

				expectedFilename := filepath.Join(tt.expectedConfigDir, fileInfo.Name())
				expectedVarFileBytes, err := os.ReadFile(expectedFilename)
				require.NoError(t, err)
				expectedVarFile := string(expectedVarFileBytes)

				assert.Equal(t, expectedVarFile, varFile)

				return nil
			})
		})
	}
}

type TestData struct {
	name                    string
	args                    []string
	packageManifest         string
	configDir               string
	releases                map[string]string
	expectError             bool
	expectedPackageManifest string
	expectedConfigDir       string
}

type GitHubReleasesMock struct {
	LatestReleases map[string]string
}

func (g *GitHubReleasesMock) GetLatestReleases() (map[string]string, error) {
	return g.LatestReleases, nil
}

func copyTestdataToTempDir(t *testing.T, tt TestData, rootDir string) {
	if tt.packageManifest != "" {
		srcPath := tt.packageManifest
		dstPath := filepath.Join(rootDir, filepath.Base(tt.packageManifest))

		err := copyFile(srcPath, dstPath)
		require.NoError(t, err)
	}

	if tt.configDir != "" {
		srcDir := os.DirFS(tt.configDir)
		dstDir := filepath.Join(rootDir, filepath.Base(tt.configDir))

		err := os.CopyFS(dstDir, srcDir)
		require.NoError(t, err)
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
