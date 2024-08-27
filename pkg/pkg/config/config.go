package config

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
)

type (
	BoilerplateVariable struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description,omitempty"`
		Type        string `yaml:"type,omitempty"`
		Default     any    `yaml:"default,omitempty"`
	}

	BoilerplateDependency struct {
		Name         string `yaml:"name"`
		TemplateUrl  string `yaml:"template-url"`
		OutputFolder string `yaml:"output-folder"`
	}

	BoilerplateConfig struct {
		Variables    []BoilerplateVariable   `yaml:"variables"`
		Dependencies []BoilerplateDependency `yaml:"dependencies"`
	}
)

type GithubFileDownloader interface {
	DownloadFile(ctx context.Context, file string) ([]byte, error)
}

type GithubFileReference struct {
	Organization string
	Repository   string
	FilePath     string
	GitRef       string
}

func DownloadBoilerplateConfig(ctx context.Context, client GithubFileDownloader, filePath string) (*BoilerplateConfig, error) {
	data, err := client.DownloadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	return parseBoilerplateConfig(data)
}

func parseBoilerplateConfig(data []byte) (*BoilerplateConfig, error) {
	var config BoilerplateConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal boilerplate config: %w", err)
	}
	return &config, nil
}
