package migrate_config

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/add_apex_domain"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
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
	firstLine, err := readFirstLine(varFile)
	if err != nil {
		return fmt.Errorf("reading first line from %s: %w", varFile, err)
	}

	varFileMetadata, err := metadata.ParseMetadata(firstLine)
	if err != nil {
		return fmt.Errorf("getting metadata from var file %s: %w", varFile, err)
	}

	err = update(varFile, varFileMetadata)
	if err != nil {
		return fmt.Errorf("updating varFile %s: %w", varFile, err)
	}

	return nil
}

func update(varFile string, metadata metadata.VarFileMetadata) error {
	err := add_apex_domain.AddApexDomainSupport(varFile, metadata)
	if err != nil {
		return err
	}

	return nil
}
