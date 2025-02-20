package pkg_test

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

type TestData struct {
	name                    string
	args                    []string
	testdataRootDir         string
	packageManifestFilename string
	configDir               string
}

func copyTestdataToTempDir(t *testing.T, tt TestData, tempDir string) {
	if tt.packageManifestFilename != "" {
		srcPath := filepath.Join(tt.testdataRootDir, tt.packageManifestFilename)
		dstPath := filepath.Join(tempDir, filepath.Base(tt.packageManifestFilename))

		err := copyFile(srcPath, dstPath)
		require.NoError(t, err)
	}

	if tt.configDir != "" {
		configDir := filepath.Join(tt.testdataRootDir, tt.configDir)

		srcDir := os.DirFS(configDir)
		dstDir := filepath.Join(tempDir, filepath.Base(tt.configDir))

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
