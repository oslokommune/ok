package pk

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"gopkg.in/yaml.v3"
)

const (
	BoilerplateTerraformPath = "boilerplate/terraform"
)

// GitHubReleases is an interface for fetching GitHub releases.
type GitHubReleases interface {
	GetLatestReleases() (map[string]string, error)
}

// AddOptions contains options for the Add function.
type AddOptions struct {
	TemplateName string
	Subfolder    string
	Ref          string
	ConfigFile   string
}

// Add adds a new template to the specified config file.
func Add(okDir string, opts AddOptions, ghReleases GitHubReleases) error {
	// Determine config file
	configFile := opts.ConfigFile
	if configFile == "" {
		var err error
		configFile, err = selectOrCreateConfigFile(okDir)
		if err != nil {
			return fmt.Errorf("selecting config file: %w", err)
		}
	}

	// Load existing config
	cfg, err := loadSingleConfig(configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Always fetch releases to validate template exists
	fmt.Printf("Fetching latest releases from GitHub repository %s/%s\n",
		common.BoilerplateRepoOwner, common.BoilerplateRepoName)

	releases, err := ghReleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("fetching releases: %w", err)
	}

	// Validate template exists
	version, ok := releases[opts.TemplateName]
	if !ok {
		return fmt.Errorf("template %q not found. Available templates: %s",
			opts.TemplateName, availableTemplates(releases))
	}

	// Determine ref (version)
	ref := opts.Ref
	if ref == "" {
		ref = fmt.Sprintf("%s-%s", opts.TemplateName, version)
	} else {
		// Normalize ref: if user provided "v1.0.0", convert to "template-v1.0.0"
		ref = normalizeRef(opts.TemplateName, ref)
	}

	// Determine subfolder
	subfolder := opts.Subfolder
	if subfolder == "" {
		subfolder = opts.TemplateName
	}

	// Check for duplicate subfolder
	for _, t := range cfg.Templates {
		if t.Subfolder == subfolder {
			return fmt.Errorf("template with subfolder %q already exists", subfolder)
		}
	}

	// Build template entry
	// Only include fields that differ from common
	tpl := Template{
		Name:      opts.TemplateName,
		Path:      fmt.Sprintf("%s/%s", BoilerplateTerraformPath, opts.TemplateName),
		Ref:       ref,
		Subfolder: subfolder,
	}

	cfg.Templates = append(cfg.Templates, tpl)

	// Save config
	if err := WriteConfig(configFile, cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("\n✅ Added template %s to %s\n", ref, configFile)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  ok pk install  - Install the template\n")

	return nil
}

// availableTemplates returns a comma-separated list of template names.
func availableTemplates(releases map[string]string) string {
	names := make([]string, 0, len(releases))
	for name := range releases {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// normalizeRef ensures the ref is in the format "template-vX.Y.Z".
func normalizeRef(templateName, ref string) string {
	// If ref already starts with template name, return as-is
	if strings.HasPrefix(ref, templateName+"-") {
		return ref
	}

	// If ref starts with "v", prepend template name
	if strings.HasPrefix(ref, "v") {
		return fmt.Sprintf("%s-%s", templateName, ref)
	}

	// Otherwise, assume it's a version number and add "v" prefix
	return fmt.Sprintf("%s-v%s", templateName, ref)
}

// selectOrCreateConfigFile returns the config file to use.
// If multiple files exist, prompts the user to select one.
// If no files exist, returns the default config file path.
func selectOrCreateConfigFile(okDir string) (string, error) {
	files, err := YAMLFiles(okDir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return filepath.Join(okDir, DefaultConfigFileName), nil
	}

	if len(files) == 1 {
		return files[0], nil
	}

	// Multiple files - prompt user to select
	return selectConfigFileInteractively(files)
}

// selectConfigFileInteractively prompts the user to select a config file.
func selectConfigFileInteractively(files []string) (string, error) {
	options := make([]huh.Option[string], 0, len(files))
	for _, f := range files {
		options = append(options, huh.NewOption(filepath.Base(f), f))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Select config file").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", fmt.Errorf("selecting config file: %w", err)
	}

	return selected, nil
}

// loadSingleConfig loads a single config file.
func loadSingleConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return Config{}, fmt.Errorf("reading file: %w", err)
	}

	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decoding YAML: %w", err)
	}

	return cfg, nil
}
