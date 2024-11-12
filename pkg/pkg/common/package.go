package common

import (
	"fmt"
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
	outputFolder := fmt.Sprintf("%-*.*s", outputFolderWidth, outputFolderWidth, p.OutputFolder)
	template := fmt.Sprintf("%-*.*s", templateWidth, templateWidth, p.Template)
	varFiles := fmt.Sprintf("%-*.*s", varFilesWidth, varFilesWidth, fmt.Sprint(p.VarFiles))

	return fmt.Sprintf("%s %s %s", outputFolder, template, varFiles)
}

// Key returns a unique key for the package
func (p Package) Key() string {
	return fmt.Sprintf("outputFolder:%s___Template:%s___Ref:%s___VarFiles:%s", p.OutputFolder, p.Template, p.Ref, strings.Join(p.VarFiles, ","))
}

func ContainsPackage(packages []Package, pkg Package) bool {
	for _, p := range packages {
		if p.Key() == pkg.Key() {
			return true
		}
	}

	return false
}
