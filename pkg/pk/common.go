package pk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"dario.cat/mergo"
	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v3"
)

// RepoRoot returns the root directory of the Git repository.
func RepoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Join(fmt.Errorf("git rev-parse --show-toplevel"), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// OkDir returns the path to the ".ok" directory inside the Git root.
func OkDir(ctx context.Context) (string, error) {
	gitRoot, err := RepoRoot(ctx)
	if err != nil {
		return "", errors.Join(fmt.Errorf("RepoRoot failed"), err)
	}
	return filepath.Join(gitRoot, ".ok"), nil
}

// YAMLFiles returns a slice of paths to all YAML files in the specified directory.
func YAMLFiles(dir string) ([]string, error) {
	var yamlFiles []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort the YAML files to ensure deterministic order
	sort.Strings(yamlFiles)

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
		dec := yaml.NewDecoder(bytes.NewReader(data))
		dec.KnownFields(true) // Enable strict mode
		if err := dec.Decode(&config); err != nil {
			return nil, fmt.Errorf("decoding YAML from file %s: %w", file, err)
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
			merged := tpl

			if err := mergo.Merge(
				&merged,
				cfg.Common,
				mergo.WithAppendSlice,
			); err != nil {
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
	qs := q.Encode() // "" when ref is empty

	if qs != "" {
		return fmt.Sprintf("%s//%s?%s", repo, subPath, qs)
	}
	return fmt.Sprintf("%s//%s", repo, subPath)
}

func BuildBoilerplateArgs(tpl Template) []string {
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
	return args
}

// RunBoilerplateCommand takes arguments and a working directory as input and executes the boilerplate command.
func RunBoilerplateCommand(ctx context.Context, args []string, workingDir string) error {
	cmd := exec.CommandContext(ctx, "boilerplate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workingDir

	return cmd.Run()
}

// FilterTemplatesByWorkingDir returns templates whose output path contains the current working directory.
// If cwd is within or equals a template's output path (baseOutputFolder/subfolder), that template is included.
func FilterTemplatesByWorkingDir(templates []Template, cwd, repoRoot string) []Template {
	var matched []Template

	// Make cwd relative to repo root for comparison
	relCwd, err := filepath.Rel(repoRoot, cwd)
	if err != nil {
		return nil
	}

	// Normalize to handle "." for repo root
	if relCwd == "." {
		return nil // At repo root, no filtering - caller should show picker
	}

	for _, tpl := range templates {
		outputPath := filepath.Join(tpl.BaseOutputFolder, tpl.Subfolder)
		// Clean both paths for consistent comparison
		outputPath = filepath.Clean(outputPath)

		// Check if cwd is within or equals the output path
		if relCwd == outputPath || strings.HasPrefix(relCwd, outputPath+string(filepath.Separator)) {
			matched = append(matched, tpl)
		}
	}

	return matched
}

// SelectTemplatesInteractively shows a multi-select picker for templates.
func SelectTemplatesInteractively(templates []Template) ([]Template, error) {
	if len(templates) == 0 {
		return nil, nil
	}

	options := make([]huh.Option[string], 0, len(templates))
	templateMap := make(map[string]Template)

	for i, tpl := range templates {
		label := tpl.Subfolder
		if label == "" {
			label = tpl.Name
		}
		outputPath := filepath.Clean(filepath.Join(tpl.BaseOutputFolder, tpl.Subfolder))
		optionKey := fmt.Sprintf("%s|%s|%s|%d", outputPath, tpl.Repo, tpl.Path, i)
		displayText := fmt.Sprintf("%s (%s)", label, tpl.Path)
		options = append(options, huh.NewOption(displayText, optionKey))
		templateMap[optionKey] = tpl
	}

	var selectedKeys []string

	s := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Select template(s) to install").
		Value(&selectedKeys)

	err := huh.NewForm(huh.NewGroup(s)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, nil
		}
		return nil, fmt.Errorf("running template selection: %w", err)
	}

	var selected []Template
	for _, key := range selectedKeys {
		selected = append(selected, templateMap[key])
	}

	return selected, nil
}
