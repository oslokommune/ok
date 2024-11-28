package add

import (
	"context"
	"fmt"
	"github.com/google/go-github/v63/github"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
    "bytes"
	"slices"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/oslokommune/ok/pkg/pkg/config"
)

type AddResult struct {
	OutputFolder    string
	VarFiles        []string
	TemplateName    string
	TemplateVersion string
}

/**
 * Add Boilerplate template to packages manifest with an optional stack name.
 * The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
 * The output folder is prefixed with the stack name and added to the packages manifest.
 */

func Run(pkgManifestFilename string, templateName, outputFolder string, updateSchema bool) (*AddResult, error) {
	ctx := context.Background()

	gh, err := getGitHubClient()
	if err != nil {
		return nil, err
	}

	templateVersion, err := getTemplateVersion(templateName)
	if err != nil {
		return nil, err
	}
	gitRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, err
	}
    
	newPackage, err := createNewPackage(manifest, templateName, gitRef, outputFolder)
	if err != nil {
		return nil, err
	}

	if err := allowDuplicateOutputFolder(manifest, newPackage); err != nil {
		return nil, err
	}

	manifest.Packages = append(manifest.Packages, newPackage)
	if err := common.SavePackageManifest(pkgManifestFilename, manifest); err != nil {
		return nil, err
	}

	if err := createDefaultConfig(ctx, gh, manifest, templateName, gitRef, outputFolder); err != nil {
		return nil, err
	}

	if updateSchema {
		if err := updateSchemaConfig(ctx, gh, manifest, templateName, gitRef, outputFolder); err != nil {
			return nil, err
		}
	}

	return &AddResult{
		OutputFolder:    manifest.PackageOutputFolder(outputFolder),
		VarFiles:        newPackage.VarFiles,
		TemplateName:    templateName,
		TemplateVersion: templateVersion,
	}, nil
}

func getGitHubClient() (*github.Client, error) {
	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub client: %w", err)
	}
	return gh, nil
}

func getTemplateVersion(templateName string) (string, error) {
	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}
		return "", fmt.Errorf("failed getting latest github releases: %w", err)
	}

	templateVersion := latestReleases[templateName]
	if templateVersion == "" {
		return "", fmt.Errorf("template %s not found in latest releases", templateName)
	}
	return templateVersion, nil
}

func createNewPackage(manifest common.PackageManifest, templateName, gitRef, outputFolder string) (common.Package, error) {
	configFile := common.ConfigFile(manifest.PackageConfigPrefix(), outputFolder)
	commonConfigFile := common.ConfigFile(manifest.PackageConfigPrefix(), "common-config")

	varFiles := []string{
		commonConfigFile,
		configFile,
	}

	newPackage := common.Package{
		Template:     templateName,
		Ref:          gitRef,
		OutputFolder: manifest.PackageOutputFolder(outputFolder),
		VarFiles:     varFiles,
	}

	return newPackage, nil
}

type KeyVal struct {
    Key string
    Val interface{}
}

// Define an ordered map
type OrderedMap []KeyVal

// Implement the json.Marshaler interface
func (omap OrderedMap) MarshalYAML() ([]byte, error) {
    var buf bytes.Buffer

    for _, kv := range omap {
        key, err := yaml.Marshal(kv.Key)
        if err != nil {
            return nil, err
        }
		newKey := strings.TrimSuffix(string(key),"\n") + ": "
        buf.WriteString(newKey)
		
		val, err := yaml.Marshal(kv.Val)
        if err != nil {
            return nil, err
        }
		

		if _, ok := kv.Val.(string); ok {
			buf.Write(val)
		} else if _, ok := kv.Val.(int); ok {
			buf.Write(val)
		} else if _, ok := kv.Val.(bool); ok {
			buf.Write(val)
		} else {
			buf.WriteString(string("\n"))
			for _, line := range strings.Split(string(val), "\n") {
				buf.WriteString("    " + line + "\n")
			}
		}
    }

    return buf.Bytes(), nil
}

