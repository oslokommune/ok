package pkg_test

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

type TestData struct {
	name                        string
	args                        []string
	testdataRootDir             string
	workingDirectoryFromRootDir string
	releases                    map[string]string
	keepTempDir                 bool

	expectFiles   []string
	expectNoFiles []string

	expectError        bool
	expectErrorMessage string
}

const inputDir = "input"
const inputRootDir = "root"

func copyTestdataRootDirToTempDir(t *testing.T, tt TestData, testWorkingDirectory string, tempDir string) {
	var err error
	rootDir := filepath.Join(testWorkingDirectory, tt.testdataRootDir, inputDir, inputRootDir)

	srcDir := os.DirFS(rootDir)
	dstDir := tempDir

	err = os.CopyFS(dstDir, srcDir)
	require.NoError(t, err)
}
