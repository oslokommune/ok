package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

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
	BoilerplateStack struct {
		Path   string
		Config *BoilerplateConfig
	}
)

type FileDownloader interface {
	DownloadFile(ctx context.Context, file string) ([]byte, error)
}

type GithubFileReference struct {
	Organization string
	Repository   string
	FilePath     string
	GitRef       string
}

func DownloadBoilerplateStacksWithDependencies(ctx context.Context, client FileDownloader, stackPath string) ([]*BoilerplateStack, error) {
	stacks := make([]*BoilerplateStack, 0)
	stackPathsToDownload := []string{stackPath}
	downloadedStacks := make(map[string]bool)

	for len(stackPathsToDownload) > 0 {
		stackPath := stackPathsToDownload[0]
		stackPathsToDownload = stackPathsToDownload[1:]
		stack, err := DownloadBoilerplateStack(ctx, client, stackPath)
		if err != nil {
			return nil, fmt.Errorf("download boilerplate stack: %w", err)
		}
		stacks = append(stacks, stack)
		downloadedStacks[stackPath] = true
		// Add dependencies to download queue if not already downloaded
		for _, dep := range stack.Config.Dependencies {
			templateUrl := JoinPath(stackPath, dep.TemplateUrl)
			if _, ok := downloadedStacks[templateUrl]; !ok {
				stackPathsToDownload = append(stackPathsToDownload, templateUrl)
			}
		}
	}
	return stacks, nil
}

func DownloadBoilerplateStack(ctx context.Context, client FileDownloader, stackPath string) (*BoilerplateStack, error) {
	boilerplatePath := JoinPath(stackPath, "boilerplate.yml")
	cfg, err := DownloadBoilerplateConfig(ctx, client, boilerplatePath)
	if err != nil {
		return nil, fmt.Errorf("download boilerplate config: %w", err)
	}

	return &BoilerplateStack{
		Path:   stackPath,
		Config: cfg,
	}, nil
}

func DownloadBoilerplateConfig(ctx context.Context, client FileDownloader, filePath string) (*BoilerplateConfig, error) {
	data, err := client.DownloadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	cfg, err := parseBoilerplateConfig(data)
	if err != nil {
		return nil, fmt.Errorf("parse boilerplate config: %w", err)
	}
	return cfg, nil
}

func parseBoilerplateConfig(data []byte) (*BoilerplateConfig, error) {
	var config BoilerplateConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal boilerplate config: %w", err)
	}
	return &config, nil
}

func JoinPath(base, path string) string {
	uri, err := url.JoinPath(base, path)
	if err != nil {
		slog.Error("could not join paths", slog.String("base", base), slog.String("path", path), slog.String("error", err.Error()))
		panic(err)
	}
	return uri
}
