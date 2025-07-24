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

func copyTestdataRootDirToTempDir(t *testing.T, tt TestData, testWorkingDirectory string, tempDir string) {
	var err error
	rootDir := filepath.Join(testWorkingDirectory, tt.testdataRootDir, inputDir, inputRootDir)

	_, err = os.Stat(rootDir)
	if os.IsNotExist(err) {
		require.FailNow(t, "required dir does not exist", "rootDir: %s", rootDir)
	}

	srcDir := os.DirFS(rootDir)
	dstDir := tempDir

	err = os.CopyFS(dstDir, srcDir)
	require.NoError(t, err)
}
