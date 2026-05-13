package jit

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = "jit.yaml"

// Config holds the JIT configuration.
type Config struct {
	TenantID string `yaml:"tenant_id"`
	ClientID string `yaml:"client_id"`
	BaseURL  string `yaml:"base_url"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ok", configFileName), nil
}

// LoadConfig reads the JIT config from disk. Returns nil if it doesn't exist.
func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig writes the JIT config to disk.
func SaveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
