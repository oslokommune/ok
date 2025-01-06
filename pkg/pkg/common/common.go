package common

import "fmt"

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
