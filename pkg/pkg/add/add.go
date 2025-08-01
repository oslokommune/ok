package add

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/schema"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

const FlagNoVar = "no-var-file"

type Options struct {
	BaseUrl         string
	CurrentDir      string
	TemplateName    string
	OutputFolder    string
	AddSchema       bool
	DownloadVarFile bool
	VarFile         string
}

/**
 * Add Boilerplate template to packages manifest with an optional stack name.
 * The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
 * The output folder is prefixed with the stack name and added to the packages manifest.
 */

type Adder struct {
	ghReleases GitHubReleases
}

type GitHubReleases interface {
	GetLatestReleases() (map[string]string, error)
	DownloadGithubFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error)
}

func NewAdder(ghReleases GitHubReleases) Adder {
	return Adder{
		ghReleases: ghReleases,
	}
}

func (a Adder) Run(opts Options) error {
	consolidatedPackageStructure, err := common.UseConsolidatedPackageStructure(opts.CurrentDir)
	if err != nil {
		return fmt.Errorf("checking whether to use old or new package structure: %w", err)
	}

	var packagesManifestFilename string
	if consolidatedPackageStructure {
		packagesManifestFilename = common.PackagesManifestFilename
	} else {
		packagesManifestFilename = filepath.Join(opts.OutputFolder, common.PackagesManifestFilename)
	}

	manifest, err := common.LoadPackageManifest(packagesManifestFilename)
	if err != nil {
		return err
	}

	err = createErrorIfOutputFolderExists(manifest, opts.OutputFolder)
	if err != nil {
		return err
	}

	var templateVersion string
	if len(opts.BaseUrl) == 0 {
		templateVersion, err = a.getTemplateVersion(opts.TemplateName)
		if err != nil {
			return fmt.Errorf("getting template version: %w", err)
		}
	} else {
		templateVersion = "using-base-url-not-relevant"
	}

	pkgRef := fmt.Sprintf("%s-%s", opts.TemplateName, templateVersion)

	newPackage, err := createNewPackage(manifest, opts.TemplateName, pkgRef, opts.OutputFolder, consolidatedPackageStructure)
	if err != nil {
		return fmt.Errorf("creating new package: %w", err)
	}

	err = createErrorIfPackageExistsInManifest(manifest, packagesManifestFilename, newPackage)
	if err != nil {
		return err
	}

	manifest.Packages = append(manifest.Packages, newPackage)

	fmt.Printf("Creating package manifest %s\n", packagesManifestFilename)
	err = common.SavePackageManifest(packagesManifestFilename, manifest)
	if err != nil {
		return fmt.Errorf("saving package manifest: %w", err)
	}

	varFilePath := getVarFilePath(consolidatedPackageStructure, manifest, opts.OutputFolder)

	if opts.DownloadVarFile {
		err = a.downloadVarFile(manifest, newPackage, opts.BaseUrl, opts.VarFile, varFilePath, opts.OutputFolder)
		if err != nil {
			return fmt.Errorf("downloading var file: %w", err)
		}
	}

	if opts.DownloadVarFile && opts.AddSchema {
		err = schema.SetSchemaDeclarationInVarFile(varFilePath, newPackage.Ref)
		if err != nil {
			return fmt.Errorf("creating or updating configuration file: %w", err)
		}
	}

	var partMsg string
	if consolidatedPackageStructure {
		partMsg = "with output"
	} else {
		partMsg = "to"
	}

	fmt.Println()
	fmt.Printf("✅ Successfully added package %s %s directory %s.\n",
		common.StyleHighlight.Render(
			fmt.Sprintf("%s-%s", opts.TemplateName, templateVersion),
		),
		partMsg,
		common.StyleHighlight.Render(manifest.PackageOutputFolder(opts.OutputFolder)),
	)
	fmt.Println()
	fmt.Printf("%sOpen %s to configure your stack.\n",
		common.StyleHighlight.Render("Next step: "),
		common.StyleHighlight.Render(varFilePath),
	)

	return nil
}

func createErrorIfOutputFolderExists(manifest common.PackageManifest, outputFolder string) error {
	// If we are generating GHA there is no restriction on output folder
	if manifest.PackagePrefix() == common.BoilerplatePackageGitHubActionsPath {
		return nil
	}

	info, err := os.Stat(outputFolder)
	if err == nil && info.IsDir() {
		return fmt.Errorf("folder already exists: %s", outputFolder)
	}

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to verify folder existence: %w", err)
	}

	return nil
}

func (a Adder) getTemplateVersion(templateName string) (string, error) {
	fmt.Printf("Fetching latest releases from GitHub repository %s/%s\n", common.BoilerplateRepoOwner, common.BoilerplateRepoName)

	latestReleases, err := a.ghReleases.GetLatestReleases()
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

func getVarFilePath(consolidatedPackageStructure bool, manifest common.PackageManifest, outputFolder string) string {
	if consolidatedPackageStructure {
		return common.VarFile(manifest.PackageConfigPrefix(), outputFolder)
	} else {
		return common.VarFile(outputFolder, common.DefaultVarFileName)
	}
}

func createNewPackage(manifest common.PackageManifest, templateName, gitRef, outputFolderCmdArgument string, consolidatedPackageStructure bool) (common.Package, error) {
	var mainVarFile, commonVarFile, outputFolder string
	if consolidatedPackageStructure {
		mainVarFile = common.VarFile(manifest.PackageConfigPrefix(), outputFolderCmdArgument)
		commonVarFile = common.VarFile(manifest.PackageConfigPrefix(), "common-config")
		outputFolder = manifest.PackageOutputFolder(outputFolderCmdArgument)
	} else {
		mainVarFile = common.VarFile("", common.DefaultVarFileName)
		commonVarFile = common.VarFile(common.GenerateRelativePath(outputFolderCmdArgument), "common-config")
		outputFolder = "."
	}

	varFiles := []string{
		commonVarFile,
		mainVarFile,
	}

	newPackage := common.Package{
		Template:     templateName,
		Ref:          gitRef,
		OutputFolder: outputFolder,
		VarFiles:     varFiles,
	}

	return newPackage, nil
}

func createErrorIfPackageExistsInManifest(
	manifest common.PackageManifest, packagesManifestFilename string, newPackage common.Package,
) error {
	// If we are generating GHA there is no restriction on output folder
	if manifest.PackagePrefix() == common.BoilerplatePackageGitHubActionsPath {
		return nil
	}

	for _, pkg := range manifest.Packages {
		if pkg.OutputFolder == newPackage.OutputFolder {
			return fmt.Errorf(
				"output folder %s already exists in package manifest %s",
				newPackage.OutputFolder,
				packagesManifestFilename,
			)
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
