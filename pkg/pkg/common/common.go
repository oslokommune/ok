package common

import (
	"fmt"
	"os"
	"strings"
)

const BaseUrlEnvName = "BASE_URL"

const PackagesManifestFilename = "packages.yml"
const BoilerplateRepoOwner = "oslokommune"
const BoilerplateRepoName = "golden-path-boilerplate"

// boilerplate terraform packages
const BoilerplatePackageTerraformPath = "boilerplate/terraform"
const BoilerplatePackageTerraformConfigPrefix = "_config"

// boilerplate github actions packages
const BoilerplatePackageGitHubActionsPath = "boilerplate/github-actions"
const BoilerplatePackageGitHubActionsConfigPrefix = ""
const BoilerplatePackageGitHubActionsOutputFolder = "../.."

const DefaultVarFileName = "package-config"

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//"
const DefaultPackagePathPrefix = BoilerplatePackageTerraformPath
const DefaultPackageConfigPrefix = BoilerplatePackageTerraformConfigPrefix

func VarFile(prefix, varFileName string) string {
	if prefix == "" {
		return fmt.Sprintf("%s.yml", varFileName)
	}
	return fmt.Sprintf("%s/%s.yml", prefix, varFileName)
}

func PrintProcessedPackages(update []Package, action string) {
	fmt.Println()
	if len(update) == 0 {
		fmt.Printf("No packages were %s.\n", action)
		return
	}
	fmt.Printf("✅ Successfully %s packages:\n", action)

	for _, pkg := range update {
		fmt.Printf("  - %s\n", pkg.String())
	}
}

// GenerateRelativePath creates a path to navigate back to the root directory
// from the specified outputFolder
func GenerateRelativePath(outputFolder string) string {
	outputFolder = strings.TrimRight(outputFolder, "/")
	dirCount := strings.Count(outputFolder, "/") + 1
	path := ""
	for i := 0; i < dirCount; i++ {
		path += "../"
	}
	return strings.TrimRight(path, "/")
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
