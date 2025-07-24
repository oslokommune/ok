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
	releases        map[string]string
	keepTempDir     bool

	expectError   bool
	expectedFiles []string
}

const inputDir = "input"
const inputRootDir = "root"

func copyTestdataRootDirToTempDir(t *testing.T, tt TestData, tempDir string) {
	configDir := filepath.Join(tt.testdataRootDir, inputDir, inputRootDir)

	_, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		require.FailNow(t, "required test data dir does not exist", "configDir: %s", configDir)
	}

	srcDir := os.DirFS(configDir)
	dstDir := tempDir

	err = os.CopyFS(dstDir, srcDir)
	require.NoError(t, err)
}
