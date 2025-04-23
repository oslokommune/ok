package pk

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

// RepoRoot returns the root directory of the Git repository.
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// OkDir returns the path to the ".ok" directory inside the Git root.
func OkDir() (string, error) {
	gitRoot, err := RepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitRoot, ".ok"), nil
}

// YAMLFiles returns a slice of paths to all YAML files in the specified directory.
func YAMLFiles(dir string) ([]string, error) {
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
	yamlFiles, err := YAMLFiles(dir)
	if err != nil {
		return nil, err
	}

	var configs []Config
	for _, file := range yamlFiles {
		data, err := os.ReadFile(file) // Updated from ioutil.ReadFile to os.ReadFile
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

// ApplyCommon applies cfg.Common to each cfg.Templates entry and
// returns the fully-resolved templates.
func ApplyCommon(cfgs []Config) ([]Template, error) {
	var out []Template

	for _, cfg := range cfgs {
		for _, tpl := range cfg.Templates {
			merged := tpl // start with the template

			// copy non-zero fields from cfg.Common into merged
			if err := mergo.Merge(&merged, cfg.Common); err != nil {
				return nil, err
			}

			out = append(out, merged)
		}
	}
	return out, nil
}

// BuildBoilerplateArgs takes a Template and constructs the arguments for the boilerplate command.
func BuildBoilerplateArgs(tpl Template) []string {
	args := []string{
		"--template-url", tpl.Repo + "//" + tpl.Path + "?ref=" + tpl.Ref,
		"--output-folder", filepath.Join(tpl.BaseOutputFolder, tpl.Subfolder),
	}

	if tpl.NonInteractive {
		args = append(args, "--non-interactive")
	}

	for _, varFile := range tpl.VarFiles {
		args = append(args, "--var-file", varFile)
	}

	return args
}

// RunBoilerplateCommand takes arguments as input and executes the boilerplate command.
func RunBoilerplateCommand(args []string) error {
	cmd := exec.Command("boilerplate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
