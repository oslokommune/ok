package pk

import (
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
