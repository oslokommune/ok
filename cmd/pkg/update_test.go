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
	args            []string // []string{"out/app-common"},
	packageManifest string
	configDir       string
	expectError     bool
}

func TestUpdateCommand(t *testing.T) {
	cmd := pkg.UpdateCommand

	tests := []TestData{
		{
			name:            "Should work with no arguments",
			args:            []string{},
			packageManifest: "packages.yml",
			configDir:       "config",
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
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
