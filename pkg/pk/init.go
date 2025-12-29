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

// Init creates the .ok directory and a default config.yaml file.
// Returns an error if the .ok directory already exists.
func Init(okDir string) error {
	if _, err := os.Stat(okDir); err == nil {
		return fmt.Errorf(".ok directory already exists at %s", okDir)
	}

	if err := os.MkdirAll(okDir, 0755); err != nil {
		return fmt.Errorf("creating .ok directory: %w", err)
	}

	configPath := filepath.Join(okDir, DefaultConfigFileName)
	return WriteConfig(configPath, DefaultConfig())
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
