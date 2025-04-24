package pk

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

// RepoRoot returns the root directory of the Git repository.
func RepoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// OkDir returns the path to the ".ok" directory inside the Git root.
func OkDir(ctx context.Context) (string, error) {
	gitRoot, err := RepoRoot(ctx)
	if err != nil {
		return "", fmt.Errorf("RepoRoot failed: %w", err)
	}
	return filepath.Join(gitRoot, ".ok"), nil
}

// YAMLFiles returns a slice of paths to all YAML files in the specified directory.
func YAMLFiles(dir string) ([]string, error) {
	var yamlFiles []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking directory %s: %w", dir, err)
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
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", file, err)
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("unmarshalling YAML from file %s: %w", file, err)
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

// buildGitSource constructs a Git source URL with a subpath and optional query parameters.
func buildGitSource(repo, subPath, ref string) string {
	// Ensure exactly one “//” separator between repo and subPath.
	repo = strings.TrimSuffix(repo, "/")
	subPath = strings.TrimPrefix(subPath, "/")

	q := url.Values{}
	if ref != "" {
		q.Set("ref", ref)
	}

	return fmt.Sprintf("%s//%s?%s", repo, subPath, q.Encode())
}

func BuildBoilerplateArgs(tpl Template) ([]string, error) {
	source := buildGitSource(tpl.Repo, tpl.Path, tpl.Ref)

	args := []string{
		"--template-url", source,
		"--output-folder", filepath.Join(tpl.BaseOutputFolder, tpl.Subfolder),
	}

	if tpl.NonInteractive {
		args = append(args, "--non-interactive")
	}
	for _, vf := range tpl.VarFiles {
		args = append(args, "--var-file", vf)
	}
	return args, nil
}

// RunBoilerplateCommand takes arguments and a working directory as input and executes the boilerplate command.
func RunBoilerplateCommand(ctx context.Context, args []string, workingDir string) error {
	cmd := exec.CommandContext(ctx, "boilerplate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workingDir

	return cmd.Run()
}
