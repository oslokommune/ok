package common

import (
	"fmt"
	"github.com/Masterminds/semver"
	"strings"
)

const (
	outputFolderWidth = 45
	templateWidth     = 40
	varFilesWidth     = 80
)

type Package struct {
	OutputFolder string   `yaml:"OutputFolder"` // Example: app-hello
	Template     string   `yaml:"Template"`     // Example: app
	Ref          string   `yaml:"Ref"`          // Example: app-v1.2.3
	VarFiles     []string `yaml:"VarFiles"`     // Example: [ "common.yml", "package-config.yml" ]
}

func (p Package) String() string {
	outputFolder := fmt.Sprintf("%-*.*s", outputFolderWidth, outputFolderWidth, p.OutputFolder)
	ref := fmt.Sprintf("%-*.*s", templateWidth, templateWidth, p.Ref)
	varFiles := fmt.Sprintf("%-*.*s", varFilesWidth, varFilesWidth, fmt.Sprint(p.VarFiles))

	return fmt.Sprintf("%s %s %s", outputFolder, ref, varFiles)
}

// Key returns a unique key for the package
func (p Package) Key() string {
	return fmt.Sprintf("outputFolder:%s___Template:%s___Ref:%s___VarFiles:%s", p.OutputFolder, p.Template, p.Ref, strings.Join(p.VarFiles, ","))
}

// PackageVersion returns a semver.Version of the package's Ref, or an error if it fails to parse.
// Example: If Ref is "app-v1.5.0", the returned version will be 1.5.0.
func (p Package) PackageVersion() (*semver.Version, error) {
	parts := strings.Split(p.Ref, "-")
	versionString := parts[len(parts)-1]

	version, err := semver.NewVersion(versionString)
	if err != nil {
		return nil, fmt.Errorf("parsing semantic version: %w", err)
	}

	return version, nil
}

func ContainsPackage(packages []Package, pkg Package) bool {
	for _, p := range packages {
		if p.Key() == pkg.Key() {
			return true
		}
	}

	return false
}