func createConfigFromBoilerplate(ctx context.Context, downloader config.FileDownloader, stackPath, gitRef string, stackName string) ([]byte, error) {
	stacks, err := config.DownloadBoilerplateStacksWithDependencies(ctx, downloader, stackPath)
	if err != nil {
		return []byte(""), fmt.Errorf("downloading boilerplate stacks: %w", err)
	}

	modules := schema.BuildModuleVariables(stacks)

	var  variables OrderedMap

	ignore_fields := []string{"AccountId", "Team", "Environment", "TemplateVersion", "TerraformVersion", "AwsProviderVersion", "Region"}

	var prefix string
	for _, module := range modules {
		fmt.Println(module.Namespace)
		if module.Namespace != "" {
			prefix = module.Namespace + "."
		} else {
			prefix = ""
		}

		for _, variable := range module.Variables {
			if module.Namespace != "" && variable.Name != "StackName" {
				continue
			}

			fmt.Println(variable.Name, variable)
			if slices.Contains(ignore_fields, variable.Name) {
				continue
			}
			
			if variable.Default == nil {
				if variable.Name == "StackName" || variable.Name == "AppName"  {
					value := stackName
					var valuePostfix string

					if variable.Name == "AppName" {
						value = strings.TrimPrefix(stackName, "app-")
					}

					if strings.Contains(prefix, "data") {
						valuePostfix = "-data"
					} 

					variables = append(variables, KeyVal{Key: prefix + variable.Name, Val: value + valuePostfix})
				} else {
					variables = append(variables, KeyVal{Key: prefix + variable.Name, Val: "<fill this in>"})
				}
			} else {
				variables = append(variables, KeyVal{Key: prefix + variable.Name, Val: variable.Default})
			}
		}
	}

	defaultConfig, err := variables.MarshalYAML()
	if err != nil {
		return []byte(""), fmt.Errorf("marshalling modules: %w", err)
	}

	return defaultConfig, nil
}

func createDefaultConfig(ctx context.Context, gh *github.Client, manifest common.PackageManifest, templateName, gitRef, outputFolder string) error {
	downloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, gitRef)
	stackPath := githubreleases.GetTemplatePath(manifest.PackagePrefix(), templateName)
	defaultConfig, err := createConfigFromBoilerplate(ctx, downloader, stackPath, gitRef, outputFolder)
	if err != nil {
		return fmt.Errorf("generating config for: %w", err)
	}
	configFile := common.ConfigFile(manifest.PackageConfigPrefix(), outputFolder)

	if _, err := os.Stat(configFile); err == nil {
		return nil
	}

	if err := os.MkdirAll(manifest.PackageConfigPrefix(), 0755); err != nil {
		return fmt.Errorf("creating folder: %w", err)
	}

	err = os.WriteFile(configFile, defaultConfig, 0644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	fmt.Println("Created configuration file: ", configFile)

	return nil
}

func updateSchemaConfig(ctx context.Context, gh *github.Client, manifest common.PackageManifest, templateName, gitRef, outputFolder string) error {
	downloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, gitRef)
	stackPath := githubreleases.GetTemplatePath(manifest.PackagePrefix(), templateName)
	generatedSchema, err := schema.GenerateJsonSchemaForApp(ctx, downloader, stackPath, gitRef)
	if err != nil {
		return fmt.Errorf("generating json schema for app: %w", err)
	}
	configFile := common.ConfigFile(manifest.PackageConfigPrefix(), outputFolder)
	_, err = schema.CreateOrUpdateConfigurationFile(configFile, gitRef, generatedSchema)
	if err != nil {
		return fmt.Errorf("creating or updating configuration file: %w", err)
	}
	return nil
}

func allowDuplicateOutputFolder(manifest common.PackageManifest, newPackage common.Package) error {
	// If we are generating GHA there is no restriction on output folder
	if manifest.PackagePrefix() == common.BoilerplatePackageGitHubActionsPath {
		return nil
	}
	for _, pkg := range manifest.Packages {
		if pkg.OutputFolder == newPackage.OutputFolder {
			return fmt.Errorf("output folder %s already exists in packages manifest", newPackage.OutputFolder)
		}
	}
	return nil
}
