package common

import "path/filepath"

// UseConsolidatedPackageStructure returns true if the current file system is either
// - using the old way of organizing package files for Terraform, as described here:
// https://github.com/oslokommune/ok/pull/429
// - the package manifest is using DefaultPackagePrefix: "boilerplate/github-actions"
//
// Otherwise, false is returned.
//
// The "old" means using a centralized package manifest together with a separate var file directory such as "_config".
func UseConsolidatedPackageStructure(dir string) (bool, error) {
	packageManifestPath := filepath.Join(dir, PackagesManifestFilename)
	varFileDir := filepath.Join(dir, BoilerplatePackageTerraformConfigPrefix)

	packageManifestExists, err := fileExists(packageManifestPath)
	if err != nil {
		return false, err
	}

	varFileDirExists, err := dirExists(varFileDir)
	if err != nil {
		return false, err
	}

	if packageManifestExists && varFileDirExists {
		return true, nil
	}

	if packageManifestExists {
		manifest, err := LoadPackageManifest(packageManifestPath)
		if err != nil {
			return false, err
		}

		if manifest.DefaultPackagePathPrefix == BoilerplatePackageGitHubActionsPath {
			return true, nil
		}
	}

	return false, nil
}
