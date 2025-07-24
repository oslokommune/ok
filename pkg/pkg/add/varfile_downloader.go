package add

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"path/filepath"
)

func (a Adder) downloadVarFile(newPackage common.Package, varFile string, varFilePath string) error {
	// TODO
	// - [x] ok pkg add app hello --var-file default
	// - [ ] ok pkg add databases --var-file non-serverless

	// TODO: Validate this earlier.
	if _, err := os.Stat(varFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", varFilePath)
	}

	// https://github.com/oslokommune/golden-path-boilerplate/tree/databases-v4.0.3/boilerplate/terraform/databases/package-config-default.yml
	varFileDownloadFilename := fmt.Sprintf("package-config-%s.yml", varFile)

	// boilerplate/terraform/{pkg.Template}/package-config-default.yml
	path := filepath.Join(common.BoilerplatePackageTerraformPath, newPackage.Template, varFileDownloadFilename)

	// Show URL to the user
	varFileUrl := fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s/%s/%s",
		common.BoilerplateRepoOwner,
		common.BoilerplateRepoName,
		newPackage.Ref,
		common.BoilerplatePackageTerraformPath,
		newPackage.Template,
		varFileDownloadFilename,
	)
	fmt.Printf("Creating var file %s from %s\n", varFilePath, varFileUrl)

	fileBytes, err := a.ghReleases.DownloadGithubFile(
		context.Background(),
		"oslokommune",
		"golden-path-boilerplate",
		path,
		newPackage.Ref,
	)
	if err != nil {
		return fmt.Errorf("downloading file from GitHub: %w", err)
	}

	// TODO bør sette StackName til outputfolder

	err = os.WriteFile(varFilePath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	return err
}
