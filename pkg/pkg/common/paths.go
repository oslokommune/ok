package common

import "path/filepath"

// UseOldPackageStructure returns true if the current file system is using the old
// way of organizing package files, as described here:
//
// https://github.com/oslokommune/ok/pull/429
//
// Otherwise, false is returned.
//
// The "old" means using a centralized package manifest together with a separate var file directory such as "_config".
func UseOldPackageStructure(dir string) (bool, error) {
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
