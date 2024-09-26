package add

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
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

	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub client: %w", err)
	}

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}
		return nil, fmt.Errorf("failed getting latest github releases: %w", err)
	}

	templateVersion := latestReleases[templateName]
	if templateVersion == "" {
		return nil, fmt.Errorf("template %s not found in latest releases", templateName)
	}
	gitRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, err
	}

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

	_err := allowDuplicateOutputFolder(manifest, newPackage)
	if _err != nil {
		return nil, _err
	}

	manifest.Packages = append(manifest.Packages, newPackage)
	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return nil, err
	}

	if updateSchema {
		downloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, gitRef)
		stackPath := githubreleases.GetTemplatePath(manifest.PackagePrefix(), templateName)
		schema, err := config.GenerateJsonSchemaForApp(ctx, downloader, stackPath, gitRef)
		if err != nil {
			return nil, fmt.Errorf("generating json schema for app: %w", err)
		}
		_, err = config.CreateOrUpdateConfigurationFile(configFile, gitRef, schema)
		if err != nil {
			return nil, fmt.Errorf("creating or updating configuration file: %w", err)
		}
	}

	return &AddResult{
		OutputFolder:    manifest.PackageOutputFolder(outputFolder),
		VarFiles:        varFiles,
		TemplateName:    templateName,
		TemplateVersion: templateVersion,
	}, nil
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
