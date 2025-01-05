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

func TestUpdateCommand(t *testing.T) {
	tests := []TestData{
		{
			name:            "Should work with no arguments",
			args:            []string{"app-hello", "load-balancing-alb-main"},
			packageManifest: "packages.yml",
			configDir:       "config",
			expectError:     false,
			releases: map[string]string{
				"app":                "v8.0.0",
				"load-balancing-alb": "v4.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up dependencies
			ghReleases := &GitHubReleasesMock{
				LatestReleases: tt.releases,
			}
			cmd := pkg.NewUpdateCommand(ghReleases)

			// More setup code
			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)
			defer os.RemoveAll(tempDir)
			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataToTempDir(t, tt, tempDir)
			cmd.SetArgs(tt.args)

			err = os.Chdir(tempDir)
			require.NoError(t, err)

			// When
			err = cmd.Execute()

			// Then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func copyTestdataToTempDir(t *testing.T, tt TestData, rootDir string) {
	if tt.packageManifest != "" {
		srcPath := filepath.Join("testdata", tt.packageManifest)
		dstPath := filepath.Join(rootDir, tt.packageManifest)

		err := copyFile(srcPath, dstPath)
		require.NoError(t, err)
	}

	if tt.configDir != "" {
		srcDir := os.DirFS(filepath.Join("testdata", tt.configDir))
		dstDir := filepath.Join(rootDir, tt.configDir)

		err := os.CopyFS(dstDir, srcDir)
		require.NoError(t, err)
	}
}
