package install

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func FindPackagesFromOutputFolders(packages []common.Package, outputFolders []string) []common.Package {
	packagesFound := []common.Package{}

	for _, pkg := range packages {
		for _, outputFolder := range outputFolders {
			if pkg.OutputFolder == outputFolder {
				packagesFound = append(packagesFound, pkg)
			}
		}
	}

	return packagesFound
}
