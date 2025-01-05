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
			name:            "Should bump the Ref for some packages",
			args:            []string{"app-hello", "load-balancing-alb-main"},
			packageManifest: "testdata/input/packages.yml",
			configDir:       "testdata/input/config",
			expectError:     false,
			releases: map[string]string{
				"app":                "v8.0.0",
				"load-balancing-alb": "v4.0.0",
			},
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

			// Compare package file
			actualBytes, err := os.ReadFile(filepath.Join(tempDir, "packages.yml"))
			require.NoError(t, err)
			actual := string(actualBytes)

			expectedBytes, err := os.ReadFile(filepath.Join("testdata", "expected", "packages.yml"))
			require.NoError(t, err)
			expected := string(expectedBytes)

			assert.Equal(t, expected, actual)

			// Compare config files
			err = filepath.Walk(filepath.Join("testdata", "expected", "config"), func(testdataPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				expectedBytes, err := os.ReadFile(testdataPath)
				require.NoError(t, err)
				expected := string(expectedBytes)

				actualBytes, err := os.ReadFile(filepath.Join(tempDir, "config", filepath.Base(testdataPath)))
				require.NoError(t, err)
				actual := string(actualBytes)

				assert.Equal(t, expected, actual)

				return nil
			})

			require.NoError(t, err)
		})
	}
}

type TestData struct {
	name            string
	args            []string
	packageManifest string
	configDir       string
	expectError     bool
	releases        map[string]string
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
