package migrate_config

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/add_apex_domain"
)

func UpdatePackageConfig(packagesToUpdate []common.Package) error {
	for _, pkg := range packagesToUpdate {
		for _, varFile := range pkg.VarFiles {
			if err := updateVarFile(varFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateVarFile(varFile string) error {
	err := add_apex_domain.AddApexDomainSupport(varFile)
	if err != nil {
		return err
	}

	return nil
}
