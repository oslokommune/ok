package common

type PackageManifest struct {
	DefaultPackagePathPrefix string    `yaml:"DefaultPackagePathPrefix,omitempty"`
	Packages                 []Package `yaml:"Packages"`
}

func (pm *PackageManifest) PackagePrefix() string {
	if pm.DefaultPackagePathPrefix != "" {
		return pm.DefaultPackagePathPrefix
	}
	return DefaultPackagePathPrefix
}

// Different package types have different config prefixes
func (pm *PackageManifest) PackageConfigPrefix() string {
	prefix := pm.PackagePrefix()
	if prefix == BoilerplatePackageGitHubActionsPath {
		return BoilerplatePackageGitHubActionsConfigPrefix
	}
	if prefix == BoilerplatePackageTerraformPath {
		return BoilerplatePackageTerraformConfigPrefix
	}
	return DefaultPackageConfigPrefix
}

func (pm *PackageManifest) PackageOutputFolder(outputFolder string) string {
	prefix := pm.PackagePrefix()
	// All GHA must have the same output folder, so we can't use the outputFolder argument
	// that the user is supplying to decide where to configure the output
	if prefix == BoilerplatePackageGitHubActionsPath {
		return BoilerplatePackageGitHubActionsOutputFolder
	}
	return outputFolder
}
