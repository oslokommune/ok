package common

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// CheckVarFiles returns an error if a var file of the packages does not exist. Var file paths are relative to
// workingDirectory, which is also the directory of manifestFile. If workingDirectory is in a git repository, the error
// suggests a file with the same name in a parent directory, up to the repository root.
func CheckVarFiles(manifestFile string, packages []Package, workingDirectory string) error {
	prefix, levelsToRoot, inGitRepo := gitLocation(workingDirectory)

	displayManifestFile := manifestFile
	if inGitRepo {
		displayManifestFile = path.Join(prefix, filepath.Base(manifestFile))
	}

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

			message := fmt.Sprintf("%s: var file %q does not exist", displayManifestFile, varFile)

			if inGitRepo {
				suggestion, ok := findVarFileSuggestion(varFile, workingDirectory, levelsToRoot)
				if ok {
					message += fmt.Sprintf(". Did you mean %q?", suggestion)
				}
			}

			messages = append(messages, message)
		}
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}

// findVarFileSuggestion looks for varFile, without its leading "../" parts, in workingDirectory and in each of the
// levels parent directories above it. It returns the nearest match, relative to workingDirectory.
//
// Example: if varFile is "../common-config.yml", it checks "common-config.yml", "../../common-config.yml" and so on.
// It skips "../common-config.yml", as that file does not exist.
func findVarFileSuggestion(varFile string, workingDirectory string, levels int) (string, bool) {
	if filepath.IsAbs(varFile) {
		return "", false
	}

	name := stripLeadingParentDirs(varFile)
	if name == "" {
		return "", false
	}

	for i := 0; i <= levels; i++ {
		candidate := filepath.Join(strings.Repeat("../", i), name)
		if candidate == filepath.Clean(varFile) {
			continue
		}

		exists, err := fileExists(filepath.Join(workingDirectory, candidate))
		if err == nil && exists {
			return candidate, true
		}
	}

	return "", false
}

func stripLeadingParentDirs(varFile string) string {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(varFile)), "/")

	for len(parts) > 0 && (parts[0] == ".." || parts[0] == ".") {
		parts = parts[1:]
	}

	return filepath.Join(parts...)
}

// gitLocation returns the path of dir relative to the root of its git repository, for example "stacks/prod/", and the
// number of directory levels from dir up to the root. ok is false if dir is not in a git repository.
func gitLocation(dir string) (prefix string, levelsToRoot int, ok bool) {
	cmd := exec.Command("git", "rev-parse", "--show-prefix", "--show-cdup")
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return "", 0, false
	}

	// The output has one line for each flag. The second line is "../" for each level, for example "../../".
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(lines) != 2 {
		return "", 0, false
	}

	return lines[0], strings.Count(lines[1], "../"), true
}
