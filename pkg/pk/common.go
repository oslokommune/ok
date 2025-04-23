package pk

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetGitRoot returns the root directory of the Git repository.
func GetGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetOkDirPath returns the path to the ".ok" directory inside the Git root.
func GetOkDirPath() (string, error) {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitRoot, ".ok"), nil
}

// FindYamlFiles returns a slice of paths to all YAML files in the specified directory.
func FindYamlFiles(dir string) ([]string, error) {
	var yamlFiles []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return yamlFiles, nil
}
