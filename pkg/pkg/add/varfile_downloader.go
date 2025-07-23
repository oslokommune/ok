package add

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"os"
	"path/filepath"
)

func downloadVarFile(
	newPackage common.Package, varFile string, outputFolder string, consolidatedPackageStructure bool) error {
	var varFilePath string
	if consolidatedPackageStructure {
		//varFilePath = common.VarFile(manifest.PackageConfigPrefix(), outputFolder)
	} else {
		varFilePath = common.VarFile(outputFolder, common.DefaultVarFileName)
	}

	// TODO: Validate this earlier.
	if _, err := os.Stat(varFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", varFilePath)
	}

	client, err := githubreleases.GetGitHubClient()
	if err != nil {
		return fmt.Errorf("getting GitHub client: %w", err)
	}

	// https://github.com/oslokommune/golden-path-boilerplate/tree/databases-v4.0.3/boilerplate/terraform/databases/package-config-default.yml
	filename := fmt.Sprintf("package-config-%s.yml", varFile)
	varFileUrl := fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s/%s/%s",
		common.BoilerplateRepoOwner,
		common.BoilerplateRepoName,
		newPackage.Ref,
		common.BoilerplatePackageTerraformPath,
		newPackage.Template,
		filename,
	)

	// boilerplate/terraform/{pkg.Template}/package-config-default.yml
	path := filepath.Join(common.BoilerplatePackageTerraformPath, newPackage.Template, filename)

	fmt.Printf("Downloading var file %s to %s\n", varFileUrl, varFilePath)

	fileBytes, err := githubreleases.DownloadGithubFile(
		context.Background(),
		client,
		"oslokommune",
		"golden-path-boilerplate",
		path,
		"move-lock-files", // newPackage.Ref,
	)
	if err != nil {
		return fmt.Errorf("downloading file from GitHub: %w", err)
	}

	err = os.WriteFile(varFilePath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	// TODO consolidatedPackageStructure logikk må komme et annet sted. Alt denne funksjonen ønsker er å vite
	// hvor varfila skal legges. Og outputFolder er nyttig for å sette StackName antakelig. Men da er det StackName
	// som bør være input.
	return err
}
