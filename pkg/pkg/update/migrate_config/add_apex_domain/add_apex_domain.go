package add_apex_domain

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"log/slog"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

func AddApexDomainSupport(varFile string, varFileJsonSchema metadata.JsonSchema) error {
	slog.Debug("adding apex support",
		slog.String("varFile", varFile),
		slog.Any("varFileJsonSchema", varFileJsonSchema),
	)

	//
	// Check version
	// If json schema version is 9.0.0 or greater, we know that the config file supports Apex domain routing.
	// https://github.com/oslokommune/golden-path-boilerplate/releases/tag/app-v9.0.0
	//

	versionOfApexDomainSupport, err := semver.NewVersion("9.0.0")
	if err != nil {
		return fmt.Errorf("creating semver: %w", err)
	}

	if varFileJsonSchema.Version.LessThan(versionOfApexDomainSupport) {
		slog.Debug("not updating, package less than required version",
			slog.String("requiredVersion", versionOfApexDomainSupport.String()),
			slog.String("varFile", varFile))
		return nil
	}

	if varFileJsonSchema.Template != "app" {
		slog.Debug("not updating, var file template is not 'app'")
		return nil
	}

	// Check if already transformed
	migrated, err := isMigrated(varFile)
	if err != nil {
		return fmt.Errorf("checking if YAML already is transformed: %w", err)
	}

	if migrated {
		slog.Debug("not updating, is already transformed")
		return nil
	}

	return migrate(varFile)
}

// isMigrated returns true if AlbHostRouting.Subdomain or AlbHostRouting.Apex is set. This most likely means
// that the YAML file already has support for Apex domain routing.
func isMigrated(varFile string) (bool, error) {
	args := []string{
		"eval",
		".AlbHostRouting.Subdomain != null or .AlbHostRouting.ApexDomain != null",
		varFile,
	}

	output, err := sh.Output("yq", args...)
	if err != nil {
		return false, fmt.Errorf("error checking AlbHostRouting fields: %w", err)
	}

	isAlreadyTransformed, err := strconv.ParseBool(strings.TrimSpace(output))
	if err != nil {
		return false, fmt.Errorf("error parsing yq output: %w", err)
	}

	return isAlreadyTransformed, nil
}

func migrate(varFile string) error {
	fmt.Printf("Transforming configuration to add support for Apex domain routing. File: %s\n", varFile)

	// Proceed with the transformation
	args := []string{
		"-i",
		`
    .AlbHostRouting = {
        "Enable": (.AlbHostRouting.Enable // false),
        "Internal": (.AlbHostRouting.Internal // false),
        "Subdomain": {
            "Enable": (.AlbHostRouting.Enable // false),
            "TargetGroupTargetStickiness": (.AlbHostRouting.TargetGroupTargetStickiness // false)
        },
        "ApexDomain": {
            "Enable": false,
            "TargetGroupTargetStickiness": false
        }
    } |
	del(.AlbHostRouting.TargetGroupTargetStickiness)
    `,
		varFile,
	}

	err := sh.RunV("yq", args...)
	if err != nil {
		return fmt.Errorf("error transforming YAML: %w", err)
	}

	return nil
}
