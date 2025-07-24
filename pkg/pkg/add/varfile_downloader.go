package add

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"path/filepath"
	"regexp"
)

func (a Adder) downloadVarFile(newPackage common.Package, varFile string, varFilePath string, outputFolder string) error {
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

	fileString := string(fileBytes)

	/*
		StackName: "load-balancing-alb-main"
		load-balancing-alb-data.StackName: "load-balancing-alb-main-data"

		StackName: "backup"

		StackName: "app-km"
		app-data.StackName: "app-km-data"
	*/

	var replacement string

	replacement = fmt.Sprintf(`StackName: "%s"`, outputFolder)
	fileString = regexStackName.ReplaceAllString(fileString, replacement)

	fileString = regexStackNameForData.ReplaceAllStringFunc(fileString, func(m string) string {
		submatches := regexStackNameForData.FindStringSubmatch(m)
		if len(submatches) >= 2 {
			templateNameWithData := submatches[1]
			return fmt.Sprintf(`%s.StackName: "%s"`, templateNameWithData, outputFolder+"-data")
		}
		return m
	})

	fileBytes = []byte(fileString)
	err = os.WriteFile(varFilePath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	return err
}

var regexStackName = regexp.MustCompile(`StackName:\s*"[\w\-]*"`)
var regexStackNameForData = regexp.MustCompile(`([\w\-]+)\.StackName:\s*"[\w\-]*"`)
