package common

import (
	"fmt"
	"strings"
)

type Package struct {
	OutputFolder string   `yaml:"OutputFolder"`
	Template     string   `yaml:"Template"`
	Ref          string   `yaml:"Ref"`
	VarFiles     []string `yaml:"VarFiles"`
}

func (p Package) String() string {
	return fmt.Sprintf("%s (%s)", p.OutputFolder, p.Ref)
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
