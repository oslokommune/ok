package common

import "fmt"

type PackageManifest struct {
	Packages []Package `yaml:"Packages"`
}

type Package struct {
	Template     string   `yaml:"Template"`
	Ref          string   `yaml:"Ref"`
	OutputFolder string   `yaml:"OutputFolder"`
	VarFiles     []string `yaml:"VarFiles"`
}

func (p Package) String() string {
	return fmt.Sprintf("%s (%s)", p.OutputFolder, p.Ref)
}
