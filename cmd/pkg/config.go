package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v63/github"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ConfigCommand = &cobra.Command{
	Use: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("config command")
		token, err := getGHToken()
		if err != nil {
			return fmt.Errorf("error getting token: %w", err)
		}
		cmd.Printf("token: %s\n", token)
		//return doStuff()
		return nil
	},
}

func mustOpenFileWrite(path string) *os.File {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

type DownloadedBoilThingy struct {
	IsRootCfg    bool
	Path         string
	Namespaces   []string
	Config       *BoilerplateConfig
	Dependencies []BoilerplateDependency
	Variables    []BoilerplateVariable
}

var DownloadCommand = &cobra.Command{
	Use: "download",
	RunE: func(cmd *cobra.Command, args []string) error {
		gh, err := githubreleases.GetGitHubClient()
		if err != nil {
			return fmt.Errorf("getting GitHub client: %w", err)
		}
		releases, err := githubreleases.GetLatestReleases()
		if err != nil {
			return fmt.Errorf("getting latest releases: %w", err)
		}

		templateName := args[0]
		templateVersion := releases[templateName]
		githubRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

		templatePath := "boilerplate/terraform"

		type boilDep struct {
			namespace string
			path      string
		}

		var downloaded = make(map[string]*DownloadedBoilThingy)

		var toDownload []boilDep = []boilDep{
			// {
			// 	namespace: "",
			// 	path:      mustJoinUri(templatePath, templateName),
			// },
		}
		var addDependencyToDownload = func(namespace, path string, name string) {
			dep := boilDep{
				namespace: namespace,
				path:      mustJoinUri(path, name),
			}
			cmd.Printf("-> Adding dependency [%s] %s\n", dep.namespace, dep.path)
			toDownload = append(toDownload, dep)
		}

		addDependencyToDownload("", templatePath, templateName)

		for i := 0; i < len(toDownload); i++ {
			repoDirPath := toDownload[i]
			//		for _, repoDirPath := range toDownload {
			log.Printf("downloading %s\n", repoDirPath)
			var conf *BoilerplateConfig = nil
			if cached, ok := downloaded[repoDirPath.path]; ok {
				conf = cached.Config
				log.Printf("-> Already downloaded %s\n", repoDirPath.path)
			} else {
				config, err := downloadBoilerplateConfig(cmd.Context(), gh, repoDirPath.path, githubRef)
				if err != nil {
					return fmt.Errorf("downloading boilerplate config: %s:  %w", repoDirPath, err)
				}
				conf = config
			}
			if cached, ok := downloaded[repoDirPath.path]; ok {
				cached.Namespaces = append(cached.Namespaces, repoDirPath.namespace)
			} else {
				downloaded[repoDirPath.path] = &DownloadedBoilThingy{
					IsRootCfg:    i == 0,
					Path:         repoDirPath.path,
					Namespaces:   []string{repoDirPath.namespace},
					Config:       conf,
					Dependencies: conf.Dependencies,
					Variables:    conf.Variables,
				}
			}
			// downloaded = append(downloaded, DownloadedBoilThingy{
			// 	Namespace:    repoDirPath.namespace,
			// 	Config:       conf,
			// 	Dependencies: conf.Dependencies,
			// 	Variables:    conf.Variables,
			// })
			for _, d := range conf.Dependencies {
				depNamespace := joinNamespacePath(repoDirPath.namespace, d.Name)
				// if toDownload does not contain d.TemplateUrl add it to the list
				// dependencyURL := mustJoinUri(repoDirPath.path, d.TemplateUrl)
				// cmd.Printf("-> Adding dependency [%s] %s\n", depNamespace, dependencyURL)
				addDependencyToDownload(depNamespace, repoDirPath.path, d.TemplateUrl)
				// toDownload = append(toDownload, boilDep{
				// 	namespace: depNamespace,
				// 	path:      dependencyURL,
				// })
			}
			//time.Sleep(1 * time.Second)
		}

		fileSchema := mustOpenFileWrite(fmt.Sprintf("%s-%s.schema.json", templateName, templateVersion))
		fileDependencies := mustOpenFileWrite(fmt.Sprintf("%s-%s.dependencies.yml", templateName, templateVersion))
		defer fileSchema.Close()
		enc := yaml.NewEncoder(fileDependencies)
		if err := enc.Encode(downloaded); err != nil {
			return fmt.Errorf("encoding downloaded: %w", err)
		}
		defer enc.Close()

		/*
			boilerplateYmlPath, err := url.JoinPath(templatePath, templateName, "boilerplate.yml")
			if err != nil {
				return fmt.Errorf("joining path: %w", err)
			}

			cmd.Printf("downloading boilerplate template %s@%s\n", templateName, templateVersion)
			data, res, err := gh.Repositories.DownloadContents(cmd.Context(), "oslokommune", "golden-path-boilerplate", boilerplateYmlPath, &github.RepositoryContentGetOptions{
				Ref: githubRef,
			})
			if err != nil {
				return fmt.Errorf("downloading boilerplate.yml: %w", err)
			}
			defer res.Body.Close()
			defer data.Close()

			dec := yaml.NewDecoder(data)
			var config BoilerplateConfig
			if err := dec.Decode(&config); err != nil {
				return fmt.Errorf("decoding boilerplate.yml: %w", err)
			}
			enc := yaml.NewEncoder(os.Stdout)
			if err := enc.Encode(config); err != nil {
				return fmt.Errorf("encoding boilerplate.yml: %w", err)
			}
			enc.Close()

			//fmt.Printf("\n\n%#v\n", config)

			githubreleases.GetLatestReleases()

			cmd.Println("download command")
			return nil
		*/

		_, err = makeJsonSchemaFromDependencies(downloaded)
		return err
		//		return nil
	},
}

func makeJsonSchemaFromDependencies(configs map[string]*DownloadedBoilThingy) ([]byte, error) {
	for app, config := range configs {
		for _, namespace := range config.Namespaces {
			for _, variable := range config.Variables {

				log.Printf("->[%s] \t%s [%s]\n", app, joinNamespacePath(namespace, variable.Name), variable.Description)
			}
		}
	}

	return nil, fmt.Errorf("not implemented")
}

func mustJoinUri(base, path string) string {
	uri, err := url.JoinPath(base, path)
	if err != nil {
		panic(err)
	}
	return uri
}

func joinNamespacePath(namespace, path string) string {
	if namespace == "" {
		return path
	}
	return fmt.Sprintf("%s.%s", namespace, path)
}

type (
	BoilerplateVariable struct {
		Name        string `yaml:"name,omitempty"`
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

func downloadBoilerplateConfig(ctx context.Context, gh *github.Client, directory, githubRef string) (*BoilerplateConfig, error) {
	boilerplateYmlPath, err := url.JoinPath(directory, "boilerplate.yml")
	if err != nil {
		return nil, fmt.Errorf("joining path: %w", err)
	}
	data, res, err := gh.Repositories.DownloadContents(ctx, "oslokommune", "golden-path-boilerplate", boilerplateYmlPath, &github.RepositoryContentGetOptions{
		Ref: githubRef,
	})

	if err != nil {
		return nil, fmt.Errorf("downloading boilerplate.yml: %w", err)
	}
	defer res.Body.Close()
	defer data.Close()

	dec := yaml.NewDecoder(data)
	var config BoilerplateConfig
	if err := dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("decoding boilerplate.yml: %w", err)
	}
	return &config, nil
}

func init() {
	ConfigCommand.AddCommand(DownloadCommand)
}

func getGHToken() (string, error) {
	output, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", fmt.Errorf("getting token from gh cli: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func doStuff() error {
	f, err := os.Open("some-boilerplate.yml")
	if err != nil {
		return err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	type YamlStruct struct {
		Variables []struct {
			Name        string `yaml:"name" json:"name"`
			Description string `yaml:"description" json:"description"`
			Type        string `yaml:"type" json:"type"`
			Default     any    `yaml:"default" json:"default"`
		} `yaml:"variables"`
	}

	var someStruct YamlStruct
	if err := dec.Decode(&someStruct); err != nil {
		return err
	}
	fmt.Printf("%#v\n\n", someStruct)

	type JSONProps struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		Default     any    `json:"default"`
	}

	type JSONStruct struct {
		Schema     string               `json:"$schema"`
		Title      string               `json:"title"`
		Type       string               `json:"type"`
		Properties map[string]JSONProps `json:"properties"`
	}

	var outProps = make(map[string]JSONProps)
	for _, v := range someStruct.Variables {
		if strings.Contains(v.Description, "do NOT edit") {
			continue
		}
		outProps[v.Name] = JSONProps{
			Description: v.Description,
			Type:        mapTypeToJsonSchema(v.Type),
			Default:     v.Default,
		}
	}
	schema := JSONStruct{
		Schema:     "http://json-schema.org/draft-07/schema#",
		Title:      "Some Boilerplate",
		Type:       "object",
		Properties: outProps,
	}
	bts, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bts))
	return nil
}

func mapTypeToJsonSchema(t string) string {
	switch t {
	case "map":
		return "object"
	case "int":
		return "integer"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	default:
		return "string"
	}
}
