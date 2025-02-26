package pkg_test

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

type TestData struct {
	name            string
	args            []string
	testdataRootDir string
}

func copyTestdataRootDirToTempDir(t *testing.T, tt TestData, tempDir string) {
	configDir := filepath.Join(tt.testdataRootDir, "input", "root")

	srcDir := os.DirFS(configDir)
	dstDir := tempDir

	err := os.CopyFS(dstDir, srcDir)
	require.NoError(t, err)
}
