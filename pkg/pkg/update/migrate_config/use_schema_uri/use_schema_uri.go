package use_schema_uri

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"log/slog"
	"regexp"
)

func ReplaceDirWithUri(varFile string, varFileJsonSchema metadata.JsonSchema) error {
	slog.Debug("replaing dir with uri for json schema",
		slog.String("varFile", varFile),
		slog.Any("varFileJsonSchema", varFileJsonSchema),
	)

	// Check if already transformed
	should, err := shouldMigrate(varFile)
	if err != nil {
		return fmt.Errorf("checking if YAML already is transformed: %w", err)
	}

	if !should {
		slog.Debug("not updating, is already transformed")
		return nil
	}

	return migrate(varFile, varFileJsonSchema)
}

// shouldMigrate returns true if json schema declaration uses URI instead of directory
func shouldMigrate(varFile string) (bool, error) {
	firstLine, err := metadata.ReadFirstLine(varFile)
	if err != nil {
		return false, fmt.Errorf("reading file %s: %w", varFile, err)
	}

	// This regexp should match on 1, but not 2:
	// 1: # yaml-language-server: $schema=.schemas/app-v9.0.0.schema.json
	// 2: # yaml-language-server: $schema=https://raw.githubusercontent.com/oslokommune/golden-path-boilerplate-schemas/refs/heads/main/schemas/app-v9.0.0.schema.json
	re := regexp.MustCompile(`\$schema=\.schemas`)

	return re.MatchString(firstLine), nil
}

func migrate(varFile string, varFileJsonSchema metadata.JsonSchema) error {
	fmt.Printf("Replacing reference to JSON schema from local directory .schemas to an URL. File: %s\n", varFile)

	err := schema.SetSchemaDeclarationInVarFile(varFile, varFileJsonSchema.Ref())
	if err != nil {
		return fmt.Errorf("creating or updating configuration file: %w", err)
	}

	return nil
}
