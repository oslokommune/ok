package install

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

var ErrPackageNotFound = fmt.Errorf("no package found")

func FindPackageFromOutputFolders(packages []common.Package, outputFolders []string) (common.Package, error) {
	for _, pkg := range packages {
		for _, outputFolder := range outputFolders {
			if pkg.OutputFolder == outputFolder {
				return pkg, nil
			}
		}
	}

	return common.Package{}, ErrPackageNotFound
}
