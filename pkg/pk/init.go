package pk

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFileName = "config.yaml"
	DefaultRepo           = "git@github.com:oslokommune/golden-path-boilerplate.git"
)

// InitOptions contains options for the Init function.
type InitOptions struct {
	ConfigFileName   string // e.g., "dev.yaml" or "config.yaml"
	BaseOutputFolder string // e.g., "stacks/dev" or "."
}

// DefaultConfig returns a Config with sensible default values for the common section.
func DefaultConfig() Config {
	return Config{
		Common: Template{
			Repo:             DefaultRepo,
			NonInteractive:   true,
			BaseOutputFolder: ".",
		},
		Templates: []Template{},
	}
}

// ConfigWithBase returns a Config with a custom base_output_folder.
func ConfigWithBase(baseOutputFolder string) Config {
	return Config{
		Common: Template{
			Repo:             DefaultRepo,
			NonInteractive:   true,
			BaseOutputFolder: baseOutputFolder,
		},
		Templates: []Template{},
	}
}

// Init creates the .ok directory and a config file.
func Init(okDir string, opts InitOptions) error {
	// Create .ok directory if it doesn't exist
	if err := os.MkdirAll(okDir, 0755); err != nil {
		return fmt.Errorf("creating .ok directory: %w", err)
	}

	// Determine config filename
	configFileName := opts.ConfigFileName
	if configFileName == "" {
		configFileName = DefaultConfigFileName
	}

	configPath := filepath.Join(okDir, configFileName)

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	// Determine base output folder
	baseOutputFolder := opts.BaseOutputFolder
	if baseOutputFolder == "" {
		baseOutputFolder = "."
	}

	return WriteConfig(configPath, ConfigWithBase(baseOutputFolder))
}

// WriteConfig writes a Config to the specified file path.
func WriteConfig(path string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
