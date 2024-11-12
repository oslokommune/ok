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

func ConfigFile(prefix, configName string) string {
	if prefix == "" {
		return fmt.Sprintf("%s.yml", configName)
	}
	return fmt.Sprintf("%s/%s.yml", prefix, configName)
}

func PrintProcessedPackages(update []Package, action string) {
	fmt.Println()
	fmt.Printf("✅ Successfully %s packages:\n", action)

	for _, pkg := range update {
		fmt.Printf("  - %s\n", pkg.String())
	}
}
