package common

import (
	"fmt"
	"strings"
)

type PackageManifest struct {
	DefaultPackagePathPrefix string    `yaml:"DefaultPackagePathPrefix,omitempty"`
	Packages                 []Package `yaml:"Packages"`
}

type Package struct {
	OutputFolder string   `yaml:"OutputFolder"`
	Template     string   `yaml:"Template"`
	Ref          string   `yaml:"Ref"`
	VarFiles     []string `yaml:"VarFiles"`
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

func (p Package) String() string {
	return fmt.Sprintf("%s (%s)", p.OutputFolder, p.Ref)
}

// Key returns a unique key for the package
func (p Package) Key() string {
	return fmt.Sprintf("outputFolder:%s___Template:%s___Ref:%s___VarFiles:%s", p.OutputFolder, p.Template, p.Ref, strings.Join(p.VarFiles, ","))
}
