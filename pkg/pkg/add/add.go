package add

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/schema"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

type AddOptions struct {
	CurrentDir                   string
	TemplateName                 string
	OutputFolder                 string
	ConsolidatedPackageStructure bool
	AddSchema                    bool
	DownloadVarFile              bool
	VarFile                      string
}

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

type Adder struct {
}

func NewAdder() Adder {
	return Adder{}
}

func (a Adder) Run(opts AddOptions) (*AddResult, error) {
	oldPackageStructure, err := common.UseOldPackageStructure(opts.CurrentDir)
	if err != nil {
		return nil, fmt.Errorf("checking whether to use old or new package structure: %w", err)
	}

	var packagesManifestFilename = common.PackagesManifestFilename
	if !oldPackageStructure {
		packagesManifestFilename = filepath.Join(opts.OutputFolder, packagesManifestFilename)
	}

	templateVersion, err := getTemplateVersion(opts.TemplateName)
	if err != nil {
		return nil, fmt.Errorf("getting template version: %w", err)
	}

	pkgRef := fmt.Sprintf("%s-%s", opts.TemplateName, templateVersion)

	manifest, err := common.LoadPackageManifest(packagesManifestFilename)
	if err != nil {
		return nil, err
	}

	newPackage, err := createNewPackage(manifest, opts.TemplateName, pkgRef, opts.OutputFolder, opts.ConsolidatedPackageStructure)
	if err != nil {
		return nil, fmt.Errorf("creating new package: %w", err)
	}

	err = allowDuplicateOutputFolder(manifest, newPackage)
	if err != nil {
		return nil, err
	}

	manifest.Packages = append(manifest.Packages, newPackage)
	err = common.SavePackageManifest(packagesManifestFilename, manifest)
	if err != nil {
		return nil, fmt.Errorf("saving package manifest: %w", err)
	}

	if opts.DownloadVarFile {
		// ok pkg add app hello --var-file default
		// ok pkg add databases --var-file non-serverless
		// TODO implementer denne
		err := downloadVarFile(newPackage, opts.VarFile, opts.OutputFolder, opts.ConsolidatedPackageStructure)
		if err != nil {
			return &AddResult{}, fmt.Errorf("downloading var file: %w", err)
		}
	}

	if opts.AddSchema {
		err := a.addSchemaConfig(manifest, newPackage, opts.OutputFolder, opts.ConsolidatedPackageStructure)
		if err != nil {
			return &AddResult{}, fmt.Errorf("adding schema config: %w", err)
		}
	}

	if oldPackageStructure {
		nonExistingConfigFiles := findNonExistingConfigurationFiles(newPackage.VarFiles)
		if len(nonExistingConfigFiles) > 0 {
			fmt.Println("Create the following configuration files:")
			for _, configFile := range nonExistingConfigFiles {
				fmt.Printf("- %s\n", configFile)
			}
		}
	}

	return &AddResult{
		OutputFolder:    manifest.PackageOutputFolder(opts.OutputFolder),
		VarFiles:        newPackage.VarFiles,
		TemplateName:    opts.TemplateName,
		TemplateVersion: templateVersion,
	}, nil
}

func getTemplateVersion(templateName string) (string, error) {
	fmt.Printf("Fetching latest releases from GitHub repository %s/%s\n", common.BoilerplateRepoOwner, common.BoilerplateRepoName)

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

func createNewPackage(manifest common.PackageManifest, templateName, gitRef, outputFolder string, consolidatedPackageStructure bool) (common.Package, error) {
	var configFile, commonConfigFile, out string
	if consolidatedPackageStructure {
		configFile = common.VarFile(manifest.PackageConfigPrefix(), outputFolder)
		commonConfigFile = common.VarFile(manifest.PackageConfigPrefix(), "common-config")
		out = manifest.PackageOutputFolder(outputFolder)
	} else {
		configFile = common.VarFile("", common.DefaultVarFileName)
		commonConfigFile = common.VarFile(common.GenerateRelativePath(outputFolder), "common-config")
		out = "."
	}

	varFiles := []string{
		commonConfigFile,
		configFile,
	}

	newPackage := common.Package{
		Template:     templateName,
		Ref:          gitRef,
		OutputFolder: out,
		VarFiles:     varFiles,
	}

	return newPackage, nil
}

func (a Adder) addSchemaConfig(manifest common.PackageManifest, pkg common.Package, outputFolder string, consolidatedPackageStructure bool) error {
	var varFilePath string
	if consolidatedPackageStructure {
		varFilePath = common.VarFile(manifest.PackageConfigPrefix(), outputFolder)
	} else {
		varFilePath = common.VarFile(outputFolder, common.DefaultVarFileName)
	}

	err := schema.SetSchemaDeclarationInVarFile(varFilePath, pkg.Ref)
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

func findNonExistingConfigurationFiles(varFiles []string) []string {
	var nonExisting []string
	for _, varFile := range varFiles {
		_, err := os.Stat(varFile)
		notExists := errors.Is(err, os.ErrNotExist)
		if notExists {
			nonExisting = append(nonExisting, varFile)
		}
	}
	return nonExisting
}
