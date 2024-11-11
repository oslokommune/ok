package migrate_config

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/rogpeppe/go-internal/semver"
)

/**
	Packages:
    - OutputFolder: app-veileder-api
      Template: app
      Ref: app-v8.0.5
      VarFiles:
        - _config/common-config.yml
        - _config/app-veileder-api.yml

    - OutputFolder: task-hvv-pg-backup-data-stores
      Template: scaffold
      Ref: scaffold-v2.1.2
      VarFiles:
        - _config/common-config.yml
        - _config/task-hvv-pg-backup-data-stores.yml

*/

func UpdatePackageConfig(packagesToUpdate []common.Package) error {
	for _, pkg := range packagesToUpdate {
		for _, varFile := range pkg.VarFiles {
			// Get current package version from header of file:
			// # yaml-language-server: $schema=.schemas/app-v8.0.5.schema.json

			// template: app
			// version: v8.0.5

			template := "app"
			version := "v8.0.5"

			err := update(template, version)
			if err != nil {
				return fmt.Errorf("updating package config. Filename: %s. Error: %w", varFile, err)
			}

		}
	}

	return nil
}

func update(template string, version string, varFile string) error {
	// Read file
	if template == "app" && semver.Compare(version, "9.0.0") > 9 {
		err := addApexDomainSupport(varFile)
		if err != nil {
			return fmt.Errorf("adding apex domain support: %w", err)
		}
	}

	return nil
}
