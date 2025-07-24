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

func TestInstallCommand(t *testing.T) {
	tests := []TestData{
		{
			name:            "Should install ok packages recursively",
			args:            []string{"--recursive"},
			testdataRootDir: "testdata/install/recursive",
			expectedFiles: []string{
				"app-hello/.boilerplate/_template_app.json",
				"networking/.boilerplate/_template_networking.json",
			},
			keepTempDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testDir, err := os.Getwd()
			require.NoError(t, err)

			command := pkg.NewInstallCommand() // figure out which parameters to pass here, if any

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

			err = os.Setenv("BASE_URL", "../boilerplate-repo")
			require.NoError(t, err)

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
