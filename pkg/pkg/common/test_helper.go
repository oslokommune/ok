package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetTestdataFilepath(testDataFilename string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}

	return filepath.Join(cwd, "testdata", testDataFilename), nil
}
