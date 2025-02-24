package add

import (
	"context"
	"fmt"
	"github.com/google/go-github/v63/github"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
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

type Adder struct {
	schemaGenerator common.SchemaGenerator
}

func NewAdder(schemaGenerator common.SchemaGenerator) Adder {
	return Adder{
		schemaGenerator: schemaGenerator,
	}
}

func (a Adder) Run(pkgManifestFilename string, templateName, outputFolder string, updateSchema bool) (*AddResult, error) {
	ctx := context.Background()

	templateVersion, err := getTemplateVersion(templateName)
	if err != nil {
		return nil, err
	}
	pkgRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, err
	}

	newPackage, err := createNewPackage(manifest, templateName, pkgRef, outputFolder)
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

	if updateSchema {
		if err := a.updateSchemaConfig(ctx, manifest, newPackage, outputFolder); err != nil {
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
	configFile := common.VarFile(manifest.PackageConfigPrefix(), outputFolder)
	commonConfigFile := common.VarFile(manifest.PackageConfigPrefix(), "common-config")

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

func (a Adder) updateSchemaConfig(
	ctx context.Context, manifest common.PackageManifest, pkg common.Package, outputFolder string) error {
	varFile := common.VarFile(manifest.PackageConfigPrefix(), outputFolder)

	// TODO replace
	jsonSchemaData, err := a.schemaGenerator.CreateJsonSchemaFile(ctx, manifest.PackagePrefix(), pkg)
	if err != nil {
		return fmt.Errorf("creating json schema file: %w", err)
	}

	schemaFilePath := schema.GetSchemaFilePath(varFile, pkg.Ref)
	err = schema.WriteSchemaToFile(schemaFilePath, jsonSchemaData)
	if err != nil {
		return fmt.Errorf("writing schema to file %s: %w", schemaFilePath, err)
	}

	_, err = schema.SetVarFileSchemaDeclaration(varFile, pkg.Ref)
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
