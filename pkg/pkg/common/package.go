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
	OutputFolder string   `yaml:"OutputFolder"`
	Template     string   `yaml:"Template"`
	Ref          string   `yaml:"Ref"`
	VarFiles     []string `yaml:"VarFiles"`
}

func (p Package) String() string {
	var reversedVarFiles []string
	for i := len(p.VarFiles) - 1; i >= 0; i-- {
		reversedVarFiles = append(reversedVarFiles, p.VarFiles[i])
	}

	outputFolder := fmt.Sprintf("%-*.*s", outputFolderWidth, outputFolderWidth, p.OutputFolder)
	if len(outputFolder) > outputFolderWidth {
		outputFolder = fmt.Sprintf("%s...", outputFolder[:outputFolderWidth-3])
		fmt.Println(outputFolder)
	}

	template := fmt.Sprintf("%-*.*s", templateWidth, templateWidth, p.Template)
	if len(template) > templateWidth {
		template = fmt.Sprintf("%s...", template[:templateWidth-3])
	}

	varFiles := fmt.Sprintf("%-*.*s", varFilesWidth, varFilesWidth, fmt.Sprint(reversedVarFiles))
	if len(varFiles) > varFilesWidth {
		varFiles = fmt.Sprintf("%s...", varFiles[:varFilesWidth-3])
	}

	//return fmt.Sprintf("OutputFolder:%s\n\tTemplate: %s\n\tVarFiles: %s", outputFolder, template, varFiles)
	return fmt.Sprintf("%s %s %s", outputFolder, template, varFiles)
}

// Key returns a unique key for the package
func (p Package) Key() string {
	return fmt.Sprintf("outputFolder:%s___Template:%s___Ref:%s___VarFiles:%s", p.OutputFolder, p.Template, p.Ref, strings.Join(p.VarFiles, ","))
}

// PackageVersion returns a semver.Version of the package's Ref, or an error if it fails to parse.
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
