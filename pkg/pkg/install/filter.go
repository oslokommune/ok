package install

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func FindPackagesFromOutputFolders(packages []common.Package, outputFolders []string) []common.Package {
	result := make([]common.Package, 0)

	for _, pkg := range packages {
		for _, outputFolder := range outputFolders {

			if pkg.OutputFolder == outputFolder {
				result = append(result, pkg)
				break
			}

		}
	}

	return result
}
