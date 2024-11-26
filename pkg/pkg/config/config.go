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
	BoilerplateTemplate struct {
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

func DownloadBoilerplateTemplatesWithDependencies(ctx context.Context, client FileDownloader, initialTemplatePath string) ([]*BoilerplateTemplate, error) {
	templates := make([]*BoilerplateTemplate, 0)
	// stackPath example: boilerplate/terraform/app
	// TODO: Add support for: git@github.com:oslokommune/golden-path-boilerplate//boilerplate/terraform/versions?ref=versions-v3.0.6
	templatePathsOrUrisToDownload := []string{initialTemplatePath}
	downloadedTemplates := make(map[string]bool)

	for len(templatePathsOrUrisToDownload) > 0 {
		templatePath := templatePathsOrUrisToDownload[0]
		templatePathsOrUrisToDownload = templatePathsOrUrisToDownload[1:]

		template, err := DownloadBoilerplateTemplate(ctx, client, templatePath)
		if err != nil {
			return nil, fmt.Errorf("download boilerplate stack: %w", err)
		}

		templates = append(templates, template)
		downloadedTemplates[templatePath] = true

		// Add dependencies to download queue if not already downloaded
		for _, dep := range template.Config.Dependencies {
			templateUrl := JoinPath(templatePath, dep.TemplateUrl)
			if _, ok := downloadedTemplates[templateUrl]; !ok {
				templatePathsOrUrisToDownload = append(templatePathsOrUrisToDownload, templateUrl)
			}
		}
	}

	return templates, nil
}

func DownloadBoilerplateTemplate(ctx context.Context, client FileDownloader, stackPath string) (*BoilerplateTemplate, error) {
	boilerplatePath := JoinPath(stackPath, "boilerplate.yml")
	cfg, err := DownloadBoilerplateConfig(ctx, client, boilerplatePath)
	if err != nil {
		return nil, fmt.Errorf("download boilerplate config: %w", err)
	}

	return &BoilerplateTemplate{
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
