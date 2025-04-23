package pk

import (
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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

// LoadConfigs loads YAML files from the specified directory into a slice of Config structures.
func LoadConfigs(dir string) ([]Config, error) {
	yamlFiles, err := FindYamlFiles(dir)
	if err != nil {
		return nil, err
	}

	var configs []Config
	for _, file := range yamlFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	return configs, nil
}
