package common

import "fmt"

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

func (p Package) String() string {
	return fmt.Sprintf("%s (%s)", p.OutputFolder, p.Ref)
}
