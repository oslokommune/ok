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

// TODO Rewrite copyTestdataToTempDir to work like this. I.e. test structure for update tests must be the same. We
// want all the test files to be copied, not select which ones.
func copyTestdataRootDirToTempDir(t *testing.T, tt TestData, tempDir string) {
	configDir := filepath.Join(tt.testdataRootDir, "input", "root")

	srcDir := os.DirFS(configDir)
	dstDir := tempDir

	err := os.CopyFS(dstDir, srcDir)
	require.NoError(t, err)
}
