package common

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckVarFiles returns an error if a var file of the packages does not exist. Var file paths are relative to
// workingDirectory. If workingDirectory is in a git repository, the error suggests a file with the same name in a
// parent directory, up to the repository root.
func CheckVarFiles(manifestFile string, packages []Package, workingDirectory string) error {
	return checkVarFiles(manifestFile, packages, workingDirectory, gitRepoRoot(workingDirectory))
}

func checkVarFiles(manifestFile string, packages []Package, workingDirectory string, repoRoot string) error {
	var messages []string
	checked := make(map[string]bool)

	for _, pkg := range packages {
		for _, varFile := range pkg.VarFiles {
			if checked[varFile] {
				continue
			}
			checked[varFile] = true

			varFilePath := varFile
			if !filepath.IsAbs(varFilePath) {
				varFilePath = filepath.Join(workingDirectory, varFilePath)
			}

			exists, err := fileExists(varFilePath)
			if err != nil {
				return fmt.Errorf("checking var file %s: %w", varFile, err)
			}

			if exists {
				continue
			}

			message := fmt.Sprintf("%s: var file %q does not exist", displayPath(manifestFile, repoRoot), varFile)

			suggestion, ok := findVarFileSuggestion(varFile, workingDirectory, repoRoot)
			if ok {
				message += fmt.Sprintf(". Did you mean %q?", suggestion)
			}

			messages = append(messages, message)
		}
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}

// findVarFileSuggestion looks for varFile, without its leading "../" parts, in workingDirectory and each parent
// directory up to repoRoot. It returns the nearest match, relative to workingDirectory.
//
// Example: if varFile is "../common-config.yml", it checks "common-config.yml", "../common-config.yml",
// "../../common-config.yml" and so on, and skips "../common-config.yml" as it does not exist.
func findVarFileSuggestion(varFile string, workingDirectory string, repoRoot string) (string, bool) {
	if repoRoot == "" || filepath.IsAbs(varFile) {
		return "", false
	}

	name := stripLeadingParentDirs(varFile)
	if name == "" {
		return "", false
	}

	start, err := resolvePath(workingDirectory)
	if err != nil {
		return "", false
	}

	missing := filepath.Join(start, varFile)
	dir := start

	for isInDir(dir, repoRoot) {
		candidate := filepath.Join(dir, name)

		if candidate != missing {
			exists, err := fileExists(candidate)
			if err == nil && exists {
				suggestion, err := filepath.Rel(start, candidate)
				if err != nil {
					return "", false
				}

				return suggestion, true
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false
}

func stripLeadingParentDirs(path string) string {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")

	for len(parts) > 0 && (parts[0] == ".." || parts[0] == ".") {
		parts = parts[1:]
	}

	return filepath.Join(parts...)
}

// isInDir returns true if path is dir or a subdirectory of dir.
func isInDir(path string, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}

	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// displayPath returns path relative to repoRoot, or path as is if this is not possible.
func displayPath(path string, repoRoot string) string {
	if repoRoot == "" {
		return path
	}

	resolved, err := resolvePath(path)
	if err != nil {
		return path
	}

	rel, err := filepath.Rel(repoRoot, resolved)
	if err != nil || !isInDir(resolved, repoRoot) {
		return path
	}

	return rel
}

// gitRepoRoot returns the root of the git repository that contains dir, or an empty string if dir is not in a git
// repository.
func gitRepoRoot(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	root, err := resolvePath(strings.TrimSpace(string(out)))
	if err != nil {
		return ""
	}

	return root
}

// resolvePath returns the absolute path of path, with symbolic links resolved. Symbolic links make paths from git and
// from the working directory different, for example /var and /private/var on macOS.
func resolvePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(abs)
}
